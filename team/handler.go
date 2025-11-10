package team

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
)

type TeamHandler struct {
    service TeamService
}

func NewTeamHandler(service TeamService) *TeamHandler {
    return &TeamHandler{service: service}
}

// ============================================
// Send Team Request
// POST /api/team-requests
// ============================================

func (h *TeamHandler) SendTeamRequest(c *gin.Context) {
    var req SendTeamRequestRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // ✅ Get student ID from JWT claims (set by AuthMiddleware)
    studentID, exists := c.Get("student_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "student id not found in token"})
        return
    }

    res, err := h.service.SendTeamRequest(studentID.(uint), req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, res)
}

// ============================================
// Get Received Requests
// GET /api/team-requests/received
// ============================================

func (h *TeamHandler) GetReceivedRequests(c *gin.Context) {
    // ✅ Get student ID from JWT claims
    studentID, exists := c.Get("student_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "student id not found in token"})
        return
    }

    res, err := h.service.GetReceivedRequests(studentID.(uint))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, res)
}

// ============================================
// Get Sent Requests
// GET /api/team-requests/sent
// ============================================

func (h *TeamHandler) GetSentRequests(c *gin.Context) {
    // ✅ Get student ID from JWT claims
    studentID, exists := c.Get("student_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "student id not found in token"})
        return
    }

    res, err := h.service.GetSentRequests(studentID.(uint))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, res)
}

// ============================================
// Accept Team Request
// POST /api/team-requests/:id/accept
// ============================================

func (h *TeamHandler) AcceptRequest(c *gin.Context) {
    // Get request ID from URL
    idStr := c.Param("id")
    requestID, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
        return
    }

    // ✅ Get student ID from JWT claims
    studentID, exists := c.Get("student_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "student id not found in token"})
        return
    }

    res, err := h.service.AcceptRequest(uint(requestID), studentID.(uint))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, res)
}

// ============================================
// Reject Team Request
// POST /api/team-requests/:id/reject
// ============================================

func (h *TeamHandler) RejectRequest(c *gin.Context) {
    // Get request ID from URL
    idStr := c.Param("id")
    requestID, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
        return
    }

    // ✅ Get student ID from JWT claims
    studentID, exists := c.Get("student_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "student id not found in token"})
        return
    }

    res, err := h.service.RejectRequest(uint(requestID), studentID.(uint))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, res)
}

// ============================================
// Cancel Team Request
// DELETE /api/team-requests/:id
// ============================================

func (h *TeamHandler) CancelRequest(c *gin.Context) {
    // Get request ID from URL
    idStr := c.Param("id")
    requestID, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
        return
    }

    // ✅ Get student ID from JWT claims
    studentID, exists := c.Get("student_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "student id not found in token"})
        return
    }

    res, err := h.service.CancelRequest(uint(requestID), studentID.(uint))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, res)
}

// ============================================
// Get My Team
// GET /api/my-team
// ============================================

func (h *TeamHandler) GetMyTeam(c *gin.Context) {
    // ✅ Get student ID from JWT claims
    studentID, exists := c.Get("student_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "student id not found in token"})
        return
    }

    res, err := h.service.GetMyTeam(studentID.(uint))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, res)
}


type SupervisorRequestHandler struct {
    service SupervisorRequestService
}

func NewSupervisorRequestHandler(service SupervisorRequestService) *SupervisorRequestHandler {
    return &SupervisorRequestHandler{service: service}
}

// ============================================
// STUDENT: Send Supervisor Request
// POST /api/supervisor-requests
// ============================================

func (h *SupervisorRequestHandler) SendSupervisorRequest(c *gin.Context) {
    var req SendSupervisorRequestRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // ✅ Get student ID from JWT claims
    studentID, exists := c.Get("student_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "student id not found in token"})
        return
    }

    res, err := h.service.SendSupervisorRequest(studentID.(uint), req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, res)
}

// ============================================
// STUDENT: Get My Team's Supervisor Requests
// GET /api/supervisor-requests/my
// ============================================

func (h *SupervisorRequestHandler) GetMySupervisorRequests(c *gin.Context) {
    // ✅ Get student ID from JWT claims
    studentID, exists := c.Get("student_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "student id not found in token"})
        return
    }

    res, err := h.service.GetMySupervisorRequests(studentID.(uint))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, res)
}

// ============================================
// SUPERVISOR: Get Pending Requests
// GET /api/supervisor-requests/pending
// ============================================

func (h *SupervisorRequestHandler) GetPendingRequests(c *gin.Context) {
    // ✅ Get supervisor ID from JWT claims
    supervisorID, exists := c.Get("supervisor_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "supervisor id not found in token"})
        return
    }

    res, err := h.service.GetPendingRequests(supervisorID.(uint))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, res)
}

// ============================================
// SUPERVISOR: Accept Request
// POST /api/supervisor-requests/:id/accept
// ============================================

func (h *SupervisorRequestHandler) AcceptRequest(c *gin.Context) {
    // Get request ID from URL
    idStr := c.Param("id")
    requestID, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
        return
    }

    // ✅ Get supervisor ID from JWT claims
    supervisorID, exists := c.Get("supervisor_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "supervisor id not found in token"})
        return
    }

    res, err := h.service.AcceptRequest(uint(requestID), supervisorID.(uint))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, res)
}

// ============================================
// SUPERVISOR: Reject Request
// POST /api/supervisor-requests/:id/reject
// ============================================

func (h *SupervisorRequestHandler) RejectRequest(c *gin.Context) {
    var req RejectSupervisorRequestRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Get request ID from URL
    idStr := c.Param("id")
    requestID, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
        return
    }

    // ✅ Get supervisor ID from JWT claims
    supervisorID, exists := c.Get("supervisor_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "supervisor id not found in token"})
        return
    }

    res, err := h.service.RejectRequest(uint(requestID), supervisorID.(uint), req)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, res)
}