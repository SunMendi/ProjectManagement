package user

import (
	"time"

	"gorm.io/gorm"
)

type UserRepository interface {
	EmailExists(email string) (bool, error)
	CreateUser(user *User) error
	FindByID(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
}

type StudentRepository interface {
	RegistrationNumberExists(regNumber string) (bool, error)
	CreateStudent(student *Student) error
	FindByID(id uint) (*Student, error)
	UpdateStudent(student *Student) error 
	FindByDepartmentAndSession(department, session string) ([]Student, error)
	FindByUserID(userID uint) (*Student, error)
    HasTeam(studentID uint) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}


func NewUserRepository(db *gorm.DB) UserRepository {
	 return &userRepository{db: db}
}

func(r *userRepository) EmailExists(email string) (bool, error) {
	 var count int64 
	 err := r.db.Model(&User{}).Where("email = ?",email).Count(&count).Error 
	 return count>0, err 
}

func (r *userRepository) CreateUser(user *User) error {
    return r.db.Create(user).Error
}

type studentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) StudentRepository {
    return &studentRepository{db: db}
}

func (r *studentRepository) RegistrationNumberExists(regNumber string) (bool, error) {
    var count int64
    err := r.db.Model(&Student{}).Where("registration_number = ?", regNumber).Count(&count).Error
    return count > 0, err
}

func (r *studentRepository) CreateStudent(student *Student) error {
    return r.db.Create(student).Error
}

func(r *userRepository) FindByID(id uint) (*User, error) {
	 var user User 
	 err := r.db.Where("id = ?", id).First(&user).Error 
	 if err != nil {
		 return nil, err 
	 }
	 return &user, nil 
}

func (r *studentRepository) FindByID(id uint) (*Student, error) {
    var student Student
    err := r.db.Where("id = ?", id).First(&student).Error
    if err != nil {
        return nil, err
    }
    return &student, nil
}

func(r *studentRepository) UpdateStudent(student *Student) error {
	 student.UpdatedAt= time.Now()
	 return r.db.Save(student).Error 
}

func (r *studentRepository) FindByDepartmentAndSession(department, session string) ([]Student, error) {
    var students []Student
    err := r.db.
        Preload("User").
        Where("department = ? AND session = ?", department, session).
        Find(&students).Error
    if err != nil {
        return nil, err
    }
    return students, nil
}


func (r *userRepository) FindByEmail(email string) (*User, error) {
    var user User
    err := r.db.Where("email = ?", email).First(&user).Error
    if err != nil {
        return nil, err
    }
    return &user, nil
}

func (r *studentRepository) FindByUserID(userID uint) (*Student, error) {
    var student Student
    err := r.db.Preload("User").Where("user_id = ?", userID).First(&student).Error
    if err != nil {
        return nil, err
    }
    return &student, nil
}

func (r *studentRepository) HasTeam(studentID uint) (bool, error) {
    var count int64
    err := r.db.Table("teams").
        Where("student1_id = ? OR student2_id = ?", studentID, studentID).
        Count(&count).Error
    return count > 0, err
}

// ... existing repositories ...

type SupervisorRepository interface {
    CreateSupervisor(supervisor *Supervisor) error
    FindByID(id uint) (*Supervisor, error)
    FindByUserID(userID uint) (*Supervisor, error)
    UpdateSupervisor(supervisor *Supervisor) error
    FindByDepartment(department string) ([]Supervisor, error)
}

type supervisorRepository struct {
    db *gorm.DB
}

func NewSupervisorRepository(db *gorm.DB) SupervisorRepository {
    return &supervisorRepository{db: db}
}

func (r *supervisorRepository) CreateSupervisor(supervisor *Supervisor) error {
    return r.db.Create(supervisor).Error
}

func (r *supervisorRepository) FindByID(id uint) (*Supervisor, error) {
    var supervisor Supervisor
    err := r.db.Preload("User").Where("id = ?", id).First(&supervisor).Error
    if err != nil {
        return nil, err
    }
    return &supervisor, nil
}

func (r *supervisorRepository) FindByUserID(userID uint) (*Supervisor, error) {
    var supervisor Supervisor
    err := r.db.Preload("User").Where("user_id = ?", userID).First(&supervisor).Error
    if err != nil {
        return nil, err
    }
    return &supervisor, nil
}

func (r *supervisorRepository) UpdateSupervisor(supervisor *Supervisor) error {
    supervisor.UpdatedAt = time.Now()
    return r.db.Save(supervisor).Error
}

func (r *supervisorRepository) FindByDepartment(department string) ([]Supervisor, error) {
    var supervisors []Supervisor
    err := r.db.Preload("User").Where("department = ?", department).Find(&supervisors).Error
    if err != nil {
        return nil, err
    }
    return supervisors, nil
}