package chat

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
)

// ============================================
// Message Handler
// ============================================

type MessageHandler struct {
    service MessageService
}

func NewMessageHandler(service MessageService) *MessageHandler {
    return &MessageHandler{
        service: service,
    }
}

// ============================================
// SUPERVISOR: Send Message to Team
// ============================================

func (h *MessageHandler) SendMessageAsSupervisor(c *gin.Context) {
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
    var req SendMessageRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Send message
    response, err := h.service.SendMessageAsSupervisor(supervisorID.(uint), uint(teamID), req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, response)
}

// ============================================
// SUPERVISOR: Get Team Messages
// ============================================

func (h *MessageHandler) GetTeamMessages(c *gin.Context) {
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

    // Get messages
    response, err := h.service.GetTeamMessages(supervisorID.(uint), uint(teamID))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, response)
}

// ============================================
// STUDENT: Send Message to Supervisor
// ============================================

func (h *MessageHandler) SendMessageAsStudent(c *gin.Context) {
    // Get student ID from JWT middleware
    studentID, exists := c.Get("student_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // Bind request
    var req SendMessageRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Send message
    response, err := h.service.SendMessageAsStudent(studentID.(uint), req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, response)
}

// ============================================
// STUDENT: Get My Team Messages
// ============================================

func (h *MessageHandler) GetMyTeamMessages(c *gin.Context) {
    // Get student ID from JWT middleware
    studentID, exists := c.Get("student_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    // Get messages
    response, err := h.service.GetMyTeamMessages(studentID.(uint))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, response)
}