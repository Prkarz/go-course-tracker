package apilayer

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Prkarz/course-tracker/models"
	"github.com/Prkarz/course-tracker/repository"
)

func (s *APIServer) Course_Creation_Handler(w http.ResponseWriter, r *http.Request) {
	var req models.CourseCreationRequest
	OwnerID := r.Context().Value("userID").(int)
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "[400_INVALID_REQUEST] Failed to parse course creation request. Ensure title and URL are provided.", http.StatusBadRequest)
		return
	}

	tx, err := s.DB.Begin()
	if err != nil {
		http.Error(w, "[500_DB_TRANSACTION_FAILED] Unable to initiate database transaction for course creation.", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	parsed_url, err := url.Parse(req.URL)
	if err != nil {
		http.Error(w, "Error parsing url", http.StatusBadRequest)
		return
	}
	if !strings.Contains(parsed_url.Host, "youtube.com") && !strings.Contains(parsed_url.Host, "youtu.be") {
		http.Error(w, "[400_INVALID_DOMAIN] Submitted link is not a valid YouTube domain.", http.StatusBadRequest)
		return
	}

	ListID := parsed_url.Query().Get("list")
	if ListID == "" {
		http.Error(w, "[400_INVALID_URL] Course tracker requires a YouTube playlist link, not a single video.", http.StatusBadRequest)
		return
	}

	courseID, _, err := repository.Create_course(tx, OwnerID, req.URL, req.Title)
	if err != nil {
		http.Error(w, "[500_COURSE_CREATION_FAILED] Course creation failed. Duplicate course or insufficient permissions.", http.StatusInternalServerError)
		return
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "[500_DB_COMMIT_FAILED] Failed to commit course creation changes to database.", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	respose := map[string]interface{}{
		"message":   "Course created successfully!",
		"course_id": courseID,
	}
	json.NewEncoder(w).Encode((respose))

}

func (s *APIServer) Start_Course_Handler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	contxt, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	var req models.StartCourseRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "[400_INVALID_REQUEST] Failed to parse course start request. Ensure courseID is provided.", http.StatusBadRequest)
		return
	}
	tx, err := s.DB.Begin()
	if err != nil {
		http.Error(w, "[500_DB_TRANSACTION_FAILED] Unable to initiate database transaction for course start.", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	err = repository.Start_course(contxt, tx, userID, req.CourseID)
	if err != nil {
		http.Error(w, "[500_COURSE_START_FAILED] Failed to start course. Course not found or you don't have access.", http.StatusInternalServerError)
		return
	}
	err = tx.Commit()
	if err != nil {
		http.Error(w, "[500_DB_COMMIT_FAILED] Failed to commit course start transaction.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Course Started Successfully"))

}
