package models

import (
	"time"
)

type File struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	TaskID      *uint      `json:"task_id"`
	GroupTaskID *uint      `json:"group_task_id"`
	UserID      uint       `json:"user_id" gorm:"index"`
	FileName    string     `json:"file_name"`
	FilePath    string     `json:"file_path"`
	FileSize    int64      `json:"file_size"`
	FileType    string     `json:"file_type"`
	UploadedAt  time.Time  `json:"uploaded_at"`
	ApprovedBy  *uint      `json:"approved_by"`
	ApprovedAt  *time.Time `json:"approved_at"`

	// Relations
	Task      *Task      `json:"task,omitempty" gorm:"foreignKey:TaskID"`
	GroupTask *GroupTask `json:"group_task,omitempty" gorm:"foreignKey:GroupTaskID"`
	User      *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (File) TableName() string {
	return "files"
}
