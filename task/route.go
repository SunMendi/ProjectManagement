package task

import (
    "ProjectManagement/middleware"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
    // Initialize repositories
    taskRepo := NewTaskRepository(db)
    submissionRepo := NewTaskSubmissionRepository(db)

    // Initialize service
    taskService := NewTaskService(db, taskRepo, submissionRepo)

    // Initialize handler
    taskHandler := NewTaskHandler(taskService)

    // ✅ Protected routes (auth required)
    protected := router.Group("/api")
    protected.Use(middleware.AuthMiddleware())
    {
        // ========== SUPERVISOR ROUTES ==========
        protected.GET("/supervisor/my-teams", taskHandler.GetMyTeams)
        protected.POST("/teams/:team_id/tasks", taskHandler.CreateTask)
        protected.GET("/teams/:team_id/tasks", taskHandler.GetTeamTasks)
        protected.DELETE("/tasks/:task_id", taskHandler.DeleteTask)
        protected.POST("/submissions/:submission_id/review", taskHandler.ReviewSubmission)

        // ========== STUDENT ROUTES ==========
        protected.GET("/my-team/tasks", taskHandler.GetMyTeamTasks)
        protected.POST("/tasks/:task_id/submit", taskHandler.SubmitTask)
    }
}