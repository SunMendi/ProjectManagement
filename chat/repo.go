package chat

import (
    "gorm.io/gorm"
)

// ============================================
// Message Repository Interface
// ============================================

type MessageRepository interface {
    // Create message
    Create(message *Message) error
    
    // Get all messages for a team
    FindByTeamID(teamID uint) ([]Message, error)
    
    // Count messages for a team
    CountByTeamID(teamID uint) (int64, error)
}

// ============================================
// Repository Implementation
// ============================================

type messageRepository struct {
    db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
    return &messageRepository{db: db}
}

func (r *messageRepository) Create(message *Message) error {
    return r.db.Create(message).Error
}

func (r *messageRepository) FindByTeamID(teamID uint) ([]Message, error) {
    var messages []Message
    err := r.db.Where("team_id = ?", teamID).
        Order("created_at ASC").
        Find(&messages).Error
    return messages, err
}

func (r *messageRepository) CountByTeamID(teamID uint) (int64, error) {
    var count int64
    err := r.db.Model(&Message{}).
        Where("team_id = ?", teamID).
        Count(&count).Error
    return count, err
}