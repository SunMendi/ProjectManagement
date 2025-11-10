package upload

import (
    "ProjectManagement/middleware"
    "github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
    handler := NewUploadHandler()

    // Protected upload endpoint (requires login)
    protected := router.Group("/api")
    protected.Use(middleware.AuthMiddleware())
    {
        protected.POST("/upload", handler.UploadFile)
    }

    // Serve uploaded files publicly
    router.Static("/uploads", "./uploads")
}