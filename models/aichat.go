package models

import "time"

type AIChatMessage struct {
	UserID    string    `json:"userid"`
	Response  string    `gorm:"type:text" json:"response"`
	Request   string    `gorm:"type:text" json:"request"`
	CreatedAt time.Time `json:"createdAt"`
	ChatID    string    `json:"chatid"`
}

type AIChatHistory struct {
	UserID    string    `json:"userid"`
	ChatID    string    `json:"chatid"`
	CreatedAt time.Time `json:"createdAt"`
	History   string    `gorm:"type:text" json:"history"`
}
