package user

import (
	"ProjectManagement/utils"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type StudentService interface {
	 RegisterStudent(req RegisterStudentRequest) (*RegisterStudentResponse, error)
	 GetStudentByID(id uint) (*GetStudentResponse, error) 
	 UpdateStudentProfile(id uint, req UpdateStudentProfileRequest) ( error)
	 GetStudentsByFilter(department, session string) (*GetStudentsByFilterResponse, error)
	 Login(req LoginRequest) (*LoginResponse, error)
}

type studentService struct {
	 db *gorm.DB 
	 userRepo UserRepository
	 studentRepo StudentRepository
}

func NewStudentService(db *gorm.DB, userRepo UserRepository, studentRepo StudentRepository) StudentService {
    return &studentService{
        db:          db,
        userRepo:    userRepo,
        studentRepo: studentRepo,
    }
}

func(s *studentService) RegisterStudent(req RegisterStudentRequest) (*RegisterStudentResponse, error) {
	 exists, err := s.userRepo.EmailExists(req.Email)
	 if err != nil {
		 return nil, err 
	 }
	 if exists {
        return nil, errors.New("email already exists")
    }
	exists, err = s.studentRepo.RegistrationNumberExists(req.RegistrationNumber)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, errors.New("registration_number already exists")
    }
	hashedPassword , err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
	   return nil, err 
	}

	tx:= s.db.Begin()
	if tx.Error != nil {
		 return nil, tx.Error 
	}
	defer func() {
		 if r:= recover() ; r!= nil {
			 tx.Rollback()
		 }
	}()

	var user User 
	var student Student 

	user = User {
		 Email: req.Email,
		 Password: string(hashedPassword),
		 Role: "student",
		 Status: "pending",
		 CreatedAt: time.Now(),
		 UpdatedAt: time.Now(),
	}

	if err:= tx.Create(&user).Error ; err != nil {
		 tx.Rollback()
		 return nil, err 
	}

	student = Student{
        UserID:             user.ID, // Use the ID from created user
        FirstName:          req.FirstName,
        LastName:           req.LastName,
        Department:         req.Department,
        Session:            req.Session,
        RegistrationNumber: req.RegistrationNumber,
        Batch:              req.Batch,
        CreatedAt:          time.Now(),
        UpdatedAt:          time.Now(),
    }

    if err := tx.Create(&student).Error; err != nil {
        tx.Rollback() 
        return nil, err
    }
	if err := tx.Commit().Error; err != nil {
        return nil, err
    }
	return &RegisterStudentResponse{
        Message:   "Registration successful. Waiting for admin approval.",
        UserID:    user.ID,
        StudentID: student.ID,
    }, nil
}

func (s *studentService) GetStudentByID(id uint) (*GetStudentResponse, error) {
    student, err := s.studentRepo.FindByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("student not found")
        }
        return nil, err
    }

    user, err := s.userRepo.FindByID(student.UserID)
    if err != nil {
        return nil, err
    }

    return &GetStudentResponse{
        ID:                 student.ID,
        UserID:             student.UserID,
        Email:              user.Email,
        FirstName:          student.FirstName,
        LastName:           student.LastName,
        Department:         student.Department,
        Session:            student.Session,
        RegistrationNumber: student.RegistrationNumber,
        Batch:              student.Batch,
		Status:             user.Status,
        Role:               user.Role,
        Image:              student.Image,
    }, nil
}

func (s *studentService) UpdateStudentProfile(id uint, req UpdateStudentProfileRequest) (error) {
    // Find student
    student, err := s.studentRepo.FindByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return errors.New("student not found")
        }
        return err
    }

    // Update fields
    if req.FirstName != "" {
        student.FirstName = req.FirstName
    }
    
    if req.LastName != "" {
        student.LastName = req.LastName
    }
    
    if req.Image != "" {
        student.Image = req.Image
    }
    
    if req.Department != "" {
        student.Department = req.Department
    }

    err = s.studentRepo.UpdateStudent(student)
    if err != nil {
        return err
    }
	return nil 
}

func(s *studentService) GetStudentsByFilter(department, session string) (*GetStudentsByFilterResponse, error) {
	students, err := s.studentRepo.FindByDepartmentAndSession(department,session)
	if err != nil {
		 return nil, err 
	}
	var studentList []StudentListItem
	for _, student := range students {
        if student.User == nil{
            continue
        }
        hasTeam, _ := s.studentRepo.HasTeam(student.ID)
        studentList = append(studentList, StudentListItem{
            ID:                 student.ID,
            UserID:             student.UserID,
            Email:              student.User.Email,
            FirstName:          student.FirstName,
            LastName:           student.LastName,
            Department:         student.Department,
            Session:            student.Session,
            RegistrationNumber: student.RegistrationNumber,
            Batch:              student.Batch,
            Image:              student.Image,
            Status:             student.User.Status,
            HasTeam:            hasTeam,
        })
    }
	return &GetStudentsByFilterResponse{
        Total:    len(studentList),
        Students: studentList,
    }, nil
}



func (s *studentService) Login(req LoginRequest) (*LoginResponse, error) {
    // Find user by email
    user, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("invalid email or password")
        }
        return nil, err
    }

    // Check if user is a student
    if user.Role != "student" {
        return nil, errors.New("invalid credentials for student login")
    }

    // Verify password
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
    if err != nil {
        return nil, errors.New("invalid email or password")
    }

    // Check account status
    if user.Status == "pending" {
        return nil, errors.New("account is pending approval")
    }

    if user.Status == "rejected" {
        return nil, errors.New("account has been rejected")
    }

    if user.Status != "active" {
        return nil, errors.New("account is not active")
    }

    // Get student details
    student, err := s.studentRepo.FindByUserID(user.ID)
    if err != nil {
        return nil, errors.New("student profile not found")
    }

	token, err := utils.GenerateStudentToken(user.ID, student.ID, user.Email)
    if err != nil {
        return nil, err
    }

    // Prepare response
    studentData := GetStudentResponse{
        ID:                 student.ID,
        UserID:             student.UserID,
        Email:              user.Email,
        FirstName:          student.FirstName,
        LastName:           student.LastName,
        Department:         student.Department,
        Session:            student.Session,
        RegistrationNumber: student.RegistrationNumber,
        Batch:              student.Batch,
        Status:             user.Status,
        Role:               user.Role,
    }

    return &LoginResponse{
        Message:   "Login successful",
		Token: token,
        UserID:    user.ID,
        StudentID: student.ID,
        Role:      user.Role,
        Status:    user.Status,
        Student:   studentData,
    }, nil
}

type SupervisorService interface {
    RegisterSupervisor(req RegisterSupervisorRequest) (*RegisterSupervisorResponse, error)
    GetSupervisorByID(id uint) (*GetSupervisorResponse, error)
    UpdateSupervisorProfile(id uint, req UpdateSupervisorProfileRequest) error
    GetSupervisorsByDepartment(department string) (*GetSupervisorsByDepartmentResponse, error)
    Login(req LoginRequest) (*SupervisorLoginResponse, error)
}

type supervisorService struct {
    db             *gorm.DB
    userRepo       UserRepository
    supervisorRepo SupervisorRepository
}

func NewSupervisorService(db *gorm.DB, userRepo UserRepository, supervisorRepo SupervisorRepository) SupervisorService {
    return &supervisorService{
        db:             db,
        userRepo:       userRepo,
        supervisorRepo: supervisorRepo,
    }
}

func (s *supervisorService) RegisterSupervisor(req RegisterSupervisorRequest) (*RegisterSupervisorResponse, error) {
    exists, err := s.userRepo.EmailExists(req.Email)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, errors.New("email already exists")
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }

    tx := s.db.Begin()
    if tx.Error != nil {
        return nil, tx.Error
    }

    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    user := User{
        Email:     req.Email,
        Password:  string(hashedPassword),
        Role:      "supervisor",
        Status:    "active",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    if err := tx.Create(&user).Error; err != nil {
        tx.Rollback()
        return nil, err
    }

    supervisor := Supervisor{
        UserID:       user.ID,
        Name:         req.Name,
        Designation:  req.Designation,
        Department:   req.Department,
        ResearchArea: req.ResearchArea,
        Phone:        req.Phone,
        CreatedAt:    time.Now(),
        UpdatedAt:    time.Now(),
    }

    if err := tx.Create(&supervisor).Error; err != nil {
        tx.Rollback()
        return nil, err
    }

    if err := tx.Commit().Error; err != nil {
        return nil, err
    }

    return &RegisterSupervisorResponse{
        Message:      "Account Created Successfully",
        UserID:       user.ID,
        SupervisorID: supervisor.ID,
    }, nil
}

func (s *supervisorService) GetSupervisorByID(id uint) (*GetSupervisorResponse, error) {
    supervisor, err := s.supervisorRepo.FindByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("supervisor not found")
        }
        return nil, err
    }

    if supervisor.User == nil || supervisor.User.ID == 0 {
        return nil, errors.New("user data not found")
    }

    return &GetSupervisorResponse{
        ID:           supervisor.ID,
        UserID:       supervisor.UserID,
        Email:        supervisor.User.Email,
        Name:         supervisor.Name,
        Designation:  supervisor.Designation,
        Department:   supervisor.Department,
        ResearchArea: supervisor.ResearchArea,
        Phone:        supervisor.Phone,
        Image:        supervisor.Image,
        Status:       supervisor.User.Status,
        Role:         supervisor.User.Role,
    }, nil
}

func (s *supervisorService) UpdateSupervisorProfile(id uint, req UpdateSupervisorProfileRequest) error {
    supervisor, err := s.supervisorRepo.FindByID(id)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return errors.New("supervisor not found")
        }
        return err
    }

    if req.Name != "" {
        supervisor.Name = req.Name
    }
    
    if req.Designation != "" {
        supervisor.Designation = req.Designation
    }
    
    if req.Department != "" {
        supervisor.Department = req.Department
    }
    
    if req.ResearchArea != "" {
        supervisor.ResearchArea = req.ResearchArea
    }
    
    if req.Phone != "" {
        supervisor.Phone = req.Phone
    }
    
    if req.Image != "" {
        supervisor.Image = req.Image
    }

    return s.supervisorRepo.UpdateSupervisor(supervisor)
}

func (s *supervisorService) GetSupervisorsByDepartment(department string) (*GetSupervisorsByDepartmentResponse, error) {
    supervisors, err := s.supervisorRepo.FindByDepartment(department)
    if err != nil {
        return nil, err
    }

    var supervisorList []SupervisorListItem
    for _, supervisor := range supervisors {
        if supervisor.User == nil || supervisor.User.ID == 0 {
            continue
        }

        supervisorList = append(supervisorList, SupervisorListItem{
            ID:           supervisor.ID,
            UserID:       supervisor.UserID,
            Email:        supervisor.User.Email,
            Name:         supervisor.Name,
            Designation:  supervisor.Designation,
            Department:   supervisor.Department,
            ResearchArea: supervisor.ResearchArea,
            Phone:        supervisor.Phone,
            Image:        supervisor.Image,
            Status:       supervisor.User.Status,
        })
    }

    return &GetSupervisorsByDepartmentResponse{
        Total:       len(supervisorList),
        Supervisors: supervisorList,
    }, nil
}

func (s *supervisorService) Login(req LoginRequest) (*SupervisorLoginResponse, error) {
    user, err := s.userRepo.FindByEmail(req.Email)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("invalid email or password")
        }
        return nil, err
    }

    if user.Role != "supervisor" {
        return nil, errors.New("invalid credentials for supervisor login")
    }

    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
    if err != nil {
        return nil, errors.New("invalid email or password")
    }

    if user.Status == "pending" {
        return nil, errors.New("account is pending approval")
    }

    if user.Status == "rejected" {
        return nil, errors.New("account has been rejected")
    }

    if user.Status != "active" {
        return nil, errors.New("account is not active")
    }

    supervisor, err := s.supervisorRepo.FindByUserID(user.ID)
    if err != nil {
        return nil, errors.New("supervisor profile not found")
    }

    token, err := utils.GenerateSupervisorToken(user.ID, supervisor.ID, user.Email)
    if err != nil {
        return nil, err
    }

    supervisorData := GetSupervisorResponse{
        ID:           supervisor.ID,
        UserID:       supervisor.UserID,
        Email:        user.Email,
        Name:         supervisor.Name,
        Designation:  supervisor.Designation,
        Department:   supervisor.Department,
        ResearchArea: supervisor.ResearchArea,
        Phone:        supervisor.Phone,
        Image:        supervisor.Image,
        Status:       user.Status,
        Role:         user.Role,
    }

    return &SupervisorLoginResponse{
        Message:      "Login successful",
        Token:        token,
        UserID:       user.ID,
        SupervisorID: supervisor.ID,
        Role:         user.Role,
        Status:       user.Status,
        Supervisor:   supervisorData,
    }, nil
}


// ============================================
// ADMIN SERVICE (Super Simple)
// ============================================

type AdminService interface {
    GetPendingUsers() (*GetPendingUsersResponse, error) // ✅ Changed return type
    ApproveUser(userID uint) error
    RejectUser(userID uint) error
}

type adminService struct {
    db *gorm.DB
}

func NewAdminService(db *gorm.DB) AdminService {
    return &adminService{db: db}
}

// GetPendingUsers - Get all pending students only
// GetPendingUsers - Get all pending students with their details
func (s *adminService) GetPendingUsers() (*GetPendingUsersResponse, error) {
    // ✅ Query: Join users with students table
    var results []struct {
        UserID             uint
        Email              string
        Role               string
        Status             string
        CreatedAt          time.Time
        StudentID          uint
        FirstName          string
        LastName           string
        Department         string
        Session            string
        RegistrationNumber string
        Batch              string
    }

    err := s.db.Table("users").
        Select(`
            users.id as user_id,
            users.email,
            users.role,
            users.status,
            users.created_at,
            students.id as student_id,
            students.first_name,
            students.last_name,
            students.department,
            students.session,
            students.registration_number,
            students.batch
        `).
        Joins("LEFT JOIN students ON students.user_id = users.id").
        Where("users.status = ? AND users.role = ?", "pending", "student").
        Order("users.created_at DESC").
        Scan(&results).Error

    if err != nil {
        return nil, err
    }

    // ✅ Build response
    var userList []PendingUserItem
    for _, r := range results {
        userList = append(userList, PendingUserItem{
            UserID:             r.UserID,
            Email:              r.Email,
            Role:               r.Role,
            Status:             r.Status,
            CreatedAt:          r.CreatedAt,
            StudentID:          r.StudentID,
            FirstName:          r.FirstName,
            LastName:           r.LastName,
            Department:         r.Department,
            Session:            r.Session,
            RegistrationNumber: r.RegistrationNumber,
            Batch:              r.Batch,
        })
    }

    return &GetPendingUsersResponse{
        Total: len(userList),
        Users: userList,
    }, nil
}

// ApproveUser - Change status from pending to active
func (s *adminService) ApproveUser(userID uint) error {
    result := s.db.Model(&User{}).
        Where("id = ? AND status = ? AND role = ?", userID, "pending", "student").
        Update("status", "active")
    
    if result.Error != nil {
        return result.Error
    }
    
    if result.RowsAffected == 0 {
        return errors.New("user not found or already approved")
    }
    
    return nil
}

func (s *adminService) RejectUser(userID uint) error {
    result := s.db.Model(&User{}).
        Where("id = ? AND status = ? AND role = ?", userID, "pending", "student").
        Update("status", "rejected")

    if result.Error != nil {
        return result.Error
    }

    if result.RowsAffected == 0 {
        return errors.New("user not found or already processed")
    }

    return nil
}