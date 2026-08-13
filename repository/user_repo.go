package repository

import (
	"database/sql"
	"errors"
)

// Create_user function attempts to create a new user in the database with the provided username and email.
// It uses a SQL transaction to ensure atomicity.
// If a user with the same email already exists, it retrieves the existing user's ID instead.
// The function returns the user ID, a boolean indicating whether a new user was created, and an error if any occurred during the process.
func Create_user(tx *sql.Tx, username, email string) (int, bool, error) {
	query := "INSERT INTO users (username, email) VALUES ($1, $2) ON CONFLICT(email)DO NOTHING RETURNING id;"
	var userId int

	err := tx.QueryRow(query, username, email).Scan(&userId)
	// If the insertion fails due to a conflict (i.e., the user already exists), it will return sql.ErrNoRows.
	if err == sql.ErrNoRows {
		fallback_query := "SELECT id FROM users WHERE email = $1"
		err = tx.QueryRow(fallback_query, email).Scan(&userId)
		//Scan is used to retrieve the user ID of the existing user with the provided email.
		if err != nil {
			return 0, false, err // Safety check in case the fallback fails
		}
		return userId, false, nil
	}

	return userId, true, nil
}

// RowsAffected is used to determine how many rows were affected by the delete operation.
func Delete_user_by_id(tx *sql.Tx, userId int) error {
	query := "Update users SET deleted_at=CURRENT_TIMESTAMP Where id=$1 AND deleted_at IS NULL;"
	result, err := tx.Exec(query, userId)
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

func Login_User_byemail(db *sql.DB, email string) (int, error) {
	var resultID int
	query := `SELECT id FROM users
	        WHERE email=$1 AND deleted_at IS NULL `

	err := db.QueryRow(query, email).Scan(&resultID)
	if err == sql.ErrNoRows {
		return 0, errors.New("invalid credentials")
	}
	if err != nil {
		return 0, err
	}
	return resultID, nil
}
