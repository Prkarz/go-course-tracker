package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	aiintegration "github.com/Prkarz/course-tracker/Ai_integration"
	apilayer "github.com/Prkarz/course-tracker/api_layer"
	"github.com/Prkarz/course-tracker/middleWare"
	"github.com/Prkarz/course-tracker/repository"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

type spaHandler struct {
	staticPath string
	indexPath  string
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// If an unmatched API route is requested, return JSON 404 instead of serving HTML
	if strings.HasPrefix(r.URL.Path, "/api/") || r.URL.Path == "/api" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":{"code":"404_NOT_FOUND","message":"API endpoint not found"}}`))
		return
	}

	// For non-GET / non-HEAD requests that reached the SPA fallback, return 404 JSON
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":{"code":"404_METHOD_NOT_ALLOWED","message":"Resource not found"}}`))
		return
	}

	path := filepath.Join(h.staticPath, r.URL.Path)

	// Check if the file exists and is not a directory
	fi, err := os.Stat(path)
	if os.IsNotExist(err) || fi.IsDir() {
		// File does not exist, serve the React index.html
		fallbackPath := filepath.Join(h.staticPath, h.indexPath)
		if _, err := os.Stat(fallbackPath); os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, fallbackPath)
		return
	}
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading from system environment variables")
	}
	// Database connection string
	connStr := os.Getenv("DB_URL")
	// Open a connection to the PostgreSQL database
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("[DB_CONFIG_ERROR] Failed to configure database connection. Check driver and connection string. Error: %v", err)
	}
	defer db.Close()
	// Verify the connection to the database
	err = db.Ping()
	if err != nil {
		log.Fatalf("[DB_CONNECTION_FAILED] Cannot connect to PostgreSQL database. Verify server is running and credentials are correct. Error: %v", err)
	}
	fmt.Println("[DB_CONNECTED] Successfully connected to PostgreSQL.")

	// Auto-ensure all required tables and schema migrations
	if err := repository.EnsureAllTables(db); err != nil {
		log.Fatalf("[DB_SCHEMA_ERROR] Failed to initialize database tables: %v", err)
	}

	client, err := aiintegration.Client_Ai_Init()
	if err != nil {
		log.Printf("[AI_CONFIG_WARNING] Gemini client initialization warning: %v", err)
	}

	server := &apilayer.APIServer{DB: db, AIClient: client}

	mux := http.NewServeMux()

	// Helper to register both /api/... and /... routes to support all frontend configurations
	registerPost := func(pattern string, handler http.HandlerFunc) {
		mux.HandleFunc("POST /api"+pattern, handler)
		mux.HandleFunc("POST "+pattern, handler)
	}
	registerGet := func(pattern string, handler http.HandlerFunc) {
		mux.HandleFunc("GET /api"+pattern, handler)
		mux.HandleFunc("GET "+pattern, handler)
	}

	// Administration Routes
	registerPost("/users/create", server.User_creation_Handler)
	registerPost("/users/login", server.User_Login_Handler)
	registerPost("/courses/create", middleWare.JWT_Middleware(server.Course_Creation_Handler))
	registerPost("/users/delete", middleWare.JWT_Middleware(server.User_delete_Handler))

	// Course & Gamification Routes
	registerPost("/courses/mycourses", middleWare.JWT_Middleware(server.List_myCourses_Handler))
	registerPost("/courses/detail", middleWare.JWT_Middleware(server.Course_Detail_Handler))
	registerPost("/courses/delete", middleWare.JWT_Middleware(server.Course_Delete_Handler))
	registerPost("/courses/start", middleWare.JWT_Middleware(server.Start_Course_Handler))
	registerPost("/courses/progress", middleWare.JWT_Middleware(server.Update_progress_Handler))
	registerGet("/stats/me", middleWare.JWT_Middleware(server.User_stats_Handler))
	registerPost("/courses/videos/viewed", middleWare.JWT_Middleware(server.Video_Viewed_Handler))

	// Determine static path (support running from root or frontend folder)
	staticDir := os.Getenv("STATIC_DIR")
	if staticDir == "" {
		if _, err := os.Stat("frontend/dist"); err == nil {
			staticDir = "frontend/dist"
		} else if _, err := os.Stat("dist"); err == nil {
			staticDir = "dist"
		} else {
			staticDir = "frontend/dist"
		}
	}

	spa := spaHandler{staticPath: staticDir, indexPath: "index.html"}
	mux.Handle("/", spa)

	handlerWithCORS := middleWare.CORSMiddleware(mux)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("[SERVER_STARTED] Server is running on port %s.\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handlerWithCORS))
}
