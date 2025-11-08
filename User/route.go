package user

import (
    "ProjectManagement/middleware"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB) {
    userRepo := NewUserRepository(db)
    studentRepo := NewStudentRepository(db)
    supervisorRepo := NewSupervisorRepository(db)

    // Student
    studentService := NewStudentService(db, userRepo, studentRepo)
    studentHandler := NewStudentHandler(studentService)

    // Supervisor
    supervisorService := NewSupervisorService(db, userRepo, supervisorRepo)
    supervisorHandler := NewSupervisorHandler(supervisorService)

    // Public routes (no auth required)
    router.POST("/api/auth/register/student", studentHandler.RegisterStudent)
    router.POST("/api/auth/login/student", studentHandler.Login)
    router.POST("/api/auth/register/supervisor", supervisorHandler.RegisterSupervisor)
    router.POST("/api/auth/login/supervisor", supervisorHandler.Login)

    // Protected routes (auth required)
    protected := router.Group("/api")
    protected.Use(middleware.AuthMiddleware())
    {
        // Student routes
        protected.GET("/students/:id", studentHandler.GetStudentByID)
        protected.PATCH("/students/:id", studentHandler.UpdateStudentProfile)
        protected.GET("/students", studentHandler.GetStudentsByFilter)

        // Supervisor routes
        protected.GET("/supervisors/:id", supervisorHandler.GetSupervisorByID)
        protected.PATCH("/supervisors/:id", supervisorHandler.UpdateSupervisorProfile)
        protected.GET("/supervisors", supervisorHandler.GetSupervisorsByDepartment)
    }
}