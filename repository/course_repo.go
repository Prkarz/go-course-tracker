package repository

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/Prkarz/course-tracker/models"
)

// Creating a Course
func Create_course(tx *sql.Tx, owner_id int, url, title string) (int, bool, error) {
	// 1. The Correct Ghost-Proof Insert for COURSES
	query := `
    INSERT INTO courses (created_by, playlist_url, title)
    SELECT $1, $2, $3
    FROM users 
    WHERE id = $1 AND deleted_at IS NULL
    ON CONFLICT (playlist_url) DO NOTHING
    RETURNING id;
    `

	var course_id int
	err := tx.QueryRow(query, owner_id, url, title).Scan(&course_id)

	if err == sql.ErrNoRows {
		// In this specific query, ErrNoRows happens for ONE of TWO reasons:
		// Reason A: ON CONFLICT triggered (The URL already exists)
		// Reason B: The SELECT returned nothing (The user is deleted/doesn't exist)

		// Let's check Reason A first by looking for the URL:
		fallback_query := "SELECT id FROM courses WHERE playlist_url = $1"
		err = tx.QueryRow(fallback_query, url).Scan(&course_id)

		if err == nil {
			// Success! The course already existed. Return its ID and false (not newly created)
			return course_id, false, nil
		}

		if err == sql.ErrNoRows {
			// Reason B confirmed! The URL doesn't exist, which means the user must be deleted.
			return 0, false, errors.New("blocked: course creator is deleted or does not exist")
		}

		// Some other database error occurred
		return 0, false, err
	}

	if err != nil {
		return 0, false, err
	}

	// Success! A brand new course was created.
	return course_id, true, nil
}

// List of all Courses
func List_my_courses(db *sql.DB, userID int) ([]models.Course_data, error) {
	var reports []models.Course_data
	query := `SELECT users.id, courses.owner_id, users.username, users.email, courses.playlist_url, courses.title 
              FROM users
              LEFT JOIN courses ON users.id = courses.owner_id
              WHERE users.id = $1 AND users.deleted_at IS NULL;`
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	// The defer statement ensures that the rows are closed after the function completes, preventing resource leaks.
	defer rows.Close()
	for rows.Next() {
		var item models.Course_data
		//struct is used to hold the data for each row returned by the query.
		//The Scan method maps the columns from the query result to the fields of the struct.
		//When multiple users are present, the loop iterates through each row,
		// creating a new struct instance for each user and appending it to the reports slice.
		err := rows.Scan(&item.ID, &item.OwnerID, &item.Name, &item.Email, &item.URL, &item.Title)
		if err != nil {
			return nil, err
		}
		reports = append(reports, item)
	}
	return reports, nil
}

// Starting a course
// this function is called when a user starts a course,
// it inserts a new record into the user_progress table with the user ID, course ID, initial completion percentage (0), and timestamps for when the course was started and last accessed.
func Start_course(tx *sql.Tx, userID, courseID int) error {
	current_Time := time.Now()

	// Combines SELECT block and ON CONFLICT
	query := `
    INSERT INTO user_progress (user_id, course_id, completion_percentage, started_at, last_accessed_at)
    SELECT $1, $2, $3, $4, $5
    FROM users 
    WHERE id = $1 AND deleted_at IS NULL
    ON CONFLICT (user_id, course_id) DO NOTHING;
    `

	result, err := tx.Exec(query, userID, courseID, 0, current_Time, current_Time)
	if err != nil {
		return err
	}

	// Check the result
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If rows == 0, it means EITHER the user is deleted, OR they already started this course.
	// Both are safe states, so we don't necessarily need to return a fatal server error.
	if rows == 0 {
		log.Printf("Start skipped: User %d is deleted or already enrolled in Course %d", userID, courseID)
		// You can choose to return an error here, or just return nil because the end result
		// (the user being enrolled) is already true!
	}

	return nil
}
