package models

import (
	"time"
)

type User_stats struct {
	User_id          int
	Streak_count     int
	Total_points     int
	Last_active_date time.Time
}

// Course_data struct represents the structure of a course in the system.
// It contains fields for the course ID, owner ID, name, email, URL, and title.
// The OwnerID and URL fields are pointers to allow for null values.
type Course_data struct {
	ID      int
	OwnerID *int
	Name    string
	Email   string
	URL     *string
	Title   *string
}

// used by the List_my_courses function to return a list of courses along with user information.
type Progress_data struct {
	id                 int
	User_id            int
	Course_id          int
	Completion_percent int
	Started_At         time.Time
	Last_accessed_date time.Time
}

// CreateUserRequest struct is used to represent the request payload for creating a new user.
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

// CourseCreationRequest struct is used to represent the request payload for creating a new course.
type CourseCreationRequest struct {
	OwnerID int    `json:"user_id"`
	URL     string `json:"url"`
	Title   string `json:"title"`
}

type StartCourseRequest struct {
	UserID   int
	CourseID int
}

type UpdateProgress struct {
	UserID        int
	CourseID      int
	NewPercentage int
}

type PointsTOUpdate struct {
	UserID             int
	PointstoAdd        int
	IsFirstActionToday bool
}

type DeleteUserRequest struct {
	UserID int
}
