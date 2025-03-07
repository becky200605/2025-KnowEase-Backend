package dao

import (
	"KnowEase/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 连接数据库
func NewDB(addr string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(addr), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	db.AutoMigrate(&models.User{}, &models.Emailverify{}, &models.Comment{}, &models.Message{}, &models.PostMessage{}, &models.Reply{}, &models.UserLikeHistory{}, &models.UserSaveHistory{}, &models.UserViewHistory{}, &models.FollowMessage{}, &models.QAs{}, &models.SearchRecord{}, &models.AIChatMessage{},&models.AIChatHistory{})
	if err != nil {
		panic(err)
	}
	return db
}
