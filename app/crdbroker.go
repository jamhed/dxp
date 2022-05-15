package app

import (
	"dxp/broker"
	"dxp/crd"
	"net/http"
)

type CrdBroker struct {
	crd    []*crd.Crd
	broker *broker.Broker
}

func (a *App) addBroker(sessionId, name string, broker *CrdBroker) *CrdBroker {
	nameMap, ok := a.brokers[sessionId]
	if !ok {
		nameMap = make(map[string]*CrdBroker)
		a.brokers[sessionId] = nameMap
	}
	nameMap[name] = broker
	return broker
}

func (a *App) deleteBroker(sessionId, name string) {
	nameMap, ok := a.brokers[sessionId]
	if !ok {
		return
	}
	delete(nameMap, name)
}

func (a *App) newBroker(sessionId, name string) *CrdBroker {
	return a.addBroker(sessionId, name, NewCrdBroker(name))
}

func (a *App) getBroker(sessionId, name string) *CrdBroker {
	if nameMap, ok := a.brokers[sessionId]; ok {
		return nameMap[name]
	}
	return nil
}

func (a *App) getSubsetBroker(sessionId, id string) *CrdBroker {
	if nameMap, ok := a.subset[sessionId]; ok {
		return nameMap[id]
	}
	return nil
}

func (a *App) maybeNewSubsetBroker(sessionId string, crd *crd.Crd) *CrdBroker {
	return a.maybeNewIdSubsetBroker(sessionId, crd.Id()).AddCrd(crd)
}

func (a *App) maybeNewIdSubsetBroker(sessionId string, id string) *CrdBroker {
	m, ok := a.subset[sessionId]
	if !ok {
		m = make(map[string]*CrdBroker)
		a.subset[sessionId] = m
	}
	cb, ok := m[id]
	if ok {
		return cb
	}
	cb = NewCrdBroker(id)
	m[id] = cb
	return cb
}

func NewCrdBroker(id string) *CrdBroker {
	cb := new(CrdBroker)
	cb.broker = broker.New(id)
	return cb
}

func (cb *CrdBroker) AddCrd(crd *crd.Crd) *CrdBroker {
	id := crd.Id()
	for _, i := range cb.crd {
		if id == i.Id() {
			return cb
		}
	}
	cb.crd = append(cb.crd, crd)
	go func() {
		for msg := range crd.Notifier() {
			cb.broker.Notifier <- msg
		}
	}()
	crd.Watch()
	return cb
}

func (cb *CrdBroker) Stop() {
	for _, crd := range cb.crd {
		crd.Stop()
	}
	close(cb.broker.Notifier)
}

func (cb *CrdBroker) Broker() *broker.Broker {
	return cb.broker
}

func (cb *CrdBroker) ServeHTTP(w http.ResponseWriter, r *http.Request) *appError {
	flusher, ok := w.(http.Flusher)
	if !ok {
		return makeError(http.StatusInternalServerError, "Streaming unsupported!")
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Transfer-Encoding", "identity")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()
	cb.broker.Serve(w, flusher)
	return nil
}
