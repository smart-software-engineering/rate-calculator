package web

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"smart-software-engineering/rate-calculator/rates"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewHandler(rates rates.RateCalculator, staticFS embed.FS, templateFS embed.FS) *Handler {
	h := &Handler{
		Mux:            chi.NewRouter(),
		RateCalculator: rates,
		StaticFS:       staticFS,
		TemplateFS:     templateFS,
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

	// Mount the admin sub-router
	h.Mount("/admin", adminRouter())

	// Serve static files from /static directory at the root URL path
	h.ServeStatic()

	return h
}

type Handler struct {
	*chi.Mux

	rates.RateCalculator
	TemplateFS embed.FS
	StaticFS   embed.FS
}

func (h *Handler) Home() http.HandlerFunc {
	type data struct {
		LayoutData
		Schedules rates.Schedule
		Inspect   string
	}
	tmpl := template.Must(template.ParseFS(h.TemplateFS, "templates/layout.html", "templates/home.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		schedules, err := h.RateCalculator.Schedules()
		if err != nil {
			http.Error(w, "schedule unloadable", http.StatusBadRequest)
		}
		tmpl.Execute(w, data{Schedules: schedules, Inspect: fmt.Sprintf("%#v", schedules), LayoutData: LayoutData{LoggedIn: true, Admin: true}})
	}
}

type key string

const userKey key = "user"

type LayoutData struct {
	LoggedIn bool
	Admin    bool
}

func (h *Handler) withUser(next http.Handler) http.Handler {
	// TODO move this to the new handler to store it inside our own Handler
	// var store = sessions.NewCookieStore([]byte("123456" /* os.Getenv("SESSION_KEY") */))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// session, _ := store.Get(r, "user_id")

		// TODO after login update it or something like that, will move that somewhere else!
		// err := session.Save(r, w)

		// TODO add the user here!
		ctx := context.WithValue(r.Context(), userKey, nil)
		next.ServeHTTP(w, r.WithContext(ctx))
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

func (h *Handler) ServeStatic() {
	staticSub, err := fs.Sub(h.StaticFS, "static")
	if err != nil {
		log.Fatal("Failed to get static subdirectory from embed FS: ", err)
	}
	fsys := http.FS(staticSub)

	h.FileServer("/", fsys)
}
