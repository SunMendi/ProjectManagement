package team

import (
    "errors"
    "time"

    "gorm.io/gorm"
)

// ============================================
// Team Service Interface
// ============================================

type TeamService interface {
    SendTeamRequest(fromStudentID uint, req SendTeamRequestRequest) (*SendTeamRequestResponse, error)
    GetReceivedRequests(studentID uint) (*GetReceivedRequestsResponse, error)
    GetSentRequests(studentID uint) (*GetSentRequestsResponse, error)
    AcceptRequest(requestID, acceptingStudentID uint) (*AcceptRequestResponse, error)
    RejectRequest(requestID, rejectingStudentID uint) (*RejectRequestResponse, error)
    CancelRequest(requestID, studentID uint) (*CancelRequestResponse, error)
    GetMyTeam(studentID uint) (*GetMyTeamResponse, error)
}

// ============================================
// Team Service Implementation
// ============================================

type teamService struct {
    db          *gorm.DB
    teamRepo    TeamRepository
    requestRepo TeamRequestRepository
}

func NewTeamService(db *gorm.DB, teamRepo TeamRepository, requestRepo TeamRequestRepository) TeamService {
    return &teamService{
        db:          db,
        teamRepo:    teamRepo,
        requestRepo: requestRepo,
    }
}

// ============================================
// Send Team Request
// ============================================

func (s *teamService) SendTeamRequest(fromStudentID uint, req SendTeamRequestRequest) (*SendTeamRequestResponse, error) {
    // ✅ 1. Validate: Cannot send request to yourself
    if fromStudentID == req.ToStudentID {
        return nil, errors.New("cannot send team request to yourself")
    }

    // ✅ 2. Check: Sender doesn't have team
    hasTeam, err := s.teamRepo.StudentHasTeam(fromStudentID)
    if err != nil {
        return nil, err
    }
    if hasTeam {
        return nil, errors.New("you already have a team")
    }

    // ✅ 3. Check: Receiver doesn't have team
    hasTeam, err = s.teamRepo.StudentHasTeam(req.ToStudentID)
    if err != nil {
        return nil, err
    }
    if hasTeam {
        return nil, errors.New("that student already has a team")
    }

    // ✅ 4. Check: Both students exist and are in same department/session
    var sender, receiver struct {
        ID         uint
        Department string
        Session    string
    }

    err = s.db.Table("students").Where("id = ?", fromStudentID).
        Select("id, department, session").Scan(&sender).Error
    if err != nil {
        return nil, errors.New("sender student not found")
    }

    err = s.db.Table("students").Where("id = ?", req.ToStudentID).
        Select("id, department, session").Scan(&receiver).Error
    if err != nil {
        return nil, errors.New("receiver student not found")
    }

    if sender.Department != receiver.Department {
        return nil, errors.New("students must be in same department")
    }

    if sender.Session != receiver.Session {
        return nil, errors.New("students must be in same session")
    }

    // ✅ 5. Check: No pending request already exists
    exists, err := s.requestRepo.PendingRequestExists(fromStudentID, req.ToStudentID)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, errors.New("you already sent a pending request to this student")
    }

    // ✅ 6. Check: No reverse pending request exists
    exists, err = s.requestRepo.PendingRequestExists(req.ToStudentID, fromStudentID)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, errors.New("this student already sent you a request. Please check your received requests")
    }

    // ✅ 7. Create request
    request := &TeamRequest{
        FromStudentID: fromStudentID,
        ToStudentID:   req.ToStudentID,
        TeamName:      req.TeamName,
        ProjectName:   req.ProjectName,
        Status:        "pending",
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }

    err = s.requestRepo.Create(request)
    if err != nil {
        return nil, err
    }

    return &SendTeamRequestResponse{
        Message:   "Team request sent successfully",
        RequestID: request.ID,
    }, nil
}

// ============================================
// Get Received Requests
// ============================================

func (s *teamService) GetReceivedRequests(studentID uint) (*GetReceivedRequestsResponse, error) {
    requests, err := s.requestRepo.FindReceivedRequests(studentID)
    if err != nil {
        return nil, err
    }

    var items []ReceivedRequestItem
    for _, req := range requests {
        // Get sender details
        var student struct {
            ID                 uint
            FirstName          string
            LastName           string
            RegistrationNumber string
            Email              string
            Department         string
            Session            string
        }

        err := s.db.Table("students").
            Select("students.id, students.first_name, students.last_name, students.registration_number, students.department, students.session, users.email").
            Joins("JOIN users ON users.id = students.user_id").
            Where("students.id = ?", req.FromStudentID).
            Scan(&student).Error

        if err != nil {
            continue
        }

        item := ReceivedRequestItem{
            ID:          req.ID,
            TeamName:    req.TeamName,
            ProjectName: req.ProjectName,
            Status:      req.Status,
            CreatedAt:   req.CreatedAt,
        }
        item.FromStudent.ID = student.ID
        item.FromStudent.Name = student.FirstName + " " + student.LastName
        item.FromStudent.RegistrationNumber = student.RegistrationNumber
        item.FromStudent.Email = student.Email
        item.FromStudent.Department = student.Department
        item.FromStudent.Session = student.Session

        items = append(items, item)
    }

    return &GetReceivedRequestsResponse{
        Total:    len(items),
        Requests: items,
    }, nil
}

// ============================================
// Get Sent Requests
// ============================================

func (s *teamService) GetSentRequests(studentID uint) (*GetSentRequestsResponse, error) {
    requests, err := s.requestRepo.FindSentRequests(studentID)
    if err != nil {
        return nil, err
    }

    var items []SentRequestItem
    for _, req := range requests {
        // Get receiver details
        var student struct {
            ID                 uint
            FirstName          string
            LastName           string
            RegistrationNumber string
            Email              string
            Department         string
            Session            string
        }

        err := s.db.Table("students").
            Select("students.id, students.first_name, students.last_name, students.registration_number, students.department, students.session, users.email").
            Joins("JOIN users ON users.id = students.user_id").
            Where("students.id = ?", req.ToStudentID).
            Scan(&student).Error

        if err != nil {
            continue
        }

        item := SentRequestItem{
            ID:          req.ID,
            TeamName:    req.TeamName,
            ProjectName: req.ProjectName,
            Status:      req.Status,
            CreatedAt:   req.CreatedAt,
        }
        item.ToStudent.ID = student.ID
        item.ToStudent.Name = student.FirstName + " " + student.LastName
        item.ToStudent.RegistrationNumber = student.RegistrationNumber
        item.ToStudent.Email = student.Email
        item.ToStudent.Department = student.Department
        item.ToStudent.Session = student.Session

        items = append(items, item)
    }

    return &GetSentRequestsResponse{
        Total:    len(items),
        Requests: items,
    }, nil
}

// ============================================
// Accept Team Request
// ============================================

func (s *teamService) AcceptRequest(requestID, acceptingStudentID uint) (*AcceptRequestResponse, error) {
    // ✅ 1. Get request
    request, err := s.requestRepo.FindByID(requestID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("request not found")
        }
        return nil, err
    }

    // ✅ 2. Verify the person accepting is the receiver
    if request.ToStudentID != acceptingStudentID {
        return nil, errors.New("you are not authorized to accept this request")
    }

    // ✅ 3. Check request is still pending
    if request.Status != "pending" {
        return nil, errors.New("request already processed")
    }

    // ✅ 4. Verify students are different (extra safety)
    if request.FromStudentID == request.ToStudentID {
        return nil, errors.New("invalid request: same student")
    }

    // ✅ 5. Check both students still don't have teams
    hasTeam, _ := s.teamRepo.StudentHasTeam(request.FromStudentID)
    if hasTeam {
        return nil, errors.New("sender already has a team")
    }

    hasTeam, _ = s.teamRepo.StudentHasTeam(request.ToStudentID)
    if hasTeam {
        return nil, errors.New("you already have a team")
    }

    // ✅ 6. Get student details for team creation
    var student1, student2 struct {
        Department string
        Session    string
    }

    s.db.Table("students").Where("id = ?", request.FromStudentID).
        Select("department, session").Scan(&student1)

    s.db.Table("students").Where("id = ?", request.ToStudentID).
        Select("department, session").Scan(&student2)

    // ✅ 7. Verify same department and session
    if student1.Department != student2.Department || student1.Session != student2.Session {
        return nil, errors.New("students must be in same department and session")
    }

    // ✅ 8. Create team in transaction
    tx := s.db.Begin()
    if tx.Error != nil {
        return nil, tx.Error
    }

    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    team := &Team{
        Name:        request.TeamName,
        ProjectName: request.ProjectName,
        Department:  student1.Department,
        Session:     student1.Session,
        Student1ID:  request.FromStudentID,
        Student2ID:  request.ToStudentID,
        Status:      "pending_supervisor",
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }

    if err := tx.Create(team).Error; err != nil {
        tx.Rollback()
        return nil, err
    }

    // ✅ 9. Update request status
    if err := tx.Model(&TeamRequest{}).Where("id = ?", requestID).
        Updates(map[string]interface{}{
            "status":     "accepted",
            "updated_at": time.Now(),
        }).Error; err != nil {
        tx.Rollback()
        return nil, err
    }

    // ✅ 10. Reject all other pending requests for both students
    tx.Model(&TeamRequest{}).
        Where("(from_student_id = ? OR to_student_id = ?) AND status = 'pending' AND id != ?",
            request.FromStudentID, request.FromStudentID, requestID).
        Updates(map[string]interface{}{
            "status":     "rejected",
            "updated_at": time.Now(),
        })

    tx.Model(&TeamRequest{}).
        Where("(from_student_id = ? OR to_student_id = ?) AND status = 'pending' AND id != ?",
            request.ToStudentID, request.ToStudentID, requestID).
        Updates(map[string]interface{}{
            "status":     "rejected",
            "updated_at": time.Now(),
        })

    if err := tx.Commit().Error; err != nil {
        return nil, err
    }

    return &AcceptRequestResponse{
        Message: "Team request accepted successfully",
        Team: TeamOverview{
            ID:          team.ID,
            Name:        team.Name,
            ProjectName: team.ProjectName,
            Status:      team.Status,
        },
    }, nil
}

// ============================================
// Reject Team Request
// ============================================

func (s *teamService) RejectRequest(requestID, rejectingStudentID uint) (*RejectRequestResponse, error) {
    // Get request
    request, err := s.requestRepo.FindByID(requestID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("request not found")
        }
        return nil, err
    }

    // Verify the person rejecting is the receiver
    if request.ToStudentID != rejectingStudentID {
        return nil, errors.New("you are not authorized to reject this request")
    }

    // Check request is still pending
    if request.Status != "pending" {
        return nil, errors.New("request already processed")
    }

    // Update status to rejected
    err = s.requestRepo.UpdateStatus(requestID, "rejected")
    if err != nil {
        return nil, err
    }

    return &RejectRequestResponse{
        Message: "Team request rejected successfully",
    }, nil
}

// ============================================
// Cancel Request (Delete sent request)
// ============================================

func (s *teamService) CancelRequest(requestID, studentID uint) (*CancelRequestResponse, error) {
    // Get request
    request, err := s.requestRepo.FindByID(requestID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("request not found")
        }
        return nil, err
    }

    // Verify the person canceling is the sender
    if request.FromStudentID != studentID {
        return nil, errors.New("you are not authorized to cancel this request")
    }

    // Check request is still pending
    if request.Status != "pending" {
        return nil, errors.New("cannot cancel a processed request")
    }

    // Delete request
    err = s.requestRepo.Delete(requestID)
    if err != nil {
        return nil, err
    }

    return &CancelRequestResponse{
        Message: "Team request cancelled successfully",
    }, nil
}

// ============================================
// Get My Team
// ============================================

func (s *teamService) GetMyTeam(studentID uint) (*GetMyTeamResponse, error) {
    // Check if student has team
    team, err := s.teamRepo.FindByStudentID(studentID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return &GetMyTeamResponse{
                HasTeam: false,
                Team:    nil,
            }, nil
        }
        return nil, err
    }

    // Get both team members
    var members []TeamMember

    // Get student1
    var student1 struct {
        ID                 uint
        FirstName          string
        LastName           string
        RegistrationNumber string
        Email              string
        Department         string
        Session            string
        Image              string
    }

    s.db.Table("students").
        Select("students.id, students.first_name, students.last_name, students.registration_number, students.department, students.session, students.image, users.email").
        Joins("JOIN users ON users.id = students.user_id").
        Where("students.id = ?", team.Student1ID).
        Scan(&student1)

    members = append(members, TeamMember{
        ID:                 student1.ID,
        Name:               student1.FirstName + " " + student1.LastName,
        Email:              student1.Email,
        RegistrationNumber: student1.RegistrationNumber,
        Department:         student1.Department,
        Session:            student1.Session,
        Image:              student1.Image,
    })

    // Get student2
    var student2 struct {
        ID                 uint
        FirstName          string
        LastName           string
        RegistrationNumber string
        Email              string
        Department         string
        Session            string
        Image              string
    }

    s.db.Table("students").
        Select("students.id, students.first_name, students.last_name, students.registration_number, students.department, students.session, students.image, users.email").
        Joins("JOIN users ON users.id = students.user_id").
        Where("students.id = ?", team.Student2ID).
        Scan(&student2)

    members = append(members, TeamMember{
        ID:                 student2.ID,
        Name:               student2.FirstName + " " + student2.LastName,
        Email:              student2.Email,
        RegistrationNumber: student2.RegistrationNumber,
        Department:         student2.Department,
        Session:            student2.Session,
        Image:              student2.Image,
    })

    return &GetMyTeamResponse{
        HasTeam: true,
        Team: &TeamDetails{
            ID:          team.ID,
            Name:        team.Name,
            ProjectName: team.ProjectName,
            Department:  team.Department,
            Session:     team.Session,
            Status:      team.Status,
            Members:     members,
            CreatedAt:   team.CreatedAt,
        },
    }, nil
}