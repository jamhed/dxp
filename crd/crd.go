package crd

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"dxp/kube"
	"dxp/msg"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

type Crd struct {
	group         string
	version       string
	resource      string
	labelSelector string
	fieldSelector string
	notify        chan *msg.Msg
	stop          chan struct{}
	informer      informers.GenericInformer
}

func (crd *Crd) Notify(action string, obj interface{}) {
	mObj := obj.(metav1.Object)
	log.Debugf("CRD: %s %s %s@%s uid:%s", crd.Id(), action, mObj.GetName(), mObj.GetNamespace(), mObj.GetUID())
	crd.notify <- msg.New(action, obj)
}

func New(group, version, resource string) *Crd {
	crd := new(Crd)
	crd.group = group
	crd.version = version
	crd.resource = resource
	crd.notify = make(chan *msg.Msg)
	crd.stop = make(chan struct{})
	return crd
}

func (crd *Crd) SetLabelSelector(labels string) *Crd {
	crd.labelSelector = labels
	return crd
}

func (crd *Crd) SetFieldSelector(fields string) *Crd {
	crd.fieldSelector = fields
	return crd
}

func (crd *Crd) Id() string {
	return fmt.Sprintf("%s|%s|%s|%s|%s", crd.group, crd.version, crd.resource, crd.labelSelector, crd.fieldSelector)
}

func (crd *Crd) Stop() {
	log.Debugf("CRD: %s stop", crd.Id())

	close(crd.notify)
	close(crd.stop)
}

func (crd *Crd) tweakListOptions(opts *metav1.ListOptions) {
	opts.LabelSelector = crd.labelSelector
	opts.FieldSelector = crd.fieldSelector
}

func (crd *Crd) Watch() *Crd {
	log.Debugf("CRD: %s start", crd.Id())

	cfg, err := kube.GetConfig()
	if err != nil {
		log.Fatalf("CRD: %s could not configure kubernetes access, error:%s", crd.Id(), err)
	}
	dc, err := dynamic.NewForConfig(cfg)
	if err != nil {
		log.Fatalf("CRD: %s could not generate dynamic client for config, error:%s", crd.Id(), err)
	}
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dc, 0, metav1.NamespaceAll, crd.tweakListOptions)
	gvr := schema.GroupVersionResource{Group: crd.group, Version: crd.version, Resource: crd.resource}
	crd.informer = factory.ForResource(gvr)

	crd.informer.Informer().AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				crd.Notify("add", obj)
			},
			DeleteFunc: func(obj interface{}) {
				crd.Notify("delete", obj)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				crd.Notify("update", newObj)
			},
		},
	)
	go crd.informer.Informer().Run(crd.stop)
	return crd
}

func (crd *Crd) Notifier() chan *msg.Msg {
	return crd.notify
}

func GetByKind(kind, namespace, name string) (crd *Crd, err error) {
	err = nil
	switch kind {
	case "pod":
		crd = New("", "v1", "pods")
	case "job":
		crd = New("batch", "v1", "jobs")
	case "pvc":
		crd = New("", "v1", "persistentvolumeclaims")
	case "service":
		crd = New("", "v1", "services")
	case "workflow":
		crd = New("argoproj.io", "v1alpha1", "workflows")
	case "ingress":
		crd = New("extensions", "v1beta1", "ingresses")
	default:
		return nil, fmt.Errorf("Can't create crd by kind:%s", kind)
	}
	crd.SetFieldSelector(fmt.Sprintf("metadata.name=%s,metadata.namespace=%s", name, namespace))
	return
}
