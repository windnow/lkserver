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
	repo            *repository.Repo
	reportsService  *services.ReportService
	usersService    *services.UserService
	catalogsService *services.CatalogsService
	fileStore       repository.FileProvider
	router          *mux.Router
	config          *config.Config
	logger          *logrus.Logger
	sessionStore    sessions.Store
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
	sessionStore.Options.Secure = false
	sessionStore.Options.SameSite = http.SameSiteLaxMode

	s := &lkserver{
		config:          config,
		repo:            r,
		fileStore:       fileStore,
		logger:          logger,
		sessionStore:    sessionStore,
		router:          mux.NewRouter(),
		reportsService:  services.NewReportService(r),
		usersService:    services.NewUsersService(r),
		catalogsService: services.NewCatalogsService(r),
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
	users.HandleFunc("/", s.handleGetUserList()).Methods("GET")
	users.HandleFunc("/{guid}", s.handleGetUserInfo()).Methods("GET")

	s.catalogsRoutes(private)
	s.reportsRoutes(private)

	s.router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, s.config.StaticFilesPath+"/index.html")
	})
	s.handleFileServerIfExists()
}

func (s *lkserver) catalogsRoutes(route *mux.Router) {

	catalogs := route.PathPrefix("/cat").Subrouter()
	catalogs.HandleFunc("/cato", s.handleCatoList()).Methods("GET")
	catalogs.HandleFunc("/cato/{guid}", s.handleGetCato()).Methods("GET")
	catalogs.HandleFunc("/vus", s.handleVusList()).Methods("GET")
	catalogs.HandleFunc("/vus/{guid}", s.handleGetVus()).Methods("GET")
	catalogs.HandleFunc("/orgs", s.handleOrgsList()).Methods("GET")
	catalogs.HandleFunc("/orgs/{guid}", s.handleGetOrganization()).Methods("GET")
	catalogs.HandleFunc("/devision", s.handleDevisionsList()).Methods("GET")
	catalogs.HandleFunc("/devision/{guid}", s.handleGetDevision()).Methods("GET")
	catalogs.HandleFunc("/order-source", s.handleOrderSourceList()).Methods("GET")
	catalogs.HandleFunc("/order-source/{guid}", s.handleGetOrderSource()).Methods("GET")
}

func (s *lkserver) reportsRoutes(route *mux.Router) {

	reports := route.PathPrefix("/reports").Subrouter()
	reports.HandleFunc("/", s.handleReportsList()).Methods("GET") // Список рапортов текущего пользователя х
	reports.HandleFunc("/types", s.handleGetReportTypes()).Methods("GET")
	reports.HandleFunc("/types/{guid}", s.handleGetReportType()).Methods("GET")
	reports.HandleFunc("/{type}/new", s.handleNewReport()).Methods("GET")
	reports.HandleFunc("/{type}/save", s.handleSaveReport()).Methods("POST")

	reports.HandleFunc("/{guid}", s.handleReportData()).Methods("GET") // Данные рапорта
	reports.HandleFunc("/approvals", nil).Methods("GET")               // Список рапортов для согласования
}
