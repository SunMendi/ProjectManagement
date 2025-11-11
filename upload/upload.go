package upload

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "time"

    "github.com/cloudinary/cloudinary-go/v2"
    "github.com/cloudinary/cloudinary-go/v2/api/uploader"
    "github.com/gin-gonic/gin"
)

// UploadFile handles file uploads to Cloudinary
func UploadFile(c *gin.Context) {
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
        ".gif":  true,
        ".doc":  true,
        ".docx": true,
        ".mp4":  true,
        ".zip":  true,
        ".txt":  true,
    }

    if !allowedExts[ext] {
        c.JSON(http.StatusBadRequest, gin.H{"error": "File type not allowed"})
        return
    }

    // Validate file size (max 10MB)
    if file.Size > 10*1024*1024 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "File size must be less than 10MB"})
        return
    }

    // Open file
    fileHeader, err := file.Open()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
        return
    }
    defer fileHeader.Close()

    // Get Cloudinary credentials from environment
    cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
    apiKey := os.Getenv("CLOUDINARY_API_KEY")
    apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

    if cloudName == "" || apiKey == "" || apiSecret == "" {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Cloudinary credentials not configured"})
        return
    }

    // Initialize Cloudinary
    cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize Cloudinary"})
        return
    }

    // Determine resource type based on extension
    resourceType := "auto"
    if ext == ".pdf" || ext == ".doc" || ext == ".docx" || ext == ".txt" || ext == ".zip" {
        resourceType = "raw"
    } else if ext == ".mp4" {
        resourceType = "video"
    } else {
        resourceType = "image"
    }

    // Generate unique filename
    timestamp := time.Now().Unix()
    publicID := fmt.Sprintf("%d_%s", timestamp, filepath.Base(file.Filename))

    // Upload to Cloudinary
    ctx := context.Background()
    uploadResult, err := cld.Upload.Upload(ctx, fileHeader, uploader.UploadParams{
        Folder:       "project-management",
        ResourceType: resourceType,
        PublicID:     publicID,
    })

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload file: %v", err)})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":  "File uploaded successfully",
        "file_url": uploadResult.SecureURL,
    })
}