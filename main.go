package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"path/filepath"

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
	path := filepath.Join(h.staticPath, r.URL.Path)

	// Check if the file exists and is not a directory
	fi, err := os.Stat(path)
	if os.IsNotExist(err) || fi.IsDir() {
		// File does not exist, serve the React index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
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
	if err := repository.EnsureCourseVideoProgressTable(db); err != nil {
		log.Fatalf("[DB_SCHEMA_ERROR] Failed to initialize video progress table: %v", err)
	}

	client, err := aiintegration.Client_Ai_Init()
	if err != nil {
		log.Fatalf("[AI_CONFIG_ERROR] Failed to initialize Gemini client: %v", err)
	}

	server := &apilayer.APIServer{DB: db, AIClient: client}

	mux := http.NewServeMux()
	//Administration Routes
	mux.HandleFunc("POST /api/users/create", server.User_creation_Handler)
	mux.HandleFunc("POST /api/users/login", server.User_Login_Handler)
	mux.HandleFunc("POST /api/courses/create", middleWare.JWT_Middleware(server.Course_Creation_Handler))
	mux.HandleFunc("POST /api/users/delete", middleWare.JWT_Middleware(server.User_delete_Handler))
	// Course & Gamification Routes
	mux.HandleFunc("POST /api/courses/mycourses", middleWare.JWT_Middleware(server.List_myCourses_Handler))
	mux.HandleFunc("POST /api/courses/detail", middleWare.JWT_Middleware(server.Course_Detail_Handler))
	mux.HandleFunc("POST /api/courses/delete", middleWare.JWT_Middleware(server.Course_Delete_Handler))
	mux.HandleFunc("POST /api/courses/start", middleWare.JWT_Middleware(server.Start_Course_Handler))
	mux.HandleFunc("POST /api/courses/progress", middleWare.JWT_Middleware(server.Update_progress_Handler))
	mux.HandleFunc("GET /api/stats/me", middleWare.JWT_Middleware(server.User_stats_Handler))
	mux.HandleFunc("POST /api/courses/videos/viewed", middleWare.JWT_Middleware(server.Video_Viewed_Handler))

	handlerWithCORS := middleWare.CORSMiddleware(mux)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	spa := spaHandler{staticPath: "frontend/dist", indexPath: "index.html"}
	mux.Handle("/", spa)
	fmt.Printf("[SERVER_STARTED] Server is running on port %s.\n", port)
	log.Println(http.ListenAndServe(":"+port, handlerWithCORS))
}
