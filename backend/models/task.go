package models

import (
	"time"
)

type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted  TaskStatus = "completed"
	StatusExpired    TaskStatus = "expired"
)

type Task struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	UserID      uint       `json:"user_id" gorm:"index"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      string     `json:"status" gorm:"default:'pending'"`
	Priority    string     `json:"priority" gorm:"default:'medium'"`
	DueDate     *time.Time `json:"due_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relations
	User     *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Progress []Progress `json:"progress,omitempty" gorm:"foreignKey:TaskID"`
	Files    []File     `json:"files,omitempty" gorm:"foreignKey:TaskID"`
}

type TaskAssignment struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID    uint      `gorm:"not null" json:"task_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Completed bool      `gorm:"default:false" json:"completed"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	Task      Task      `gorm:"foreignKey:TaskID" json:"task"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
