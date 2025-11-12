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

// UploadFile handles file uploads to Cloudinary following official documentation
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

	// Generate unique filename with extension (critical for PDFs)
	timestamp := time.Now().Unix()
	publicID := fmt.Sprintf("%d%s", timestamp, ext) // Include extension in public ID
	ctx := context.Background()

	// Upload based on file type following Cloudinary documentation
	var uploadResult *uploader.UploadResult
	var params uploader.UploadParams

	// Set common parameters
	params.PublicID = publicID
	params.Folder = "project-management"

	switch ext {
	case ".pdf":
		// ✅ CORRECT APPROACH: PDFs should be uploaded as "image" resource type
		// According to documentation: "Cloudinary treats PDF files the same as any other image file"
		params.ResourceType = "image"
		params.Folder = "project-management/documents"
		
	case ".doc", ".docx", ".txt", ".zip":
		// Office documents and other files should be "raw" for download
		params.ResourceType = "raw"
		params.Folder = "project-management/documents"
		
	case ".mp4":
		params.ResourceType = "video"
		params.Folder = "project-management/videos"
		
	default:
		// Images
		params.ResourceType = "image"
		params.Folder = "project-management/images"
	}

	uploadResult, err = cld.Upload.Upload(ctx, fileHeader, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to upload file: %v", err)})
		return
	}

	// ✅ Transform URL for PDFs to enable inline viewing - CORRECTED APPROACH
	finalURL := transformCloudinaryURL(uploadResult.SecureURL, ext, params.ResourceType)

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"file_url": finalURL,
		"pages":    uploadResult.Pages, // Include page count for PDFs
	})
}

// ✅ CORRECTED transformCloudinaryURL based on Cloudinary documentation
func transformCloudinaryURL(url string, ext string, resourceType string) string {
	if ext == ".pdf" {
		if resourceType == "image" {
			// For PDFs uploaded as image type (recommended approach)
			// This enables inline viewing in browsers
			if strings.Contains(url, "/image/upload/") {
				// Add fl_attachment:false for inline viewing
				return strings.Replace(url, "/upload/", "/upload/fl_attachment:false/", 1)
			}
		} else if resourceType == "raw" {
			// For PDFs uploaded as raw (fallback approach)
			return strings.Replace(url, "/upload/", "/upload/fl_attachment:false/", 1)
		}
	}
	
	// For other raw files that should download
	if resourceType == "raw" && (ext == ".doc" || ext == ".docx" || ext == ".zip") {
		if !strings.Contains(url, "fl_attachment") {
			return strings.Replace(url, "/upload/", "/upload/fl_attachment/", 1)
		}
	}
	
	return url
}