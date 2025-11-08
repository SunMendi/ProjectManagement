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

    // Initialize service
    teamService := NewTeamService(db, teamRepo, requestRepo)

    // Initialize handler
    teamHandler := NewTeamHandler(teamService)

    // ✅ Protected routes (auth required)
    protected := router.Group("/api")
    protected.Use(middleware.AuthMiddleware()) // All routes need authentication
    {
        // Team request routes
        protected.POST("/team-requests", teamHandler.SendTeamRequest)
        protected.GET("/team-requests/received", teamHandler.GetReceivedRequests)
        protected.GET("/team-requests/sent", teamHandler.GetSentRequests)
        protected.POST("/team-requests/:id/accept", teamHandler.AcceptRequest)
        protected.POST("/team-requests/:id/reject", teamHandler.RejectRequest)
        protected.DELETE("/team-requests/:id", teamHandler.CancelRequest)

        // Team routes
        protected.GET("/my-team", teamHandler.GetMyTeam)
    }
}