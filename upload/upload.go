package upload

import (
    "fmt"
    "net/http"
    "path/filepath"
    "time"

    "github.com/gin-gonic/gin"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
    return &UploadHandler{}
}

// UploadFile handles file uploads
func (h *UploadHandler) UploadFile(c *gin.Context) {
    // Get file from form
    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
        return
    }

    // Validate file extension
    ext := filepath.Ext(file.Filename)
    allowedExts := map[string]bool{
        ".pdf":  true,
        ".jpg":  true,
        ".jpeg": true,
        ".png":  true,
        ".doc":  true,
        ".docx": true,
    }

    if !allowedExts[ext] {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Only PDF, JPG, PNG, DOC, DOCX files are allowed"})
        return
    }

    // Validate file size (max 10MB)
    if file.Size > 10*1024*1024 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "File size must be less than 10MB"})
        return
    }

    // Generate unique filename with timestamp
    timestamp := time.Now().Unix()
    newFilename := fmt.Sprintf("%d_%s", timestamp, file.Filename)

    // Save file to uploads folder
    savePath := filepath.Join("uploads", newFilename)
    if err := c.SaveUploadedFile(file, savePath); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
        return
    }

    // Return public URL
    fileURL := fmt.Sprintf("http://localhost:8081/uploads/%s", newFilename)

    c.JSON(http.StatusOK, gin.H{
        "message":  "File uploaded successfully",
        "file_url": fileURL,
    })
}