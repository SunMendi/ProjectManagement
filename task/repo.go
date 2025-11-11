package task

import (
    "gorm.io/gorm"
)

// ============================================
// Task Repository Interface
// ============================================

type TaskRepository interface {
    // Task CRUD
    Create(task *Task) error
    FindByID(id uint) (*Task, error)
    FindByTeamID(teamID uint) ([]Task, error)
    Delete(id uint) error
    
    // Get tasks with submissions
    FindByTeamIDWithSubmissions(teamID uint) ([]Task, error)
    
    // Count tasks for a team
    CountTasksByTeamID(teamID uint) (int64, error)
}

// ============================================
// TaskSubmission Repository Interface
// ============================================

type TaskSubmissionRepository interface {
    // Submission CRUD
    Create(submission *TaskSubmission) error
    FindByID(id uint) (*TaskSubmission, error)           // ✅ ADD
    Update(submission *TaskSubmission) error
    FindByTaskID(taskID uint) ([]TaskSubmission, error)
    FindByStudentAndTask(studentID, taskID uint) (*TaskSubmission, error)
    StudentHasSubmitted(studentID, taskID uint) (bool, error)
}

// ============================================
// Task Repository Implementation
// ============================================

type taskRepository struct {
    db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
    return &taskRepository{db: db}
}

func (r *taskRepository) Create(task *Task) error {
    return r.db.Create(task).Error
}

func (r *taskRepository) FindByID(id uint) (*Task, error) {
    var task Task
    err := r.db.Where("id = ?", id).First(&task).Error
    if err != nil {
        return nil, err
    }
    return &task, nil
}

func (r *taskRepository) FindByTeamID(teamID uint) ([]Task, error) {
    var tasks []Task
    err := r.db.Where("team_id = ?", teamID).
        Order("created_at DESC").
        Find(&tasks).Error
    return tasks, err
}

func (r *taskRepository) Delete(id uint) error {
    return r.db.Delete(&Task{}, id).Error
}

func (r *taskRepository) FindByTeamIDWithSubmissions(teamID uint) ([]Task, error) {
    var tasks []Task
    err := r.db.Where("team_id = ?", teamID).
        Order("created_at DESC").
        Find(&tasks).Error
    return tasks, err
}

func (r *taskRepository) CountTasksByTeamID(teamID uint) (int64, error) {
    var count int64
    err := r.db.Model(&Task{}).
        Where("team_id = ?", teamID).
        Count(&count).Error
    return count, err
}

// ============================================
// TaskSubmission Repository Implementation
// ============================================

type taskSubmissionRepository struct {
    db *gorm.DB
}

func NewTaskSubmissionRepository(db *gorm.DB) TaskSubmissionRepository {
    return &taskSubmissionRepository{db: db}
}

func (r *taskSubmissionRepository) Create(submission *TaskSubmission) error {
    return r.db.Create(submission).Error
}

func (r *taskSubmissionRepository) FindByTaskID(taskID uint) ([]TaskSubmission, error) {
    var submissions []TaskSubmission
    err := r.db.Where("task_id = ?", taskID).
        Order("submitted_at DESC").
        Find(&submissions).Error
    return submissions, err
}

func (r *taskSubmissionRepository) FindByStudentAndTask(studentID, taskID uint) (*TaskSubmission, error) {
    var submission TaskSubmission
    err := r.db.Where("student_id = ? AND task_id = ?", studentID, taskID).
        First(&submission).Error
    if err != nil {
        return nil, err
    }
    return &submission, nil
}

func (r *taskSubmissionRepository) StudentHasSubmitted(studentID, taskID uint) (bool, error) {
    var count int64
    err := r.db.Model(&TaskSubmission{}).
        Where("student_id = ? AND task_id = ?", studentID, taskID).
        Count(&count).Error
    return count > 0, err
}


func (r *taskSubmissionRepository) FindByID(id uint) (*TaskSubmission, error) {
    var submission TaskSubmission
    err := r.db.Where("id = ?", id).First(&submission).Error
    if err != nil {
        return nil, err
    }
    return &submission, nil
}

func (r *taskSubmissionRepository) Update(submission *TaskSubmission) error {
    return r.db.Save(submission).Error
}