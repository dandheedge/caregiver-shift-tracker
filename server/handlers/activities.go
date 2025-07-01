package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"visit-tracker-api/database"
	"visit-tracker-api/models"

	"github.com/gin-gonic/gin"
)

// GetActivityByID godoc
// @Summary Get activity by ID
// @Description Get a specific activity by its ID
// @Tags activities
// @Accept json
// @Produce json
// @Param id path int true "Activity ID"
// @Success 200 {object} models.Activity
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /activities/{id} [get]
func GetActivityByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	var activity models.Activity
	var createdAt, updatedAt string
	var reason sql.NullString

	query := `
		SELECT id, schedule_id, title, description, is_resolved, reason, created_at, updated_at
		FROM activities
		WHERE id = ?`

	err = database.DB.QueryRow(query, id).Scan(
		&activity.ID, &activity.ScheduleID, &activity.Title, &activity.Description,
		&activity.IsResolved, &reason, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
			return
		}
		log.Printf("Database query error in GetActivityByID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activity"})
		return
	}

	// Handle nullable reason field
	if reason.Valid {
		activity.Reason = reason.String
	}

	// Parse time strings using the same function from schedules.go
	activity.CreatedAt = parseTime(createdAt)
	activity.UpdatedAt = parseTime(updatedAt)

	c.JSON(http.StatusOK, activity)
}

// GetActivitiesBySchedule godoc
// @Summary Get activities by schedule ID
// @Description Get all activities for a specific schedule
// @Tags activities
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {array} models.Activity
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /schedules/{id}/activities [get]
func GetActivitiesBySchedule(c *gin.Context) {
	scheduleIDParam := c.Param("id")
	scheduleID, err := strconv.Atoi(scheduleIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}

	query := `
		SELECT id, schedule_id, title, description, is_resolved, reason, created_at, updated_at
		FROM activities
		WHERE schedule_id = ?
		ORDER BY created_at ASC`

	rows, err := database.DB.Query(query, scheduleID)
	if err != nil {
		log.Printf("Database query error in GetActivitiesBySchedule: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activities"})
		return
	}
	defer rows.Close()

	var activities []models.Activity
	for rows.Next() {
		var activity models.Activity
		var createdAt, updatedAt string
		var reason sql.NullString

		err := rows.Scan(
			&activity.ID, &activity.ScheduleID, &activity.Title, &activity.Description,
			&activity.IsResolved, &reason, &createdAt, &updatedAt,
		)
		if err != nil {
			log.Printf("Error scanning activity data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse activity data"})
			return
		}

		// Handle nullable reason field
		if reason.Valid {
			activity.Reason = reason.String
		}

		// Parse time strings
		activity.CreatedAt = parseTime(createdAt)
		activity.UpdatedAt = parseTime(updatedAt)

		activities = append(activities, activity)
	}

	c.JSON(http.StatusOK, activities)
}

// CreateActivity godoc
// @Summary Create a new activity
// @Description Create a new activity for a specific schedule
// @Tags activities
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Param activity body models.CreateActivityRequest true "Activity data"
// @Success 201 {object} models.Activity
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /schedules/{id}/activities [post]
func CreateActivity(c *gin.Context) {
	scheduleIDParam := c.Param("id")
	scheduleID, err := strconv.Atoi(scheduleIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}

	var req models.CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify that the schedule exists
	var exists int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM schedules WHERE id = ?", scheduleID).Scan(&exists)
	if err != nil {
		log.Printf("Database error checking schedule existence: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if exists == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Schedule not found"})
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	result, err := database.DB.Exec(`
		INSERT INTO activities (schedule_id, title, description, is_resolved, created_at, updated_at)
		VALUES (?, ?, ?, 0, ?, ?)`,
		scheduleID, req.Title, req.Description, now, now)
	if err != nil {
		log.Printf("Database error creating activity: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create activity"})
		return
	}

	activityID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting last insert ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get activity ID"})
		return
	}

	// Return the created activity
	activity := models.Activity{
		ID:          int(activityID),
		ScheduleID:  scheduleID,
		Title:       req.Title,
		Description: req.Description,
		IsResolved:  false,
		Reason:      "",
		CreatedAt:   parseTime(now),
		UpdatedAt:   parseTime(now),
	}

	c.JSON(http.StatusCreated, activity)
}

// UpdateActivity godoc
// @Summary Update activity progress
// @Description Update the resolution status of an activity
// @Tags activities
// @Accept json
// @Produce json
// @Param id path int true "Activity ID"
// @Param activity body models.UpdateActivityRequest true "Activity update data"
// @Success 200 {object} models.Activity
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /activities/{id} [put]
func UpdateActivity(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	var req models.UpdateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate that if is_resolved is false, reason is required
	if !req.IsResolved && req.Reason == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reason is required when activity is not resolved"})
		return
	}

	// Check if activity exists
	var exists int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM activities WHERE id = ?", id).Scan(&exists)
	if err != nil {
		log.Printf("Database error checking activity existence: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if exists == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	_, err = database.DB.Exec(`
		UPDATE activities 
		SET is_resolved = ?, reason = ?, updated_at = ?
		WHERE id = ?`,
		req.IsResolved, req.Reason, now, id)
	if err != nil {
		log.Printf("Database error updating activity: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update activity"})
		return
	}

	// Fetch and return the updated activity
	var activity models.Activity
	var createdAt, updatedAt string
	var reason sql.NullString

	query := `
		SELECT id, schedule_id, title, description, is_resolved, reason, created_at, updated_at
		FROM activities
		WHERE id = ?`

	err = database.DB.QueryRow(query, id).Scan(
		&activity.ID, &activity.ScheduleID, &activity.Title, &activity.Description,
		&activity.IsResolved, &reason, &createdAt, &updatedAt,
	)
	if err != nil {
		log.Printf("Database error fetching updated activity: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated activity"})
		return
	}

	// Handle nullable reason field
	if reason.Valid {
		activity.Reason = reason.String
	}

	// Parse time strings
	activity.CreatedAt = parseTime(createdAt)
	activity.UpdatedAt = parseTime(updatedAt)

	c.JSON(http.StatusOK, activity)
} 