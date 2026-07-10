package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Prkarz/course-tracker/repository"

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

	//Start of transaction
	tx1, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to Start transaction: %v", err)
	}
	defer tx1.Rollback()

	//CREATING USER
	user_id, isNewUser, err := repository.Create_user(tx1, "Debandhu_Mukherjee", "debandhumukherjee56@gmail.com")
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	if isNewUser {
		fmt.Printf("👤 New user created successfully! User ID: %d\n", user_id)
		err = repository.InitUserStats(tx1, user_id)
		if err != nil {
			log.Fatalf("Error Gamifying Experience : %v", err)
		} else {
			fmt.Println("Get Ready for a Gamified Experience")
		}

	} else {
		fmt.Printf("👤 User already exists! User ID: %d\n", user_id)
	}
	err = tx1.Commit()
	if err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	fmt.Println("✅ User and Stats successfully saved to the database!")

	//CREATING COURSE FK USER

	//Start of transaction
	tx2, err := db.Begin()
	if err != nil {
		log.Fatalf("Failed to Start transaction: %v", err)
	}
	defer tx2.Rollback()

	course_id, isNewCourse, err := repository.Create_course(tx2, user_id, "https://youtube.com/playlist?list=ML_VIEW_984", "How to win an election")
	if err != nil {
		log.Fatalf("Failed to create course: %v", err)
	}
	if isNewCourse {
		fmt.Printf("📚 New course created successfully! Course ID: %d\n", course_id)
	} else {
		fmt.Printf("📚 Course already exists! Course ID: %d\n", course_id)
	}
	err = tx2.Commit()
	if err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
	}

	//DELETING USERS
	err = repository.Delete_user_by_id(db, 15)
	if err != nil {
		log.Fatalln("USER DOESNOT EXISTS OR ALREADY DELETED")
	} else {
		fmt.Println("USER DELETED SUCCESSFULLY")
	}

	//LIST OF USERS
	report, err := repository.List_my_courses(db)
	if err != nil {
		log.Fatalf("Error listing courses: %v", err)
	}

	fmt.Println("\n=========================================================================================")
	fmt.Printf("%-4s | %-20s | %-30s | %s | %s\n", "ID", "USERNAME", "EMAIL ADDRESS", "ASSIGNED COURSE TRACK", "LINK")
	fmt.Println("-----------------------------------------------------------------------------------------")

	for _, value := range report {
		courseTitle := "NO COURSES ASSIGNED YET"
		playlistURL := "N/A"

		if value.Title != nil {
			courseTitle = *value.Title
		}
		if value.URL != nil {
			playlistURL = *value.URL
		}

		fmt.Printf("%-4d | %-20s | %-30s | %s | %s\n", value.ID, value.Name, value.Email, courseTitle, playlistURL)
	}
	fmt.Println("=========================================================================================")

}
