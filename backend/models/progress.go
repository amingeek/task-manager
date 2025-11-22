package models

import (
	"time"

	"gorm.io/gorm"
)

type TaskProgress struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	TaskID      uint           `gorm:"not null" json:"task_id"`
	Task        Task           `gorm:"foreignKey:TaskID" json:"task"`
	UserID      uint           `gorm:"not null" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"user"`
	Progress    int            `gorm:"default:0" json:"progress"` // 0-100
	Notes       string         `json:"notes"`
	IsCompleted bool           `gorm:"default:false" json:"is_completed"`
	CompletedAt *time.Time     `json:"completed_at"`
}

type GroupTaskProgress struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	TaskID      uint           `gorm:"not null" json:"task_id"`
	Task        Task           `gorm:"foreignKey:TaskID" json:"task"`
	UserID      uint           `gorm:"not null" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"user"`
	AssignedBy  uint           `gorm:"not null" json:"assigned_by"` // User ID of admin who assigned/updated
	Progress    int            `gorm:"default:0" json:"progress"`   // 0-100
	Notes       string         `json:"notes"`
	IsCompleted bool           `gorm:"default:false" json:"is_completed"`
	CompletedAt *time.Time     `json:"completed_at"`
	Approved    bool           `gorm:"default:false" json:"approved"`
	ApprovedBy  *uint          `json:"approved_by"`
	ApprovedAt  *time.Time     `json:"approved_at"`
}
