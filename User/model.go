package user

import (
	"time"

	"gorm.io/gorm"
)


type UserRole string 

type UserStatus string 



const (
	 RoleAdmin UserRole = "admin"
	 RoleStudent UserRole = "student"
	 RoleSupervisor UserRole = "supervisor"
)

const (
    StatusPending  UserStatus = "pending" 
    StatusActive   UserStatus = "active"  
    StatusRejected UserStatus = "rejected" 
)


type User struct {
	 ID uint `json:"id" gorm:"primaryKey"`
	 Email string `json:"email" gorm:"uniqueIndex;not null"`
	 Password string `json:"-" gorm:"not null"`
	 Role UserRole `json:"user_role"`
	 Status UserStatus `json:"user_status"`
	 CreatedAt time.Time `json:"created_at"`
	 UpdatedAt time.Time `json:"updated_at"`
	 DeletedAt gorm.DeletedAt `json:"-"`

	 //Relationship
	 
    Student    *Student    `gorm:"foreignKey:UserID" json:"student,omitempty"`
    Supervisor *Supervisor `gorm:"foreignKey:UserID" json:"supervisor,omitempty"`
}

type Student struct {
    ID                 uint      `json:"id" gorm:"primaryKey"`
    UserID             uint      `json:"user_id" gorm:"not null;uniqueIndex"`
    FirstName          string    `json:"first_name" gorm:"not null"`
    LastName           string    `json:"last_name" gorm:"not null"`
	Image              string    `json:"image"`
    Department         string    `json:"department" gorm:"not null"`
    StudentID          string    `json:"student_id" gorm:"uniqueIndex;not null"`          // Roll
    Session            string    `json:"session" gorm:"not null"`
    RegistrationNumber string    `json:"registration_number" gorm:"uniqueIndex;not null"` // Registration
    Batch              string    `json:"batch" gorm:"not null"`
    CreatedAt          time.Time `json:"created_at"`
    UpdatedAt          time.Time `json:"updated_at"`

    // Relationship
    User User `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type Supervisor struct {
    ID           uint      `json:"id" gorm:"primaryKey"`
    UserID       uint      `json:"user_id" gorm:"not null;uniqueIndex"`
    Name         string    `json:"name" gorm:"not null"`
    Image        string    `json:"image"`
    Designation  string    `json:"designation" gorm:"not null"`
    ResearchArea string    `json:"research_area"`
    Department   string    `json:"department" gorm:"not null"`
    Phone        string    `json:"phone"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`

    // Relationship
    User User `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}