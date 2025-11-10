package task

import "time"

// Task - Supervisor assigns tasks to team
type Task struct {
    ID           uint       `json:"id" gorm:"primaryKey"`
    TeamID       uint       `json:"team_id" gorm:"not null;index"`
    SupervisorID uint       `json:"supervisor_id" gorm:"not null"`
    Title        string     `json:"title" gorm:"not null"`
    Description  string     `json:"description" gorm:"type:text"`
    Deadline     *time.Time `json:"deadline"`
    Status       string     `json:"status" gorm:"default:'pending'"` // pending, in_progress, completed
    CreatedAt    time.Time  `json:"created_at"`
    UpdatedAt    time.Time  `json:"updated_at"`
}

// TaskSubmission - Students submit their work
type TaskSubmission struct {
    ID             uint      `json:"id" gorm:"primaryKey"`
    TaskID         uint      `json:"task_id" gorm:"not null;index"`
    StudentID      uint      `json:"student_id" gorm:"not null"` // Who submitted
    SubmissionType string    `json:"submission_type"`            // 'file', 'link', 'text'
    FileURL        string    `json:"file_url"`                   // S3/local path
    LinkURL        string    `json:"link_url"`                   // External URL
    TextContent    string    `json:"text_content" gorm:"type:text"` // Text submission
    SubmittedAt    time.Time `json:"submitted_at"`
}