package models

import (
	"time"
)

type Group struct {
	ID          uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string        `gorm:"not null" json:"name"`
	Description string        `json:"description"`
	CreatorID   uint          `gorm:"not null" json:"creator_id"`
	Creator     User          `gorm:"foreignKey:CreatorID" json:"creator"`
	Members     []GroupMember `json:"members"`
	Tasks       []Task        `json:"tasks"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

type GroupMember struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	GroupID   uint      `gorm:"not null" json:"group_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	Accepted  bool      `gorm:"default:false" json:"accepted"`
	Role      string    `gorm:"default:'member'" json:"role"` // 'member', 'admin'
	CreatedAt time.Time `json:"created_at"`
}
