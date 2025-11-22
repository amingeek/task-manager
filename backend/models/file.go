package models

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Filename  string         `gorm:"not null" json:"filename"`
	Filepath  string         `gorm:"not null" json:"filepath"`
	FileSize  int64          `json:"file_size"`
	MimeType  string         `json:"mime_type"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"user"`
	TaskID    uint           `gorm:"not null" json:"task_id"`
	Task      Task           `gorm:"foreignKey:TaskID" json:"task"`
}
