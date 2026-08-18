package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Prkarz/course-tracker/middleWare"

	apilayer "github.com/Prkarz/course-tracker/api_layer"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("%v", err)
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
	fmt.Println("🚀 Successfully connected to the PostgreSQL database!")

	server := &apilayer.APIServer{DB: db}

	mux := http.NewServeMux()
	//Administration Routes
	mux.HandleFunc("POST /users/create", server.User_creation_Handler)
	mux.HandleFunc("POST /users/login", server.User_Login_Handler)
	mux.HandleFunc("POST /courses/create", middleWare.JWT_Middleware(server.Course_Creation_Handler))
	mux.HandleFunc("POST /users/delete", middleWare.JWT_Middleware(server.User_delete_Handler))
	// Course & Gamification Routes
	mux.HandleFunc("POST /courses/mycourses", middleWare.JWT_Middleware(server.List_myCourses_Handler))
	mux.HandleFunc("POST /courses/start", middleWare.JWT_Middleware(server.Start_Course_Handler))
	mux.HandleFunc("POST /courses/progress", middleWare.JWT_Middleware(server.Update_progress_Handler))

	handlerWithCORS := middleWare.CORSMiddleware(mux)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("🌐 Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":"+port, handlerWithCORS))
}
