package repository

import (
	"database/sql"
	"time"

	"github.com/Prkarz/course-tracker/models"
)

// Creating a Course
func Create_course(tx *sql.Tx, owner_id int, url, title string) (int, bool, error) {
	query := "INSERT INTO courses(owner_id,playlist_url,title) VALUES($1,$2,$3) ON CONFLICT(playlist_url) DO NOTHING RETURNING id;"
	//On conflict, the query does nothing and returns no rows, which is handled by checking for sql.ErrNoRows.
	//RETURNING id is used to get the ID of the newly created course, which is useful for further operations.
	var course_id int
	err := tx.QueryRow(query, owner_id, url, title).Scan(&course_id)
	if err == sql.ErrNoRows {
		fallback_query := "SELECT id FROM courses WHERE playlist_url = $1"
		err = tx.QueryRow(fallback_query, url).Scan(&course_id)
		if err != nil {
			return 0, false, err // Safety check in case fallback fails
		}
		return course_id, false, nil
	}
	if err != nil {
		return 0, false, err
	}

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
	query := "INSERT INTO user_progress(user_id,course_id, completion_percentage, started_at, last_accessed_at) VALUES($1,$2,$3,$4,$5)"
	_, err := tx.Exec(query, userID, courseID, 0, current_Time, current_Time)
	if err != nil {
		return err
	}
	return nil
}
