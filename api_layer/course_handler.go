package apilayer

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Prkarz/course-tracker/models"
	"github.com/Prkarz/course-tracker/repository"
)

func (s *APIServer) Course_Creation_Handler(w http.ResponseWriter, r *http.Request) {
	var req models.CourseCreationRequest
	OwnerID := r.Context().Value("UserID").(int)
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Couldn't create Course", http.StatusBadRequest)
		return
	}

	tx, err := s.DB.Begin()
	if err != nil {
		http.Error(w, "SERVER ERROR", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	_, _, err = repository.Create_course(tx, OwnerID, req.URL, req.Title)
	if err != nil {
		http.Error(w, "SERVER ERROR", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "SERVER ERROR", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Course successfully created!"))

}

func (s *APIServer) Start_Course_Handler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	contxt, cancel := context.WithTimeout(r.Context(), http.DefaultClient.Timeout)
	defer cancel()
	var req models.StartCourseRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Unable to Decode Data", http.StatusBadRequest)
		return
	}
	tx, err := s.DB.Begin()
	if err != nil {
		http.Error(w, "SERVER ERROR", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	err = repository.Start_course(contxt, tx, userID, req.CourseID)
	if err != nil {
		http.Error(w, "Unable to Start Course", http.StatusInternalServerError)
		return
	}
	err = tx.Commit()
	if err != nil {
		http.Error(w, "SERVER ERROR", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Course Started Successfully"))
}
