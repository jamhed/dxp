package args

import (
	"flag"
	"fmt"
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"
)

// Args type
type Args struct {
	verboseLevel     string
	port             int
	bindAddr         string
	dir              string
	args             []string
	oidcProvider     string
	oidcClientID     string
	oidcClientSecret string
	oidcRedirectURL  string
	oidcScopes       string
	oidcUserId       string
	uiRootURL        string
	uiRootDomain     string
	baseDomain       string
	sessionKey       string
	redisAddr        string
}

// New type
func New() *Args {
	return new(Args).Parse()
}

func getEnvOrDefault(name, def string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}
	return def
}

// Parse parameters
func (a *Args) Parse() *Args {
	flag.StringVar(&a.verboseLevel, "verbose", getEnvOrDefault("VERBOSE", "info"), "Set verbosity level")
	flag.IntVar(&a.port, "port", 8080, "Listen port")
	flag.StringVar(&a.bindAddr, "bind", os.Getenv("BIND_ADDR"), "Bind address")
	flag.StringVar(&a.dir, "dir", getEnvOrDefault("UI_DIR", "ui/dist"), "Static files folder")
	flag.StringVar(&a.oidcProvider, "oidc-provider", os.Getenv("OIDC_PROVIDER"), "OIDC provider")
	flag.StringVar(&a.oidcClientID, "oidc-client-id", os.Getenv("OIDC_CLIENT_ID"), "OIDC client id")
	flag.StringVar(&a.oidcClientSecret, "oidc-client-secret", os.Getenv("OIDC_CLIENT_SECRET"), "OIDC client secret")
	flag.StringVar(&a.oidcRedirectURL, "oidc-redirect-url", os.Getenv("OIDC_REDIRECT_URL"), "OIDC redirect url")
	flag.StringVar(&a.oidcScopes, "oidc-scopes", getEnvOrDefault("OIDC_SCOPES", "roles"), "OIDC scopes")
	flag.StringVar(&a.oidcUserId, "oidc-user-id", getEnvOrDefault("OIDC_USER_ID", "preferred_username"), "OIDC user id field")
	flag.StringVar(&a.uiRootURL, "ui-root-url", getEnvOrDefault("UI_ROOT_URL", "http://localhost:8080/ui/#/"), "UI root url for redirects")
	flag.StringVar(&a.baseDomain, "base-domain", getEnvOrDefault("BASE_DOMAIN", ""), "Base domain to base ingress on")
	flag.StringVar(&a.sessionKey, "session-key", getEnvOrDefault("SESSION_KEY", "0123456789abcdef"), "HTTP Session encryption key")

	url, _ := url.Parse(a.uiRootURL)
	a.uiRootDomain = fmt.Sprintf("%s://%s", url.Scheme, url.Host)

	flag.Parse()
	a.args = flag.Args()
	return a
}

// Dir to serve files from
func (a *Args) Dir() string {
	return a.dir
}

// BindAddr to bind to Web Server
func (a *Args) BindAddr() string {
	return a.bindAddr
}

// Port to bind to Web Server
func (a *Args) Port() int {
	return a.port
}

// UIRootURL returns UI root url
func (a *Args) UIRootURL() string {
	return a.uiRootURL
}

func (a *Args) UIRootDomain() string {
	return a.uiRootDomain
}

func (a *Args) SessionKey() string {
	return a.sessionKey
}

func (a *Args) OIDC() (string, string, string, string, string) {
	return a.oidcProvider, a.oidcClientID, a.oidcClientSecret, a.oidcRedirectURL, a.oidcScopes
}

func (a *Args) OidcUserId() string {
	return a.oidcUserId
}

func (a *Args) BaseDomain() string {
	return a.baseDomain
}

func (a *Args) Redis() string {
	return a.redisAddr
}

// LogLevel set loglevel
func (a *Args) LogLevel() *Args {
	switch a.verboseLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	return a
}
