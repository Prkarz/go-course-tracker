package apilayer

import (
	"encoding/json"
	"net/http"

	"github.com/Prkarz/course-tracker/models"
	"github.com/Prkarz/course-tracker/repository"
)

func (s *APIServer) Update_progress_Handler(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateProgress
	userID := r.Context().Value("userID").(int)
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "[400_INVALID_REQUEST] Failed to parse progress update request. Ensure courseID and percentage are provided.", http.StatusBadRequest)
		return
	}
	tx, err := s.DB.Begin()
	if err != nil {
		http.Error(w, "[500_DB_TRANSACTION_FAILED] Unable to initiate database transaction for progress update.", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	err = repository.Update_progress(tx, userID, req.CourseID, req.NewPercentage)
	if err != nil {
		http.Error(w, "[500_PROGRESS_UPDATE_FAILED] Failed to update course progress. Course not found or invalid percentage.", http.StatusInternalServerError)
		return
	}
	err = tx.Commit()
	if err != nil {
		http.Error(w, "[500_DB_COMMIT_FAILED] Failed to commit progress update to database.", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Progress Updated"))
}

func (s *APIServer) Points_Streak_toUpdate_Handler(w http.ResponseWriter, r *http.Request) {
	var req models.PointsTOUpdate
	userID := r.Context().Value("userID").(int)
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "[400_INVALID_REQUEST] Failed to parse points update request. Ensure pointsToAdd is provided.", http.StatusBadRequest)
		return
	}
	tx, err := s.DB.Begin()
	if err != nil {
		http.Error(w, "[500_DB_TRANSACTION_FAILED] Unable to initiate database transaction for points update.", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	err = repository.Points_update(tx, userID, req.PointstoAdd)
	if err != nil {
		http.Error(w, "[500_POINTS_UPDATE_FAILED] Failed to update points. User not found or invalid point value.", http.StatusInternalServerError)
		return
	}

	if req.IsFirstActionToday {
		err = repository.Streak_update(tx, userID)
		if err != nil {
			http.Error(w, "[500_STREAK_UPDATE_FAILED] Failed to update streak. Please try again.", http.StatusInternalServerError)
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		http.Error(w, "[500_DB_COMMIT_FAILED] Failed to commit points and streak update to database.", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Points Updated"))
}
