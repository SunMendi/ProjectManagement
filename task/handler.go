package task

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
)

// ============================================
// Task Handler
// ============================================

type TaskHandler struct {
    service TaskService
}

func NewTaskHandler(service TaskService) *TaskHandler {
    return &TaskHandler{
        service: service,
    }
}

// ============================================
// SUPERVISOR: Get My Teams
// ============================================

func (h *TaskHandler) GetMyTeams(c *gin.Context) {
    // Get supervisor ID from JWT middleware
    supervisorID, exists := c.Get("supervisor_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    response, err := h.service.GetMyTeams(supervisorID.(uint))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, response)
}

// ============================================
// SUPERVISOR: Create Task
// ============================================

func (h *TaskHandler) CreateTask(c *gin.Context) {
    // Get supervisor ID from JWT middleware
    supervisorID, exists := c.Get("supervisor_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // Get team ID from URL parameter
    teamIDParam := c.Param("team_id")
    teamID, err := strconv.ParseUint(teamIDParam, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
        return
    }

    // Bind request
    var req CreateTaskRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Create task
    response, err := h.service.CreateTask(supervisorID.(uint), uint(teamID), req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, response)
}

// ============================================
// SUPERVISOR: Get Team Tasks
// ============================================

func (h *TaskHandler) GetTeamTasks(c *gin.Context) {
    // Get supervisor ID from JWT middleware
    supervisorID, exists := c.Get("supervisor_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // Get team ID from URL parameter
    teamIDParam := c.Param("team_id")
    teamID, err := strconv.ParseUint(teamIDParam, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID"})
        return
    }

    // Get tasks
    response, err := h.service.GetTeamTasks(supervisorID.(uint), uint(teamID))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, response)
}

// ============================================
// SUPERVISOR: Delete Task
// ============================================

func (h *TaskHandler) DeleteTask(c *gin.Context) {
    // Get supervisor ID from JWT middleware
    supervisorID, exists := c.Get("supervisor_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // Get task ID from URL parameter
    taskIDParam := c.Param("task_id")
    taskID, err := strconv.ParseUint(taskIDParam, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
        return
    }

    // Delete task
    err = h.service.DeleteTask(supervisorID.(uint), uint(taskID))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// ============================================
// STUDENT: Get My Team Tasks
// ============================================

func (h *TaskHandler) GetMyTeamTasks(c *gin.Context) {
    // Get student ID from JWT middleware
    studentID, exists := c.Get("student_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // Get tasks
    response, err := h.service.GetMyTeamTasks(studentID.(uint))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, response)
}

// ============================================
// STUDENT: Submit Task
// ============================================

func (h *TaskHandler) SubmitTask(c *gin.Context) {
    // Get student ID from JWT middleware
    studentID, exists := c.Get("student_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // Get task ID from URL parameter
    taskIDParam := c.Param("task_id")
    taskID, err := strconv.ParseUint(taskIDParam, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
        return
    }

    // Bind request
    var req SubmitTaskRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Submit task
    response, err := h.service.SubmitTask(studentID.(uint), uint(taskID), req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, response)
}

// ...existing code...

// ============================================
// SUPERVISOR: Review Submission
// ============================================

func (h *TaskHandler) ReviewSubmission(c *gin.Context) {
    supervisorID, exists := c.Get("supervisor_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    submissionIDParam := c.Param("submission_id")
    submissionID, err := strconv.ParseUint(submissionIDParam, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
        return
    }

    var req ReviewSubmissionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    response, err := h.service.ReviewSubmission(supervisorID.(uint), uint(submissionID), req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, response)
}