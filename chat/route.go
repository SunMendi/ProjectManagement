package chat

import (
    "ProjectManagement/middleware"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
    // Initialize repository
    msgRepo := NewMessageRepository(db)

    // Initialize service
    msgService := NewMessageService(db, msgRepo)

    // Initialize handler
    msgHandler := NewMessageHandler(msgService)

    // ✅ Protected routes (auth required)
    protected := router.Group("/api")
    protected.Use(middleware.AuthMiddleware())
    {
        // ========== SUPERVISOR ROUTES ==========
        protected.POST("/teams/:team_id/messages", msgHandler.SendMessageAsSupervisor)
        protected.GET("/teams/:team_id/messages", msgHandler.GetTeamMessages)

        // ========== STUDENT ROUTES ==========
        protected.POST("/my-team/messages", msgHandler.SendMessageAsStudent)
        protected.GET("/my-team/messages", msgHandler.GetMyTeamMessages)
    }
}