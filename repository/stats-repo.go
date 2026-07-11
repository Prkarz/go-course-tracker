package repository

import (
	"database/sql"
	"errors"
)

// InitUserStats initializes the user stats for a newly created user.
// It inserts a new record into the user_stats table with the provided user ID.
// This function is typically called after a new user is created to set up their initial statistics.
func InitUserStats(tx *sql.Tx, userID int) error {
	query := "INSERT INTO user_stats(user_id) VALUES($1)"
	_, err := tx.Exec(query, userID)
	return err
}

// UpdateUserStats updates the user stats for a given user.
// It takes the user ID, points to add, and a boolean indicating whether to update the streak.
// The function updates the total points and optionally updates the streak count based on the last active date.
func Points_update(tx *sql.Tx, userID, pointsToAdd int) error {
	query := `UPDATE user_stats 
	SET total_points = total_points + $1 WHERE user_id = $2
	AND EXISTS (
        SELECT 1 FROM users 
        WHERE users.id = $1 
        AND users.deleted_at IS NULL
    )`
	_, err := tx.Exec(query, pointsToAdd, userID)
	if err != nil {
		return err
	}
	return nil
}

// Streak_update updates the streak count for a user based on their last active date.
// If the user was active yesterday, the streak is incremented; if they missed a day, the streak resets to 1.
// The last active date is updated to the current date.
func Streak_update(tx *sql.Tx, userID int) error {
	query := `UPDATE user_stats 
	SET streak_count=CASE
	WHEN CURRENT_DATE - last_active_date::DATE = 1 THEN streak_count + 1
	WHEN CURRENT_DATE - last_active_date::DATE > 1 THEN 1
	ELSE streak_count 
    END,
	last_active_date=CURRENT_DATE
	WHERE user_id = $1 AND EXISTS (
        SELECT 1 FROM users 
        WHERE users.id = $1 
        AND users.deleted_at IS NULL
    )`
	_, err := tx.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}

// Update_progress updates the completion percentage and last accessed date for a user's progress in a specific course.
// It takes the user ID, course ID, and the new completion percentage as parameters.
// The function updates the corresponding record in the user_progress table.
func Update_progress(tx *sql.Tx, userID, courseID, newPercentage int) error {
	query := `UPDATE user_progress 
              SET completion_percentage = COALESCE(completion_percentage, 0) + $1, 
                  last_accessed_at = CURRENT_TIMESTAMP 
              WHERE user_id = $2 AND course_id = $3 
			  AND EXISTS (
        SELECT 1 FROM users 
        WHERE users.id = $2 
        AND users.deleted_at IS NULL
    )`

	result, err := tx.Exec(query, newPercentage, userID, courseID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err // Error retrieving the affected row count
	}
	if rowsAffected == 0 {
		return errors.New("no progress record found: has the user started this course?")
	}
	return nil
}
