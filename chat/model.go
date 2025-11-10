package chat

import "time"

// Message - Chat between supervisor and team
type Message struct {
    ID         uint      `json:"id" gorm:"primaryKey"`
    TeamID     uint      `json:"team_id" gorm:"not null;index"`
    SenderID   uint      `json:"sender_id" gorm:"not null"`
    SenderType string    `json:"sender_type" gorm:"not null"` // 'student' or 'supervisor'
    Message    string    `json:"message" gorm:"type:text;not null"`
    CreatedAt  time.Time `json:"created_at"`
}