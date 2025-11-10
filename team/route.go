package team

import (
    "ProjectManagement/middleware"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
    // Initialize repositories
    teamRepo := NewTeamRepository(db)
    requestRepo := NewTeamRequestRepository(db)
    supervisorReqRepo := NewSupervisorRequestRepository(db) // ✅ NEW

    // Initialize services
    teamService := NewTeamService(db, teamRepo, requestRepo)
    supervisorReqService := NewSupervisorRequestService(db, teamRepo, supervisorReqRepo) // ✅ NEW

    // Initialize handlers
    teamHandler := NewTeamHandler(teamService)
    supervisorReqHandler := NewSupervisorRequestHandler(supervisorReqService) // ✅ NEW

    // ✅ Protected routes (auth required)
    protected := router.Group("/api")
    protected.Use(middleware.AuthMiddleware())
    {
        // ========== Team Request Routes ==========
        protected.POST("/team-requests", teamHandler.SendTeamRequest)
        protected.GET("/team-requests/received", teamHandler.GetReceivedRequests)
        protected.GET("/team-requests/sent", teamHandler.GetSentRequests)
        protected.POST("/team-requests/:id/accept", teamHandler.AcceptRequest)
        protected.POST("/team-requests/:id/reject", teamHandler.RejectRequest)
        protected.DELETE("/team-requests/:id", teamHandler.CancelRequest)
        
        // ========== Team Routes ==========
        protected.GET("/my-team", teamHandler.GetMyTeam)

        // ========== Supervisor Request Routes (STUDENT) ========== ✅ NEW
        protected.POST("/supervisor-requests", supervisorReqHandler.SendSupervisorRequest)
        protected.GET("/supervisor-requests/my", supervisorReqHandler.GetMySupervisorRequests)

        // ========== Supervisor Request Routes (SUPERVISOR) ========== ✅ NEW
        protected.GET("/supervisor-requests/pending", supervisorReqHandler.GetPendingRequests)
        protected.POST("/supervisor-requests/:id/accept", supervisorReqHandler.AcceptRequest)
        protected.POST("/supervisor-requests/:id/reject", supervisorReqHandler.RejectRequest)
    }
}