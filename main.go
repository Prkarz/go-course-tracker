package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type course_data struct {
	ID      int
	OwnerID *int
	Name    string
	Email   string
	URL     *string
	Title   *string
}

func list_my_courses(db *sql.DB) ([]course_data, error) {
	var reports []course_data
	query := "SELECT users.id,courses.owner_id,users.username,users.email,courses.playlist_url,courses.title FROM users LEFT JOIN courses ON users.id=courses.owner_id;"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item course_data
		err := rows.Scan(&item.ID, &item.OwnerID, &item.Name, &item.Email, &item.URL, &item.Title)
		if err != nil {
			return nil, err
		}
		reports = append(reports, item)
	}
	return reports, nil
}

func create_course(db *sql.DB, owner_id int, url, title string) (int, bool, error) {
	query := "INSERT INTO courses(owner_id,playlist_url,title) VALUES($1,$2,$3) ON CONFLICT(playlist_url) DO NOTHING RETURNING id;"
	var course_id int
	err := db.QueryRow(query, owner_id, url, title).Scan(&course_id)
	if err == sql.ErrNoRows {
		fallback_query := "SELECT id FROM courses WHERE playlist_url = $1"
		err = db.QueryRow(fallback_query, url).Scan(&course_id)
		if err != nil {
			return 0, false, err // Safety check in case fallback fails
		}
		return course_id, false, err
	}
	if err != nil {
		return 0, false, err
	}

	return course_id, true, nil
}

func create_user(db *sql.DB, username, email string) (int, bool, error) {
	query := "INSERT INTO users (username, email) VALUES ($1, $2) ON CONFLICT(email)DO NOTHING RETURNING id;"
	var userId int
	err := db.QueryRow(query, username, email).Scan(&userId)
	if err == sql.ErrNoRows {
		fallback_query := "SELECT id FROM users WHERE email = $1"
		err = db.QueryRow(fallback_query, email).Scan(&userId)
		if err != nil {
			return 0, false, err // Safety check in case the fallback fails
		}
		return userId, false, err
	}
	return userId, true, nil
}

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

	//CREATING USER
	user_id, Bool, err := create_user(db, "sarnavo_chakraborty", "sarnavochakra02@gmail.com")
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	if Bool == true {
		fmt.Println("👤 User processed successfully (Created or Found)")
	}

	//CREATING COURSE FK USER
	course_id, Bool, err := create_course(db, user_id, "https://youtube.com/playlist?list=PI_MOCK_777", "NEET exam preparation v2.0")
	if err != nil {
		log.Fatalf("Failed to create course: %v", err)
	}
	if Bool == true {
		fmt.Printf("📚 Course processed successfully! ID is: %d\n", course_id)
	}

	//LIST OF USERS
	report, err := list_my_courses(db)
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
