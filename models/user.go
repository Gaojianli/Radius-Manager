package models

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	Salt      string    `json:"-" gorm:"not null"`
	IsAdmin   bool      `json:"is_admin" gorm:"default:false"`
	Email     string    `json:"email" gorm:"unique"`
	Status    bool      `json:"status" gorm:"default:true"`
	Banned    bool      `json:"banned" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) generateSalt() (string, error) {
	saltBytes := make([]byte, 32)
	_, err := rand.Read(saltBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(saltBytes), nil
}

func (u *User) HashPassword(password string) error {
	salt, err := u.generateSalt()
	if err != nil {
		return err
	}
	u.Salt = salt

	h := sha256.New()
	h.Write([]byte(password + salt))
	u.Password = hex.EncodeToString(h.Sum(nil))
	return nil
}

func (u *User) CheckPassword(password string) bool {
	h := sha256.New()
	h.Write([]byte(password + u.Salt))
	hashedPassword := hex.EncodeToString(h.Sum(nil))
	return hashedPassword == u.Password
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Password != "" {
		return u.HashPassword(u.Password)
	}
	return nil
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	IsAdmin   bool      `json:"is_admin"`
	Status    bool      `json:"status"`
	Banned    bool      `json:"banned"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		IsAdmin:   u.IsAdmin,
		Status:    u.Status,
		Banned:    u.Banned,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
