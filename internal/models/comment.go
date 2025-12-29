package models

import (
	"time"

	"gorm.io/gorm"
)

// Comment represents a comment on a task
type Comment struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Content   string         `json:"content" gorm:"type:text;not null"` // Comment text
	TaskID    uint           `json:"task_id" gorm:"not null;index"`     // ID of the task this comment belongs to
	UserID    uint           `json:"user_id" gorm:"not null;index"`      // ID of the user who created the comment
	Task      Task           `json:"task,omitempty" gorm:"foreignKey:TaskID"`
	User      User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

