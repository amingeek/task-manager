// backend/models/notification.go

package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	UserID    uint           `gorm:"not null" json:"user_id"`
	Title     string         `gorm:"not null" json:"title"`
	Message   string         `json:"message"`
	IsRead    bool           `gorm:"default:false" json:"is_read"`
	Type      string         `json:"type"` // task_assigned, group_invitation, file_uploaded, etc.
	RelatedID uint           `json:"related_id"`
}

func (Notification) TableName() string {
	return "notifications"
}
