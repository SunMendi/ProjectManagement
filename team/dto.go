package team

import "time"



type SendTeamRequestRequest struct {
    ToStudentID uint   `json:"to_student_id" binding:"required"`
    TeamName    string `json:"team_name" binding:"required,min=3,max=100"`
    ProjectName string `json:"project_name" binding:"required,min=5,max=200"`
}

type SendTeamRequestResponse struct {
    Message   string `json:"message"`
    RequestID uint   `json:"request_id"`
}

// ============================================
// Get Received Requests (Requests sent TO me)
// ============================================

type ReceivedRequestItem struct {
    ID          uint      `json:"id"`
    FromStudent struct {
        ID                 uint   `json:"id"`
        Name               string `json:"name"`
        RegistrationNumber string `json:"registration_number"`
        Email              string `json:"email"`
        Department         string `json:"department"`
        Session            string `json:"session"`
    } `json:"from_student"`
    TeamName    string    `json:"team_name"`
    ProjectName string    `json:"project_name"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
}

type GetReceivedRequestsResponse struct {
    Total    int                   `json:"total"`
    Requests []ReceivedRequestItem `json:"requests"`
}

// ============================================
// Get Sent Requests (Requests I sent)
// ============================================

type SentRequestItem struct {
    ID        uint      `json:"id"`
    ToStudent struct {
        ID                 uint   `json:"id"`
        Name               string `json:"name"`
        RegistrationNumber string `json:"registration_number"`
        Email              string `json:"email"`
        Department         string `json:"department"`
        Session            string `json:"session"`
    } `json:"to_student"`
    TeamName    string    `json:"team_name"`
    ProjectName string    `json:"project_name"`
    Status      string    `json:"status"`
    CreatedAt   time.Time `json:"created_at"`
}

type GetSentRequestsResponse struct {
    Total    int               `json:"total"`
    Requests []SentRequestItem `json:"requests"`
}

// ============================================
// Accept Team Request
// ============================================

type AcceptRequestResponse struct {
    Message string       `json:"message"`
    Team    TeamOverview `json:"team"`
}

type TeamOverview struct {
    ID          uint   `json:"id"`
    Name        string `json:"name"`
    ProjectName string `json:"project_name"`
    Status      string `json:"status"`
}

// ============================================
// Reject Team Request
// ============================================

type RejectRequestResponse struct {
    Message string `json:"message"`
}

// ============================================
// Get My Team
// ============================================

type GetMyTeamResponse struct {
    HasTeam bool         `json:"has_team"`
    Team    *TeamDetails `json:"team"` // null if no team
}

type TeamDetails struct {
    ID          uint         `json:"id"`
    Name        string       `json:"name"`
    ProjectName string       `json:"project_name"`
    Department  string       `json:"department"`
    Session     string       `json:"session"`
    Status      string       `json:"status"`
    Members     []TeamMember `json:"members"`
    CreatedAt   time.Time    `json:"created_at"`
}

type TeamMember struct {
    ID                 uint   `json:"id"`
    Name               string `json:"name"`
    Email              string `json:"email"`
    RegistrationNumber string `json:"registration_number"`
    Department         string `json:"department"`
    Session            string `json:"session"`
    Image              string `json:"image"`
}

type CancelRequestResponse struct {
    Message string `json:"message"`
}