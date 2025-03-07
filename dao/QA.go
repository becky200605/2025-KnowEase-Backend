package dao

import (
	"KnowEase/models"
	"time"

	"gorm.io/gorm"
)

type QADao struct {
	db *gorm.DB
}

func NewQADao(db *gorm.DB) *QADao {
	return &QADao{db: db}
}
func ProvideQADao(db *gorm.DB) QADaoInterface {
	return NewQADao(db)
}

type QADaoInterface interface {
	SyncQABodyToDB(Body *models.QAs) error
	DeleteQABody(PostID string) error
	SearchAllQA() ([]models.QAs, error)
	SearchUnviewedQAByTag(UserViewRecordID []string, Tag string) ([]models.QAs, error)
	SearchQAByID(PostID string) (models.QAs, error)
	SearchUpdateTime(LastUpdateTime time.Time) (int64, error)
	GetUserQA(UserID string) ([]models.QAs, error)
	GetCommentCount(PostID string) (int64, error)
}

// 将问答内容写入数据库
func (qd *QADao) SyncQABodyToDB(Body *models.QAs) error {
	Body.CreatedAt = time.Now()
	Body.UpdatedAT = time.Now()
	return qd.db.Create(Body).Error

}

// 将问答内容从数据库里删除
func (qd *QADao) DeleteQABody(PostID string) error {
	return qd.db.Model(&models.QAs{}).Delete("post_id = ?", PostID).Error
}

// 查询所有问答帖子
func (qd *QADao) SearchAllQA() ([]models.QAs, error) {
	var Posts []models.QAs
	err := qd.db.Model(&models.QAs{}).Not("tag = ?","text").Find(&Posts).Error
	if err != nil {
		return nil, err
	}
	return Posts, nil
}

// 查询某一tag的未浏览帖子
func (qd *QADao) SearchUnviewedQAByTag(UserViewRecordID []string, Tag string) ([]models.QAs, error) {
	var posts []models.QAs
	query := qd.db.Where("post_id NOT IN ?", UserViewRecordID).Not("tag = ?","text")
	if Tag != "" {
		query = query.Where("tag = ?", Tag)
	}
	err := query.Find(&posts).Error
	return posts, err
}

// 根据id查找帖子
func (qd *QADao) SearchQAByID(PostID string) (models.QAs, error) {
	var QA models.QAs
	err := qd.db.Where("post_id = ?", PostID).First(&QA).Error
	return QA, err
}

// 根据时间查询未写入的问答内容
func (qd *QADao) SearchUpdateTime(LastUpdateTime time.Time) (int64, error) {
	var FromID int
	err := qd.db.Model(&models.QAs{}).Where("created_at > ?", LastUpdateTime).Select("id").First(&FromID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, err
	}
	return int64(FromID), nil

}

// 查询用户已发布的所有问答帖子
func (qd *QADao) GetUserQA(UserID string) ([]models.QAs, error) {
	var QAS []models.QAs
	err := qd.db.Order("created_at DESC").Where("author_id = ?", UserID).Find(&QAS).Error
	return QAS, err
}

// 查询某个问答帖的评论总数
func (qd *QADao) GetCommentCount(PostID string) (int64, error) {
	var count1, count2 int64
	err := qd.db.Model(&models.Comment{}).Where("post_id = ?", PostID).Count(&count1).Error
	if err != nil {
		return 0, err
	}
	err = qd.db.Model(&models.Reply{}).Where("post_id = ?", PostID).Count(&count2).Error
	if err != nil {
		return count1, nil
	}
	return count1 + count2, nil
}
