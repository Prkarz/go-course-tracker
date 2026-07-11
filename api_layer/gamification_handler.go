package apilayer

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Prkarz/course-tracker/models"
	"github.com/Prkarz/course-tracker/repository"
)

func (s *APIServer) Update_progress_Handler(w http.ResponseWriter, r *http.Request) {
	var req models.UpdateProgress
	err := json.NewDecoder(r.Body).Decode(&req)
	log.Printf("DEBUG: Received UserID: %d, CourseID: %d, NewPercentage: %d", req.UserID, req.CourseID, req.NewPercentage)
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
	err = repository.Update_progress(tx, req.UserID, req.CourseID, req.NewPercentage)
	if err != nil {
		http.Error(w, "Receiver Progress", http.StatusInternalServerError)
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
	err = repository.Points_update(tx, req.UserID, req.PointstoAdd)
	if err != nil {
		http.Error(w, "Unable to  Update Points", http.StatusInternalServerError)
		return
	}

	if req.IsFirstActionToday {
		err = repository.Streak_update(tx, req.UserID)
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
