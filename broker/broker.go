package broker

import (
	"dxp/msg"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// Broker defines broker data
type Broker struct {
	id             string
	Notifier       chan *msg.Msg
	newClients     chan chan *msg.Msg
	closingClients chan chan *msg.Msg
	cache          map[types.UID]interface{}
	clients        map[chan *msg.Msg]bool
}

func (broker *Broker) Find(name, namespace string) interface{} {
	for _, obj := range broker.cache {
		mObj := obj.(v1.Object)
		if mObj.GetName() == name && mObj.GetNamespace() == namespace {
			return obj
		}
	}
	return nil
}

func sendText(w http.ResponseWriter, flusher http.Flusher, text string) {
	w.Write([]byte(fmt.Sprintf("%s\n\n", text)))
	flusher.Flush()
}

func sendMsg(w http.ResponseWriter, flusher http.Flusher, m *msg.Msg) {
	jsonMsg, err := json.Marshal(m)
	if err != nil {
		log.Errorf("Broker can't encode message:%s", m)
		return
	}
	w.Write([]byte(fmt.Sprintf("data: %s\n\n", jsonMsg)))
	flusher.Flush()
}

// New creates broker instance
func New(id string) (broker *Broker) {
	broker = &Broker{
		id:             id,
		Notifier:       make(chan *msg.Msg),
		newClients:     make(chan chan *msg.Msg),
		closingClients: make(chan chan *msg.Msg),
		cache:          make(map[types.UID]interface{}),
		clients:        make(map[chan *msg.Msg]bool),
	}
	go broker.listen()
	return
}

// Serve forwards events to HTTP client
func (broker *Broker) Serve(w http.ResponseWriter, flusher http.Flusher) {
	messageChan := make(chan *msg.Msg)
	broker.newClients <- messageChan

	notify := w.(http.CloseNotifier).CloseNotify()
	go func() {
		<-notify
		close(messageChan)
	}()

	defer func() {
		broker.closingClients <- messageChan
	}()

	for _, obj := range broker.cache {
		sendMsg(w, flusher, msg.New("add", obj))
	}

	pingTicker := time.NewTicker(5 * time.Second)

loop:
	for {
		select {
		case <-pingTicker.C:
			sendText(w, flusher, "event: ping")
		case m, open := <-messageChan:
			if !open {
				break loop
			}
			sendMsg(w, flusher, m)
		}
	}
}

func (broker *Broker) updateCache(m *msg.Msg) {
	mObj := m.Content.(v1.Object)
	if m.Action == "delete" {
		delete(broker.cache, mObj.GetUID())
	} else {
		broker.cache[mObj.GetUID()] = mObj
	}
}

const patience time.Duration = time.Second * 1

func (broker *Broker) Id() string {
	return broker.id
}

func (broker *Broker) listen() {
	log.Debugf("Broker: %s start", broker.id)
	for {
		select {
		case s := <-broker.newClients:
			broker.clients[s] = true
			log.Debugf("Broker: %s add client, total:%d", broker.id, len(broker.clients))
		case s := <-broker.closingClients:
			delete(broker.clients, s)
			log.Debugf("Broker: %s remove client, total:%d", broker.id, len(broker.clients))
		case msg, ok := <-broker.Notifier:
			if !ok {
				log.Debugf("Broker: %s stop", broker.id)
				return
			}
			broker.updateCache(msg)
			for clientMessageChan := range broker.clients {
				select {
				case clientMessageChan <- msg:
				case <-time.After(patience):
					log.Print("Broker: %s client timeout", broker.id)
				}
			}
		}
	}
}
