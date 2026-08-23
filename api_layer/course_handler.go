package apilayer

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	aiintegration "github.com/Prkarz/course-tracker/Ai_integration"
	"github.com/Prkarz/course-tracker/models"
	"github.com/Prkarz/course-tracker/repository"
	"github.com/Prkarz/course-tracker/service"
)

func (s *APIServer) Course_Creation_Handler(w http.ResponseWriter, r *http.Request) {
	var req models.CourseCreationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "400_INVALID_REQUEST", "Failed to parse course creation request. Ensure title and URL are provided.")
		return
	}
	if strings.TrimSpace(req.URL) == "" {
		req.URL = req.PlaylistURL
	}
	req.URL = strings.TrimSpace(req.URL)
	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" || req.URL == "" {
		writeError(w, http.StatusBadRequest, "400_INVALID_REQUEST", "Course title and playlist URL are required.")
		return
	}

	parsed_url, err := url.Parse(req.URL)
	if err != nil {
		writeError(w, http.StatusBadRequest, "400_INVALID_URL", "Failed to parse the playlist URL.")
		return
	}
	if parsed_url.Scheme != "https" && parsed_url.Scheme != "http" {
		writeError(w, http.StatusBadRequest, "400_INVALID_URL", "Playlist URL must use http or https.")
		return
	}
	host := strings.ToLower(parsed_url.Hostname())
	if host != "youtube.com" && !strings.HasSuffix(host, ".youtube.com") && host != "youtu.be" {
		writeError(w, http.StatusBadRequest, "400_INVALID_DOMAIN", "Submitted link is not a valid YouTube domain.")
		return
	}
	ListID := parsed_url.Query().Get("list")
	if ListID == "" {
		writeError(w, http.StatusBadRequest, "400_INVALID_URL", "Course tracker requires a YouTube playlist link, not a single video.")
		return
	}

	videos, err := service.FetchPlayListVideos(ListID)
	if err != nil {
		log.Printf("[YOUTUBE_API_FAILED] listID=%q error=%v", ListID, err)
		writeError(w, http.StatusInternalServerError, "500_YOUTUBE_FETCH_FAILED", "Failed to fetch playlist videos from YouTube.")
		return
	}

	OwnerID, ok := r.Context().Value("userID").(int)
	if !ok || OwnerID <= 0 {
		writeError(w, http.StatusUnauthorized, "401_INVALID_USER", "Authenticated user is invalid.")
		return
	}

	response, err := aiintegration.CourseInsights(r.Context(), s.AIClient, req.Title)
	if err != nil {
		log.Printf("[AI_SUMMARY_FAILED] title=%q error=%v", req.Title, err)
		writeError(w, http.StatusInternalServerError, "500_AI_SUMMARY_FAILED", "Unable to generate course insights.")
		return
	}

	tx, err := s.DB.Begin()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "500_DB_TRANSACTION_FAILED", "Unable to initiate database transaction for course creation.")
		return
	}
	defer tx.Rollback()

	courseID, isNewCourse, err := repository.Create_course(tx, OwnerID, req.URL, req.Title, response.Summary, response.Tags)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "500_COURSE_CREATION_FAILED", "Unable to save the course.")
		return
	}

	if isNewCourse {
		err = repository.Insert_course_videos(tx, courseID, videos)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "500_VIDEO_SAVE_FAILED", "Failed to save the playlist videos.")
			return
		}
	} else {
		var videoCount int
		_ = tx.QueryRow("SELECT COUNT(*) FROM course_videos WHERE course_id = $1", courseID).Scan(&videoCount)
		if videoCount == 0 && len(videos) > 0 {
			_ = repository.Insert_course_videos(tx, courseID, videos)
		}
	}

	err = tx.Commit()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "500_DB_COMMIT_FAILED", "Failed to commit course creation changes to database.")
		return
	}
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"course_id":      courseID,
		"created":        true,
		"already_exists": false,
		"message":        "[201_COURSE_CREATED] Course added to your courses successfully.",
	})

}

func (s *APIServer) Start_Course_Handler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	contxt, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	var req models.StartCourseRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		writeError(w, http.StatusBadRequest, "400_INVALID_REQUEST", "Failed to parse course start request. Ensure course_id is provided.")
		return
	}
	if req.CourseID <= 0 {
		writeError(w, http.StatusBadRequest, "400_INVALID_REQUEST", "course_id must be a positive integer.")
		return
	}
	tx, err := s.DB.Begin()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "500_DB_TRANSACTION_FAILED", "Unable to initiate database transaction for course start.")
		return
	}
	defer tx.Rollback()
	err = repository.Start_course(contxt, tx, userID, req.CourseID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "500_COURSE_START_FAILED", "Failed to start course. Course not found or you do not have access.")
		return
	}
	err = tx.Commit()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "500_DB_COMMIT_FAILED", "Failed to commit course start transaction.")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Course started successfully.",
	})

}

func (s *APIServer) Course_Detail_Handler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	var req models.StartCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.CourseID <= 0 {
		writeError(w, http.StatusBadRequest, "400_INVALID_REQUEST", "course_id must be a positive integer.")
		return
	}

	tx, err := s.DB.Begin()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "500_DB_TRANSACTION_FAILED", "Unable to initiate course lookup.")
		return
	}
	defer tx.Rollback()

	detail, err := repository.Fetch_IndiCourse(tx, userID, req.CourseID)
	if err != nil {
		log.Printf("[COURSE_DETAIL_FAILED] user_id=%d course_id=%d error=%v", userID, req.CourseID, err)
		writeError(w, http.StatusInternalServerError, "500_COURSE_DETAIL_FAILED", "Unable to load course videos.")
		return
	}

	// Auto-heal missing or placeholder durations for existing course records
	var missingVideoIDs []string
	for _, vid := range detail.VideoInfo {
		if vid.Duration == "" || vid.Duration == "TBD" || vid.Duration == "10:00" {
			missingVideoIDs = append(missingVideoIDs, vid.VideoId)
		}
	}
	if len(missingVideoIDs) > 0 {
		fetchedDurations := service.FetchVideoDurations(missingVideoIDs)
		for i, vid := range detail.VideoInfo {
			if realDur, ok := fetchedDurations[vid.VideoId]; ok && realDur != "" {
				detail.VideoInfo[i].Duration = realDur
				_, _ = tx.Exec("UPDATE course_videos SET duration = $1 WHERE youtube_video_id = $2 AND course_id = $3", realDur, vid.VideoId, req.CourseID)
			}
		}
	}
	_ = tx.Commit()

	writeJSON(w, http.StatusOK, detail)
}

func (s *APIServer) Course_Delete_Handler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	var req models.StartCourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.CourseID <= 0 {
		writeError(w, http.StatusBadRequest, "400_INVALID_REQUEST", "course_id must be a positive integer.")
		return
	}

	tx, err := s.DB.Begin()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "500_DB_TRANSACTION_FAILED", "Unable to initiate course deletion.")
		return
	}
	defer tx.Rollback()
	if err := repository.Delete_course(tx, userID, req.CourseID); err != nil {
		writeError(w, http.StatusNotFound, "404_COURSE_NOT_FOUND", "Course not found.")
		return
	}
	if err := tx.Commit(); err != nil {
		writeError(w, http.StatusInternalServerError, "500_DB_COMMIT_FAILED", "Failed to commit course deletion.")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"message": "Course deleted successfully. Your earned XP is permanent and preserved.",
	})
}
