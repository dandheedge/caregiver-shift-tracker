package models

import (
	"time"
)

// Schedule represents a caregiver's assigned shift
type Schedule struct {
	ID          int       `json:"id" db:"id"`
	ClientName  string    `json:"client_name" db:"client_name"`
	ShiftStart  time.Time `json:"shift_start" db:"shift_start"`
	ShiftEnd    time.Time `json:"shift_end" db:"shift_end"`
	Latitude    float64   `json:"latitude" db:"latitude"`
	Longitude   float64   `json:"longitude" db:"longitude"`
	Status      string    `json:"status" db:"status"` // upcoming, in_progress, completed, missed
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Task represents a care activity assigned to a schedule
type Task struct {
	ID          int    `json:"id" db:"id"`
	ScheduleID  int    `json:"schedule_id" db:"schedule_id"`
	Description string `json:"description" db:"description"`
	Status      string `json:"status" db:"status"` // pending, completed, not_completed
	Reason      string `json:"reason,omitempty" db:"reason"` // reason if not completed
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// Visit represents the actual visit log with timestamps and location
type Visit struct {
	ID         int        `json:"id" db:"id"`
	ScheduleID int        `json:"schedule_id" db:"schedule_id"`
	StartTime  *time.Time `json:"start_time,omitempty" db:"start_time"`
	EndTime    *time.Time `json:"end_time,omitempty" db:"end_time"`
	StartLat   *float64   `json:"start_lat,omitempty" db:"start_lat"`
	StartLng   *float64   `json:"start_lng,omitempty" db:"start_lng"`
	EndLat     *float64   `json:"end_lat,omitempty" db:"end_lat"`
	EndLng     *float64   `json:"end_lng,omitempty" db:"end_lng"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at" db:"updated_at"`
}

// ScheduleWithTasks represents a schedule with its associated tasks
type ScheduleWithTasks struct {
	Schedule
	Tasks []Task `json:"tasks"`
	Visit *Visit `json:"visit,omitempty"`
}

// StartVisitRequest represents the request payload for starting a visit
type StartVisitRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

// EndVisitRequest represents the request payload for ending a visit
type EndVisitRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

// UpdateTaskRequest represents the request payload for updating a task
type UpdateTaskRequest struct {
	Status string `json:"status" binding:"required,oneof=completed not_completed"`
	Reason string `json:"reason,omitempty"`
}

// Activity represents a care activity assigned to a schedule
type Activity struct {
	ID          int       `json:"id" db:"id"`
	ScheduleID  int       `json:"schedule_id" db:"schedule_id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	IsResolved  bool      `json:"is_resolved" db:"is_resolved"`
	Reason      string    `json:"reason,omitempty" db:"reason"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// UpdateActivityRequest represents the request payload for updating an activity
type UpdateActivityRequest struct {
	IsResolved bool   `json:"is_resolved"`
	Reason     string `json:"reason,omitempty"`
}

// CreateActivityRequest represents the request payload for creating an activity
type CreateActivityRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

// StatsResponse represents the dashboard statistics
type StatsResponse struct {
	TotalSchedules    int `json:"total_schedules"`
	MissedSchedules   int `json:"missed_schedules"`
	UpcomingToday     int `json:"upcoming_today"`
	CompletedToday    int `json:"completed_today"`
} 