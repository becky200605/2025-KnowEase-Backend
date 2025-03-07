package dao

import (
	"KnowEase/models"
	"time"

	"gorm.io/gorm"
)

type PostDao struct {
	db *gorm.DB
}

func NewPostDao(db *gorm.DB) *PostDao {
	return &PostDao{db: db}
}

func ProvidePostDao(db *gorm.DB) PostDaoInterface {
	return NewPostDao(db)
}

type PostDaoInterface interface {
	SyncPostBodyToDB(Body *models.PostMessage) error
	DeletePostBody(PostID string) error
	SyncCommentBodyToDB(Body *models.Comment) error
	DeleteComment(CommentID string) error
	DeleteReply(ReplyID string) error
	SyncReplyToDB(Body *models.Reply) error
	DeleteAllComment(PostID string) error
	DeleteAllReply(CommentID string) error
	SearchAllPost() ([]string, error)
	SearchUnviewedPostByTag(UserViewRecordID []string, Tag string) ([]models.PostMessage, error)
	SearchCountOfTag(PostID []string, Tag string) (int, error)
	SearchAllComment(PostID string) ([]models.Comment, error)
	SearchALLReply(CommentID string) ([]models.Reply, error)
	SearchCommentByID(CommentID string) (models.Comment, error)
	SearchReplyByID(ReplyID string) (models.Reply, error)
	SearchPostByID(PostID string) (models.PostMessage, error)
	GetUserPosts(UserID string) ([]models.PostMessage, error)
}

// 将帖子内容写入数据库
func (pd *PostDao) SyncPostBodyToDB(Body *models.PostMessage) error {
	Body.CreateAt = time.Now()
	r := pd.db.Create(Body)
	return r.Error
}

// 将帖子内容从数据库中删除
func (pd *PostDao) DeletePostBody(PostID string) error {
	err := pd.db.Model(&models.PostMessage{}).Delete("post_id = ?", PostID).Error
	return err
}

// 将评论信息写入数据库
func (pd *PostDao) SyncCommentBodyToDB(Body *models.Comment) error {
	Body.CreatAt = time.Now()
	r := pd.db.Create(Body)
	return r.Error
}

// 将评论删除
func (pd *PostDao) DeleteComment(CommentID string) error {
	err := pd.db.Delete(&models.Comment{}, "comment_id = ?", CommentID).Error
	return err
}

// 将回复删除
func (pd *PostDao) DeleteReply(ReplyID string) error {
	err := pd.db.Delete(&models.Reply{}, "reply_id = ?", ReplyID).Error
	return err
}

// 将回复内容写入数据库
func (pd *PostDao) SyncReplyToDB(Body *models.Reply) error {
	Body.CreatAt = time.Now()
	r := pd.db.Create(Body)
	return r.Error
}

// 删除帖子所有评论
func (pd *PostDao) DeleteAllComment(PostID string) error {
	err := pd.db.Model(&models.Comment{}).Delete("post_id = ?", PostID).Error
	return err
}

// 删除评论所有回复
func (pd *PostDao) DeleteAllReply(CommentID string) error {
	err := pd.db.Delete(&models.Reply{}, "comment_id = ?", CommentID).Error
	return err
}

// 查询所有帖子
func (pd *PostDao) SearchAllPost() ([]string, error) {
	var Posts []string
	err := pd.db.Model(&models.PostMessage{}).Select("post_id").Find(&Posts).Error
	if err != nil {
		return nil, err
	}
	return Posts, nil
}

// 查询某一tag的未浏览帖子
func (pd *PostDao) SearchUnviewedPostByTag(UserViewRecordID []string, Tag string) ([]models.PostMessage, error) {
	var posts []models.PostMessage
	query := pd.db.Where("post_id NOT IN ?", UserViewRecordID)
	if Tag != "" {
		query = query.Where("tag = ?", Tag)
	}
	err := query.Find(&posts).Error
	return posts, err
}

// 查询所有包含某个tag的数据
func (pd *PostDao) SearchCountOfTag(PostID []string, Tag string) (int, error) {
	var Count int64
	err := pd.db.Model(&models.PostMessage{}).Where("post_id IN ? AND tag= ?", PostID, Tag).Count(&Count).Error

	return int(Count), err
}

// 查询帖子的所有评论
func (pd *PostDao) SearchAllComment(PostID string) ([]models.Comment, error) {
	var Comments []models.Comment
	if PostID == "" {
		err := pd.db.Find(&Comments).Error
		return Comments, err
	}
	err := pd.db.Where("post_id = ?", PostID).Find(&Comments).Error
	return Comments, err
}

// 查询评论的所有回复
func (pd *PostDao) SearchALLReply(CommentID string) ([]models.Reply, error) {
	var Replys []models.Reply
	if CommentID == "" {
		err := pd.db.Find(&Replys).Error
		return Replys, err
	}
	err := pd.db.Where("comment_id = ?", CommentID).Find(&Replys).Error
	return Replys, err
}

// 查询评论
func (pd *PostDao) SearchCommentByID(CommentID string) (models.Comment, error) {
	var Comment models.Comment
	err := pd.db.Where("comment_id = ?", CommentID).First(&Comment).Error
	return Comment, err
}

// 查询回复
func (pd *PostDao) SearchReplyByID(ReplyID string) (models.Reply, error) {
	var Reply models.Reply
	err := pd.db.Where("reply_id = ?", ReplyID).First(&Reply).Error
	return Reply, err
}

// 查询帖子
func (pd *PostDao) SearchPostByID(PostID string) (models.PostMessage, error) {
	var Post models.PostMessage
	err := pd.db.Where("post_id = ?", PostID).First(&Post).Error
	return Post, err
}

// 查询用户已发布的所有帖子
func (pd *PostDao) GetUserPosts(UserID string) ([]models.PostMessage, error) {
	var PostS []models.PostMessage
	err := pd.db.Order("create_at DESC").Where("poster_id = ?", UserID).Find(&PostS).Error
	return PostS, err
}
