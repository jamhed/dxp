package app

import (
	"dxp/args"
	"dxp/auth"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"regexp"
	"sync"

	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type App struct {
	args   *args.Args
	store  *sessions.FilesystemStore
	auth   *auth.Auth
	wg     sync.WaitGroup
	groups map[string]string
	ver    map[string]interface{}
}

func New() *App {
	a := new(App)
	a.args = args.New().LogLevel()
	a.store = sessions.NewFilesystemStore("sessions", []byte(a.args.SessionKey()))
	a.store.MaxLength(0)
	a.store.Options.Domain = a.args.BaseDomain()
	gob.Register(map[string]interface{}{})
	go a.Serve()
	a.auth = auth.New(a.args.OIDC())
	a.ver = make(map[string]interface{})
	return a
}

func (a *App) SetVersion(version, commit, builddate string) *App {
	a.ver["version"] = version
	a.ver["commit"] = commit
	a.ver["builddate"] = builddate
	return a
}

var bypassAuth []*regexp.Regexp = []*regexp.Regexp{
	regexp.MustCompile("^/auth"),
	regexp.MustCompile("^/logout$"),
	regexp.MustCompile("^/callback.*$"),
	regexp.MustCompile("^/ui/.*$"),
}

func makeError(code int, format string, args ...interface{}) *appError {
	return &appError{Error: fmt.Sprintf(format, args...), Code: code}
}

func makeStringError(err error) *appError {
	return &appError{Error: fmt.Sprintf("%s", err)}
}

type appError struct {
	Error string
	Code  int
}

type httpErrHandler func(w http.ResponseWriter, r *http.Request) *appError

func (fn httpErrHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		log.Debug(err.Error)
		http.Error(w, err.Error, err.Code)
	}
}

func (a *App) Profile(w http.ResponseWriter, r *http.Request) {
	session, err := a.store.Get(r, "auth-session")
	if err != nil {
		log.Debugf("Can't get session, error:%s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	obj, err := json.Marshal(session.Values["profile"])
	w.Write([]byte(obj))
}

func (a *App) Token(w http.ResponseWriter, r *http.Request) {
	session, err := a.store.Get(r, "auth-session")
	if err != nil {
		log.Debugf("Can't get session, error:%s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	obj, err := json.Marshal(session.Values["access_token"])
	if err != nil {
		log.Debugf("Can't get token from session, error:%s", err)
	}
	w.Write([]byte(obj))
}

func (a *App) Serve() {
	a.wg.Add(1)
	defer a.wg.Done()
	bindAddr := fmt.Sprintf("%s:%d", a.args.BindAddr(), a.args.Port())
	log.Infof("HTTP: at %s static:%s start", bindAddr, a.args.Dir())
	r := mux.NewRouter()
	r.PathPrefix("/ui/").Handler(http.StripPrefix("/ui/", http.FileServer(http.Dir(a.args.Dir()))))
	r.HandleFunc("/auth", a.AuthInitiate)
	r.HandleFunc("/profile", a.Profile)
	r.HandleFunc("/token", a.Token)
	r.Handle("/callback", httpErrHandler(a.AuthCallback))
	r.HandleFunc("/logout", a.Logout)
	srv := &http.Server{
		Handler: r,
		Addr:    bindAddr,
	}
	srv.SetKeepAlivesEnabled(true)
	log.Fatal(srv.ListenAndServe())
}

func (a *App) Wait() {
	a.wg.Wait()
}
