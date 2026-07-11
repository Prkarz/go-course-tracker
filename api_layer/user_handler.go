package apilayer

import (
	"encoding/json"
	"net/http"

	"github.com/Prkarz/course-tracker/models"
	"github.com/Prkarz/course-tracker/repository"
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
	//write is used to send a success message back to the client indicating that the user was successfully created.

}

func (s *APIServer) User_delete_Handler(w http.ResponseWriter, r *http.Request) {
	var req models.DeleteUserRequest
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
	err = repository.Delete_user_by_id(tx, req.UserID)
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
	var req models.ListMyCourses
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Couldn't Fetch USER's Courses", http.StatusBadRequest)
		return
	}
	reports, err := repository.List_my_courses(s.DB, req.UserID)
	if err != nil {
		http.Error(w, "Couldn't Fetch USER's Courses", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(reports)
}
