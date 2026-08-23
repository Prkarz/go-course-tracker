package apilayer

import (
	"encoding/json"
	"net/http"

	"github.com/Prkarz/course-tracker/models"
	"github.com/Prkarz/course-tracker/repository"
)

func (s *APIServer) User_stats_Handler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	stats, err := repository.Get_user_stats(s.DB, userID)
	if err != nil {
		writeError(w, http.StatusNotFound, "404_STATS_NOT_FOUND", "User stats not found.")
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

func (s *APIServer) Video_Viewed_Handler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	var req models.VideoViewedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.CourseID <= 0 || req.VideoID == "" {
		writeError(w, http.StatusBadRequest, "400_INVALID_REQUEST", "course_id and video_id are required.")
		return
	}

	tx, err := s.DB.Begin()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "500_DB_TRANSACTION_FAILED", "Unable to record video activity.")
		return
	}
	defer tx.Rollback()
	newVideo, err := repository.Insert_Video_Progress(tx, userID, req.CourseID, req.VideoID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "400_VIDEO_NOT_FOUND", "Video was not found in this course.")
		return
	}
	pointsEarned := 0
	if newVideo {
		var totalVideos, completedVideos int
		_ = tx.QueryRow("SELECT COUNT(*) FROM course_videos WHERE course_id = $1", req.CourseID).Scan(&totalVideos)
		_ = tx.QueryRow("SELECT COUNT(*) FROM course_video_progress WHERE course_id = $1 AND user_id = $2", req.CourseID, userID).Scan(&completedVideos)

		newPercent := 100.0
		if totalVideos > 0 {
			newPercent = (float64(completedVideos) / float64(totalVideos)) * 100.0
			if newPercent > 100.0 {
				newPercent = 100.0
			}
		}

		_, _ = tx.Exec(`
			INSERT INTO user_progress (user_id, course_id, completion_percentage, started_at, last_accessed_at)
			VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			ON CONFLICT (user_id, course_id) DO UPDATE
			SET completion_percentage = EXCLUDED.completion_percentage,
			    last_accessed_at = CURRENT_TIMESTAMP
		`, userID, req.CourseID, newPercent)

		pointsEarned = 10
		if newPercent >= 100.0 {
			pointsEarned = 100
		}

		if err := repository.Points_update(tx, userID, pointsEarned); err != nil {
			writeError(w, http.StatusInternalServerError, "500_POINTS_UPDATE_FAILED", "Unable to update points.")
			return
		}
	}
	if err := repository.Streak_update(tx, userID); err != nil {
		writeError(w, http.StatusInternalServerError, "500_STREAK_UPDATE_FAILED", "Unable to update streak.")
		return
	}
	if err := tx.Commit(); err != nil {
		writeError(w, http.StatusInternalServerError, "500_DB_COMMIT_FAILED", "Failed to commit video activity.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"message": "Video activity recorded.", "points_earned": pointsEarned, "new_video": newVideo})
}

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
