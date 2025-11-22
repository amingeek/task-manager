package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Username  string         `gorm:"uniqueIndex;not null;size:255" json:"username"`
	Email     string         `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	FullName  string         `json:"full_name"`
	Bio       string         `json:"bio"`
	AvatarURL string         `json:"avatar_url"`
}

type Streak struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	UserID         uint      `gorm:"uniqueIndex;not null" json:"user_id"`
	CurrentStreak  int       `gorm:"default:0" json:"current_streak"`
	LongestStreak  int       `gorm:"default:0" json:"longest_streak"`
	LastActivityAt time.Time `json:"last_activity_at"`
}
