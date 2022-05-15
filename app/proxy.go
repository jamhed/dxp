package app

import (
	"dxp/constant"
	"dxp/kube"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func authObj(kind, name, namespace string, owner string, groups []string) *appError {
	obj, err := kube.GetByKind(kind, name, namespace)
	if err != nil {
		return makeError(http.StatusNotFound, "Can't find object by kind %s/%s/%s, err:%s", kind, namespace, name, err)
	}
	if !Authorize(obj, owner, groups) {
		return makeError(http.StatusForbidden, "Not authorized to access object %s/%s/%s", kind, namespace, name)
	}
	return nil
}

func Authorize(obj metav1.Object, owner string, groups []string) bool {
	labels := obj.GetLabels()
	if groupLabel, ok := labels[constant.GroupLabel]; ok && len(groupLabel) > 0 {
		for _, group := range groups {
			if group == groupLabel {
				return true
			}
		}
	}
	if idLabel, ok := labels[constant.IdLabel]; ok && len(idLabel) > 0 && idLabel == owner {
		return true
	}
	return false
}

func serviceKey(auth, namespace, name string) string {
	return fmt.Sprintf("%s-%s-%s", namespace, name, auth)
}

func maybeGetBearer(r *http.Request) string {
	reqToken := r.Header.Get("Authorization")
	if len(reqToken) == 0 {
		return ""
	}
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) != 2 {
		log.Warnf("Invalid authorization header value:%s", reqToken)
		return ""
	}
	return splitToken[1]
}

func (a *App) proxyDomain(w http.ResponseWriter, r *http.Request) *appError {
	re := regexp.MustCompile(`^(.+?)\.(.+?)\.(.+?)\.svc\.cluster`)
	tmp := re.FindAllStringSubmatch(r.Host, -1)
	if !(len(tmp) == 1 && len(tmp[0]) == 4) {
		return makeError(http.StatusForbidden, "Proxy: can't parse host:%s", r.Host)
	}
	vars := mux.Vars(r)
	r = mux.SetURLVars(r, map[string]string{
		"namespace": tmp[0][1],
		"name":      tmp[0][2],
		"port":      tmp[0][3],
		"rest":      vars["rest"],
		"keep":      "true",
	})
	return a.proxyService(w, r)
}

func (a *App) proxyService(w http.ResponseWriter, r *http.Request) *appError {
	v := mux.Vars(r)
	name, namespace, port, rest, keep := v["name"], v["namespace"], v["port"], v["rest"], v["keep"]
	var session *sessions.Session
	var svc *v1.Service
	var err error

	// verify service exists
	svc, err = kube.GetService(name, namespace)
	if err != nil {
		return makeError(http.StatusForbidden, "Proxy: no service %s/%s, access denied, error:%s", namespace, name, err)
	}

	session, err = a.store.Get(r, "auth-session")
	if err != nil {
		return makeError(http.StatusInternalServerError, "Can't get session, error:%s", err)
	}

	pf, ok := session.Values["profile"].(map[string]interface{})
	if !ok {
		fwd := fmt.Sprintf("%s://%s/%s", r.Header.Get("X-Forwarded-Proto"), r.Header.Get("X-Forwarded-Host"), rest)
		redirect := fmt.Sprintf("https://%s/auth?redirect=%s", a.args.BaseDomain(), url.PathEscape(fwd))
		log.Debugf("Proxy: %s/%s, no profile, access denied, redirecting to:%s", namespace, name, redirect)
		http.Redirect(w, r, redirect, http.StatusFound)
		return nil
	}

	if !Authorize(svc, profileUser(pf), profileGroups(pf)) {
		return makeError(http.StatusForbidden, "Proxy: %s/%s, no auth, access denied", namespace, name)
	}

	schema := "http"
	if port == "443" {
		schema = "https"
	}

	target := fmt.Sprintf("%s://%s.%s.svc.cluster.local:%s", schema, name, namespace, port)

	if keep != "true" {
		rest = fmt.Sprintf("/proxy/%s/%s/%s/%s", namespace, name, port, rest)
	}

	log.Debugf("Proxy from:%s target:%s url:%s", r.RemoteAddr, target, rest)

	url, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(url)
	r.URL.Host = url.Host
	r.URL.Scheme = url.Scheme
	r.URL.Path = rest
	r.Header.Set("Host", r.Host)
	r.Header.Set("X-Forwarded-Host", r.Host)
	r.Header.Set("X-Forwarded-Proto", "https")
	proxy.ServeHTTP(w, r)
	return nil
}
