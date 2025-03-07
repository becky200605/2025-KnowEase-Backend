package models

import "time"

type SearchRecord struct {
	SearchMessage string    `gorm:"not null" json:"searchmessage"`
	UserID        string    `gorm:"not null" json:"userid"`
	CreatedAt     time.Time 
}
type Outer struct {
	Data Outer1 `json:"data"`
}
type Outer1 struct {
	Realtime []Hot `json:"realtime"`
}
type Hot struct {
	Label string `json:"label_name"`
	Rank  int    `json:"rank"`
	Note  string `json:"note"`
}
