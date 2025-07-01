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

// parseTime safely parses time strings using standard Go time formats
func parseTime(timeStr string) time.Time {
	layouts := []string{
		time.DateTime,    // "2006-01-02 15:04:05"
		time.RFC3339,     // "2006-01-02T15:04:05Z07:00"
		time.RFC3339Nano, // "2006-01-02T15:04:05.999999999Z07:00"
		time.DateOnly,    // "2006-01-02"
	}
	
	for _, layout := range layouts {
		if t, err := time.Parse(layout, timeStr); err == nil {
			return t
		}
	}
	
	return time.Time{}
}

// GetAllSchedules godoc
// @Summary Get all schedules
// @Description Get a list of all caregiver schedules
// @Tags schedules
// @Accept json
// @Produce json
// @Success 200 {array} models.Schedule
// @Failure 500 {object} map[string]string
// @Router /schedules [get]
func GetAllSchedules(c *gin.Context) {
	query := `
		SELECT s.id, s.client_name, s.shift_start, s.shift_end, s.latitude, s.longitude, s.status, s.created_at, s.updated_at
		FROM schedules s
		ORDER BY s.shift_start ASC`

	rows, err := database.DB.Query(query)
	if err != nil {
		log.Printf("Database query error in GetAllSchedules: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedules"})
		return
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var schedule models.Schedule
		var shiftStart, shiftEnd, createdAt, updatedAt string

		err := rows.Scan(
			&schedule.ID, &schedule.ClientName, &shiftStart, &shiftEnd,
			&schedule.Latitude, &schedule.Longitude, &schedule.Status, &createdAt, &updatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse schedule data"})
			return
		}

		// Parse time strings using flexible parsing
		schedule.ShiftStart = parseTime(shiftStart)
		schedule.ShiftEnd = parseTime(shiftEnd)
		schedule.CreatedAt = parseTime(createdAt)
		schedule.UpdatedAt = parseTime(updatedAt)

		schedules = append(schedules, schedule)
	}

	c.JSON(http.StatusOK, schedules)
}

// GetTodaySchedules godoc
// @Summary Get today's schedules
// @Description Get a list of today's caregiver schedules
// @Tags schedules
// @Accept json
// @Produce json
// @Success 200 {array} models.Schedule
// @Failure 500 {object} map[string]string
// @Router /schedules/today [get]
func GetTodaySchedules(c *gin.Context) {
	today := time.Now().Format("2006-01-02")
	
	query := `
		SELECT s.id, s.client_name, s.shift_start, s.shift_end, s.latitude, s.longitude, s.status, s.created_at, s.updated_at
		FROM schedules s
		WHERE DATE(s.shift_start) = ?
		ORDER BY s.shift_start ASC`

	rows, err := database.DB.Query(query, today)
	if err != nil {
		log.Printf("Database query error in GetTodaySchedules: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch today's schedules"})
		return
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var schedule models.Schedule
		var shiftStart, shiftEnd, createdAt, updatedAt string

		err := rows.Scan(
			&schedule.ID, &schedule.ClientName, &shiftStart, &shiftEnd,
			&schedule.Latitude, &schedule.Longitude, &schedule.Status, &createdAt, &updatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse schedule data"})
			return
		}

		// Parse time strings using flexible parsing
		schedule.ShiftStart = parseTime(shiftStart)
		schedule.ShiftEnd = parseTime(shiftEnd)
		schedule.CreatedAt = parseTime(createdAt)
		schedule.UpdatedAt = parseTime(updatedAt)

		schedules = append(schedules, schedule)
	}

	c.JSON(http.StatusOK, schedules)
}

// GetScheduleByID godoc
// @Summary Get schedule by ID
// @Description Get a specific schedule with its tasks and visit information
// @Tags schedules
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {object} models.ScheduleWithTasks
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /schedules/{id} [get]
func GetScheduleByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}

	// Get schedule
	var scheduleWithTasks models.ScheduleWithTasks
	var shiftStart, shiftEnd, createdAt, updatedAt string

	scheduleQuery := `
		SELECT id, client_name, shift_start, shift_end, latitude, longitude, status, created_at, updated_at
		FROM schedules
		WHERE id = ?`

	err = database.DB.QueryRow(scheduleQuery, id).Scan(
		&scheduleWithTasks.ID, &scheduleWithTasks.ClientName, &shiftStart, &shiftEnd,
		&scheduleWithTasks.Latitude, &scheduleWithTasks.Longitude, &scheduleWithTasks.Status, &createdAt, &updatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
			return
		}
		log.Printf("Database query error in GetScheduleByID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedule"})
		return
	}

	// Parse time strings using flexible parsing
	scheduleWithTasks.ShiftStart = parseTime(shiftStart)
	scheduleWithTasks.ShiftEnd = parseTime(shiftEnd)
	scheduleWithTasks.CreatedAt = parseTime(createdAt)
	scheduleWithTasks.UpdatedAt = parseTime(updatedAt)

	// Get tasks
	tasksQuery := `
		SELECT id, description, status, reason, created_at, updated_at
		FROM tasks
		WHERE schedule_id = ?
		ORDER BY id ASC`

	rows, err := database.DB.Query(tasksQuery, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var reason sql.NullString
		var taskCreatedAt, taskUpdatedAt string

		err := rows.Scan(
			&task.ID, &task.Description, &task.Status, &reason,
			&taskCreatedAt, &taskUpdatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse task data"})
			return
		}

		task.ScheduleID = id
		if reason.Valid {
			task.Reason = reason.String
		}
		task.CreatedAt = parseTime(taskCreatedAt)
		task.UpdatedAt = parseTime(taskUpdatedAt)

		tasks = append(tasks, task)
	}
	scheduleWithTasks.Tasks = tasks

	// Get visit information
	visitQuery := `
		SELECT id, start_time, end_time, start_lat, start_lng, end_lat, end_lng, created_at, updated_at
		FROM visits
		WHERE schedule_id = ?`

	var visit models.Visit
	var startTime, endTime sql.NullString
	var startLat, startLng, endLat, endLng sql.NullFloat64
	var visitCreatedAt, visitUpdatedAt string

	err = database.DB.QueryRow(visitQuery, id).Scan(
		&visit.ID, &startTime, &endTime, &startLat, &startLng, &endLat, &endLng,
		&visitCreatedAt, &visitUpdatedAt,
	)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch visit data"})
		return
	}

	if err != sql.ErrNoRows {
		visit.ScheduleID = id
		if startTime.Valid {
			if t := parseTime(startTime.String); !t.IsZero() {
				visit.StartTime = &t
			}
		}
		if endTime.Valid {
			if t := parseTime(endTime.String); !t.IsZero() {
				visit.EndTime = &t
			}
		}
		if startLat.Valid {
			visit.StartLat = &startLat.Float64
		}
		if startLng.Valid {
			visit.StartLng = &startLng.Float64
		}
		if endLat.Valid {
			visit.EndLat = &endLat.Float64
		}
		if endLng.Valid {
			visit.EndLng = &endLng.Float64
		}
		visit.CreatedAt = parseTime(visitCreatedAt)
		visit.UpdatedAt = parseTime(visitUpdatedAt)

		scheduleWithTasks.Visit = &visit
	}

	c.JSON(http.StatusOK, scheduleWithTasks)
}

// GetStats godoc
// @Summary Get dashboard statistics
// @Description Get statistics for the dashboard including total, missed, upcoming, and completed schedules
// @Tags stats
// @Accept json
// @Produce json
// @Success 200 {object} models.StatsResponse
// @Failure 500 {object} map[string]string
// @Router /stats [get]
func GetStats(c *gin.Context) {
	today := time.Now().Format("2006-01-02")
	
	var stats models.StatsResponse

	// Total schedules
	err := database.DB.QueryRow("SELECT COUNT(*) FROM schedules").Scan(&stats.TotalSchedules)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch total schedules"})
		return
	}

	// Missed schedules
	err = database.DB.QueryRow("SELECT COUNT(*) FROM schedules WHERE status = 'missed'").Scan(&stats.MissedSchedules)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch missed schedules"})
		return
	}

	// Upcoming today
	err = database.DB.QueryRow(
		"SELECT COUNT(*) FROM schedules WHERE DATE(shift_start) = ? AND status = 'upcoming'",
		today).Scan(&stats.UpcomingToday)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch upcoming schedules"})
		return
	}

	// Completed today
	err = database.DB.QueryRow(
		"SELECT COUNT(*) FROM schedules WHERE DATE(shift_start) = ? AND status = 'completed'",
		today).Scan(&stats.CompletedToday)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch completed schedules"})
		return
	}

	c.JSON(http.StatusOK, stats)
} 