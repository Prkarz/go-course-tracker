package service

import (
	"database/sql"
	"errors"

	"github.com/Prkarz/course-tracker/repository"
	"golang.org/x/crypto/bcrypt"
)

func Register_user(db *sql.Tx, username, email, plain_password string) (int, bool, error) {
	hashed_Password, err := bcrypt.GenerateFromPassword([]byte(plain_password), bcrypt.DefaultCost)
	if err != nil {
		return 0, false, errors.New("Unable to hash password")
	}
	userID, isNewUser, err := repository.Create_user(db, username, email, string(hashed_Password))
	if err != nil {
		return 0, false, errors.New("Unable to Registering User")
	}
	return userID, isNewUser, nil
}
func LoginUser(db *sql.DB, email string, typedPassword string) (int, error) {
	UserID, password, err := repository.Login_User_byemail(db, email)
	if err != nil {
		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(typedPassword))
	if err != nil {
		return 0, repository.ErrInvalidCredentials
	}
	return UserID, nil
    
}