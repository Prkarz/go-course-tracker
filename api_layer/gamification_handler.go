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
		writeError(w, http.StatusBadRequest, "400_INVALID_REQUEST", "Failed to parse progress update request. Ensure course_id and percentage_to_add are provided.")
		return
	}
	if req.CourseID <= 0 || req.NewPercentage <= 0 {
		writeError(w, http.StatusBadRequest, "400_INVALID_REQUEST", "course_id and a positive percentage increment are required.")
		return
	}
	tx, err := s.DB.Begin()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "500_DB_TRANSACTION_FAILED", "Unable to initiate database transaction for progress update.")
		return
	}
	defer tx.Rollback()
	progresspercentage, err := repository.Update_progress(tx, userID, req.CourseID, req.NewPercentage)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "500_PROGRESS_UPDATE_FAILED", "Failed to update course progress.")
		return
	}
	if progresspercentage == 100.00 {
		pointstoReward = 100
	} else {
		pointstoReward = 10
	}
	err = repository.Points_update(tx, userID, pointstoReward)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "500_POINTS_UPDATE_FAILED", "Unable to update points.")
		return
	}

	err = repository.Streak_update(tx, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "500_STREAK_UPDATE_FAILED", "Unable to update streak.")
		return
	}
	err = tx.Commit()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "500_DB_COMMIT_FAILED", "Failed to commit progress update to the database.")
		return
	}
	// 1. Tell the browser we are sending JSON
	w.Header().Set("Content-Type", "application/json")

	// 2. Pack the envelope
	response := ProgressResponse{
		Message:      "[200_PROGRESS_UPDATED] Progress updated successfully.",
		PointsEarned: pointstoReward,
	}

	// 3. Send it
	json.NewEncoder(w).Encode(response)

}
