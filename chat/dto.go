package chat

import "time"

// ============================================
// Send Message (Both Supervisor & Student)
// ============================================

type SendMessageRequest struct {
    Message string `json:"message" binding:"required,min=1,max=1000"`
}

type SendMessageResponse struct {
    Message string `json:"message"`
}

// ============================================
// Get Messages (Both Supervisor & Student)
// ============================================

type MessageItem struct {
    ID         uint      `json:"id"`
    SenderName string    `json:"sender_name"`
    SenderType string    `json:"sender_type"` // "student" or "supervisor"
    Message    string    `json:"message"`
    CreatedAt  time.Time `json:"created_at"`
}

type GetMessagesResponse struct {
    Total    int           `json:"total"`
    Messages []MessageItem `json:"messages"`
}