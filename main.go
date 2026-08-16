package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Prkarz/course-tracker/middleWare"

	apilayer "github.com/Prkarz/course-tracker/api_layer"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	// Database connection string
	connStr := "postgres://postgres:10102006@localhost:5432/course_tracker?sslmode=disable"
	// Open a connection to the PostgreSQL database
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("[DB_CONFIG_ERROR] Failed to configure database connection. Check driver and connection string. Error: %v", err)
	}
	defer db.Close()
	// Verify the connection to the database
	err = db.Ping()
	if err != nil {
		log.Fatalf("[DB_CONNECTION_FAILED] Cannot connect to PostgreSQL database. Verify server is running on localhost:5432 and credentials are correct. Error: %v", err)
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
	mux.HandleFunc("POST /courses/points", middleWare.JWT_Middleware(server.Points_Streak_toUpdate_Handler))
	mux.HandleFunc("POST /courses/progress", middleWare.JWT_Middleware(server.Update_progress_Handler))

	fmt.Println("🌐 Server is running on port 8080...")
	http.ListenAndServe(":8080", mux)
}
