package lkserver

import (
	"lkserver/internal/repo"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

type Config struct {
	BindAddr    string `toml:"bind_addr"`
	SessionsKey string `toml:"session_key"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr:    ":8833",
		SessionsKey: "2342340234234-2234234",
	}
}

type lkserver struct {
	repo         repo.DataProvider
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

func New(r repo.DataProvider, config *Config) *lkserver {
	logger := logrus.New()
	sessionStore := sessions.NewCookieStore([]byte(config.SessionsKey))

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

func (s *lkserver) configureRouter() {
	s.router.Use(handlers.CORS(
		//handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedOriginValidator(func(a string) bool { return true }),
		handlers.AllowedHeaders([]string{"X-Requested-With", "X-Request-ID", "Accept", "Content-Type"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}),
		handlers.AllowCredentials(),
	))

	s.HandleFunc("/session", s.handleSessionCreate()).Methods("POST")
}
