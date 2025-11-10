package chat

import (
    "errors"
    "time"

    "gorm.io/gorm"
)

// ============================================
// Message Service Interface
// ============================================

type MessageService interface {
    // Supervisor operations
    SendMessageAsSupervisor(supervisorID, teamID uint, req SendMessageRequest) (*SendMessageResponse, error)
    GetTeamMessages(supervisorID, teamID uint) (*GetMessagesResponse, error)
    
    // Student operations
    SendMessageAsStudent(studentID uint, req SendMessageRequest) (*SendMessageResponse, error)
    GetMyTeamMessages(studentID uint) (*GetMessagesResponse, error)
}

// ============================================
// Service Implementation
// ============================================

type messageService struct {
    db      *gorm.DB
    msgRepo MessageRepository
}

func NewMessageService(db *gorm.DB, msgRepo MessageRepository) MessageService {
    return &messageService{
        db:      db,
        msgRepo: msgRepo,
    }
}

// ============================================
// SUPERVISOR: Send Message to Team
// ============================================

func (s *messageService) SendMessageAsSupervisor(supervisorID, teamID uint, req SendMessageRequest) (*SendMessageResponse, error) {
    // Verify team belongs to this supervisor
    var team struct {
        ID           uint
        SupervisorID *uint
    }
    
    err := s.db.Table("teams").
        Select("id, supervisor_id").
        Where("id = ?", teamID).
        Scan(&team).Error
    
    if err != nil || team.ID == 0 {
        return nil, errors.New("team not found")
    }
    
    if team.SupervisorID == nil || *team.SupervisorID != supervisorID {
        return nil, errors.New("you are not the supervisor of this team")
    }
    
    // Create message
    message := &Message{
        TeamID:     teamID,
        SenderID:   supervisorID,
        SenderType: "supervisor",
        Message:    req.Message,
        CreatedAt:  time.Now(),
    }
    
    err = s.msgRepo.Create(message)
    if err != nil {
        return nil, err
    }
    
    return &SendMessageResponse{
        Message: "Message sent successfully",
    }, nil
}

// ============================================
// SUPERVISOR: Get Team Messages
// ============================================

func (s *messageService) GetTeamMessages(supervisorID, teamID uint) (*GetMessagesResponse, error) {
    // Verify team belongs to this supervisor
    var team struct {
        ID           uint
        SupervisorID *uint
    }
    
    err := s.db.Table("teams").
        Select("id, supervisor_id").
        Where("id = ?", teamID).
        Scan(&team).Error
    
    if err != nil || team.ID == 0 {
        return nil, errors.New("team not found")
    }
    
    if team.SupervisorID == nil || *team.SupervisorID != supervisorID {
        return nil, errors.New("you are not the supervisor of this team")
    }
    
    // Get messages
    messages, err := s.msgRepo.FindByTeamID(teamID)
    if err != nil {
        return nil, err
    }
    
    var items []MessageItem
    
    for _, msg := range messages {
        var senderName string
        
        if msg.SenderType == "supervisor" {
            // Get supervisor name
            s.db.Table("supervisors").
                Select("name").
                Where("id = ?", msg.SenderID).
                Scan(&senderName)
        } else {
            // Get student name
            s.db.Table("students").
                Select("CONCAT(first_name, ' ', last_name)").
                Where("id = ?", msg.SenderID).
                Scan(&senderName)
        }
        
        items = append(items, MessageItem{
            ID:         msg.ID,
            SenderName: senderName,
            SenderType: msg.SenderType,
            Message:    msg.Message,
            CreatedAt:  msg.CreatedAt,
        })
    }
    
    return &GetMessagesResponse{
        Total:    len(items),
        Messages: items,
    }, nil
}

// ============================================
// STUDENT: Send Message to Supervisor
// ============================================

func (s *messageService) SendMessageAsStudent(studentID uint, req SendMessageRequest) (*SendMessageResponse, error) {
    // Get student's team
    var team struct {
        ID           uint
        Student1ID   uint
        Student2ID   uint
        SupervisorID *uint
    }
    
    err := s.db.Table("teams").
        Select("id, student1_id, student2_id, supervisor_id").
        Where("student1_id = ? OR student2_id = ?", studentID, studentID).
        Scan(&team).Error
    
    if err != nil || team.ID == 0 {
        return nil, errors.New("you don't have a team yet")
    }
    
    // Verify student is in this team
    if team.Student1ID != studentID && team.Student2ID != studentID {
        return nil, errors.New("you are not a member of this team")
    }
    
    // Verify team has supervisor
    if team.SupervisorID == nil {
        return nil, errors.New("your team doesn't have a supervisor yet")
    }
    
    // Create message
    message := &Message{
        TeamID:     team.ID,
        SenderID:   studentID,
        SenderType: "student",
        Message:    req.Message,
        CreatedAt:  time.Now(),
    }
    
    err = s.msgRepo.Create(message)
    if err != nil {
        return nil, err
    }
    
    return &SendMessageResponse{
        Message: "Message sent successfully",
    }, nil
}

// ============================================
// STUDENT: Get My Team Messages
// ============================================

func (s *messageService) GetMyTeamMessages(studentID uint) (*GetMessagesResponse, error) {
    // Get student's team
    var team struct {
        ID         uint
        Student1ID uint
        Student2ID uint
    }
    
    err := s.db.Table("teams").
        Select("id, student1_id, student2_id").
        Where("student1_id = ? OR student2_id = ?", studentID, studentID).
        Scan(&team).Error
    
    if err != nil || team.ID == 0 {
        return &GetMessagesResponse{
            Total:    0,
            Messages: []MessageItem{},
        }, nil
    }
    
    // Verify student is in this team
    if team.Student1ID != studentID && team.Student2ID != studentID {
        return nil, errors.New("you are not a member of this team")
    }
    
    // Get messages
    messages, err := s.msgRepo.FindByTeamID(team.ID)
    if err != nil {
        return nil, err
    }
    
    var items []MessageItem
    
    for _, msg := range messages {
        var senderName string
        
        if msg.SenderType == "supervisor" {
            // Get supervisor name
            s.db.Table("supervisors").
                Select("name").
                Where("id = ?", msg.SenderID).
                Scan(&senderName)
        } else {
            // Get student name
            s.db.Table("students").
                Select("CONCAT(first_name, ' ', last_name)").
                Where("id = ?", msg.SenderID).
                Scan(&senderName)
        }
        
        items = append(items, MessageItem{
            ID:         msg.ID,
            SenderName: senderName,
            SenderType: msg.SenderType,
            Message:    msg.Message,
            CreatedAt:  msg.CreatedAt,
        })
    }
    
    return &GetMessagesResponse{
        Total:    len(items),
        Messages: items,
    }, nil
}