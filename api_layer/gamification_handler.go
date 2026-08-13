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
		http.Error(w, "Unable to Decode File", http.StatusBadRequest)
		return
	}
	tx, err := s.DB.Begin()
	if err != nil {
		http.Error(w, "Unable to  Begin Transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	err = repository.Update_progress(tx, userID, req.CourseID, req.NewPercentage)
	if err != nil {
		http.Error(w, "USER doesnot Exist", http.StatusInternalServerError)
		return
	}
	err = tx.Commit()
	if err != nil {
		http.Error(w, "Commit ERROR", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Progress Updated"))
}

func (s *APIServer) Points_Streak_toUpdate_Handler(w http.ResponseWriter, r *http.Request) {
	var req models.PointsTOUpdate
	userID := r.Context().Value("userID").(int)
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Unable to Decode File", http.StatusBadRequest)
		return
	}
	tx, err := s.DB.Begin()
	if err != nil {
		http.Error(w, "Unable to  Begin Transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	err = repository.Points_update(tx, userID, req.PointstoAdd)
	if err != nil {
		http.Error(w, "Unable to  Update Points", http.StatusInternalServerError)
		return
	}

	if req.IsFirstActionToday {
		err = repository.Streak_update(tx, userID)
		if err != nil {
			http.Error(w, "Unable to Update Streak", http.StatusInternalServerError)
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		http.Error(w, "SERVER ERROR", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Points Updated"))
}
