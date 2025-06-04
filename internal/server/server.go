package server

import (
	"embed"
	"io"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/smart-software-engineering/rate-calculator/internal/session"
	"github.com/smart-software-engineering/rate-calculator/internal/storage"
	tmpl "github.com/smart-software-engineering/rate-calculator/internal/template"
)

//go:embed static
var staticFiles embed.FS

type ServerOptions struct {
	DevMode bool
}

type Server struct {
	Addr         string
	Template     tmpl.Manager
	SessionStore session.Store
	Options      *ServerOptions
}

func New(addr string, tm tmpl.Manager, ss storage.ScheduleStorage, sessionStore session.Store, options *ServerOptions) *Server {
	if options == nil {
		options = &ServerOptions{
			DevMode: false,
		}
	}

	return &Server{
		Addr:         addr,
		Template:     tm,
		SessionStore: sessionStore,
		Options:      options,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()

	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return err
	}

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		// Landing page - no user data needed
		if err := s.Template.Render(w, "index.html", tmpl.TemplateData{}); err != nil {
			http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("GET /clear-session", func(w http.ResponseWriter, r *http.Request) {
		s.SessionStore.Clear(w, r)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// Serve Svelte SPA - unified handler for /app, /app/, and /app/*
	mux.HandleFunc("GET /app", func(w http.ResponseWriter, r *http.Request) {
		if s.Options.DevMode {
			// In development, redirect to Vite dev server
			target := "http://localhost:5173"
			if r.URL.Path != "/app" {
				target += r.URL.Path
			}
			http.Redirect(w, r, target, http.StatusTemporaryRedirect)
			return
		}

		// Serve SPA index.html directly
		w.Header().Set("Content-Type", "text/html")
		indexFile, err := staticFS.Open("index.html")
		if err != nil {
			http.Error(w, "SPA not found", http.StatusNotFound)
			return
		}
		defer indexFile.Close()

		_, err = io.Copy(w, indexFile)
		if err != nil {
			log.Printf("Error serving SPA: %v", err)
		}
	})

	// Handle /app/ and /app/* routes
	mux.HandleFunc("GET /app/", func(w http.ResponseWriter, r *http.Request) {
		if s.Options.DevMode {
			// In development, redirect to Vite dev server
			target := "http://localhost:5173" + r.URL.Path
			http.Redirect(w, r, target, http.StatusTemporaryRedirect)
			return
		}

		// Handle different path patterns
		path := r.URL.Path
		filePath := strings.TrimPrefix(path, "/app/")

		if filePath == "" {
			// /app/ - serve SPA index.html (same as /app)
			w.Header().Set("Content-Type", "text/html")
			indexFile, err := staticFS.Open("index.html")
			if err != nil {
				http.Error(w, "SPA not found", http.StatusNotFound)
				return
			}
			defer indexFile.Close()

			_, err = io.Copy(w, indexFile)
			if err != nil {
				log.Printf("Error serving SPA: %v", err)
			}
			return
		}

		// Check if it's an asset file (has extension)
		if strings.Contains(filePath, ".") {
			// It's an asset file, serve it with proper MIME type
			ext := filepath.Ext(filePath)
			contentType := mime.TypeByExtension(ext)
			if contentType == "" {
				// Fallback for common web file types
				switch ext {
				case ".css":
					contentType = "text/css"
				case ".js":
					contentType = "application/javascript"
				case ".html":
					contentType = "text/html"
				case ".json":
					contentType = "application/json"
				default:
					contentType = "application/octet-stream"
				}
			}
			w.Header().Set("Content-Type", contentType)

			// Serve the file directly from embedded filesystem
			file, err := staticFS.Open(filePath)
			if err != nil {
				http.Error(w, "File not found", http.StatusNotFound)
				return
			}
			defer file.Close()

			_, err = io.Copy(w, file)
			if err != nil {
				log.Printf("Error serving file %s: %v", filePath, err)
			}
			return
		} else {
			// It's a SPA route (no extension), serve index.html
			w.Header().Set("Content-Type", "text/html")
			indexFile, err := staticFS.Open("index.html")
			if err != nil {
				http.Error(w, "SPA not found", http.StatusNotFound)
				return
			}
			defer indexFile.Close()

			_, err = io.Copy(w, indexFile)
			if err != nil {
				log.Printf("Error serving SPA: %v", err)
			}
			return
		}
	})

	// Fallback route for /assets/* -> /app/assets/* (for backwards compatibility)
	mux.HandleFunc("GET /assets/", func(w http.ResponseWriter, r *http.Request) {
		// Redirect to the proper /app/assets/ path
		newPath := "/app" + r.URL.Path
		http.Redirect(w, r, newPath, http.StatusMovedPermanently)
	})

	// API routes (future JSON API)
	mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","message":"API is working","timestamp":"` +
			time.Now().Format(time.RFC3339) + `"}`))
	})

	var handler http.Handler = mux

	server := &http.Server{
		Addr:    s.Addr,
		Handler: handler,
	}

	log.Printf("Visit http://localhost%s to view the rate calculator\n", s.Addr)

	return server.ListenAndServe()
}
