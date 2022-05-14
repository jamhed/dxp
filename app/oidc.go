package app

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"

	"github.com/coreos/go-oidc"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type BrokerDef struct {
	Name     string
	Group    string
	Version  string
	Resource string
}

func (a *App) AuthInitiate(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	state := base64.StdEncoding.EncodeToString(b)
	session, _ := a.store.Get(r, "auth-session")
	session.Values["state"] = state

	redirect := r.URL.Query().Get("redirect")
	if len(redirect) > 0 {
		unescape, err := url.PathUnescape(redirect)
		if err == nil {
			log.Debugf("AUTH: keep redirect value:%s", unescape)
			session.Values["redirect"] = unescape
		} else {
			log.Debugf("AUTH: error unescape path:%s, error:%s", redirect, err)
		}
	}

	if err = session.Save(r, w); err != nil {
		log.Errorf("save session:%s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, a.auth.Config.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func (a *App) AuthCallback(w http.ResponseWriter, r *http.Request) *appError {
	session, err := a.store.Get(r, "auth-session")
	if err != nil {
		return makeError(http.StatusInternalServerError, "Can't get session, error:%s", err.Error())
	}

	if r.URL.Query().Get("state") != session.Values["state"] {
		return makeError(http.StatusBadRequest, "Can't get proper state value from request")
	}

	if len(r.URL.Query().Get("error")) > 0 {
		return makeError(http.StatusInternalServerError, "Error reply from OIDC provider:%s", r.URL.Query().Get("error_description"))
	}

	token, err := a.auth.Config.Exchange(context.TODO(), r.URL.Query().Get("code"))
	if err != nil {
		return makeError(http.StatusUnauthorized, "Can't get token by code value")
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return makeError(http.StatusInternalServerError, "Can't get id_token value from token:%s", token)
	}

	idToken, err := a.auth.Provider.Verifier(&oidc.Config{ClientID: a.auth.Config.ClientID}).Verify(context.TODO(), rawIDToken)
	if err != nil {
		return makeError(http.StatusInternalServerError, "Can't verify token, error:%s", err)
	}

	var idTokenClaims map[string]interface{}
	if err := idToken.Claims(&idTokenClaims); err != nil {
		return makeError(http.StatusInternalServerError, "Can't decode token claims, error:%s", err)
	}

	log.Debugf("OIDC: id token claims: %s", idTokenClaims)

	redirectUrl := session.Values["redirect"]
	delete(session.Values, "redirect")

	if err = session.Save(r, w); err != nil {
		return makeError(http.StatusInternalServerError, "Can't save session, error:%s", err)
	}
	// this is to set cookie for the api domain
	redirect := `<html><head><script type="text/javascript">window.location.href="%s"</script></head><body></body></html>`
	if re, ok := redirectUrl.(string); ok {
		log.Debugf("AUTH: redirecting to:%s", re)
		fmt.Fprintf(w, redirect, re)
	} else {
		fmt.Fprintf(w, redirect, a.args.UIRootURL())
	}
	return nil
}

func (a *App) Logout(w http.ResponseWriter, r *http.Request) {
	session, err := a.store.Get(r, "auth-session")
	if err != nil {
		http.Redirect(w, r, a.args.UIRootURL(), http.StatusFound)
		return
	}
	delete(session.Values, "state")
	delete(session.Values, "auth-session")
	delete(session.Values, "profile")
	session.Options.MaxAge = -1
	if err = session.Save(r, w); err != nil {
		log.Errorf("Can't delete session, error:%s", err)
	}
	http.Redirect(w, r, a.args.UIRootURL(), http.StatusFound)
}

func (a *App) userInfo(token *oauth2.Token) (map[string]interface{}, error) {
	var claims map[string]interface{}
	userinfo, err := a.auth.Provider.UserInfo(context.TODO(), a.auth.Config.TokenSource(context.TODO(), token))
	if err != nil {
		log.Errorf("Can't request user info, error:%s", err)
		return nil, err
	}
	err = userinfo.Claims(&claims)
	if err != nil {
		log.Errorf("Can't decode claims, error:%s", err)
		return nil, err
	}
	return claims, nil
}
