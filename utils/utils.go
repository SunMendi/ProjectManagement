package utils

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key-change-this-in-production")

// ✅ Updated Claims with SupervisorID
type Claims struct {
    UserID       uint   `json:"user_id"`
    StudentID    uint   `json:"student_id,omitempty"`    // Only for students
    SupervisorID uint   `json:"supervisor_id,omitempty"` // Only for supervisors
    Email        string `json:"email"`
    Role         string `json:"role"`
    jwt.RegisteredClaims
}

// ✅ Generate token for Student
func GenerateStudentToken(userID, studentID uint, email string) (string, error) {
    claims := Claims{
        UserID:    userID,
        StudentID: studentID,
        Email:     email,
        Role:      "student",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

// ✅ Generate token for Supervisor
func GenerateSupervisorToken(userID, supervisorID uint, email string) (string, error) {
    claims := Claims{
        UserID:       userID,
        SupervisorID: supervisorID,
        Email:        email,
        Role:         "supervisor",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

// ✅ Keep the old function for backward compatibility (deprecated)
// You can remove this after updating all Login calls
func GenerateToken(userID, studentID uint, email, role string) (string, error) {
    if role == "student" {
        return GenerateStudentToken(userID, studentID, email)
    }
    // For supervisor, studentID param will be supervisorID
    return GenerateSupervisorToken(userID, studentID, email)
}

// Validate JWT token (no changes needed)
func ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}

// ✅ Helper to get JWT secret (for external use)
func GetJWTSecret() string {
    return string(jwtSecret)
}