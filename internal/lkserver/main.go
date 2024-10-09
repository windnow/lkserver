package lkserver

import (
	"errors"
	"lkserver/internal/lkserver/config"
	"lkserver/internal/repository"
	"lkserver/internal/services"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

var errUnautorized = errors.New("NOT AUTORIZED")
var errNotFound = errors.New("NOT Found")

type lkserver struct {
	repo           *repository.Repo
	reportsService *services.ReportService
	usersService   *services.UserService
	fileStore      repository.FileProvider
	router         *mux.Router
	config         *config.Config
	logger         *logrus.Logger
	sessionStore   sessions.Store
}

func (s *lkserver) Start() error {
	log.Printf("Starting server: %s\n\t-> %s", s.config.BindAddr, s.config.SessionsKey)
	defer log.Printf("Server stopped")

	return http.ListenAndServe(s.config.BindAddr, s)

}

func New(r *repository.Repo, fileStore repository.FileProvider, config *config.Config) *lkserver {
	logger := logrus.New()
	sessionStore := sessions.NewCookieStore([]byte(config.SessionsKey))
	sessionStore.Options.MaxAge = config.SessionMaxAge

	s := &lkserver{
		config:         config,
		repo:           r,
		fileStore:      fileStore,
		logger:         logger,
		sessionStore:   sessionStore,
		router:         mux.NewRouter(),
		reportsService: services.NewReportService(r),
		usersService:   services.NewUsersService(r),
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
	private.HandleFunc("/edu/{iin}", s.handleEducationByIIN()).Methods("GET")

	users := private.PathPrefix("/users").Subrouter()
	users.HandleFunc("/{guid}", s.handleGetUserInfo()).Methods("GET")

	reports := private.PathPrefix("/reports").Subrouter()
	reports.HandleFunc("/types", s.handleGetReportTypes()).Methods("GET")
	reports.HandleFunc("/types/{guid}", s.handleGetReportType()).Methods("GET")
	reports.HandleFunc("/{type}/save", s.handleSaveReport()).Methods("POST")
	reports.HandleFunc("/", s.handleReportsList()).Methods("GET") //  Список рапортов текущего пользователя х
	reports.HandleFunc("/approvals", nil).Methods("GET")          // Список рапортов для согласования

	s.router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, s.config.StaticFilesPath+"/index.html")
	})
	s.handleFileServerIfExists()
}
