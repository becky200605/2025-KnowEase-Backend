package dao

import (
	"KnowEase/models"

	"gorm.io/gorm"
)

type SearchDao struct {
	db *gorm.DB
}

func NewSearchDao(db *gorm.DB) *SearchDao {
	return &SearchDao{db: db}
}
func ProvideSearchDao(db *gorm.DB) SearchDaoInterface {
	return NewSearchDao(db)
}

type SearchDaoInterface interface {
	SyncSearchRecord(Record models.SearchRecord) error
	SearchUserRecord(UserID string) ([]models.SearchRecord, error)
	SearchPost(Message string) ([]models.PostMessage, error)
	SearchQA(Message string) ([]models.QAs, error)
	UpdateRecord(UserID string, Record []string) error
	DeleteRecord(UserID, Message string) error
	DeleteAllRecord(UserID string) error
	SearchMessage(Message, UserID string) error
	UpdateCreatedTime(Record models.SearchRecord) error
}

// 存入搜索记录
func (sd *SearchDao) SyncSearchRecord(Record models.SearchRecord) error {
	return sd.db.Create(&Record).Error
}

// 更新时间
func (sd *SearchDao) UpdateCreatedTime(Record models.SearchRecord) error {
	return sd.db.Model(&models.SearchRecord{}).Where("search_message = ? AND user_id = ?", Record.SearchMessage, Record.UserID).Update("created_at", Record.CreatedAt).Error
}

// 查询用户搜索记录
func (sd *SearchDao) SearchUserRecord(UserID string) ([]models.SearchRecord, error) {
	var Records []models.SearchRecord
	err := sd.db.Where("user_id = ?", UserID).Limit(10).Find(&Records).Error
	return Records, err
}

// 模糊搜索帖子内容
func (sd *SearchDao) SearchPost(Message string) ([]models.PostMessage, error) {
	var Records []models.PostMessage
	err := sd.db.Where("title LIKE ? OR body LIKE ?", "%"+Message+"%", "%"+Message+"%").Find(&Records).Error
	return Records, err
}

// 模糊搜索问答内容
func (sd *SearchDao) SearchQA(Message string) ([]models.QAs, error) {
	var Records []models.QAs
	err := sd.db.Where("question Like ?", "%"+Message+"%").Find(&Records).Error
	return Records, err
}

// 更新搜索库
func (sd *SearchDao) UpdateRecord(UserID string, Record []string) error {
	err := sd.db.Where("message NOT IN ? AND user_id = ?", Record, UserID).Delete(&models.SearchRecord{}).Error
	return err
}

// 删除搜索记录
func (sd *SearchDao) DeleteRecord(UserID, Message string) error {
	err := sd.db.Where("search_message = ? AND user_id = ?", Message, UserID).Delete(&models.SearchRecord{}).Error
	return err
}

// 清空搜索记录
func (sd *SearchDao) DeleteAllRecord(UserID string) error {
	err := sd.db.Where("user_id = ?", UserID).Delete(&models.SearchRecord{}).Error
	return err
}

// 查询用户搜索记录
func (sd *SearchDao) SearchMessage(Message, UserID string) error {
	var Record models.SearchRecord
	err := sd.db.Where("user_id = ? AND search_message = ?", UserID, Message).First(&Record).Error
	return  err
}
