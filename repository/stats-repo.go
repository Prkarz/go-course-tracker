package repository

import (
	"database/sql"
	"errors"
	"log"
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
	if pointsToAdd < 0 {
		return errors.New("points to add cannot be negative")
	}

	query := `UPDATE user_stats 
	SET total_points = COALESCE(user_stats.total_points, 0) + $1
	FROM users
	WHERE user_stats.user_id = $2
	  AND users.id = $2
	  AND users.deleted_at IS NULL`
	log.Printf("[POINTS_UPDATE] user_id=%d points=%d", userID, pointsToAdd)
	result, err := tx.Exec(query, pointsToAdd, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("blocked: user is deleted or does not exist")
	}
	return nil
}

// Streak_update updates the streak count for a user based on their last active date.
// If the user was active yesterday, the streak is incremented; if they missed a day, the streak resets to 1.
// The last active date is updated to the current date.
func Streak_update(tx *sql.Tx, userID int) error {
	query := `UPDATE user_stats 
    SET streak_count = CASE
        -- 1. First time ever? Start at 1
        WHEN last_active_date IS NULL THEN 1 
        
        -- 2. Logged in yesterday? Increment it
        WHEN CURRENT_DATE - last_active_date::DATE = 1 THEN COALESCE(streak_count, 0) + 1 
        
        -- 3. Missed a day? Reset to 1
        WHEN CURRENT_DATE - last_active_date::DATE > 1 THEN 1 
        
        -- 4. Same day login? Keep current streak, but ensure it's at least 1 (not 0)
        ELSE COALESCE(NULLIF(streak_count, 0), 1)
    END,
    last_active_date = CURRENT_DATE
    WHERE user_id = $1 AND EXISTS (
        SELECT 1 FROM users 
        WHERE users.id = $1 
        AND users.deleted_at IS NULL
    )`
	result, err := tx.Exec(query, userID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("blocked: user is deleted or does not exist")
	}
	return nil
}

// Update_progress updates the completion percentage and last accessed date for a user's progress in a specific course.
// It takes the user ID, course ID, and the new completion percentage as parameters.
// The function updates the corresponding record in the user_progress table.
func Update_progress(tx *sql.Tx, userID, courseID, newPercentage int) (float64, error) {
	if newPercentage < 0 {
		return 0, errors.New("percentage increment cannot be negative")
	}

	var completion_percentage float64
	query := `UPDATE user_progress 
              SET completion_percentage = LEAST(COALESCE(completion_percentage, 0) + $1,100), 
                  last_accessed_at = CURRENT_TIMESTAMP 
              WHERE user_id = $2 AND course_id = $3 
			  AND EXISTS (
        SELECT 1 FROM users 
        WHERE users.id = $2 
        AND users.deleted_at IS NULL
    )AND EXISTS (
        SELECT 1 FROM courses 
        WHERE courses.id = $3 
        AND courses.deleted_at IS NULL
    )RETURNING completion_percentage`

	err := tx.QueryRow(query, newPercentage, userID, courseID).Scan(&completion_percentage)
	if err == sql.ErrNoRows {
		return 0, errors.New("no progress row was updated")
	}
	if err != nil {
		return 0, err
	}

	return completion_percentage, nil
}
