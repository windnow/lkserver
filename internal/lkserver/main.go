package lkserver

import (
	"errors"
	"lkserver/internal/repository"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

var errUnautorized = errors.New("NOT AUTORIZED")
var errNotFound = errors.New("NOT Found")

type Config struct {
	BindAddr        string `toml:"bind_addr"`
	SessionsKey     string `toml:"session_key"`
	StaticFilesPath string `toml:"files"`
	SessionMaxAge   int    `toml:"session_max_age"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr:        ":8833",
		SessionsKey:     "2342340234234-2234234",
		StaticFilesPath: "data",
		SessionMaxAge:   60 * 30,
	}
}

type lkserver struct {
	repo         *repository.Repo
	router       *mux.Router
	config       *Config
	logger       *logrus.Logger
	sessionStore sessions.Store
}

func (s *lkserver) Start() error {
	log.Printf("Starting server: %s\n\t-> %s", s.config.BindAddr, s.config.SessionsKey)
	defer log.Printf("Server stopped")

	return http.ListenAndServe(s.config.BindAddr, s)

}

func New(r *repository.Repo, config *Config) *lkserver {
	logger := logrus.New()
	sessionStore := sessions.NewCookieStore([]byte(config.SessionsKey))
	sessionStore.Options.MaxAge = config.SessionMaxAge

	s := &lkserver{
		config:       config,
		repo:         r,
		logger:       logger,
		sessionStore: sessionStore,
		router:       mux.NewRouter(),
	}
	s.configureRouter()
	return s
}

func (s *lkserver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *lkserver) HandleFunc(path string, f func(http.ResponseWriter,
	*http.Request)) *mux.Route {
	return s.router.HandleFunc(path, f)
}

func (s *lkserver) Use(mwf mux.MiddlewareFunc) {
	s.router.Use(mwf)
}

func (s *lkserver) PathPrefix(tpl string) *mux.Route {
	return s.router.PathPrefix(tpl)
}

func (s *lkserver) configureRouter() {

	s.Use(s.setRequestID) // Присвоим уникальный идентификатор каждому запросу
	s.Use(s.authUser)
	s.Use(s.logRequest)

	s.Use(handlers.CORS(
		//handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedOriginValidator(func(a string) bool { return true }),
		handlers.AllowedHeaders([]string{"X-Requested-With", "X-Request-ID", "Accept", "Content-Type"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}),
		handlers.AllowCredentials(),
	))

	s.HandleFunc("/session", s.handleSessionCreate()).Methods("POST")
	s.HandleFunc("/wai", s.handleWhoAmI()).Methods("GET")

	private := s.PathPrefix("/i").Subrouter()
	private.Use(s.checkUser)
	private.HandleFunc("/destroy", s.handleSessionDestroy())
	private.HandleFunc("/file/{id}", s.handleFile()).Methods("GET")
	private.HandleFunc("/ind/{iin}", s.handleIndividualsByIIN()).Methods("GET")

	s.router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, s.config.StaticFilesPath+"/index.html")
	})
	s.handleFileServerIfExists()
}
