package task

import (
    "errors"
    "time"

    "gorm.io/gorm"
)

// ============================================
// Task Service Interface
// ============================================

type TaskService interface {
    // Supervisor operations
    GetMyTeams(supervisorID uint, session string) (*GetMyTeamsResponse, error)
    CreateTask(supervisorID, teamID uint, req CreateTaskRequest) (*CreateTaskResponse, error)
    GetTeamTasks(supervisorID, teamID uint, session string) (*GetTasksResponse, error)
    DeleteTask(supervisorID, taskID uint) error
    ReviewSubmission(supervisorID, submissionID uint, req ReviewSubmissionRequest) (*ReviewSubmissionResponse, error) // ✅ ADD

    
    // Student operations
    GetMyTeamTasks(studentID uint) (*GetTasksResponse, error)
    SubmitTask(studentID, taskID uint, req SubmitTaskRequest) (*SubmitTaskResponse, error)
}

// ============================================
// Service Implementation
// ============================================

type taskService struct {
    db                *gorm.DB
    taskRepo          TaskRepository
    submissionRepo    TaskSubmissionRepository
}

func NewTaskService(
    db *gorm.DB,
    taskRepo TaskRepository,
    submissionRepo TaskSubmissionRepository,
) TaskService {
    return &taskService{
        db:             db,
        taskRepo:       taskRepo,
        submissionRepo: submissionRepo,
    }
}

// ============================================
// SUPERVISOR: Get My Teams
// ============================================

// ============================================
// SUPERVISOR: Get My Teams
// ============================================

func (s *taskService) GetMyTeams(supervisorID uint, session string) (*GetMyTeamsResponse, error) {
    // Get all teams assigned to this supervisor
   var teams []struct {
        ID          uint
        Name        string
        ProjectName string
        Student1ID  uint
        Student2ID  uint
        Session     string  // ✅ NEW
    }
    
    // ✅ Build query
    query := s.db.Table("teams").
        Select("teams.id, teams.name, teams.project_name, teams.student1_id, teams.student2_id, teams.session").
        Where("teams.supervisor_id = ? AND teams.status = ?", supervisorID, "active")
    
    // ✅ Add session filter if provided
    if session != "" {
        query = query.Where("teams.session = ?", session)
    }
    
    err := query.Scan(&teams).Error
    
    if err != nil {
        return nil, err
    }
    
    var items []MyTeamItem
    
    for _, team := range teams {
        // Get student1
        var student1 struct {
            Name               string
            RegistrationNumber string
        }
        s.db.Table("students").
            Select("CONCAT(students.first_name, ' ', students.last_name) as name, students.registration_number").
            Where("students.id = ?", team.Student1ID).
            Scan(&student1)
        
        // Get student2
        var student2 struct {
            Name               string
            RegistrationNumber string
        }
        s.db.Table("students").
            Select("CONCAT(students.first_name, ' ', students.last_name) as name, students.registration_number").
            Where("students.id = ?", team.Student2ID).
            Scan(&student2)
        
        // Count tasks
        tasksCount, _ := s.taskRepo.CountTasksByTeamID(team.ID)
        
        // Build item with inline member creation
        item := MyTeamItem{
            ID:          team.ID,
            Name:        team.Name,
            ProjectName: team.ProjectName,
            Members: []struct {
                Name               string `json:"name"`
                RegistrationNumber string `json:"registration_number"`
            }{
                {Name: student1.Name, RegistrationNumber: student1.RegistrationNumber},
                {Name: student2.Name, RegistrationNumber: student2.RegistrationNumber},
            },
            TasksCount: int(tasksCount),
        }
        
        items = append(items, item)
    }
    
    return &GetMyTeamsResponse{
        Total: len(items),
        Teams: items,
    }, nil
}

// ============================================
// SUPERVISOR: Create Task
// ============================================

func (s *taskService) CreateTask(supervisorID, teamID uint, req CreateTaskRequest) (*CreateTaskResponse, error) {
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
    
    // Create task
    task := &Task{
        TeamID:       teamID,
        SupervisorID: supervisorID,
        Title:        req.Title,
        Description:  req.Description,
        Deadline:     req.Deadline,
        Status:       "pending",
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }
    
    err = s.taskRepo.Create(task)
    if err != nil {
        return nil, err
    }
    
    return &CreateTaskResponse{
        Message: "Task created successfully",
        TaskID:  task.ID,
    }, nil
}

// ============================================
// SUPERVISOR: Get Team Tasks
// ============================================

// ============================================
// SUPERVISOR: Get Team Tasks
// ============================================

func (s *taskService) GetTeamTasks(supervisorID, teamID uint, session string) (*GetTasksResponse, error) {
    // Verify team belongs to this supervisor
    var team struct {
        ID           uint
        SupervisorID *uint
        Session      string // ✅ Get session from team
    }
    
    err := s.db.Table("teams").
        Select("id, supervisor_id, session").
        Where("id = ?", teamID).
        Scan(&team).Error
    
    if err != nil || team.ID == 0 {
        return nil, errors.New("team not found")
    }
    
    if team.SupervisorID == nil || *team.SupervisorID != supervisorID {
        return nil, errors.New("you are not the supervisor of this team")
    }

    // ✅ If session filter provided, verify it matches team's session
    if session != "" && team.Session != session {
        return &GetTasksResponse{
            Total: 0,
            Tasks: []TaskItem{},
        }, nil
    }
    
    // Get tasks
    tasks, err := s.taskRepo.FindByTeamID(teamID)
    if err != nil {
        return nil, err
    }
    
    var items []TaskItem
    
    for _, task := range tasks {
        // Get submissions for this task
        submissions, _ := s.submissionRepo.FindByTaskID(task.ID)
        
        // Define submission details struct
        var submissionDetails []struct {
            ID             uint      `json:"id"`
            StudentName    string    `json:"student_name"`
            SubmissionType string    `json:"submission_type"`
            FileURL        string    `json:"file_url,omitempty"`
            LinkURL        string    `json:"link_url,omitempty"`
            TextContent    string    `json:"text_content,omitempty"`
            Status         string    `json:"status"`
            Feedback       string    `json:"feedback,omitempty"`
            SubmittedAt    time.Time `json:"submitted_at"`
        }
        
        // Populate submission details
        for _, sub := range submissions {
            var studentName string
            s.db.Table("students").
                Select("CONCAT(first_name, ' ', last_name)").
                Where("id = ?", sub.StudentID).
                Scan(&studentName)

            submissionDetails = append(submissionDetails, struct {
                ID             uint      `json:"id"`
                StudentName    string    `json:"student_name"`
                SubmissionType string    `json:"submission_type"`
                FileURL        string    `json:"file_url,omitempty"`
                LinkURL        string    `json:"link_url,omitempty"`
                TextContent    string    `json:"text_content,omitempty"`
                Status         string    `json:"status"`
                Feedback       string    `json:"feedback,omitempty"`
                SubmittedAt    time.Time `json:"submitted_at"`
            }{
                ID:             sub.ID,
                StudentName:    studentName,
                SubmissionType: sub.SubmissionType,
                FileURL:        sub.FileURL,
                LinkURL:        sub.LinkURL,
                TextContent:    sub.TextContent,
                Status:         sub.Status,
                Feedback:       sub.Feedback,
                SubmittedAt:    sub.SubmittedAt,
            })
        }
        
        item := TaskItem{
            ID:          task.ID,
            Title:       task.Title,
            Description: task.Description,
            Deadline:    task.Deadline,
            Status:      task.Status,
            CreatedAt:   task.CreatedAt,
            Submissions: submissionDetails,
        }
        
        items = append(items, item)
    }
    
    return &GetTasksResponse{
        Total: len(items),
        Tasks: items,
    }, nil
}

// ============================================
// SUPERVISOR: Delete Task
// ============================================

func (s *taskService) DeleteTask(supervisorID, taskID uint) error {
    // Get task
    task, err := s.taskRepo.FindByID(taskID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return errors.New("task not found")
        }
        return err
    }
    
    // Verify task belongs to this supervisor
    if task.SupervisorID != supervisorID {
        return errors.New("you are not authorized to delete this task")
    }
    
    return s.taskRepo.Delete(taskID)
}

// ============================================
// STUDENT: Get My Team Tasks
// ============================================

// ============================================
// STUDENT: Get My Team Tasks
// ============================================

func (s *taskService) GetMyTeamTasks(studentID uint) (*GetTasksResponse, error) {
    // Get student's team
    var team struct {
        ID uint
    }
    
    err := s.db.Table("teams").
        Select("id").
        Where("student1_id = ? OR student2_id = ?", studentID, studentID).
        Scan(&team).Error
    
    if err != nil || team.ID == 0 {
        return &GetTasksResponse{
            Total: 0,
            Tasks: []TaskItem{},
        }, nil
    }
    
    // Get tasks for this team
    tasks, err := s.taskRepo.FindByTeamID(team.ID)
    if err != nil {
        return nil, err
    }
    
    var items []TaskItem
    
    for _, task := range tasks {
        // Get submissions for this task
        submissions, _ := s.submissionRepo.FindByTaskID(task.ID)
        
        // ✅ FIX: Add missing fields to match TaskItem DTO
        var submissionDetails []struct {
            ID             uint      `json:"id"`              // ✅ ADD
            StudentName    string    `json:"student_name"`
            SubmissionType string    `json:"submission_type"`
            FileURL        string    `json:"file_url,omitempty"`
            LinkURL        string    `json:"link_url,omitempty"`
            TextContent    string    `json:"text_content,omitempty"`
            Status         string    `json:"status"`          // ✅ ADD
            Feedback       string    `json:"feedback,omitempty"` // ✅ ADD
            SubmittedAt    time.Time `json:"submitted_at"`
        }
        
        for _, sub := range submissions {
            // Get student name
            var studentName string
            s.db.Table("students").
                Select("CONCAT(first_name, ' ', last_name)").
                Where("id = ?", sub.StudentID).
                Scan(&studentName)
            
            // ✅ FIX: Include all fields
            submissionDetails = append(submissionDetails, struct {
                ID             uint      `json:"id"`
                StudentName    string    `json:"student_name"`
                SubmissionType string    `json:"submission_type"`
                FileURL        string    `json:"file_url,omitempty"`
                LinkURL        string    `json:"link_url,omitempty"`
                TextContent    string    `json:"text_content,omitempty"`
                Status         string    `json:"status"`
                Feedback       string    `json:"feedback,omitempty"`
                SubmittedAt    time.Time `json:"submitted_at"`
            }{
                ID:             sub.ID,              // ✅ ADD
                StudentName:    studentName,
                SubmissionType: sub.SubmissionType,
                FileURL:        sub.FileURL,
                LinkURL:        sub.LinkURL,
                TextContent:    sub.TextContent,
                Status:         sub.Status,          // ✅ ADD
                Feedback:       sub.Feedback,        // ✅ ADD
                SubmittedAt:    sub.SubmittedAt,
            })
        }
        
        item := TaskItem{
            ID:          task.ID,
            Title:       task.Title,
            Description: task.Description,
            Deadline:    task.Deadline,
            Status:      task.Status,
            CreatedAt:   task.CreatedAt,
            Submissions: submissionDetails,
        }
        
        items = append(items, item)
    }
    
    return &GetTasksResponse{
        Total: len(items),
        Tasks: items,
    }, nil
}
// ============================================
// STUDENT: Submit Task
// ============================================

func (s *taskService) SubmitTask(studentID, taskID uint, req SubmitTaskRequest) (*SubmitTaskResponse, error) {
    // Get task
    task, err := s.taskRepo.FindByID(taskID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errors.New("task not found")
        }
        return nil, err
    }
    
    // Verify student is in this team
    var team struct {
        Student1ID uint
        Student2ID uint
    }
    
    err = s.db.Table("teams").
        Select("student1_id, student2_id").
        Where("id = ?", task.TeamID).
        Scan(&team).Error
    
    if err != nil {
        return nil, errors.New("team not found")
    }
    
    if team.Student1ID != studentID && team.Student2ID != studentID {
        return nil, errors.New("you are not a member of this team")
    }
    
    // Check if student already submitted
    hasSubmitted, _ := s.submissionRepo.StudentHasSubmitted(studentID, taskID)
    if hasSubmitted {
        return nil, errors.New("you have already submitted this task")
    }
    
    // Validate submission based on type
    if req.SubmissionType == "file" && req.FileURL == "" {
        return nil, errors.New("file_url is required for file submission")
    }
    if req.SubmissionType == "link" && req.LinkURL == "" {
        return nil, errors.New("link_url is required for link submission")
    }
    if req.SubmissionType == "text" && req.TextContent == "" {
        return nil, errors.New("text_content is required for text submission")
    }
    
    // Create submission
    submission := &TaskSubmission{
        TaskID:         taskID,
        StudentID:      studentID,
        SubmissionType: req.SubmissionType,
        FileURL:        req.FileURL,
        LinkURL:        req.LinkURL,
        TextContent:    req.TextContent,
        SubmittedAt:    time.Now(),
    }
    
    err = s.submissionRepo.Create(submission)
    if err != nil {
        return nil, err
    }
    
    return &SubmitTaskResponse{
        Message: "Task submitted successfully",
    }, nil
}


func (s *taskService) ReviewSubmission(supervisorID, submissionID uint, req ReviewSubmissionRequest) (*ReviewSubmissionResponse, error) {
    // Get submission
    submission, err := s.submissionRepo.FindByID(submissionID)
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, errors.New("submission not found")
        }
        return nil, err
    }
    
    // Get task to verify supervisor
    task, err := s.taskRepo.FindByID(submission.TaskID)
    if err != nil {
        return nil, errors.New("task not found")
    }
    
    // Verify supervisor owns this task
    if task.SupervisorID != supervisorID {
        return nil, errors.New("unauthorized")
    }
    
    // Update submission
    if req.Action == "approve" {
        submission.Status = "approved"
    } else {
        submission.Status = "rejected"
    }
    submission.Feedback = req.Feedback
    
    err = s.submissionRepo.Update(submission)
    if err != nil {
        return nil, err
    }
    
    message := "Submission approved successfully"
    if req.Action == "reject" {
        message = "Submission rejected successfully"
    }
    
    return &ReviewSubmissionResponse{
        Message: message,
    }, nil
}