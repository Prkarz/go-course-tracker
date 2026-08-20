package apilayer

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
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
		http.Error(w, "[400_INVALID_REQUEST] Failed to parse user creation request. Please check JSON payload format.", http.StatusBadRequest)
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "[400_INVALID_REQUEST] Username, email, and password are required.", http.StatusBadRequest)
		return
	}

	tx, err := s.DB.Begin()
	if err != nil {
		http.Error(w, "[500_DB_TRANSACTION_FAILED] Unable to initiate database transaction. Please try again.", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	userID, isNewUser, err := service.Register_user(tx, req.Username, req.Email, req.Password)
	if err != nil {
		http.Error(w, "[500_USER_REGISTRATION_FAILED] User registration failed. Email may already exist or password requirements not met.", http.StatusInternalServerError)
		return
	}
	if isNewUser {
		err := repository.InitUserStats(tx, userID)
		if err != nil {
			http.Error(w, "[500_USER_STATS_INIT_FAILED] Failed to initialize user statistics. Please contact support.", http.StatusInternalServerError)
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		http.Error(w, "[500_DB_COMMIT_FAILED] Failed to commit changes to database. User registration may be incomplete.", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User account successfully created."))
	//write is used to send a success message back to the client indicating that the user was successfully created.

}

func (s *APIServer) User_Login_Handler(w http.ResponseWriter, r *http.Request) {
	var req models.LoginUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "[400_INVALID_REQUEST] Failed to parse login request. Ensure email and password fields are provided in JSON format.", http.StatusBadRequest)
		return
	}
	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" || req.Password == "" {
		http.Error(w, "[400_INVALID_REQUEST] Email and password are required.", http.StatusBadRequest)
		return
	}
	userID, err := service.LoginUser(s.DB, req.Email, req.Password)
	if err != nil {
		http.Error(w, "[401_AUTH_FAILED] Authentication failed. Invalid email or password provided.", http.StatusUnauthorized)
		return
	}
	secret_key := string(config.GetJWTSecret())
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(config.JWTExpiryDuration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret_key))
	if err != nil {
		log.Printf("Login Error: %v", err)
		http.Error(w, "[500_TOKEN_GENERATION_FAILED] Token generation failed. Please try logging in again.", http.StatusInternalServerError)
		return
	}
	w.Write([]byte(tokenString))
}

func (s *APIServer) User_delete_Handler(w http.ResponseWriter, r *http.Request) {
	UserID := r.Context().Value("userID").(int)
	tx, err := s.DB.Begin()
	if err != nil {
		http.Error(w, "[500_DB_TRANSACTION_FAILED] Unable to initiate database transaction for user deletion.", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	err = repository.Delete_user_by_id(tx, UserID)
	if err != nil {
		http.Error(w, "[500_USER_DELETE_FAILED] Failed to delete user account. Please try again or contact support.", http.StatusInternalServerError)
		return
	}
	err = tx.Commit()
	if err != nil {
		http.Error(w, "[500_DB_COMMIT_FAILED] Failed to commit user deletion. Account may still exist.", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("User account successfully deleted."))
}

func (s *APIServer) List_myCourses_Handler(w http.ResponseWriter, r *http.Request) {
	contxt, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	UserID := r.Context().Value("userID").(int)
	reports, err := repository.List_my_courses(contxt, s.DB, UserID)
	if err != nil {
		log.Printf("list courses failed for user %d: %v", UserID, err)
		http.Error(w, "[500_FETCH_COURSES_FAILED] Failed to retrieve your courses. Database operation timeout or unavailable.", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(reports)
}
