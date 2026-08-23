package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Prkarz/course-tracker/models"
)

func EnsureCourseVideoProgressTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS course_video_progress (
			user_id INTEGER NOT NULL,
			course_id INTEGER NOT NULL,
			video_id TEXT NOT NULL,
			is_completed BOOLEAN NOT NULL DEFAULT TRUE,
			PRIMARY KEY (user_id, course_id, video_id)
		)`)
	return err
}

// Creating a Course
func Create_course(tx *sql.Tx, owner_id int, url, title string, summary string, tags []string) (int, bool, error) {

	var isUserActive bool
	err := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND deleted_at IS NULL)", owner_id).Scan(&isUserActive)
	if err != nil {
		return 0, false, err
	}
	if !isUserActive {
		// If they are deleted, kick them out IMMEDIATELY before doing anything else!
		return 0, false, errors.New("blocked: course creator is deleted or does not exist")
	}

	query := `
    INSERT INTO courses (owner_id, playlist_url, title, ai_summary, course_tags, deleted_at)
    SELECT $1, $2, $3, $4, $5, NULL
    FROM users 
    WHERE id = $1 AND deleted_at IS NULL
    ON CONFLICT (playlist_url) DO UPDATE
    SET deleted_at = NULL,
        title = CASE WHEN courses.title IS NULL OR courses.title = '' THEN EXCLUDED.title ELSE courses.title END,
        ai_summary = CASE WHEN courses.ai_summary IS NULL OR courses.ai_summary = '' THEN EXCLUDED.ai_summary ELSE courses.ai_summary END,
        course_tags = CASE WHEN courses.course_tags IS NULL OR cardinality(courses.course_tags) = 0 THEN EXCLUDED.course_tags ELSE courses.course_tags END
    RETURNING id, (xmax = 0) AS is_new;
    `

	var course_id int
	var isNewCourse bool
	err = tx.QueryRow(query, owner_id, url, title, summary, tags).Scan(&course_id, &isNewCourse)

	if err != nil {
		// Fallback query if conflict without RETURNING or other dialect edge case
		fallback_query := "SELECT id FROM courses WHERE playlist_url = $1"
		err = tx.QueryRow(fallback_query, url).Scan(&course_id)
		if err != nil {
			return 0, false, err
		}
		// Un-delete if it was previously soft-deleted
		_, _ = tx.Exec("UPDATE courses SET deleted_at = NULL WHERE id = $1", course_id)
		isNewCourse = false
	}

	// Always ensure user is enrolled in user_progress when they import a course so it appears in "Your Courses"
	_, _ = tx.Exec(`
		INSERT INTO user_progress (user_id, course_id, completion_percentage, started_at, last_accessed_at)
		VALUES ($1, $2, 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id, course_id) DO UPDATE
		SET last_accessed_at = CURRENT_TIMESTAMP
	`, owner_id, course_id)

	return course_id, isNewCourse, nil
}

// List of all Courses
func List_my_courses(contxt context.Context, db *sql.DB, userID int) ([]models.Course_data, error) {

	reports := make([]models.Course_data, 0)
	query := `SELECT  
	courses.id,
	courses.owner_id, 
	courses.playlist_url, 
	courses.title,
	courses.ai_summary,
	COALESCE(to_json(courses.course_tags), '[]'::json) AS course_tags,
	COALESCE(user_progress.completion_percentage, 0) AS completion_percentage,
	user_progress.started_at,
	(user_progress.course_id IS NOT NULL) AS is_started,
	(SELECT youtube_video_id FROM course_videos WHERE course_videos.course_id = courses.id ORDER BY index_order ASC LIMIT 1) AS first_video_id
              FROM courses
			  LEFT JOIN user_progress ON courses.id = user_progress.course_id AND user_progress.user_id = $1
              WHERE (courses.owner_id = $1 OR user_progress.user_id = $1) AND courses.deleted_at IS NULL
              ORDER BY COALESCE(user_progress.started_at, courses.created_at, CURRENT_TIMESTAMP) DESC, courses.id DESC;`
	rows, err := db.QueryContext(contxt, query, userID)
	if err != nil {
		// Fallback query if created_at column does not exist on courses
		fallbackQuery := `SELECT  
		courses.id,
		courses.owner_id, 
		courses.playlist_url, 
		courses.title,
		courses.ai_summary,
		COALESCE(to_json(courses.course_tags), '[]'::json) AS course_tags,
		COALESCE(user_progress.completion_percentage, 0) AS completion_percentage,
		user_progress.started_at,
		(user_progress.course_id IS NOT NULL) AS is_started,
		(SELECT youtube_video_id FROM course_videos WHERE course_videos.course_id = courses.id ORDER BY index_order ASC LIMIT 1) AS first_video_id
				  FROM courses
				  LEFT JOIN user_progress ON courses.id = user_progress.course_id AND user_progress.user_id = $1
				  WHERE (courses.owner_id = $1 OR user_progress.user_id = $1) AND courses.deleted_at IS NULL
				  ORDER BY courses.id DESC;`
		rows, err = db.QueryContext(contxt, fallbackQuery, userID)
		if err != nil {
			return nil, err
		}
	}
	// The defer statement ensures that the rows are closed after the function completes, preventing resource leaks.
	defer rows.Close()
	for rows.Next() {
		var item models.Course_data
		// Use intermediate sql.Null* types for nullable DB columns, then map to pointer fields.
		var courseID int
		var ownerID sql.NullInt64
		var url sql.NullString
		var title sql.NullString
		var summary sql.NullString
		var tagsJSON sql.NullString
		var percentageProgress sql.NullFloat64
		var startedAt sql.NullTime
		var isStarted bool
		var firstVideoID sql.NullString
		err := rows.Scan(&courseID, &ownerID, &url, &title, &summary, &tagsJSON, &percentageProgress, &startedAt, &isStarted, &firstVideoID)
		if err != nil {
			return nil, err
		}
		item.ID = courseID
		if ownerID.Valid {
			tmp := int(ownerID.Int64)
			item.OwnerID = &tmp
		} else {
			item.OwnerID = nil
		}
		if url.Valid {
			tmp := url.String
			item.URL = &tmp
		} else {
			item.URL = nil
		}
		if title.Valid {
			tmp := title.String
			item.Title = &tmp
		} else {
			item.Title = nil
		}
		if summary.Valid {
			tmp := summary.String
			item.Summary = &tmp
		}
		item.Tags = []string{}
		if tagsJSON.Valid && tagsJSON.String != "" {
			_ = json.Unmarshal([]byte(tagsJSON.String), &item.Tags)
		}
		if percentageProgress.Valid {
			item.CompletionPercent = percentageProgress.Float64
		}
		if startedAt.Valid {
			item.StartedAt = &startedAt.Time
		}
		item.IsStarted = isStarted
		if firstVideoID.Valid && firstVideoID.String != "" {
			vid := firstVideoID.String
			item.FirstVideoID = &vid
			thumb := fmt.Sprintf("https://img.youtube.com/vi/%s/hqdefault.jpg", vid)
			item.ThumbnailURL = &thumb
		}
		reports = append(reports, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return reports, nil
}

// Starting a course
func Start_course(contxt context.Context, tx *sql.Tx, userID, courseID int) error {
	current_Time := time.Now()

	query := `
	INSERT INTO user_progress (user_id, course_id, completion_percentage, started_at, last_accessed_at)
	SELECT $1, $2, $3, $4, $5
	FROM users
	WHERE id = $1 AND deleted_at IS NULL
	AND EXISTS (SELECT 1 FROM courses WHERE id = $2 AND deleted_at IS NULL)
	ON CONFLICT (user_id, course_id) DO UPDATE
	SET last_accessed_at = EXCLUDED.last_accessed_at;
	`

	_, err := tx.ExecContext(contxt, query, userID, courseID, 0, current_Time, current_Time)
	return err
}

func Delete_course(tx *sql.Tx, userID, CourseID int) error {
	// 1. Check if user is the course owner
	var isOwner bool
	_ = tx.QueryRow("SELECT (owner_id = $1) FROM courses WHERE id = $2 AND deleted_at IS NULL", userID, CourseID).Scan(&isOwner)

	if isOwner {
		_, err := tx.Exec("UPDATE courses SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1 AND owner_id = $2 AND deleted_at IS NULL", CourseID, userID)
		if err != nil {
			return err
		}
	}

	// 2. Always remove this user's progress and video completion records for this course
	// DELETION IMMUNITY: User's total_points in user_stats is permanent and NEVER deducted on deletion.
	_, _ = tx.Exec("DELETE FROM user_progress WHERE user_id = $1 AND course_id = $2", userID, CourseID)
	_, _ = tx.Exec("DELETE FROM course_video_progress WHERE user_id = $1 AND course_id = $2", userID, CourseID)

	return nil
}

func Fetch_IndiCourse(tx *sql.Tx, userID, courseID int) (*models.CourseViewerData, error) {
	query := `
    SELECT 
	courses.title, courses.ai_summary,
	course_videos.title, course_videos.youtube_video_id, course_videos.duration, course_videos.index_order,
	COALESCE(course_video_progress.is_completed, false) AS is_completed
	FROM courses
    LEFT JOIN course_videos 
        ON courses.id = course_videos.course_id
    LEFT JOIN course_video_progress 
        ON course_videos.youtube_video_id = course_video_progress.video_id 
		AND course_video_progress.course_id = courses.id
        AND course_video_progress.user_id = $2
	WHERE courses.id = $1 
	  AND (courses.owner_id = $2 OR EXISTS (SELECT 1 FROM user_progress WHERE user_progress.user_id = $2 AND user_progress.course_id = $1))
	  AND courses.deleted_at IS NULL
	ORDER BY course_videos.index_order ASC;
`

	result, err := tx.Query(query, courseID, userID)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	courseData := &models.CourseViewerData{
		VideoInfo: make([]models.VideoData, 0),
	}
	foundCourse := false
	seenVideoIDs := make(map[string]bool)

	for result.Next() {
		var title, summary sql.NullString
		var videoTitle, videoID, duration sql.NullString
		var videoIndex sql.NullInt64
		var isCompleted bool
		err := result.Scan(
			&title, &summary,
			&videoTitle, &videoID, &duration, &videoIndex, &isCompleted,
		)
		if err != nil {
			return nil, err
		}

		foundCourse = true
		if title.Valid {
			courseData.Title = title.String
		}
		if summary.Valid {
			courseData.Summary = summary.String
		}
		if videoID.Valid && videoID.String != "" {
			vID := videoID.String
			if !seenVideoIDs[vID] {
				seenVideoIDs[vID] = true
				dur := "10:00"
				if duration.Valid && strings.TrimSpace(duration.String) != "" && strings.TrimSpace(duration.String) != "TBD" {
					dur = strings.TrimSpace(duration.String)
				}
				video := models.VideoData{
					VideoId:  vID,
					Status:   isCompleted,
					Duration: dur,
					Index:    len(courseData.VideoInfo) + 1, // Clean sequential 1-based index
				}
				if videoTitle.Valid {
					video.VideoTitle = videoTitle.String
				}
				courseData.VideoInfo = append(courseData.VideoInfo, video)
			}
		}
	}
	if err := result.Err(); err != nil {
		return nil, err
	}
	if !foundCourse {
		return nil, errors.New("course not found")
	}
	return courseData, nil
}

func Insert_course_videos(tx *sql.Tx, courseID int, vd []models.VideoData) error {
	// First clean existing video list for this course to prevent duplication
	_, _ = tx.Exec("DELETE FROM course_videos WHERE course_id = $1", courseID)

	query := `INSERT INTO course_videos(course_id,youtube_video_id,title,index_order,duration)
	        VALUES($1,$2,$3,$4,$5)
			`
	seen := make(map[string]bool)
	idx := 1
	for _, vid := range vd {
		if vid.VideoId == "" || seen[vid.VideoId] {
			continue
		}
		seen[vid.VideoId] = true
		duration := strings.TrimSpace(vid.Duration)
		if duration == "" || duration == "TBD" {
			duration = "10:00"
		}
		_, err := tx.Exec(query, courseID, vid.VideoId, vid.VideoTitle, idx, duration)
		if err != nil {
			return err
		}
		idx++
	}
	return nil
}

func Insert_Video_Progress(tx *sql.Tx, userID int, courseID int, videoID string) (bool, error) {
	var alreadyViewed bool
	err := tx.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM course_video_progress
			WHERE user_id = $1 AND course_id = $2 AND video_id = $3
		)`, userID, courseID, videoID).Scan(&alreadyViewed)
	if err != nil {
		return false, err
	}
	if alreadyViewed {
		return false, nil
	}
	query := `
		INSERT INTO course_video_progress (user_id, course_id, video_id)
		SELECT $1, $2, $3
		FROM users
		WHERE users.id = $1
		  AND users.deleted_at IS NULL
		  AND EXISTS (
			SELECT 1
			FROM courses
			WHERE courses.id = $2
			  AND courses.deleted_at IS NULL
		  )
		  AND EXISTS (
			SELECT 1
			FROM course_videos
			WHERE course_videos.course_id = $2
			  AND course_videos.youtube_video_id = $3
		  )
	`

	result, err := tx.Exec(query, userID, courseID, videoID)
	if err != nil {
		return false, err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if rows == 0 {
		return false, errors.New("blocked: user, course, or video is invalid")
	}

	return true, nil
}
