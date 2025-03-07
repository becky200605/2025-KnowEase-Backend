package services

import (
	"KnowEase/dao"
	"KnowEase/models"
	"KnowEase/utils"
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type LikeService struct {
	ld dao.LikeDaoInterface
	pd dao.PostDaoInterface
	ud dao.UserDaoInterface
	qd dao.QADaoInterface
}

func NewLikeService(ld dao.LikeDaoInterface, pd dao.PostDaoInterface, ud dao.UserDaoInterface, qd dao.QADaoInterface) *LikeService {
	return &LikeService{ld: ld, pd: pd, ud: ud, qd: qd}
}

// 获取用户总获赞数
func (ls *LikeService) GetUserLikes(UserID string) (int, error) {
	key := fmt.Sprintf("Poster:%s:likes", UserID)
	likescount, err := utils.Client.HGetAll(utils.Ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	// 遍历累加点赞数
	var count, totalLikes int
	count = 0
	totalLikes = 0
	for _, countStr := range likescount {
		count, err = strconv.Atoi(countStr)
		if err != nil {
			log.Println("failed to count this record:%w", err)
		}
		totalLikes += count
	}
	return totalLikes, nil
}

// 查询用户历史点赞记录
func (ls *LikeService) GetLikeRecord(UserID, Tag string) ([]string, error) {
	UserHistory, err := ls.ld.GetLikeRecord(UserID, Tag)
	if err != nil {
		return nil, err
	}
	var Record []string
	for _, Message := range UserHistory {
		Record = append(Record, Message.PostID)
	}
	return Record, nil
}

// 查询用户收藏记录
func (ls *LikeService) GetSaveRecord(UserID, Tag string) ([]models.UserSaveHistory, error) {
	UserHistory, err := ls.ld.GetSaveRecord(UserID, Tag)
	if err != nil {
		return nil, err
	}
	return UserHistory, nil
}

// 将帖子浏览数写入数据库
func (ls *LikeService) SyncPostViewToDB(PostID string) {
	count, err := ls.GetCount("post", PostID, "view")
	if err != nil {
		log.Printf("update post %s viewcount error:%v", PostID, err)
	}
	err = ls.ld.SyncPostCountToDB(PostID, "view_count", count)
	if err != nil {
		log.Printf("update post %s viewcount error:%v", PostID, err)
	}
}

// 将问答浏览数写入数据库
func (ls *LikeService) SyncQAViewToDB(PostID string) {
	count, err := ls.GetCount("post", PostID, "view")
	if err != nil {
		log.Printf("update post %s viewcount error:%v", PostID, err)
	}
	err = ls.ld.SyncQACountToDB(PostID, "view_count", count)
	if err != nil {
		log.Printf("update post %s viewcount error:%v", PostID, err)
	}
}

// 将帖子点赞数写入数据库
func (ls *LikeService) SyncPostLikeToDB(PostID string) {
	count, err := ls.GetCount("post", PostID, "like")
	if err != nil {
		log.Printf("update post %s likecount error:%v", PostID, err)
	}
	err = ls.ld.SyncPostCountToDB(PostID, "like_count", count)
	if err != nil {
		log.Printf("update post %s likecount error:%v", PostID, err)
	}
}

// 将问答点赞数写入数据库
func (ls *LikeService) SyncQALikeToDB(PostID string) {
	count, err := ls.GetCount("post", PostID, "like")
	if err != nil {
		log.Printf("update post %s likecount error:%v", PostID, err)
	}
	err = ls.ld.SyncQACountToDB(PostID, "like_count", count)
	if err != nil {
		log.Printf("update post %s likecount error:%v", PostID, err)
	}
}

// 将帖子收藏数写入数据库
func (ls *LikeService) SyncPostSaveToDB(PostID string) {
	count, err := ls.GetCount("post", PostID, "save")
	if err != nil {
		log.Printf("update post %s savecount error:%v", PostID, err)
	}
	err = ls.ld.SyncPostCountToDB(PostID, "save_count", count)
	if err != nil {
		log.Printf("update post %s savecount error:%v", PostID, err)
	}
}

// 将问答收藏数写入数据库
func (ls *LikeService) SyncQASaveToDB(PostID string) {
	count, err := ls.GetCount("post", PostID, "save")
	if err != nil {
		log.Printf("update post %s savecount error:%v", PostID, err)
	}
	err = ls.ld.SyncQACountToDB(PostID, "save_count", count)
	if err != nil {
		log.Printf("update post %s savecount error:%v", PostID, err)
	}
}

// 将评论点赞数写入数据库
func (ls *LikeService) SyncCommentLikeToDB(CommentID string) {
	count, err := ls.GetCount("comment", CommentID, "like")
	if err != nil {
		log.Printf("update comment %s likecount error:%v", CommentID, err)
	}
	err = ls.ld.SyncCommentLikeToDB(CommentID, count)
	if err != nil {
		log.Printf("update comment %s likecount error:%v", CommentID, err)
	}
}

// 将帖子评论数写入数据库
func (ls *LikeService) SyncReplyLikeToDB(ReplyID string) {
	count, err := ls.GetCount("reply", ReplyID, "like")
	if err != nil {
		log.Printf("update reply %s likecount error:%v", ReplyID, err)
	}
	err = ls.ld.SyncReplyLikeToDB(ReplyID, count)
	if err != nil {
		log.Printf("update reply %s likecount error:%v", ReplyID, err)
	}
}

// 查询用户浏览记录
func (ls *LikeService) GetViewRecord(UserID, Tag string) ([]models.UserViewHistory, error) {
	UserHistory, err := ls.ld.GetViewRecord(UserID, Tag)
	if err != nil {
		return nil, err
	}
	return UserHistory, nil
}

// 创建消息推送
func (ls *LikeService) InitMessage(UserID, message, Userurl, Tag, PostID string) error {
	var Message models.Message
	Message.Status = "unread"
	Message.UserID = UserID
	Message.Message = message
	Message.PosterURL = Userurl
	Message.PostID = PostID
	Message.Tag = Tag
	return ls.ld.SyncMessageToDB(&Message)
}

// 设置定时更新数据库数据
func (ls *LikeService) StartUpdateTicker(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			ls.UpdateAllCount()
		}
	}
}

// 用户关注操作
func (ls *LikeService) Follow(UserID, FolloweeID string) error {
	key := fmt.Sprintf("User:%s:followers", FolloweeID)
	//查找关注记录以避免重复关注
	if utils.Client.SIsMember(utils.Ctx, key, UserID).Val() {
		return fmt.Errorf("the user %s is already followed user %s", UserID, FolloweeID)
	}
	//添加用户到粉丝合集
	utils.Client.SAdd(utils.Ctx, key, UserID)
	//更新用户关注数
	key = "User:followees"
	utils.Client.HIncrBy(utils.Ctx, key, UserID, 1)
	Record := models.FollowMessage{
		FollowerID: UserID,
		FolloweeID: FolloweeID,
	}
	go ls.ld.SyncFollowMessageToDB(&Record)
	key = "User:followers"
	utils.Client.HIncrBy(utils.Ctx, key, FolloweeID, 1)
	return nil
}

// 用户取消关注操作
func (ls *LikeService) CancelFollow(UserID, FolloweeID string) error {
	key := fmt.Sprintf("User:%s:followers", FolloweeID)
	//查找关注记录以避免重复关注
	if !utils.Client.SIsMember(utils.Ctx, key, UserID).Val() {
		return fmt.Errorf("the user %s is not followed user %s", UserID, FolloweeID)
	}
	//用户移除粉丝合集
	utils.Client.SRem(utils.Ctx, key, UserID)
	//更新用户关注数
	key = "User:followees"
	utils.Client.HIncrBy(utils.Ctx, key, UserID, -1)
	go ls.ld.DeleteFollowMessage(UserID)
	key = "User:followers"
	utils.Client.HIncrBy(utils.Ctx, key, FolloweeID, -1)
	return nil
}

// 获取用户关注列表
func (ls *LikeService) GetUserFoloweeList(UserID string) ([]models.Usermessage, error) {
	var FolloweeMessages []models.Usermessage
	UserIDs, err := ls.ld.SearchUserFollow(UserID, "follower_id", "followee_id")
	if err != nil {
		return nil, err
	}
	for _, FolloweeID := range UserIDs {
		Followee, err := ls.ud.SearchUserByID(FolloweeID)
		if err != nil {
			fmt.Printf("get followee%s error:%v", FolloweeID, err)
		} else {
			FolloweeMessages = append(FolloweeMessages, models.Usermessage{
				Username: Followee.Username,
				UserID:   Followee.ID,
				ImageURL: Followee.ImageURL,
			})
		}
	}
	return FolloweeMessages, nil
}

// 获取用户粉丝列表
func (ls *LikeService) GetUserFolowerList(UserID string) ([]models.Usermessage, error) {
	var FollowerMessages []models.Usermessage
	UserIDs, err := ls.ld.SearchUserFollow(UserID, "followee_id", "follower_id")
	if err != nil {
		return nil, err
	}
	for _, FollowerID := range UserIDs {
		Follower, err := ls.ud.SearchUserByID(FollowerID)
		if err != nil {
			fmt.Printf("get follower%s error:%v", FollowerID, err)
		} else {
			FollowerMessages = append(FollowerMessages, models.Usermessage{
				Username: Follower.Username,
				UserID:   Follower.ID,
				ImageURL: Follower.ImageURL,
			})
		}
	}
	return FollowerMessages, nil
}

// 关注数据同步
func (ls *LikeService) SyncFollowCount(UserID string) {
	followercount, err := ls.GetCount("follower", UserID, "")
	if err != nil {
		fmt.Print(err.Error())
	}
	followeecount, err := ls.GetCount("followee", UserID, "")
	if err != nil {
		fmt.Print(err.Error())
	}
	err = ls.ld.SyncFollowCount(UserID, followeecount, followercount)
	if err != nil {
		fmt.Print("update followee count error%w", err)
	}
}

// 获赞数据同步
func (ls *LikeService) SyncLikeCount(UserID string) {
	totalLikes, err := ls.GetUserLikes(UserID)
	if err != nil {
		fmt.Print(err.Error())
	}
	err = ls.ld.SyncUserCountToDB(UserID, "like_count", totalLikes)
	if err != nil {
		fmt.Print("update user like count error%w", err)
	}
}

func (ls *LikeService) UpdateAllCount() {
	Posts, _ := ls.pd.SearchAllPost()
	QAs, _ := ls.qd.SearchAllQA()
	for _, Post := range Posts {
		ls.SyncPostLikeToDB(Post)
		ls.SyncPostSaveToDB(Post)
		ls.SyncPostViewToDB(Post)

	}
	for _, QA := range QAs {
		ls.SyncQALikeToDB(QA.PostID)
		ls.SyncQASaveToDB(QA.PostID)
		ls.SyncQAViewToDB(QA.PostID)

	}
	Comments, _ := ls.pd.SearchAllComment("")
	for _, Comment := range Comments {
		ls.SyncCommentLikeToDB(Comment.CommentID)
	}
	Replys, _ := ls.pd.SearchALLReply("")
	for _, Reply := range Replys {
		ls.SyncReplyLikeToDB(Reply.ReplyID)
	}
	Users, _ := ls.ud.SearchAllUser()
	for _, User := range Users {
		ls.SyncFollowCount(User)
		ls.SyncLikeCount(User)
	}
}

// 获取点赞等相关状态
func (ls *LikeService) GetEntityStatus(EntityType, UserID, EntityID string) bool {
	var key string
	switch EntityType {
	case "post_like":
		key = fmt.Sprintf("Post:%s:like_users", EntityID)
	case "post_save":
		key = fmt.Sprintf("Post:%s:save_users", EntityID)
	case "comment_like":
		key = fmt.Sprintf("comment:%s:like_users", EntityID)
	case "reply_like":
		key = fmt.Sprintf("reply:%s:like_users", EntityID)
	case "follow":
		key = fmt.Sprintf("User:%s:followers", EntityID)
	default:
		return false
	}
	return utils.Client.SIsMember(utils.Ctx, key, UserID).Val()
}

// 点赞收藏相关操作
func (ls *LikeService) HandleUserAction(EntityType, Action, EntityID, UserID, Type, Status string) error {
	var key1, key2 string
	var Value int64

	switch EntityType {
	case "life":
		key1 = fmt.Sprintf("Post:%s:%s_users", EntityID, Action)
		Post, err := ls.pd.SearchPostByID(EntityID)
		if err != nil {
			return fmt.Errorf("failed to find the poster:%w", err)
		}
		key2 = fmt.Sprintf("Poster:%s:%ss", Post.PosterID, Action)
	case "qa":
		key1 = fmt.Sprintf("Post:%s:%s_users", EntityID, Action)
		Post, err := ls.qd.SearchQAByID(EntityID)
		if err != nil {
			return fmt.Errorf("failed to find the poster:%w", err)
		}
		key2 = fmt.Sprintf("Poster:%s:%ss", Post.AuthorID, Action)
	case "comment":
		key1 = fmt.Sprintf("comment:%s:%s_users", EntityID, Action)
		key2 = "Comment:likes"
	case "reply":
		key1 = fmt.Sprintf("reply:%s:%s_users", EntityID, Action)
		key2 = "Reply:likes"
	default:
		return fmt.Errorf("unknown entity type: %s", EntityType)
	}

	// 更新数据库记录
	if Status == "like" {
		Value = 1
		if utils.Client.SIsMember(utils.Ctx, key1, UserID).Val() {
			return fmt.Errorf("the user %s has already %s %s %s", UserID, Action, EntityType, EntityID)
		}
		utils.Client.SAdd(utils.Ctx, key1, UserID)
		if EntityType == "post" || EntityType == "qa" {
			switch Action {
			case "like":
				Record := models.UserLikeHistory{
					UserID: UserID,
					PostID: EntityID,
					Type:   Type,
				}
				go ls.ld.SyncUserLikeHistoryToDB(&Record)
			case "save":
				Record := models.UserSaveHistory{
					UserID: UserID,
					PostID: EntityID,
					Type:   Type,
				}
				go ls.ld.SyncUserSaveHistoryToDB(&Record)
			case "view":
				Record := models.UserViewHistory{
					UserID:   UserID,
					PostID:   EntityID,
					Type:     Type,
					CreateAt: time.Now(),
				}
				go ls.ld.SyncUserViewHistoryToDB(&Record)
			default:
				return fmt.Errorf("unknown action: %s", Action)
			}
		}
	} else if Status == "cancel" {
		Value = -1
		if !utils.Client.SIsMember(utils.Ctx, key1, UserID).Val() {
			return fmt.Errorf("the user %s has not %s %s %s", UserID, Action, EntityType, EntityID)
		}
		utils.Client.SRem(utils.Ctx, key1, UserID)

		// 更新数据库记录
		if EntityType == "post" || EntityType == "qa" {
			switch Action {
			case "like":
				go ls.ld.DeleteUserLikeHistory(EntityID)
			case "save":
				go ls.ld.DeleteUserSaveHistory(EntityID)
			default:
				return fmt.Errorf("unknown action: %s", Action)
			}
		}
	} else {
		return fmt.Errorf("unknown status %s", Status)
	}
	utils.Client.HIncrBy(utils.Ctx, key2, EntityID, Value)
	return nil
}

// 获取相关数值
func (ls *LikeService) GetCount(EntityType, EntityID, Action string) (int, error) {
	var key string
	switch EntityType {
	case "post":
		PosterID, err := ls.pd.SearchPostByID(EntityID)
		if err != nil {
			return 0, fmt.Errorf("failed to find the poster:%w", err)
		}
		key = fmt.Sprintf("Poster:%s:%ss", PosterID.PosterID, Action)
	case "comment":
		key = "Comment:likes"
	case "reply":
		key = "Reply:likes"
	case "follower":
		key = "User:followers"
	case "followee":
		key = "User:followees"
	default:
		return 0, fmt.Errorf("unknown entity type: %s", EntityType)
	}
	count, err := utils.Client.HGet(utils.Ctx, key, EntityID).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	Count, err := strconv.Atoi(count)
	return Count, err
}
