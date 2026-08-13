package apilayer

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Prkarz/course-tracker/config"
	"github.com/Prkarz/course-tracker/models"
	"github.com/Prkarz/course-tracker/repository"
	"github.com/Prkarz/course-tracker/service"
	"github.com/golang-jwt/jwt/v5"
)

// User_creation function handles the creation of a new user.
// It decodes the request payload into a CreateUserRequest struct, starts a database transaction,
// and calls the Create_user function from the repository package to create the user.
// If the user is new, it initializes the user's stats using the InitUserStats function.
// Finally, it commits the transaction and sends a success response to the client.
func (s *APIServer) User_creation_Handler(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUserRequest
	//req is a variable of type CreateUserRequest, which is used to store the decoded request payload.
	err := json.NewDecoder(r.Body).Decode(&req)
	//decode req the request body into the req variable. If there is an error during decoding,
	//it sends a bad request response to the client.
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

	userID, isNewUser, err := service.Register_user(tx, req.Username, req.Email, req.Password)
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
	//write is used to send a success message back to the client indicating that the user was successfully created.

}

func (s *APIServer) User_Login_Handler(w http.ResponseWriter, r *http.Request) {
	var req models.LoginUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Couldnot Decode Request", http.StatusInternalServerError)
		return
	}
	userID, err := service.LoginUser(s.DB, req.Email, req.Password)
	if err != nil {
		http.Error(w, "Inavlid creds", http.StatusUnauthorized)
		return
	}
	secret_key := config.JWTSecret
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret_key)
	if err != nil {
		http.Error(w, "Failed to forge token", http.StatusBadRequest)
		return
	}
	w.Write([]byte(tokenString))
}

func (s *APIServer) User_delete_Handler(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteUserRequest
	UserID := r.Context().Value("userID").(int)
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Couldn't delete USER", http.StatusBadRequest)
		return
	}
	tx, err := s.DB.Begin()
	if err != nil {
		http.Error(w, "SERVER ERROR", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	err = repository.Delete_user_by_id(tx, UserID)
	if err != nil {
		http.Error(w, "SERVER ERROR", http.StatusInternalServerError)
		return
	}
	err = tx.Commit()
	if err != nil {
		http.Error(w, "SERVER ERROR", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("User Profile successfully deleted!"))
}

func (s *APIServer) List_myCourses_Handler(w http.ResponseWriter, r *http.Request) {
	contxt, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	UserID := r.Context().Value("userID").(int)
	reports, err := repository.List_my_courses(contxt, s.DB, UserID)
	if err != nil {
		http.Error(w, "Couldn't Fetch USER's Courses", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(reports)
}
