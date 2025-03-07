package dao

import (
	"KnowEase/models"
	"time"

	"gorm.io/gorm"
)

type LikeDao struct {
	db *gorm.DB
}

func NewLikeDao(db *gorm.DB) *LikeDao {
	return &LikeDao{db: db}
}
func ProvideLikeDao(db *gorm.DB) LikeDaoInterface {
	return NewLikeDao(db)
}

type LikeDaoInterface interface {
	GetLikeRecord(Userid, Tag string) ([]models.UserLikeHistory, error)
	GetSaveRecord(Userid, Tag string) ([]models.UserSaveHistory, error)
	GetViewRecord(Userid, Tag string) ([]models.UserViewHistory, error)
	SyncPostCountToDB(PostID, Filed string, Count int) error
	SyncQACountToDB(PostID, Filed string, Count int) error
	SyncCommentLikeToDB(CommentID string, LikeCount int) error
	SyncReplyLikeToDB(ReplyID string, LikeCount int) error
	SyncUserCountToDB(UserID, Filed string, Count int) error
	SyncUserLikeHistoryToDB(Record *models.UserLikeHistory) error
	DeleteUserLikeHistory(PostID string) error
	DeleteUserViewHistory(PostID string) error
	SyncUserSaveHistoryToDB(Record *models.UserSaveHistory) error
	SyncUserViewHistoryToDB(Record *models.UserViewHistory) error
	DeleteUserSaveHistory(PostID string) error
	SyncMessageToDB(Body *models.Message) error
	SearchUnreadMessage(UserID, Tag string) ([]models.Message, error)
	UpdateMessageStatus(UserID, Tag string) error
	SyncFollowMessageToDB(FollowMessage *models.FollowMessage) error
	DeleteFollowMessage(UserID string) error
	SearchUserFollow(UserID, SearchField, SelectField string) ([]string, error)
	SyncFollowCount(UserID string, FolloweeCount, FollowerCount int) error
}

// 在数据库中查询最近一个月的历史点赞记录
func (ld *LikeDao) GetLikeRecord(Userid, Tag string) ([]models.UserLikeHistory, error) {
	oneMonthAgo := time.Now().AddDate(0, 0, -31) // 获取 31天前的时间
	var LikeRecords []models.UserLikeHistory
	if err := ld.db.Order("created_at DESC").Where("user_id = ? AND created_at > ? AND type = ?", Userid, oneMonthAgo, Tag).Find(&LikeRecords).Error; err != nil {
		return nil, err
	}
	return LikeRecords, nil
}

// 在数据库中查询历史收藏记录
func (ld *LikeDao) GetSaveRecord(Userid, Tag string) ([]models.UserSaveHistory, error) {
	var SaveRecords []models.UserSaveHistory
	if err := ld.db.Order("created_at DESC").Where("user_id = ? AND type = ? ", Userid, Tag).Find(&SaveRecords).Error; err != nil {
		return nil, err
	}
	return SaveRecords, nil
}

// 在数据库中查询最近一个月的历史浏览记录
func (ld *LikeDao) GetViewRecord(Userid, Tag string) ([]models.UserViewHistory, error) {
	oneMonthAgo := time.Now().AddDate(0, 0, -31) // 获取 31天前的时间
	var ViewRecords []models.UserViewHistory
	if err := ld.db.Order("create_at DESC").Where("user_id = ? AND create_at > ? AND type = ?", Userid, oneMonthAgo, Tag).Find(&ViewRecords).Error; err != nil {
		return nil, err
	}
	return ViewRecords, nil
}

// 将帖子数据写入数据库
func (ld *LikeDao) SyncPostCountToDB(PostID, Filed string, Count int) error {
	err := ld.db.Model(&models.PostMessage{}).Where("post_id = ?", PostID).Update(Filed, Count).Error
	return err
}

// 将问答数据写入数据库
func (ld *LikeDao) SyncQACountToDB(PostID, Filed string, Count int) error {
	err := ld.db.Model(&models.QAs{}).Where("post_id = ?", PostID).Update(Filed, Count).Error
	return err
}

// 将评论点赞数据写入数据库
func (ld *LikeDao) SyncCommentLikeToDB(CommentID string, LikeCount int) error {
	err := ld.db.Model(&models.Comment{}).Where("comment_id = ?", CommentID).Update("like_count", LikeCount).Error
	return err
}

// 将回复点赞数据写入数据库
func (ld *LikeDao) SyncReplyLikeToDB(ReplyID string, LikeCount int) error {
	err := ld.db.Model(&models.Reply{}).Where("reply_id = ?", ReplyID).Update("like_count", LikeCount).Error
	return err
}

// 将用户数据写入数据库
func (ld *LikeDao) SyncUserCountToDB(UserID, Filed string, Count int) error {
	err := ld.db.Model(&models.User{}).Where("id = ?", UserID).Update(Filed, Count).Error
	return err
}

// 将用户点赞历史记录写入数据库
func (ld *LikeDao) SyncUserLikeHistoryToDB(Record *models.UserLikeHistory) error {
	Record.CreatedAt = time.Now()
	r := ld.db.Create(Record)
	return r.Error
}

// 将用户点赞历史记录从数据库中删除
func (ld *LikeDao) DeleteUserLikeHistory(PostID string) error {
	err := ld.db.Delete(&models.UserLikeHistory{}, "post_id = ?", PostID).Error
	return err
}

// 将用户浏览历史记录从数据库中删除
func (ld *LikeDao) DeleteUserViewHistory(PostID string) error {
	err := ld.db.Delete(&models.UserViewHistory{}, "post_id = ?", PostID).Error
	return err
}

// 将用户收藏历史记录写入数据库
func (ld *LikeDao) SyncUserSaveHistoryToDB(Record *models.UserSaveHistory) error {
	Record.CreatedAt = time.Now()
	r := ld.db.Create(Record)
	return r.Error
}

// 将用户浏览历史记录写入数据库
func (ld *LikeDao) SyncUserViewHistoryToDB(Record *models.UserViewHistory) error {
	Record.CreateAt = time.Now()
	err := ld.db.Model(&models.UserViewHistory{}).Where("user_id = ? AND post_id = ?", Record.UserID, Record.PostID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r := ld.db.Create(Record)
			return r.Error
		}
		return err
	}
	err = ld.db.Model(&models.UserViewHistory{}).Where("user_id = ? AND post_id = ?", Record.UserID, Record.PostID).Update("create_at", Record.CreateAt).Error
	return err
}

// 将用户收藏历史记录从数据库中删除
func (ld *LikeDao) DeleteUserSaveHistory(PostID string) error {
	err := ld.db.Delete(&models.UserSaveHistory{}, "post_id = ?", PostID).Error
	return err
}

// 根据帖子id查找帖子信息
func (ld *LikeDao) SearchPostByID(PostID string) (*models.PostMessage, error) {
	var Post models.PostMessage
	re := ld.db.Where("post_id=?", PostID).Find(&Post)
	if re.Error != nil {
		if re.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, re.Error
	}
	return &Post, nil
}

// 根据评论id查找评论
func (ld *LikeDao) SearchCommentByID(CommentID string) (*models.Comment, error) {
	var Comment models.Comment
	re := ld.db.Where("comment_id = ?", CommentID).Find(&Comment)
	if re.Error != nil {
		if re.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, re.Error
	}
	return &Comment, nil
}

// 将消息内容写入数据库
func (ld *LikeDao) SyncMessageToDB(Body *models.Message) error {
	Body.CreateAt = time.Now()
	r := ld.db.Create(Body)
	return r.Error
}

// 查询用户所有未读消息
func (ld *LikeDao) SearchUnreadMessage(UserID, Tag string) ([]models.Message, error) {
	var Messages []models.Message
	err := ld.db.Model(&models.Message{}).Where("user_id = ? AND status = ? AND tag = ?", UserID, "unread", Tag).Find(&Messages).Error
	return Messages, err
}

// 更新消息状态
func (ld *LikeDao) UpdateMessageStatus(UserID, Tag string) error {
	err := ld.db.Model(&models.Message{}).Where("user_id = ? AND tag = ?", UserID, Tag).Update("status", "read").Error
	return err
}

// 将关注信息写入数据库
func (ld *LikeDao) SyncFollowMessageToDB(FollowMessage *models.FollowMessage) error {
	return ld.db.Create(FollowMessage).Error
}

// 将关注信息从数据库里删除
func (ld *LikeDao) DeleteFollowMessage(UserID string) error {
	return ld.db.Delete(&models.FollowMessage{}, "follower_id = ?", UserID).Error
}

// 获取用户关注列表
func (ld *LikeDao) SearchUserFollow(UserID, SearchField, SelectField string) ([]string, error) {
	var FolloweeIDs []string
	err := ld.db.Model(&models.FollowMessage{}).Where(SearchField, UserID).Select(SelectField).Find(&FolloweeIDs).Error
	return FolloweeIDs, err
}

// 将关注数据写入数据库
func (ld *LikeDao) SyncFollowCount(UserID string, FolloweeCount, FollowerCount int) error {
	err := ld.db.Model(&models.User{}).Where("id", UserID).Updates(map[string]interface{}{
		"followee_count": FolloweeCount,
		"follower_count": FollowerCount,
	}).Error
	return err
}
