 package team

import (
   // "ProjectManagement/User"
    "time"

    "gorm.io/gorm"
)

// ============================================
// Team Repository Interface
// ============================================

type TeamRepository interface {
    Create(team *Team) error
    FindByID(id uint) (*Team, error)
    FindByStudentID(studentID uint) (*Team, error)
    StudentHasTeam(studentID uint) (bool, error)
    Update(team *Team) error
}

// ============================================
// TeamRequest Repository Interface
// ============================================

type TeamRequestRepository interface {
    Create(request *TeamRequest) error
    FindByID(id uint) (*TeamRequest, error)
    FindReceivedRequests(studentID uint) ([]TeamRequest, error)
    FindSentRequests(studentID uint) ([]TeamRequest, error)
    UpdateStatus(id uint, status string) error
    Delete(id uint) error
    PendingRequestExists(fromStudentID, toStudentID uint) (bool, error)
}

// ============================================
// Team Repository Implementation
// ============================================

type teamRepository struct {
    db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
    return &teamRepository{db: db}
}

func (r *teamRepository) Create(team *Team) error {
    return r.db.Create(team).Error
}

func (r *teamRepository) FindByID(id uint) (*Team, error) {
    var team Team
    err := r.db.Where("id = ?", id).First(&team).Error
    if err != nil {
        return nil, err
    }
    return &team, nil
}

func (r *teamRepository) FindByStudentID(studentID uint) (*Team, error) {
    var team Team
    err := r.db.Where("student1_id = ? OR student2_id = ?", studentID, studentID).
        First(&team).Error
    if err != nil {
        return nil, err
    }
    return &team, nil
}

func (r *teamRepository) StudentHasTeam(studentID uint) (bool, error) {
    var count int64
    err := r.db.Model(&Team{}).
        Where("student1_id = ? OR student2_id = ?", studentID, studentID).
        Count(&count).Error
    return count > 0, err
}

func (r *teamRepository) Update(team *Team) error {
    team.UpdatedAt = time.Now()
    return r.db.Save(team).Error
}

// ============================================
// TeamRequest Repository Implementation
// ============================================

type teamRequestRepository struct {
    db *gorm.DB
}

func NewTeamRequestRepository(db *gorm.DB) TeamRequestRepository {
    return &teamRequestRepository{db: db}
}

func (r *teamRequestRepository) Create(request *TeamRequest) error {
    return r.db.Create(request).Error
}

func (r *teamRequestRepository) FindByID(id uint) (*TeamRequest, error) {
    var request TeamRequest
    err := r.db.Where("id = ?", id).First(&request).Error
    if err != nil {
        return nil, err
    }
    return &request, nil
}

func (r *teamRequestRepository) FindReceivedRequests(studentID uint) ([]TeamRequest, error) {
    var requests []TeamRequest
    err := r.db.Where("to_student_id = ? AND status = 'pending'", studentID).
        Order("created_at DESC").
        Find(&requests).Error
    if err != nil {
        return nil, err
    }
    return requests, nil
}

func (r *teamRequestRepository) FindSentRequests(studentID uint) ([]TeamRequest, error) {
    var requests []TeamRequest
    err := r.db.Where("from_student_id = ?", studentID).
        Order("created_at DESC").
        Find(&requests).Error
    if err != nil {
        return nil, err
    }
    return requests, nil
}

func (r *teamRequestRepository) UpdateStatus(id uint, status string) error {
    return r.db.Model(&TeamRequest{}).
        Where("id = ?", id).
        Updates(map[string]interface{}{
            "status":     status,
            "updated_at": time.Now(),
        }).Error
}

func (r *teamRequestRepository) Delete(id uint) error {
    return r.db.Delete(&TeamRequest{}, id).Error
}

func (r *teamRequestRepository) PendingRequestExists(fromStudentID, toStudentID uint) (bool, error) {
    var count int64
    err := r.db.Model(&TeamRequest{}).
        Where("from_student_id = ? AND to_student_id = ? AND status = 'pending'", fromStudentID, toStudentID).
        Count(&count).Error
    return count > 0, err
}