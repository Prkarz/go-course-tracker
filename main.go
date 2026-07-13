package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	apilayer "github.com/Prkarz/course-tracker/api_layer"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	// Database connection string
	connStr := "postgres://postgres:10102006@localhost:5432/course_tracker?sslmode=disable"
	// Open a connection to the PostgreSQL database
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Error configuring the database connection:%v ", err)
	}
	defer db.Close()
	// Verify the connection to the database
	err = db.Ping()
	if err != nil {
		log.Fatalf("Database connection failed! Error: %v", err)
	}
	fmt.Println("🚀 Successfully connected to the PostgreSQL database!")

	server := &apilayer.APIServer{DB: db}

	mux := http.NewServeMux()
	//Administration Routes
	mux.HandleFunc("POST /users/create", server.User_creation_Handler)
	mux.HandleFunc("POST /users/login", server.User_Login_Handler)
	mux.HandleFunc("POST /courses/create", server.Course_Creation_Handler)
	mux.HandleFunc("POST /users/delete", server.User_delete_Handler)
	// Course & Gamification Routes
	mux.HandleFunc("POST /courses/mycourses", server.List_myCourses_Handler)
	mux.HandleFunc("POST /courses/start", server.Start_Course_Handler)
	mux.HandleFunc("POST /courses/points", server.Points_Streak_toUpdate_Handler)
	mux.HandleFunc("POST /courses/progress", server.Update_progress_Handler)

	fmt.Println("🌐 Server is running on port 8080...")
	http.ListenAndServe(":8080", mux)
}
