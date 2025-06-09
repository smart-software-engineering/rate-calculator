package web

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func init() {
	// mime.AddExtensionType(".js", "application/javascript")
	// mime.AddExtensionType(".css", "text/css")
}

func NewHandler() *Handler {
	h := &Handler{
		Mux: chi.NewRouter(),
	}

	h.Use(middleware.RequestID)
	h.Use(middleware.RealIP)
	h.Use(middleware.Logger)
	h.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.

	// TODO validate this in the context of AI agent requests
	h.Use(middleware.Timeout(60 * time.Second))

	// TODO: h.Use(csrf.Protect(csrfKey, csrf.Secure(false)))
	// TODO: h.Use(sessions.LoadAndSave)
	h.Use(h.withUser)

	h.Get("/", h.Home())

	h.Get("/show", h.Show())

	// Create a route along /files that will serve contents from
	// the ./data/ folder.
	// TODO refactor this
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "assets"))
	h.FileServer("/assets/", filesDir)

	// Mount the admin sub-router
	h.Mount("/admin", adminRouter())

	return h
}

type Handler struct {
	*chi.Mux
}

func (h *Handler) Home() http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/home.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		// TODO pp, err := h.store.Posts()
		tmpl.Execute(w, nil)
	}
}

func (h *Handler) Show() http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/show.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	}
}

func (h *Handler) withUser(next http.Handler) http.Handler {
	// TODO move this to the new handler to store it inside our own Handler
	// var store = sessions.NewCookieStore([]byte("123456" /* os.Getenv("SESSION_KEY") */))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// session, _ := store.Get(r, "user_id")

		// TODO after login update it or something like that, will move that somewhere else!
		// err := session.Save(r, w)

		// TODO add the user here!
		// ctx := context.WithValue(r.Context(), "user", nil)
		next.ServeHTTP(w, r /* r.WithContext(ctx) */)
	})
}

// A completely separate router for administrator routes
func adminRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(AdminOnly)
	// r.Get("/", adminIndex)
	// r.Get("/accounts", adminListAccounts)
	return r
}

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/*
			ctx := r.Context()
				perm, ok := ctx.Value("acl.permission").(YourPermissionType)
				if !ok || !perm.IsAdmin() {
					http.Error(w, http.StatusText(403), 403)
					return
				}
		*/
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) FileServer(path string, root http.FileSystem) {
	if path == "" || path[len(path)-1] != '/' {
		log.Fatalf("FileServer path must end with '/'")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	h.Get(path+"*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
