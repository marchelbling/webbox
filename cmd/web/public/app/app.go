package app

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/marchelbling/webbox/pkg/middleware"
	"github.com/rs/cors"
)

//go:embed static
var staticRootFS embed.FS

//go:embed templates
var templatesRootFS embed.FS

// Application defines the API.
type Application struct {
	router       *mux.Router
	staticServer http.Handler
	templates    fs.FS
}

type HTTPFileOnlyFS struct {
	http.FileSystem
}

func (f *HTTPFileOnlyFS) Open(path string) (http.File, error) {
	file, err := f.FileSystem.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := file.Stat()
	if s.IsDir() {
		return nil, os.ErrNotExist
	}

	return file, nil
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
		log.Fatalf("static: %v", err)
	}

	templatesFS, err := fs.Sub(templatesRootFS, "templates")
	if err != nil {
		log.Fatalf("templates: %v", err)
	}

	return &Application{
		staticServer: http.FileServer(&HTTPFileOnlyFS{http.FS(staticFS)}),
		router:       mux.NewRouter(),
		templates:    templatesFS,
	}, nil
}

// Routes registers the application routes.
func (a *Application) Routes() http.Handler {
	a.router.HandleFunc("/", a.GetMain).Methods(http.MethodOptions, http.MethodGet)
	a.router.HandleFunc("/register", a.GetRegister).Methods(http.MethodOptions, http.MethodGet)
	a.router.HandleFunc("/register", a.PostRegister).Methods(http.MethodOptions, http.MethodPost)
	a.router.PathPrefix("/").Handler(http.StripPrefix("/static/", a.staticServer))
	return a.router
}

// Main is the handler that serves the main page.
func (a *Application) GetMain(w http.ResponseWriter, req *http.Request) {
	files := []string{
		"base.tmpl",
		"partials/nav.tmpl",
		"pages/home.tmpl",
	}

	ts, err := template.ParseFS(a.templates, files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

// Register is the handler that registers a new user.
func (a *Application) GetRegister(w http.ResponseWriter, req *http.Request) {
	files := []string{
		"base.tmpl",
		"partials/nav.tmpl",
		"pages/register.tmpl",
	}

	ts, err := template.ParseFS(a.templates, files...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func (a *Application) PostRegister(w http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

	for key, value := range req.PostForm {
		log.Printf("o<<< register: %s = %s", key, value)
	}
	http.Redirect(w, req, "/", http.StatusFound)
}
