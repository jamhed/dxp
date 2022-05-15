package app

import (
	"dxp/crd"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *App) watchObject(sid string, profile map[string]interface{}, w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	kind, name, namespace := v["kind"], v["name"], v["namespace"]
	// if err := authObj(kind, name, namespace, profile["preferred_email"].(string), profile["groups"].([]string)); err != nil {
	//	return err
	//}
	crd, err := crd.GetByKind(kind, namespace, name)
	if err != nil {
		return makeError(http.StatusInternalServerError, "Can't create watcher %s/%s/%s, err:%s", kind, namespace, name, err)
	}
	return a.maybeNewSubsetBroker(sid, crd).ServeHTTP(w, r)
}

func (a *App) watchKind(sid string, profile map[string]interface{}, w http.ResponseWriter, r *http.Request) *appError {
	vars := mux.Vars(r)
	kind := vars["kind"]
	broker := a.getBroker(sid, kind)
	if broker == nil {
		return makeError(http.StatusNotFound, "Can't find broker by kind:%s", kind)
	}
	return broker.ServeHTTP(w, r)
}
