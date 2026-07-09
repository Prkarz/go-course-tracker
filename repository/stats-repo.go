package repository

import (
	"database/sql"
)

func InitUserStats(tx *sql.Tx, userID int) error {
	query := "INSERT INTO user_stats(user_id) VALUES($1)"
	_, err := tx.Exec(query, userID)
	return err
}

func Points_update(tx *sql.Tx, userID, pointsToAdd int) error {
	query := "UPDATE user_stats SET total_points = total_points + $1 WHERE user_id = $2"
	_, err := tx.Exec(query, pointsToAdd, userID)
	if err != nil {
		return err
	}
	return nil
}

func Streak_update(tx *sql.Tx, userID int) error {
	query := `UPDATE user_stats 
	SET streak_count=CASE
	WHEN CURRENT_DATE-last_active_date::DATE+1 THEN streak_count + 1
	WHEN CURRENT_DATE - last_active_date::DATE > 1 THEN 1
	ELSE streak_count 
    END,
	last_active_date=CURRENT_DATE
	WHERE user_id = $1`
	_, err := tx.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}
