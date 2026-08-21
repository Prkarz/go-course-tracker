package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/Prkarz/course-tracker/models"
)

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
    INSERT INTO courses (owner_id, playlist_url, title,ai_summary,course_tags)
    SELECT $1, $2, $3, $4, $5
    FROM users 
    WHERE id = $1 AND deleted_at IS NULL
    ON CONFLICT (playlist_url) DO NOTHING
    RETURNING id;
    `

	var course_id int
	err = tx.QueryRow(query, owner_id, url, title, summary, tags).Scan(&course_id)

	if err == sql.ErrNoRows {
		// In this specific query, ErrNoRows happens for ONE of TWO reasons:
		// Reason A: ON CONFLICT triggered (The URL already exists)
		// Reason B: The SELECT returned nothing (The user is deleted/doesn't exist)

		// Let's check Reason A first by looking for the URL:
		fallback_query := "SELECT id FROM courses WHERE playlist_url = $1  "
		err = tx.QueryRow(fallback_query, url).Scan(&course_id)

		if err == nil {

			return course_id, false, nil
		}

		if err == sql.ErrNoRows {
			return 0, false, errors.New("blocked: course creator is deleted or does not exist")
		}

		return 0, false, err
	}

	if err != nil {
		return 0, false, err
	}

	return course_id, true, nil
}

// List of all Courses
func List_my_courses(contxt context.Context, db *sql.DB, userID int) ([]models.Course_data, error) {

	var reports []models.Course_data
	query := `SELECT  
	courses.id,
	courses.owner_id, 
	courses.playlist_url, 
	courses.title,
	courses.ai_summary,
	COALESCE(to_json(courses.course_tags), '[]'::json) AS course_tags,
	COALESCE(user_progress.completion_percentage,0) AS completion_percentage,
	user_progress.started_at,
	(user_progress.course_id IS NOT NULL) AS is_started
              FROM courses
			  LEFT JOIN user_progress ON courses.id=user_progress.course_id AND user_progress.user_id=$1
              WHERE courses.owner_id = $1 AND courses.deleted_at IS NULL ;`
	rows, err := db.QueryContext(contxt, query, userID)
	if err != nil {
		return nil, err
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
		err := rows.Scan(&courseID, &ownerID, &url, &title, &summary, &tagsJSON, &percentageProgress, &startedAt, &isStarted)
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
		if tagsJSON.Valid && tagsJSON.String != "" {
			if err := json.Unmarshal([]byte(tagsJSON.String), &item.Tags); err != nil {
				return nil, err
			}
		}
		if percentageProgress.Valid {
			item.CompletionPercent = percentageProgress.Float64
		}
		if startedAt.Valid {
			item.StartedAt = &startedAt.Time
		}
		item.IsStarted = isStarted
		reports = append(reports, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return reports, nil
}

// Starting a course
// this function is called when a user starts a course,
// it inserts a new record into the user_progress table with the user ID, course ID, initial completion percentage (0), and timestamps for when the course was started and last accessed.
func Start_course(contxt context.Context, tx *sql.Tx, userID, courseID int) error {
	current_Time := time.Now()

	var exists bool
	err := tx.QueryRowContext(contxt, `
		SELECT EXISTS (
			SELECT 1
			FROM user_progress
			WHERE user_id = $1 AND course_id = $2
		)
	`, userID, courseID).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("Start skipped: User %d already enrolled in Course %d", userID, courseID)
		return errors.New("blocked: user is already enrolled in this course")
	}

	query := `
	INSERT INTO user_progress (user_id, course_id, completion_percentage, started_at, last_accessed_at)
	SELECT $1, $2, $3, $4, $5
	FROM users
	WHERE id = $1 AND deleted_at IS NULL
	AND EXISTS (SELECT 1 FROM courses WHERE id = $2 AND deleted_at IS NULL)
	`

	result, err := tx.ExecContext(contxt, query, userID, courseID, 0, current_Time, current_Time)
	if err != nil {
		return err
	}

	// Check the result
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// If rows == 0, it means EITHER the user is deleted, OR they already started this course.
	// Both are safe states, so we don't necessarily need to return a fatal server error.
	if rows == 0 {
		return errors.New("blocked: user is deleted or course does not exist")
		// You can choose to return an error here, or just return nil because the end result
		// (the user being enrolled) is already true!
	}

	return nil
}

func Delete_course(tx *sql.Tx, userID, CourseID int) error {
	query := `UPDATE courses
	            SET deleted_at=CURRENT_TIMESTAMP 
				WHERE owner_id=$1 AND id=$2 AND deleted_at IS NULL
				`

	result, err := tx.Exec(query, userID, CourseID)

	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("Unauthorized or course not found")
	}
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
	WHERE courses.id = $1 AND courses.deleted_at IS NULL
	ORDER BY course_videos.index_order ASC;
`

	result, err := tx.Query(query, courseID, userID)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	courseData := &models.CourseViewerData{}
	foundCourse := false

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
		if videoID.Valid {
			video := models.VideoData{
				VideoId: videoID.String,
				Status:  isCompleted,
			}
			if videoTitle.Valid {
				video.VideoTitle = videoTitle.String
			}
			if duration.Valid {
				video.Duration = duration.String
			}
			if videoIndex.Valid {
				video.Index = int(videoIndex.Int64)
			}
			courseData.VideoInfo = append(courseData.VideoInfo, video)
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
	query := `INSERT INTO course_videos(course_id,youtube_video_id,title,index_order,duration)
	        VALUES($1,$2,$3,$4,$5)
			`
	for _, vid := range vd {
		_, err := tx.Exec(query, courseID, vid.VideoId, vid.VideoTitle, vid.Index, vid.Duration)
		if err != nil {
			return err
		}
	}
	return nil
}

func Insert_Video_Progress(tx *sql.Tx, userID int, courseID int, videoID string) error {
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
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("blocked: user, course, or video is invalid")
	}

	return nil
}
