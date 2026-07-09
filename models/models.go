package models

import "time"

type User_stats struct {
	User_id          int
	Streak_count     int
	Total_points     int
	Last_active_date time.Time
}

type Course_data struct {
	ID      int
	OwnerID *int
	Name    string
	Email   string
	URL     *string
	Title   *string
}

type Progress_data struct {
	id                 int
	User_id            int
	Course_id          int
	Completion_percent int
	Started_At         time.Time
	Last_accessed_date time.Time
}
