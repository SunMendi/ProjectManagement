package upload

import (
    "ProjectManagement/middleware"
    "github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
    // Protected upload endpoint
    protected := router.Group("/api")
    protected.Use(middleware.AuthMiddleware())
    {
        protected.POST("/upload", UploadFile)
    }
}