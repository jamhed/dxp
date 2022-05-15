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
	"net/url"

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

func (a *App) authMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", a.args.UIRootDomain())
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		for _, re := range bypassAuth {
			if re.MatchString(r.RequestURI) {
				log.Debugf("HTTP: skip auth, from:%s %v", r.RemoteAddr, r.RequestURI)
				next.ServeHTTP(w, r)
				return
			}
		}
		session, err := a.store.Get(r, "auth-session")
		if err != nil {
			log.Debugf("Can't get session, error:%s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		p, ok := session.Values["profile"]
		if !ok {
			log.Debugf("Not authorized, request:%s", r.RequestURI)
			http.Redirect(w, r, "/auth?redirect="+url.PathEscape(r.RequestURI), http.StatusFound)
			return
		}
		log.Debugf("HTTP: %s from:%s %s", p.(map[string]interface{})["preferred_username"], r.RemoteAddr, r.RequestURI)
		next.ServeHTTP(w, r)
	})
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
	r.Use(a.authMiddleWare)
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
