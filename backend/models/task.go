package models

import (
	"time"

	"gorm.io/gorm"
)

type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted  TaskStatus = "completed"
	StatusExpired    TaskStatus = "expired"
)

type Task struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	Status      TaskStatus     `gorm:"default:'pending'" json:"status"`
	DueDate     *time.Time     `json:"due_date"`
	StartTime   *time.Time     `json:"start_time"`
	EndTime     *time.Time     `json:"end_time"`
	CreatorID   uint           `gorm:"not null" json:"creator_id"`
	Creator     User           `gorm:"foreignKey:CreatorID" json:"creator"`
	IsGroupTask bool           `gorm:"default:false" json:"is_group_task"`
	GroupID     *uint          `json:"group_id"`
	Group       *Group         `gorm:"foreignKey:GroupID" json:"group,omitempty"`

	// فیلدهای جدید
	RequireFiles bool   `gorm:"default:false" json:"require_files"`
	MaxFiles     int    `gorm:"default:5" json:"max_files"`
	AllowTypes   string `json:"allow_types"` // pdf,image,video,etc

	// روابط
	Files           []File              `json:"files,omitempty"`
	Progress        []TaskProgress      `json:"progress,omitempty"`
	GroupProgress   []GroupTaskProgress `json:"group_progress,omitempty"`
	TaskAssignments []TaskAssignment    `json:"task_assignments,omitempty"`
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
