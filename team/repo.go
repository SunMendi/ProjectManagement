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
	AssignSupervisor(teamID uint, supervisorID uint) error

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

func (r *teamRepository) AssignSupervisor(teamID uint, supervisorID uint) error {
    return r.db.Model(&Team{}).
        Where("id = ?", teamID).
        Updates(map[string]interface{}{
            "supervisor_id": supervisorID,
            "status":        "active",
        }).Error
}

// ============================================
// Repository Interface
// ============================================

type SupervisorRequestRepository interface {
    // Student operations
    Create(request *SupervisorRequest) error
    FindByTeamID(teamID uint) ([]SupervisorRequest, error)
    
    // Supervisor operations
    FindPendingBySupervisorID(supervisorID uint) ([]SupervisorRequest, error)
    FindByID(id uint) (*SupervisorRequest, error)
    UpdateStatus(id uint, status string, rejectReason string) error
    
    // Validation
    HasPendingRequest(teamID uint) (bool, error)
}

// ============================================
// Repository Implementation
// ============================================

type supervisorRequestRepository struct {
    db *gorm.DB
}

func NewSupervisorRequestRepository(db *gorm.DB) SupervisorRequestRepository {
    return &supervisorRequestRepository{db: db}
}

// ============================================
// Student Operations
// ============================================

// Create - student sends request to supervisor
func (r *supervisorRequestRepository) Create(request *SupervisorRequest) error {
    return r.db.Create(request).Error
}

// FindByTeamID - get all supervisor requests for a team
func (r *supervisorRequestRepository) FindByTeamID(teamID uint) ([]SupervisorRequest, error) {
    var requests []SupervisorRequest
    err := r.db.Where("team_id = ?", teamID).
        Order("created_at DESC").
        Find(&requests).Error
    return requests, err
}

// ============================================
// Supervisor Operations
// ============================================

// FindPendingBySupervisorID - supervisor sees pending requests
func (r *supervisorRequestRepository) FindPendingBySupervisorID(supervisorID uint) ([]SupervisorRequest, error) {
    var requests []SupervisorRequest
    err := r.db.Where("supervisor_id = ? AND status = ?", supervisorID, "pending").
        Order("created_at DESC").
        Find(&requests).Error
    return requests, err
}

// FindByID - get a specific request
func (r *supervisorRequestRepository) FindByID(id uint) (*SupervisorRequest, error) {
    var request SupervisorRequest
    err := r.db.Where("id = ?", id).First(&request).Error
    if err != nil {
        return nil, err
    }
    return &request, nil
}

// UpdateStatus - accept or reject request
func (r *supervisorRequestRepository) UpdateStatus(id uint, status string, rejectReason string) error {
    updates := map[string]interface{}{
        "status": status,
    }
    
    if rejectReason != "" {
        updates["reject_reason"] = rejectReason
    }
    
    return r.db.Model(&SupervisorRequest{}).
        Where("id = ?", id).
        Updates(updates).Error
}

// ============================================
// Validation
// ============================================

// HasPendingRequest - check if team already has a pending request
func (r *supervisorRequestRepository) HasPendingRequest(teamID uint) (bool, error) {
    var count int64
    err := r.db.Model(&SupervisorRequest{}).
        Where("team_id = ? AND status = ?", teamID, "pending").
        Count(&count).Error
    return count > 0, err
}