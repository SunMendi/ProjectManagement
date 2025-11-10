package user

type RegisterStudentRequest struct {
    Email              string `json:"email" binding:"required,email"`
    Password           string `json:"password" binding:"required,min=8"`
    FirstName          string `json:"first_name" binding:"required"`
    LastName           string `json:"last_name" binding:"required"`
    Department         string `json:"department" binding:"required"`
    Session            string `json:"session" binding:"required"`
    RegistrationNumber string `json:"registration_number" binding:"required"`
    Batch              string `json:"batch" binding:"required"`
}

type RegisterStudentResponse struct {
    Message   string `json:"message"`
    UserID    uint   `json:"user_id"`
    StudentID uint   `json:"student_id"`
}

type GetStudentResponse struct {
    ID                 uint   `json:"id"`
    UserID             uint   `json:"user_id"`
    Email              string `json:"email"`
    FirstName          string `json:"first_name"`
    LastName           string `json:"last_name"`
    Department         string `json:"department"`
    Session            string `json:"session"`
    RegistrationNumber string `json:"registration_number"`
    Batch              string `json:"batch"`
    Status             string `json:"status"`
    Role               string `json:"role"`
}

type UpdateStudentProfileRequest struct {
    FirstName  string `json:"first_name"`
    LastName   string `json:"last_name"`
    Image      string `json:"image"`
    Department string `json:"department"`
}

type GetStudentsByFilterRequest struct {
    Department string `form:"department" binding:"required"` // ✅ Changed from json to form
    Session    string `form:"session" binding:"required"`    // ✅ Changed from json to form
}

type StudentListItem struct {
	 ID uint `json:"id"`
	 UserID uint `json:"user_id"`
	 Email              string `json:"email"`
     FirstName          string `json:"first_name"`
     LastName           string `json:"last_name"`
     Department         string `json:"department"`
     Session            string `json:"session"`
     RegistrationNumber string `json:"registration_number"`
     Batch              string `json:"batch"`
     Image              string `json:"image"`
     Status             string `json:"status"`
     HasTeam            bool   `json:"has_team"`
}

type GetStudentsByFilterResponse struct {
	 Total int  `json:"total"`
	 Students []StudentListItem `json:"students"`
}

type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
    Message   string             `json:"message"`
    Token     string             `json:"token"`      // ✅ Add token
    UserID    uint               `json:"user_id"`
    StudentID uint               `json:"student_id"`
    Role      string             `json:"role"`
    Status    string             `json:"status"`
    Student   GetStudentResponse `json:"student"`
}

// ... existing student DTOs ...

// Supervisor Registration
type RegisterSupervisorRequest struct {
    Email        string `json:"email" binding:"required,email"`
    Password     string `json:"password" binding:"required,min=8"`
    Name         string `json:"name" binding:"required"`
    Designation  string `json:"designation" binding:"required"`
    Department   string `json:"department" binding:"required"`
    ResearchArea string `json:"research_area"`
    Phone        string `json:"phone"`
}

type RegisterSupervisorResponse struct {
    Message      string `json:"message"`
    UserID       uint   `json:"user_id"`
    SupervisorID uint   `json:"supervisor_id"`
}

// Get Supervisor
type GetSupervisorResponse struct {
    ID           uint   `json:"id"`
    UserID       uint   `json:"user_id"`
    Email        string `json:"email"`
    Name         string `json:"name"`
    Designation  string `json:"designation"`
    Department   string `json:"department"`
    ResearchArea string `json:"research_area"`
    Phone        string `json:"phone"`
    Image        string `json:"image"`
    Status       string `json:"status"`
    Role         string `json:"role"`
}

// Update Supervisor Profile
type UpdateSupervisorProfileRequest struct {
    Name         string `json:"name"`
    Designation  string `json:"designation"`
    Department   string `json:"department"`
    ResearchArea string `json:"research_area"`
    Phone        string `json:"phone"`
    Image        string `json:"image"`
}

// Get Supervisors by Department
type GetSupervisorsByDepartmentRequest struct {
    Department string `form:"department" binding:"required"`
}

type SupervisorListItem struct {
    ID           uint   `json:"id"`
    UserID       uint   `json:"user_id"`
    Email        string `json:"email"`
    Name         string `json:"name"`
    Designation  string `json:"designation"`
    Department   string `json:"department"`
    ResearchArea string `json:"research_area"`
    Phone        string `json:"phone"`
    Image        string `json:"image"`
    Status       string `json:"status"`
   // HasTeam            bool   `json:"has_team"`
}

type GetSupervisorsByDepartmentResponse struct {
    Total       int                  `json:"total"`
    Supervisors []SupervisorListItem `json:"supervisors"`
}

// Supervisor Login
type SupervisorLoginResponse struct {
    Message      string                `json:"message"`
    Token        string                `json:"token"`
    UserID       uint                  `json:"user_id"`
    SupervisorID uint                  `json:"supervisor_id"`
    Role         string                `json:"role"`
    Status       string                `json:"status"`
    Supervisor   GetSupervisorResponse `json:"supervisor"`
}