package repository

import (
	"database/sql"
	"errors"
)

func Create_user(tx *sql.Tx, username, email string) (int, bool, error) {
	query := "INSERT INTO users (username, email) VALUES ($1, $2) ON CONFLICT(email)DO NOTHING RETURNING id;"
	var userId int

	err := tx.QueryRow(query, username, email).Scan(&userId)
	if err == sql.ErrNoRows {
		fallback_query := "SELECT id FROM users WHERE email = $1"
		err = tx.QueryRow(fallback_query, email).Scan(&userId)
		if err != nil {
			return 0, false, err // Safety check in case the fallback fails
		}
		return userId, false, nil
	}

	return userId, true, nil
}

func Delete_user_by_id(db *sql.DB, userId int) error {
	query := "DELETE FROM users Where id=$1;"
	result, err := db.Exec(query, userId)
	if err != nil {
		return err
	}
	row_affecred, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if row_affecred == 0 {
		return errors.New("USER NOT FOUND")
	}
	return nil
}
