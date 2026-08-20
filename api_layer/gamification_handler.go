package apilayer

import (
	"encoding/json"
	"net/http"

	"github.com/Prkarz/course-tracker/models"
	"github.com/Prkarz/course-tracker/repository"
)

type ProgressResponse struct {
	Message      string `json:"message"`
	PointsEarned int    `json:"points_earned"`
}

func (s *APIServer) Update_progress_Handler(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateProgress
	var pointstoReward int
	userID := r.Context().Value("userID").(int)
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "[400_INVALID_REQUEST] Failed to parse progress update request. Ensure courseID and percentage are provided.", http.StatusBadRequest)
		return
	}
	if req.CourseID <= 0 || req.NewPercentage <= 0 {
		http.Error(w, "[400_INVALID_REQUEST] course_id and a positive percentage increment are required.", http.StatusBadRequest)
		return
	}
	tx, err := s.DB.Begin()
	if err != nil {
		http.Error(w, "[500_DB_TRANSACTION_FAILED] Unable to initiate database transaction for progress update.", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	progresspercentage, err := repository.Update_progress(tx, userID, req.CourseID, req.NewPercentage)
	if err != nil {
		http.Error(w, "[500_PROGRESS_UPDATE_FAILED] Failed to update course progress. Course not found or invalid percentage.", http.StatusInternalServerError)
		return
	}
	if progresspercentage == 100.00 {
		pointstoReward = 100
	} else {
		pointstoReward = 10
	}
	err = repository.Points_update(tx, userID, pointstoReward)
	if err != nil {
		http.Error(w, "COuldnot update points", http.StatusInternalServerError)
		return
	}

	err = repository.Streak_update(tx, userID)
	if err != nil {
		http.Error(w, "COuldnot update streak", http.StatusInternalServerError)
		return
	}
	err = tx.Commit()
	if err != nil {
		http.Error(w, "[500_DB_COMMIT_FAILED] Failed to commit progress update to database.", http.StatusInternalServerError)
		return
	}
	// 1. Tell the browser we are sending JSON
	w.Header().Set("Content-Type", "application/json")

	// 2. Pack the envelope
	response := ProgressResponse{
		Message:      "Progress Successfully Updated",
		PointsEarned: pointstoReward,
	}

	// 3. Send it
	json.NewEncoder(w).Encode(response)

}
