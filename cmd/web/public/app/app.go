package app

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/marchelbling/webbox/pkg/middleware"
	"github.com/rs/cors"
)

//go:embed static
var staticRootFS embed.FS

// Application defines the API.
type Application struct {
	router       *mux.Router
	staticServer http.Handler
}

// New creates a new application instance.
// Static assets located in the static folder are embedded
// and will be served without the static prefix.
func New() (*Application, error) {
	router := mux.NewRouter()
	router.Use(
		mux.MiddlewareFunc(middleware.Recover()), // catch panic
		mux.MiddlewareFunc(cors.Default().Handler),
	)

	staticFS, err := fs.Sub(staticRootFS, "static")
	if err != nil {
		log.Fatal(err)
	}
	staticServer := http.FileServer(http.FS(staticFS))

	return &Application{
		staticServer: staticServer,
		router:       mux.NewRouter(),
	}, nil
}

// Routes registers the application routes.
func (a *Application) Routes() http.Handler {
	a.router.HandleFunc("/", Main)
	a.router.HandleFunc("/register", Register)
	a.router.PathPrefix("/").Handler(http.StripPrefix("/static/", a.staticServer))
	return a.router
}

// Main is the handler that serves the main page.
func Main(w http.ResponseWriter, req *http.Request) {
}

// Register is the handler that registers a new user.
func Register(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}
