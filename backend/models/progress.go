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

type Progress struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TaskID    uint      `json:"task_id" gorm:"index"`
	UserID    uint      `json:"user_id" gorm:"index"`
	Status    string    `json:"status" gorm:"default:'pending'"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Task *Task `json:"task,omitempty" gorm:"foreignKey:TaskID"`
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
