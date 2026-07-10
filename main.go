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
	mux.HandleFunc("POST /courses/create", server.Course_Creation_Handler)
	mux.HandleFunc("POST /users/delete", server.User_delete_Handler)
	// Course & Gamification Routes
	mux.HandleFunc("POST /courses/start", server.Start_Course_Handler)
	mux.HandleFunc("POST /courses/points", server.Points_Streak_toUpdate_Handler)
	mux.HandleFunc("POST /courses/progress", server.Update_progress_Handler)

	fmt.Println("🌐 Server is running on port 8080...")
	http.ListenAndServe(":8080", mux)
}

// 	//Start of transaction
// 	//this is atomicity of the transaction, if any of the below fails, the whole transaction will be rolled back
// 	tx1, err := db.Begin()
// 	if err != nil {
// 		log.Fatalf("Failed to Start transaction: %v", err)
// 	}
// 	defer tx1.Rollback()

// 	//CREATING USER
// 	user_id, isNewUser, err := repository.Create_user(tx1, "Debandhu_Mukherjee", "debandhumukherjee56@gmail.com")
// 	if err != nil {
// 		log.Fatalf("Failed to create user: %v", err)
// 	}
// 	if isNewUser {
// 		fmt.Printf("👤 New user created successfully! User ID: %d\n", user_id)
// 		//InitUserStats initializes the user stats for the newly created user.
// 		// If there is an error during this process, it logs a fatal error and terminates the program.
// 		// Otherwise, it prints a message indicating that the user is ready for a gamified experience.
// 		err = repository.InitUserStats(tx1, user_id)
// 		if err != nil {
// 			log.Fatalf("Error Gamifying Experience : %v", err)
// 		} else {
// 			fmt.Println("Get Ready for a Gamified Experience")
// 		}

// 	} else {
// 		fmt.Printf("👤 User already exists! User ID: %d\n", user_id)
// 	}
// 	// Commit the transaction to save the user and stats to the database
// 	err = tx1.Commit()
// 	if err != nil {
// 		log.Fatalf("Failed to commit transaction: %v", err)
// 	}

// 	fmt.Println("✅ User and Stats successfully saved to the database!")

// 	//CREATING COURSE FK USER
// 	//Start of transaction for creating course
// 	tx2, err := db.Begin()
// 	if err != nil {
// 		log.Fatalf("Failed to Start transaction: %v", err)
// 	}
// 	defer tx2.Rollback()

// 	course_id, isNewCourse, err := repository.Create_course(tx2, user_id, "https://youtube.com/playlist?list=ML_VIEW_984", "How to win an election")
// 	if err != nil {
// 		log.Fatalf("Failed to create course: %v", err)
// 	}
// 	if isNewCourse {
// 		fmt.Printf("📚 New course created successfully! Course ID: %d\n", course_id)
// 	} else {
// 		fmt.Printf("📚 Course already exists! Course ID: %d\n", course_id)
// 	}
// 	err = tx2.Commit()
// 	if err != nil {
// 		log.Fatalf("Failed to commit transaction: %v", err)
// 	}

// 	//DELETING USERS
// 	//Delete_user_by_id deletes a user from the database based on the provided user ID.
// 	// If the user does not exist or has already been deleted, it logs a fatal error and terminates the program.
// 	// Otherwise, it prints a message indicating that the user was deleted successfully.
// 	err = repository.Delete_user_by_id(db, 15)
// 	if err != nil {
// 		log.Fatalln("USER DOESNOT EXISTS OR ALREADY DELETED")
// 	} else {
// 		fmt.Println("USER DELETED SUCCESSFULLY")
// 	}

// 	//LIST OF USERS
// 	report, err := repository.List_my_courses(db)
// 	if err != nil {
// 		log.Fatalf("Error listing courses: %v", err)
// 	}

// 	fmt.Println("\n=========================================================================================")
// 	fmt.Printf("%-4s | %-20s | %-30s | %s | %s\n", "ID", "USERNAME", "EMAIL ADDRESS", "ASSIGNED COURSE TRACK", "LINK")
// 	fmt.Println("-----------------------------------------------------------------------------------------")
// 	// Loop through the report and print each user's details along with their assigned course information
// 	for _, value := range report {
// 		// Initialize default values for course title and playlist URL
// 		courseTitle := "Untitled Course"
// 		playlistURL := "N/A"

// 		if value.Title != nil {
// 			courseTitle = *value.Title
// 		}
// 		if value.URL != nil {
// 			playlistURL = *value.URL
// 		}

// 		fmt.Printf("%-4d | %-20s | %-30s | %s | %s\n", value.ID, value.Name, value.Email, courseTitle, playlistURL)
// 	}
// 	fmt.Println("=========================================================================================")
