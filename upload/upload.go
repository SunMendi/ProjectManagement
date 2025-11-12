package upload

import (
    "context"
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "strings"
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

    // Validate file extension (case-insensitive)
    ext := strings.ToLower(filepath.Ext(file.Filename))
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

    // Get Cloudinary credentials
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

    // Generate unique timestamp
    timestamp := time.Now().Unix()
    ctx := context.Background()

    // Upload based on file type
    var uploadResult *uploader.UploadResult
    var fileURL string

    if ext == ".pdf" {
        // ✅ PDFs - upload as IMAGE type
        uploadResult, err = cld.Upload.Upload(ctx, fileHeader, uploader.UploadParams{
            PublicID:     fmt.Sprintf("%d", timestamp),
            Folder:       "project-management/documents",
            ResourceType: "image",
        })

        if err == nil {
            // ✅ Return both PDF and PNG URLs for flexibility
            pdfURL := uploadResult.SecureURL
            // Convert to PNG for guaranteed viewing (first page)
            pngURL := strings.Replace(pdfURL, ".pdf", ".png", 1)
            
            fileURL = pngURL // Use PNG for guaranteed compatibility
        }

    } else if ext == ".doc" || ext == ".docx" || ext == ".txt" || ext == ".zip" {
        // Other documents - use "raw" (will download)
        uploadResult, err = cld.Upload.Upload(ctx, fileHeader, uploader.UploadParams{
            PublicID:     fmt.Sprintf("%d", timestamp),
            Folder:       "project-management/documents",
            ResourceType: "raw",
        })
        fileURL = uploadResult.SecureURL

    } else if ext == ".mp4" {
        // Videos
        uploadResult, err = cld.Upload.Upload(ctx, fileHeader, uploader.UploadParams{
            PublicID:     fmt.Sprintf("%d", timestamp),
            Folder:       "project-management/videos",
            ResourceType: "video",
        })
        fileURL = uploadResult.SecureURL

    } else {
        // Images
        uploadResult, err = cld.Upload.Upload(ctx, fileHeader, uploader.UploadParams{
            PublicID:     fmt.Sprintf("%d", timestamp),
            Folder:       "project-management/images",
            ResourceType: "image",
        })
        fileURL = uploadResult.SecureURL
    }

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload file: %v", err)})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "message":  "File uploaded successfully",
        "file_url": fileURL,
    })
}