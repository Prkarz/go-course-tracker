package models

import (
	"time"
)

type User_stats struct {
	User_id          int       `json:"user_id"`
	Streak_count     int       `json:"streak_count"`
	Total_points     int       `json:"total_points"`
	Last_active_date time.Time `json:"last_active_date"`
}

// Course_data struct represents the structure of a course in the system.
// It contains fields for the course ID, owner ID, name, email, URL, and title.
// The OwnerID and URL fields are pointers to allow for null values.
type Course_data struct {
	ID      int     `json:"id"`
	OwnerID *int    `json:"owner_id,omitempty"`
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	URL     *string `json:"url,omitempty"`
	Title   *string `json:"title,omitempty"`
}

// used by the List_my_courses function to return a list of courses along with user information.
type Progress_data struct {
	Course_id          int       `json:"course_id"`
	Completion_percent int       `json:"completion_percent"`
	Started_At         time.Time `json:"started_at"`
	Last_accessed_date time.Time `json:"last_accessed_date"`
}

// CreateUserRequest struct is used to represent the request payload for creating a new user.
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// CourseCreationRequest struct is used to represent the request payload for creating a new course.
type CourseCreationRequest struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

type StartCourseRequest struct {
	CourseID int `json:"course_id"`
}

type UpdateProgress struct {
	CourseID      int `json:"course_id"`
	NewPercentage int `json:"percentage_to_add"`
}

type PointsTOUpdate struct {
	PointstoAdd        int  `json:"points_to_add"`
	IsFirstActionToday bool `json:"is_first_action_today"`
}

type DeleteUserRequest struct {
	DeletedAt *time.Time
}

type LoginUserRequest struct {
	Email string `json:"email"`
}
