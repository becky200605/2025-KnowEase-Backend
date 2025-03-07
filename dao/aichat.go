package dao

import (
	"KnowEase/models"

	"gorm.io/gorm"
)

type AIChatDao struct {
	db *gorm.DB
}

func NewAIChatDao(db *gorm.DB) *AIChatDao {
	return &AIChatDao{db: db}
}
func ProvideAIChatDao(db *gorm.DB) AIChatInterface {
	return NewAIChatDao(db)
}

type AIChatInterface interface {
	SyncChatBodyToDB(Body *models.AIChatMessage) error
	SearchChatsByID(ChatID, UserID string) ([]models.AIChatMessage, error)
	SyncChatHistory(History *models.AIChatHistory) error
	GetUserChatHistory(UserID string) ([]models.AIChatHistory, error)
}

// 将聊天信息写入数据库
func (acd *AIChatDao) SyncChatBodyToDB(Body *models.AIChatMessage) error {
	return acd.db.Create(Body).Error
}

// 根据聊天分组查询所有聊天记录
func (acd *AIChatDao) SearchChatsByID(ChatID, UserID string) ([]models.AIChatMessage, error) {
	var History []models.AIChatMessage
	if err := acd.db.Order("created_at ASC").Where("chat_id = ? AND user_id = ?", ChatID).Find(&History).Error; err != nil {
		return nil, err
	}
	return History, nil
}

// 将聊天记录存入库
func (acd *AIChatDao) SyncChatHistory(History *models.AIChatHistory) error {
	return acd.db.Create(History).Error
}

// 查询用户所有聊天历史
func (acd *AIChatDao) GetUserChatHistory(UserID string) ([]models.AIChatHistory, error) {
	var UserHistory []models.AIChatHistory
	err := acd.db.Order("created_at DESC").Where("user_id = ?", UserID).Find(&UserHistory).Error
	return UserHistory, err
}

//删除聊天记录
func(acd *AIChatDao)DeleteChatHistory(ChatID,UserID string)error{
	return acd.db.Where("user_id = ? AND chat_id = ?",UserID,ChatID).Delete(&models.AIChatMessage{}).Error
}

//删除聊天历史
func(acd *AIChatDao)DeleteChatList(ChatID,UserID string)error{
	return acd.db.Where("user_id = ? AND chat_id = ?",UserID,ChatID).Delete(&models.AIChatHistory{}).Error
}
