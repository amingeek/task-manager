package models

import (
	"time"
)

type Group struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name" gorm:"index"`
	Description string    `json:"description"`
	CreatedBy   uint      `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	Members []GroupMember `json:"members,omitempty" gorm:"foreignKey:GroupID"`
	Tasks   []GroupTask   `json:"tasks,omitempty" gorm:"foreignKey:GroupID"`
	Creator *User         `json:"creator,omitempty" gorm:"foreignKey:CreatedBy;references:ID"`
}

type GroupMember struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	GroupID   uint      `json:"group_id" gorm:"index"`
	UserID    uint      `json:"user_id" gorm:"index"`
	Role      string    `json:"role" gorm:"default:'member'"`
	Accepted  bool      `json:"accepted" gorm:"default:false;index"` // ✅ برای دعوت‌نامه
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations - تعریف روابط
	Group *Group `json:"group,omitempty" gorm:"foreignKey:GroupID;references:ID"`
	User  *User  `json:"user,omitempty" gorm:"foreignKey:UserID;references:ID"`
}

type GroupTask struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	GroupID     uint      `json:"group_id" gorm:"index"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status" gorm:"default:'pending'"`
	Priority    string    `json:"priority" gorm:"default:'medium'"`
	AssignedBy  uint      `json:"assigned_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relations
	Group          *Group          `json:"group,omitempty" gorm:"foreignKey:GroupID"`
	Progress       []GroupProgress `json:"progress,omitempty" gorm:"foreignKey:GroupTaskID"`
	Files          []File          `json:"files,omitempty" gorm:"foreignKey:GroupTaskID"`
	AssignedByUser *User           `json:"assigned_by_user,omitempty" gorm:"foreignKey:AssignedBy;references:ID"`
}

type GroupProgress struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	GroupTaskID    uint      `json:"group_task_id" gorm:"index"`
	UserID         uint      `json:"user_id" gorm:"index"`
	Status         string    `json:"status" gorm:"default:'pending'"`
	PercentageDone int       `json:"percentage_done" gorm:"default:0"`
	Notes          string    `json:"notes"`
	LastUpdatedBy  uint      `json:"last_updated_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Relations
	GroupTask *GroupTask `json:"group_task,omitempty" gorm:"foreignKey:GroupTaskID"`
	User      *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
