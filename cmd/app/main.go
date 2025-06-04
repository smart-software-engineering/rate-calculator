package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/smart-software-engineering/rate-calculator/internal/server"
	"github.com/smart-software-engineering/rate-calculator/internal/session"
	"github.com/smart-software-engineering/rate-calculator/internal/storage/memory"
	"github.com/smart-software-engineering/rate-calculator/internal/template"
)

func main() {
	addr := flag.String("addr", ":8080", "HTTP server address")
	devModeFlag := flag.Bool("dev", false, "Run in development mode (less secure, more convenient for local development)")
	flag.Parse()

	// Check environment variable as an alternative way to enable dev mode
	devMode := *devModeFlag

	log.Print("Dev Mode: ", devMode)

	tmplManager, err := template.NewManager()
	if err != nil {
		log.Fatalf("Failed to initialize template manager: %v", err)
	}

	schedulesPath := filepath.Join("data", "schedules")
	scheduleStorage, err := memory.NewScheduleStorage(schedulesPath)
	if err != nil {
		log.Printf("Warning: Failed to initialize schedule storage: %v", err)
		log.Println("Continuing with empty schedule storage")
	}

	authKey := getSessionKey(devMode)

	// Configure cookie options based on dev mode
	cookieOptions := &session.CookieOptions{
		Secure:   !devMode, // Secure in production, not in dev
		SameSite: getSameSiteMode(devMode),
		HttpOnly: true,
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
	}

	sessionStore := session.NewCookieStore([]byte(authKey), nil, cookieOptions)

	// Create server with dev mode configuration
	srv := server.New(*addr, tmplManager, scheduleStorage, sessionStore, &server.ServerOptions{
		DevMode: devMode,
	})

	log.Printf("Server starting on %s\n", *addr)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func getSessionKey(devMode bool) string {
	// In development mode, use a fixed key
	if devMode {
		envKey := os.Getenv("COOKIE_SESSION_KEY")
		if envKey != "" {
			return envKey
		}

		// Provide a fixed dev-only key
		log.Println("WARNING: Using fixed development session key. DO NOT use in production!")
		return "development-only-insecure-session-key-32bytes"
	}

	// In production mode, require a proper key from environment
	envKey := os.Getenv("COOKIE_SESSION_KEY")
	if envKey == "" {
		log.Fatal("Production mode requires COOKIE_SESSION_KEY environment variable to be set")
	}

	return envKey
}

// Get the appropriate SameSite mode based on environment
func getSameSiteMode(devMode bool) session.SameSiteMode {
	if devMode {
		return session.SameSiteLaxMode
	}
	return session.SameSiteStrictMode
}
