package team

import "time"

// Team - stores the actual team with 2 students
type Team struct {
    ID          uint      `json:"id" gorm:"primaryKey"`
    Name        string    `json:"name" gorm:"not null"`
    ProjectName string    `json:"project_name" gorm:"not null"`
    Department  string    `json:"department" gorm:"not null"`
    Session     string    `json:"session" gorm:"not null"`
    Status      string    `json:"status" gorm:"not null;default:'pending_supervisor'"` // pending_supervisor, active, completed
    Student1ID  uint      `json:"student1_id" gorm:"not null;uniqueIndex"`              // One student = one team
    Student2ID  uint      `json:"student2_id" gorm:"not null;uniqueIndex"`  
	SupervisorID *uint     `json:"supervisor_id" gorm:"index"` 
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// TeamRequest - student invites another student to form team
type TeamRequest struct {
    ID            uint      `json:"id" gorm:"primaryKey"`
    FromStudentID uint      `json:"from_student_id" gorm:"not null"`
    ToStudentID   uint      `json:"to_student_id" gorm:"not null"`
    TeamName      string    `json:"team_name" gorm:"not null"`
    ProjectName   string    `json:"project_name" gorm:"not null"`
    Status        string    `json:"status" gorm:"not null;default:'pending'"` // pending, accepted, rejected
    CreatedAt     time.Time `json:"created_at"`
    UpdatedAt     time.Time `json:"updated_at"`
}



type SupervisorRequest struct {
    ID           uint      `json:"id" gorm:"primaryKey"`
    TeamID       uint      `json:"team_id" gorm:"not null;index"` // ✅ Contains both students already!
    SupervisorID uint      `json:"supervisor_id" gorm:"not null;index"`
    ProjectTitle string    `json:"project_title" gorm:"not null"`
    ProjectInfo  string    `json:"project_info" gorm:"type:text"`
    Status       string    `json:"status" gorm:"not null;default:'pending'"` // pending, accepted, rejected
    RejectReason string    `json:"reject_reason" gorm:"type:text"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
}