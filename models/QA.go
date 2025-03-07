package models

import "time"

const (
	QATableName = "QAs"
)

type QAs struct {
	PostID     string
	Question   string `gorm:"type:text" json:"question"`
	AuthorName string
	Comment    []Comment `gorm:"-"`
	Weight     float64   `gorm:"-"`
	AuthorID   string
	AuthorURL  string
	LikeCount  int `gorm:"default:0"`
	SaveCount  int `gorm:"default:0"`
	ViewCount  int `gorm:"default:0"`
	Tag        string
	Answer     string
	CreatedAt  time.Time
	UpdatedAT  time.Time
}

//模型回答结构体
type Responses struct {
	Similarity float64 `json:"similarity"`
	Answer     string  `json:"answer"`
	PostIDs    string  `json:"postids"`
}
type Type struct {
	UserID string `json:"userid"`
	Tag    string `json:"tag"`
}

//接受用户消息
type ChatRequest struct {
	Message string `json:"message"`
}

func (q *QAs) TableName() string {
	return QATableName
}
