package task

import "time"

// ============================================
// SUPERVISOR: Create Task
// ============================================

type CreateTaskRequest struct {
    Title       string     `json:"title" binding:"required,min=5,max=255"`
    Description string     `json:"description" binding:"required,min=10"`
    Deadline    *time.Time `json:"deadline"`
}

type CreateTaskResponse struct {
    Message string `json:"message"`
    TaskID  uint   `json:"task_id"`
}

// ============================================
// SUPERVISOR: Get My Teams (Dashboard)
// ============================================

type MyTeamItem struct {
    ID          uint   `json:"id"`
    Name        string `json:"name"`
    ProjectName string `json:"project_name"`
    Members     []struct {
        Name               string `json:"name"`
        RegistrationNumber string `json:"registration_number"`
    } `json:"members"`
    TasksCount int `json:"tasks_count"`
}

type GetMyTeamsResponse struct {
    Total int          `json:"total"`
    Teams []MyTeamItem `json:"teams"`
}

// ============================================
// Get Team Tasks (Both Supervisor & Student)
// ============================================

type TaskItem struct {
    ID          uint       `json:"id"`
    Title       string     `json:"title"`
    Description string     `json:"description"`
    Deadline    *time.Time `json:"deadline"`
    Status      string     `json:"status"`
    CreatedAt   time.Time  `json:"created_at"`
    
    // Show submissions
    Submissions []struct {
        StudentName    string    `json:"student_name"`
        SubmissionType string    `json:"submission_type"`
        FileURL        string    `json:"file_url,omitempty"`
        LinkURL        string    `json:"link_url,omitempty"`
        TextContent    string    `json:"text_content,omitempty"`
        SubmittedAt    time.Time `json:"submitted_at"`
    } `json:"submissions"`
}

type GetTasksResponse struct {
    Total int        `json:"total"`
    Tasks []TaskItem `json:"tasks"`
}

// ============================================
// STUDENT: Submit Task
// ============================================

type SubmitTaskRequest struct {
    SubmissionType string `json:"submission_type" binding:"required,oneof=file link text"`
    FileURL        string `json:"file_url"`
    LinkURL        string `json:"link_url"`
    TextContent    string `json:"text_content"`
}

type SubmitTaskResponse struct {
    Message string `json:"message"`
}