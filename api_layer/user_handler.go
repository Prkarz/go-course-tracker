package apilayer

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Prkarz/course-tracker/models"
	"github.com/Prkarz/course-tracker/repository"
)

type APIServer struct {
	DB *sql.DB
}

func (s *APIServer) User_creation(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Couldn't create USER", http.StatusBadRequest)
		return
	}

	tx, err := s.DB.Begin()
	if err != nil {
		http.Error(w, "SERVER ERROR", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	userID, isNewUser, err := repository.Create_user(tx, req.Username, req.Email)
	if err != nil {
		http.Error(w, "SERVER ERROR", http.StatusInternalServerError)
		return
	}
	if isNewUser {
		err := repository.InitUserStats(tx, userID)
		if err != nil {
			http.Error(w, "SERVER ERROR", http.StatusInternalServerError)
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		http.Error(w, "SERVER ERROR", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("User successfully created!"))

}
