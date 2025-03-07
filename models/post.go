package models

import (
	"time"
)

type PostMessage struct {
	PostID     string    `gorm:"not null;primaryKey"`
	PosterID   string    `gorm:"not null" json:"posterid"`
	LikeCount  int       `gorm:"default:0"`
	SaveCount  int       `gorm:"default:0"`
	ViewCount  int       `gorm:"default:0"`
	Body       string    `gorm:"type:text" json:"body"`
	ImageURL   string    `grom:"type:varchar(255)" json:"urls"`
	Title      string    `json:"title"`
	Tag        string    `json:"tag"`
	Comment    []Comment `gorm:"-"`
	PosterName string
	PosterURL  string
	CreateAt   time.Time
}


type PostIDs struct {
	PostID string `json:"postids"`
}

// 帖子推荐指数
type PostRecommendLevel struct {
	PostID string
	Weight int
}
type UserViewHistory struct {
	ID       int    `gorm:"primaryKey"`
	PostID   string `gorm:"not null"`
	UserID   string `gorm:"not null"`
	Type     string
	CreateAt time.Time
}

type Comment struct {
	CommentID     string  `gorm:"not null" json:"commentid"`
	LikeCount     int     `gorm:"default:0"`
	ImageURL      string  `gorm:"type:varchar(255)" json:"imageurl"`
	CommenterID   string  `gorm:"not null"`
	PostID        string  `gorm:"not null"`
	Reply         []Reply `gorm:"-"`
	CommenterName string
	CommenterURL  string
	Body          string `gorm:"type:text" json:"body"`
	CreatAt       time.Time
}

// 需要实现发布评论的时候更新reply数组
type Reply struct {
	PostID      string `gorm:"not null" json:"postid"`
	ReplyID     string `gorm:"not null" json:"commentid"`
	LikeCount   int    `gorm:"default:0"`
	ImageURL    string `gorm:"type:varchar(255)" json:"imageurl"`
	ReplyerID   string `gorm:"not null"`
	CommentID   string `gorm:"not null"`
	ReplyerName string
	ReplyURL    string
	Reply       []Reply `gorm:"-"`
	Body        string  `gorm:"type:text" json:"body"`
	CreatAt     time.Time
}
type Delete struct {
	PostID string `gorm:"not null"`
}

type PostData struct {
	PostID    string `json:"postid"`
	LikeCount int    `json:"like_count"`
	SaveCount int    `json:"save_count"`
	ViewCount int    `json:"view_count"`
}

type QiNiuYunConfig struct {
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket_name"`
	Domain    string `mapstructure:"domain"`
}
