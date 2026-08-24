package repository

import (
	"database/sql"
	"errors"
)

// ErrInvalidCredentials is returned when login fails due to wrong email/password.
var ErrInvalidCredentials = errors.New("invalid credentials")

// Create_or_Reactivate_User creates a new user or reactivates a soft-deleted user.
// Returns: (userID, isNewOrReactivated, isReactivated, error)
func Create_or_Reactivate_User(tx *sql.Tx, username, email, passwordHash string) (int, bool, bool, error) {
	var existingID int
	var deletedAt sql.NullTime
	err := tx.QueryRow("SELECT id, deleted_at FROM users WHERE email = $1", email).Scan(&existingID, &deletedAt)

	if err == sql.ErrNoRows {
		// Brand new user
		var newID int
		err := tx.QueryRow(`
			INSERT INTO users (username, email, password_hash, deleted_at)
			VALUES ($1, $2, $3, NULL)
			RETURNING id
		`, username, email, passwordHash).Scan(&newID)
		if err != nil {
			return 0, false, false, err
		}
		return newID, true, false, nil
	} else if err != nil {
		return 0, false, false, err
	}

	// User exists and was soft-deleted -> Reactivate user and restore data
	if deletedAt.Valid {
		var reactivatedID int
		err := tx.QueryRow(`
			UPDATE users
			SET deleted_at = NULL,
			    username = $1,
			    password_hash = $2
			WHERE id = $3
			RETURNING id
		`, username, passwordHash, existingID).Scan(&reactivatedID)
		if err != nil {
			return 0, false, false, err
		}

		// Restore all courses previously owned by this user
		_, _ = tx.Exec("UPDATE courses SET deleted_at = NULL WHERE owner_id = $1", reactivatedID)

		// Ensure user_stats record exists for the reactivated user
		_, _ = tx.Exec("INSERT INTO user_stats (user_id) VALUES ($1) ON CONFLICT (user_id) DO NOTHING", reactivatedID)

		return reactivatedID, true, true, nil
	}

	// User already exists and is active
	return existingID, false, false, nil
}

func Create_user(tx *sql.Tx, username, email, passwordHash string) (int, bool, error) {
	id, isNewOrReactivated, _, err := Create_or_Reactivate_User(tx, username, email, passwordHash)
	return id, isNewOrReactivated, err
}

func Login_User_byemail(db *sql.DB, email string) (int, string, string, error) {
	var resultID int
	var username string
	var password string
	query := `SELECT id, username, password_hash FROM users
	        WHERE email=$1  AND deleted_at IS NULL `

	err := db.QueryRow(query, email).Scan(&resultID, &username, &password)
	if err == sql.ErrNoRows {
		return 0, "", "", ErrInvalidCredentials
	}
	if err != nil {
		return 0, "", "", err
	}
	return resultID, username, password, nil
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
