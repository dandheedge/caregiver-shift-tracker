package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"visit-tracker-api/database"
	"visit-tracker-api/models"
	"visit-tracker-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// StartVisit godoc
// @Summary Start a visit
// @Description Start a caregiver visit by logging timestamp and geolocation
// @Tags visits
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Param startVisitRequest body models.StartVisitRequest true "Start visit data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /schedules/{id}/start [post]
func StartVisit(c *gin.Context) {
	// Log the start of the operation
	utils.LogInfo("Starting visit", logrus.Fields{
		"request_id": c.GetString("request_id"),
		"path":       c.Request.URL.Path,
	})

	idParam := c.Param("id")
	scheduleID, err := strconv.Atoi(idParam)
	if err != nil {
		utils.HandleValidationError(c, err, "schedule_id")
		return
	}

	var req models.StartVisitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleValidationError(c, err, "request_body")
		return
	}

	// Validate coordinates
	if req.Latitude < -90 || req.Latitude > 90 || req.Longitude < -180 || req.Longitude > 180 {
		utils.HandleValidationError(c, 
			&ValidationError{Field: "coordinates", Message: "Invalid latitude or longitude"}, 
			"coordinates")
		return
	}

	// Check if schedule exists and is not already started
	var currentStatus string
	err = database.DB.QueryRow("SELECT status FROM schedules WHERE id = ?", scheduleID).Scan(&currentStatus)
	if err != nil {
		utils.HandleDatabaseError(c, err, "get_schedule_status")
		return
	}

	if currentStatus == "completed" {
		utils.HandleValidationError(c, 
			&ValidationError{Field: "status", Message: "Visit already completed"}, 
			"visit_status")
		return
	}

	if currentStatus == "in_progress" {
		utils.HandleValidationError(c, 
			&ValidationError{Field: "status", Message: "Visit already started"}, 
			"visit_status")
		return
	}

	// Start transaction
	tx, err := database.DB.Begin()
	if err != nil {
		utils.HandleDatabaseError(c, err, "begin_transaction")
		return
	}
	defer tx.Rollback()

	// Update visit record with start time and location
	now := time.Now()
	_, err = tx.Exec(`
		UPDATE visits 
		SET start_time = ?, start_lat = ?, start_lng = ?, updated_at = CURRENT_TIMESTAMP
		WHERE schedule_id = ?`,
		now.Format("2006-01-02 15:04:05"), req.Latitude, req.Longitude, scheduleID)
	if err != nil {
		utils.HandleDatabaseError(c, err, "update_visit_record")
		return
	}

	// Update schedule status to in_progress
	_, err = tx.Exec(`
		UPDATE schedules 
		SET status = 'in_progress', updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`, scheduleID)
	if err != nil {
		utils.HandleDatabaseError(c, err, "update_schedule_status")
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		utils.HandleDatabaseError(c, err, "commit_transaction")
		return
	}

	// Log successful operation
	utils.LogInfo("Visit started successfully", logrus.Fields{
		"request_id":  c.GetString("request_id"),
		"schedule_id": scheduleID,
		"latitude":    req.Latitude,
		"longitude":   req.Longitude,
	})

	// Return success response
	utils.JSONSuccess(c, gin.H{
		"message": "Visit started successfully",
		"timestamp": now,
		"location": gin.H{
			"latitude":  req.Latitude,
			"longitude": req.Longitude,
		},
	})
}

// EndVisit godoc
// @Summary End a visit
// @Description End a caregiver visit by logging timestamp and geolocation
// @Tags visits
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Param endVisitRequest body models.EndVisitRequest true "End visit data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /schedules/{id}/end [post]
func EndVisit(c *gin.Context) {
	idParam := c.Param("id")
	scheduleID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}

	var req models.EndVisitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Check if schedule exists and is in progress
	var currentStatus string
	err = database.DB.QueryRow("SELECT status FROM schedules WHERE id = ?", scheduleID).Scan(&currentStatus)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedule"})
		return
	}

	if currentStatus != "in_progress" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Visit not started yet or already completed"})
		return
	}

	// Check if visit has start time
	var startTime sql.NullString
	err = database.DB.QueryRow("SELECT start_time FROM visits WHERE schedule_id = ?", scheduleID).Scan(&startTime)
	if err != nil || !startTime.Valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Visit not properly started"})
		return
	}

	// Start transaction
	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	// Update visit record with end time and location
	now := time.Now()
	_, err = tx.Exec(`
		UPDATE visits 
		SET end_time = ?, end_lat = ?, end_lng = ?, updated_at = CURRENT_TIMESTAMP
		WHERE schedule_id = ?`,
		now.Format("2006-01-02 15:04:05"), req.Latitude, req.Longitude, scheduleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update visit record"})
		return
	}

	// Update schedule status to completed
	_, err = tx.Exec(`
		UPDATE schedules 
		SET status = 'completed', updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`, scheduleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update schedule status"})
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Calculate visit duration
	startTimeObj, _ := time.Parse("2006-01-02 15:04:05", startTime.String)
	duration := now.Sub(startTimeObj)

	c.JSON(http.StatusOK, gin.H{
		"message": "Visit ended successfully",
		"start_time": startTimeObj,
		"end_time": now,
		"duration_minutes": int(duration.Minutes()),
		"end_location": gin.H{
			"latitude":  req.Latitude,
			"longitude": req.Longitude,
		},
	})
} 