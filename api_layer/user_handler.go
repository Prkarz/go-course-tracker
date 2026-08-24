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
		writeError(w, http.StatusBadRequest, "400_INVALID_REQUEST", "Failed to parse user creation request. Please check the JSON payload format.")
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	if req.Username == "" || req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "400_INVALID_REQUEST", "Username, email, and password are required.")
		return
	}

	tx, err := s.DB.Begin()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "500_DB_TRANSACTION_FAILED", "Unable to initiate database transaction. Please try again.")
		return
	}
	defer tx.Rollback()

	userID, isNewOrReactivated, isReactivated, err := service.Register_user(tx, req.Username, req.Email, req.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "500_USER_REGISTRATION_FAILED", "User registration failed.")
		return
	}
	if !isNewOrReactivated {
		writeError(w, http.StatusConflict, "409_EMAIL_EXISTS", "An active account with this email already exists. Please log in.")
		return
	}
	if !isReactivated {
		err := repository.InitUserStats(tx, userID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "500_USER_STATS_INIT_FAILED", "Failed to initialize user statistics.")
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "500_DB_COMMIT_FAILED", "Failed to commit changes to the database.")
		return
	}

	msg := "Account created successfully."
	if isReactivated {
		msg = "Welcome back! Account reactivated and your courses have been restored."
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"message":     msg,
		"reactivated": isReactivated,
		"user_id":     userID,
	})
	//write is used to send a success message back to the client indicating that the user was successfully created.

}

func (s *APIServer) User_Login_Handler(w http.ResponseWriter, r *http.Request) {
	var req models.LoginUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "400_INVALID_REQUEST", "Failed to parse login request. Ensure email and password are provided.")
		return
	}
	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "400_INVALID_REQUEST", "Email and password are required.")
		return
	}
	userID, username, err := service.LoginUser(s.DB, req.Email, req.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "401_AUTH_FAILED", "Invalid email or password.")
		return
	}

	// Award login reward (+50 XP bonus on daily login / re-login) and update streak
	_, _ = repository.Record_Login_Reward(s.DB, userID, 50)

	secret_key := string(config.GetJWTSecret())
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"email":    req.Email,
		"exp":      time.Now().Add(config.JWTExpiryDuration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret_key))
	if err != nil {
		log.Printf("[TOKEN_GENERATION_FAILED] error=%v", err)
		writeError(w, http.StatusInternalServerError, "500_TOKEN_GENERATION_FAILED", "Token generation failed. Please try logging in again.")
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(tokenString))
}

func (s *APIServer) User_delete_Handler(w http.ResponseWriter, r *http.Request) {
	UserID := r.Context().Value("userID").(int)
	tx, err := s.DB.Begin()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "500_DB_TRANSACTION_FAILED", "Unable to initiate database transaction for user deletion.")
		return
	}
	defer tx.Rollback()
	err = repository.Delete_user_by_id(tx, UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "500_USER_DELETE_FAILED", "Failed to delete user account.")
		return
	}
	err = tx.Commit()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "500_DB_COMMIT_FAILED", "Failed to commit user deletion.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{
		"message": "User account deleted successfully.",
	})
}

func (s *APIServer) List_myCourses_Handler(w http.ResponseWriter, r *http.Request) {
	contxt, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	UserID := r.Context().Value("userID").(int)
	reports, err := repository.List_my_courses(contxt, s.DB, UserID)
	if err != nil {
		log.Printf("[COURSE_LIST_FAILED] user_id=%d error=%v", UserID, err)
		writeError(w, http.StatusInternalServerError, "500_FETCH_COURSES_FAILED", "Failed to retrieve your courses.")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reports)
}
