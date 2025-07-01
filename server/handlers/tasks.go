package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"visit-tracker-api/database"
	"visit-tracker-api/models"

	"github.com/gin-gonic/gin"
)

// UpdateTask updates the status of a specific task
func UpdateTask(c *gin.Context) {
	idParam := c.Param("taskId")
	taskID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Validate that reason is provided when marking as not_completed
	if req.Status == "not_completed" && req.Reason == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reason is required when marking task as not completed"})
		return
	}

	// Check if task exists
	var existingStatus string
	var scheduleID int
	err = database.DB.QueryRow("SELECT status, schedule_id FROM tasks WHERE id = ?", taskID).Scan(&existingStatus, &scheduleID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task"})
		return
	}

	// Check if the associated schedule is in progress (visit started)
	var scheduleStatus string
	err = database.DB.QueryRow("SELECT status FROM schedules WHERE id = ?", scheduleID).Scan(&scheduleStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedule status"})
		return
	}

	if scheduleStatus != "in_progress" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot update tasks before starting the visit"})
		return
	}

	// Update task
	_, err = database.DB.Exec(`
		UPDATE tasks 
		SET status = ?, reason = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`,
		req.Status, req.Reason, taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	// Return updated task
	var updatedTask models.Task
	var reason sql.NullString
	var createdAt, updatedAt string

	err = database.DB.QueryRow(`
		SELECT id, schedule_id, description, status, reason, created_at, updated_at
		FROM tasks WHERE id = ?`, taskID).Scan(
		&updatedTask.ID, &updatedTask.ScheduleID, &updatedTask.Description,
		&updatedTask.Status, &reason, &createdAt, &updatedAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated task"})
		return
	}

	if reason.Valid {
		updatedTask.Reason = reason.String
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Task updated successfully",
		"task": updatedTask,
	})
}

// GetTasksBySchedule returns all tasks for a specific schedule
func GetTasksBySchedule(c *gin.Context) {
	idParam := c.Param("id")
	scheduleID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedule ID"})
		return
	}

	// Check if schedule exists
	var count int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM schedules WHERE id = ?", scheduleID).Scan(&count)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify schedule"})
		return
	}

	if count == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	// Get tasks
	rows, err := database.DB.Query(`
		SELECT id, schedule_id, description, status, reason, created_at, updated_at
		FROM tasks
		WHERE schedule_id = ?
		ORDER BY id ASC`, scheduleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var reason sql.NullString
		var createdAt, updatedAt string

		err := rows.Scan(
			&task.ID, &task.ScheduleID, &task.Description, &task.Status,
			&reason, &createdAt, &updatedAt)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse task data"})
			return
		}

		if reason.Valid {
			task.Reason = reason.String
		}

		tasks = append(tasks, task)
	}

	c.JSON(http.StatusOK, tasks)
} 