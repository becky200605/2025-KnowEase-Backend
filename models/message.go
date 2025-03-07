package models

import "time"

// 消息推送
type Message struct {
	UserID    string `gorm:"not null"`
	PosterURL string `grom:"type:varchar(255)"`
	Message   string
	PostID    string
	Tag       string
	Status    string `gorm:"not null"`
	CreateAt  time.Time
}