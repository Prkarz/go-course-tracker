package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/Prkarz/course-tracker/models"
)

// Creating a Course
func Create_course(tx *sql.Tx, owner_id int, url, title string) (int, bool, error) {

	var isUserActive bool
	err := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND deleted_at IS NULL)", owner_id).Scan(&isUserActive)
	if err != nil {
		return 0, false, err
	}
	if !isUserActive {
		// If they are deleted, kick them out IMMEDIATELY before doing anything else!
		return 0, false, errors.New("blocked: course creator is deleted or does not exist")
	}
	query := `
    INSERT INTO courses (owner_id, playlist_url, title)
    SELECT $1, $2, $3
    FROM users 
    WHERE id = $1 AND deleted_at IS NULL
    ON CONFLICT (playlist_url) DO NOTHING
    RETURNING id;
    `

	var course_id int
	err = tx.QueryRow(query, owner_id, url, title).Scan(&course_id)

	if err == sql.ErrNoRows {
		// In this specific query, ErrNoRows happens for ONE of TWO reasons:
		// Reason A: ON CONFLICT triggered (The URL already exists)
		// Reason B: The SELECT returned nothing (The user is deleted/doesn't exist)

		// Let's check Reason A first by looking for the URL:
		fallback_query := "SELECT id FROM courses WHERE playlist_url = $1  "
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
func List_my_courses(contxt context.Context, db *sql.DB, userID int) ([]models.Course_data, error) {

	var reports []models.Course_data
	query := `SELECT users.id, courses.owner_id, users.username, users.email, courses.playlist_url, courses.title 
              FROM users
              LEFT JOIN courses ON users.id = courses.owner_id
              WHERE users.id = $1 AND users.deleted_at IS NULL;`
	rows, err := db.QueryContext(contxt, query, userID)
	if err != nil {
		return nil, err
	}
	// The defer statement ensures that the rows are closed after the function completes, preventing resource leaks.
	defer rows.Close()
	for rows.Next() {
		var item models.Course_data
		// Use intermediate sql.Null* types for nullable DB columns, then map to pointer fields.
		var ownerID sql.NullInt64
		var url sql.NullString
		var title sql.NullString
		err := rows.Scan(&item.ID, &ownerID, &item.Name, &item.Email, &url, &title)
		if err != nil {
			return nil, err
		}
		if ownerID.Valid {
			tmp := int(ownerID.Int64)
			item.OwnerID = &tmp
		} else {
			item.OwnerID = nil
		}
		if url.Valid {
			tmp := url.String
			item.URL = &tmp
		} else {
			item.URL = nil
		}
		if title.Valid {
			tmp := title.String
			item.Title = &tmp
		} else {
			item.Title = nil
		}
		reports = append(reports, item)
	}
	return reports, nil
}

// Starting a course
// this function is called when a user starts a course,
// it inserts a new record into the user_progress table with the user ID, course ID, initial completion percentage (0), and timestamps for when the course was started and last accessed.
func Start_course(contxt context.Context, tx *sql.Tx, userID, courseID int) error {
	current_Time := time.Now()

	var exists bool
	err := tx.QueryRowContext(contxt, `
		SELECT EXISTS (
			SELECT 1
			FROM user_progress
			WHERE user_id = $1 AND course_id = $2
		)
	`, userID, courseID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("Start skipped: User %d already enrolled in Course %d", userID, courseID)
		return errors.New("blocked: user is already enrolled in this course")
	}

	query := `
	INSERT INTO user_progress (user_id, course_id, completion_percentage, started_at, last_accessed_at)
	SELECT $1, $2, $3, $4, $5
	FROM users
	WHERE id = $1 AND deleted_at IS NULL
	AND EXISTS (SELECT 1 FROM courses WHERE id = $2 )
	`

	result, err := tx.ExecContext(contxt, query, userID, courseID, 0, current_Time, current_Time)
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
		return errors.New("blocked: user is deleted or course does not exist")
		// You can choose to return an error here, or just return nil because the end result
		// (the user being enrolled) is already true!
	}

	return nil
}

func Delete_course(tx *sql.Tx, userID, CourseID int) error {
	query := `UPDATE courses
	            SET deleted_at=CURRENT_TIMESTAMP 
				WHERE owner_id=$1 AND id=$2
				`

	result, err := tx.Exec(query, userID, CourseID)

	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("Unauthorized or course not found")
	}
	return nil
}
