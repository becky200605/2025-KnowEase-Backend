package services

import (
	"KnowEase/dao"
	"KnowEase/models"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type PostService struct {
	PostDao dao.PostDaoInterface
	LikeDao dao.LikeDaoInterface
	UserDao dao.UserDaoInterface
}

func NewPostService(PostDao dao.PostDaoInterface, LikeDao dao.LikeDaoInterface, UserDao dao.UserDaoInterface) *PostService {
	return &PostService{PostDao: PostDao, LikeDao: LikeDao, UserDao: UserDao}
}

// 发布帖子
func (ps *PostService) PublishPost(Post models.PostMessage) error {
	return ps.PostDao.SyncPostBodyToDB(&Post)
}

// 发布评论
func (ps *PostService) PublishComment(Comment models.Comment) error {
	return ps.PostDao.SyncCommentBodyToDB(&Comment)
}

// 发布评论
func (ps *PostService) PublishReply(Reply models.Reply) error {
	return ps.PostDao.SyncReplyToDB(&Reply)
}

// 删除评论
func (ps *PostService) DeleteComment(CommentID string) error {
	return ps.PostDao.DeleteComment(CommentID)
}

// 删除回复
func (ps *PostService) DeleteReply(ReplyID string) error {
	return ps.PostDao.DeleteReply(ReplyID)
}

// 删除帖子
func (ps *PostService) DeletePost(PostID string) error {
	err := ps.PostDao.DeletePostBody(PostID)
	if err != nil {
		return err
	}
	go ps.DeletePostComment(PostID)
	return nil
}

// 批量删除帖子
func (ps *PostService) DeletePosts(PostIDs []string) error {
	var Err []error
	for _, PostID := range PostIDs {
		err := ps.PostDao.DeletePostBody(PostID)
		if err != nil {
			Err = append(Err, err)
		}
		go ps.DeletePostComment(PostID)
	}
	if Err != nil {
		return fmt.Errorf("failed to delete posts")
	}
	return nil
}

// 删除帖子相关评论
func (ps *PostService) DeletePostComment(PostID string) {
	Comments, err := ps.PostDao.SearchAllComment(PostID)
	if err != nil {
		fmt.Printf("failed to find post %s comment", PostID)
		return
	}
	err = ps.PostDao.DeleteAllComment(PostID)
	if err != nil {
		fmt.Printf("failed to delete post %s comment", PostID)
		return
	}
	for _, Comment := range Comments {
		err := ps.PostDao.DeleteAllReply(Comment.CommentID)
		if err != nil {
			fmt.Printf("failed to delete comment %s reply", Comment.CommentID)
		}
	}
}

// 查找相关帖子信息-根据id
func (ps *PostService) SearchPostByID(ID []string) ([]models.PostMessage, []error) {
	var PostMessages []models.PostMessage
	var Err []error
	for _, PostID := range ID {
		PostMessage, err := ps.PostDao.SearchPostByID(PostID)
		if err != nil {
			Err = append(Err, err)
		}
		PostMessages = append(PostMessages, PostMessage)
	}
	if PostMessages == nil {
		return nil, Err
	}
	if Err != nil {
		log.Println("something went wrong while querying recommended posts:%w", Err)
	}
	return PostMessages, nil
}

// 推荐的加权
func (ps *PostService) WeightedRecommendation(UserID string) ([]models.PostRecommendLevel, error) {
	Record, err := ps.LikeDao.GetViewRecord(UserID, "生活")
	if err != nil {
		return nil, err
	}
	var PostIDS []string
	for _, Post := range Record {
		PostIDS = append(PostIDS, Post.PostID)
	}
	if PostIDS == nil {
		PostIDS = append(PostIDS, "")
	}
	var WeightedPosts []models.PostRecommendLevel
	tagCount := make(map[string]int)
	Tag := [4]string{"校园", "生活", "美食", "绘画"}
	for i := 0; i < len(Tag); i++ {
		Count, err := ps.PostDao.SearchCountOfTag(PostIDS, Tag[i])
		if err != nil {
			fmt.Printf("search tag %s error:%v", Tag[i], err)
		}
		tagCount[Tag[i]] = Count
	}
	PostMessage, err := ps.PostDao.SearchUnviewedPostByTag(PostIDS, "")
	if err != nil {
		return nil, err
	}
	//计算所有帖子的权重
	for _, Posts := range PostMessage {
		weight := (Posts.LikeCount*3+Posts.SaveCount*4+Posts.ViewCount*3)*4/100 + tagCount[Posts.Tag]*6/10
		WeightedPosts = append(WeightedPosts, models.PostRecommendLevel{
			PostID: Posts.PostID,
			Weight: weight,
		})
	}
	//根据权重排序
	for i := 0; i < len(WeightedPosts); i++ {
		for j := 0; j < len(WeightedPosts)-i; j++ {
			if WeightedPosts[i].Weight < WeightedPosts[j].Weight {
				WeightedPosts[i], WeightedPosts[j] = WeightedPosts[j], WeightedPosts[i]
			}
		}
	}
	return WeightedPosts, nil
}

// 查询某一tag的未浏览帖子
func (ps *PostService) SearchUnviewedPostsByTag(UserID, Tag string) ([]models.PostMessage, error) {
	Record, err := ps.LikeDao.GetViewRecord(UserID, "生活")
	if err != nil {
		return nil, err
	}
	var PostIDS []string
	for _, Post := range Record {
		PostIDS = append(PostIDS, Post.PostID)
	}
	if PostIDS == nil {
		PostIDS = append(PostIDS, "")
	}
	return ps.PostDao.SearchUnviewedPostByTag(PostIDS, Tag)
}
func (ps *PostService) DeleteAllReply(CommentID string) error {
	return ps.PostDao.DeleteAllReply(CommentID)
}

// 查询帖子的所有评论
func (ps *PostService) GetAllComment(PostID string) (*models.PostMessage, error) {
	PostMessage, err := ps.PostDao.SearchPostByID(PostID)
	if err != nil {
		return nil, fmt.Errorf("failed to find this Post")
	}
	//查询所有的评论
	Comments, err := ps.PostDao.SearchAllComment(PostMessage.PostID)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(Comments); i++ {
		Reply, _ := ps.PostDao.SearchALLReply(Comments[i].CommentID)
		Comments[i].Reply = Reply
		for j := 0; j < len(Reply); j++ {
			Replys, _ := ps.PostDao.SearchALLReply(Reply[j].ReplyID)
			Reply[j].Reply = Replys
		}
	}
	PostMessage.Comment = Comments
	return &PostMessage, nil

}

func (ps *PostService) GetPostByID(PostID string) (models.PostMessage, error) {
	return ps.PostDao.SearchPostByID(PostID)
}

// 查询未读消息
func (ps *PostService) SearchAllUnreadMessage(UserID, Tag string) ([]models.Message, error) {
	Messages, err := ps.LikeDao.SearchUnreadMessage(UserID, Tag)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return Messages, nil
}

// 更新消息状态
func (ps *PostService) UpdateMessageStatus(UserID, Tag string) {
	err := ps.LikeDao.UpdateMessageStatus(UserID, Tag)
	if err != nil {
		fmt.Printf("update message error:%v", err)
	}
}

func (ps *PostService) ReadConfig(filename string) (*models.QiNiuYunConfig, error) {
	v := viper.New()
	v.SetConfigFile(filename)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	var config models.QiNiuYunConfig
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &config, nil
}

func (ps *PostService) GetToken(config *models.QiNiuYunConfig) (string, error) {
	fmt.Print(config)
	accesskey := config.AccessKey
	secretkey := config.SecretKey
	bucket := config.Bucket
	mac := credentials.NewCredentials(accesskey, secretkey)
	putPolicy, err := uptoken.NewPutPolicy(bucket, time.Now().Add(1*time.Hour))
	if err != nil {
		return "", err
	}
	//获取上传凭证
	upToken, err := uptoken.NewSigner(putPolicy, mac).GetUpToken(context.Background())
	if err != nil {
		return "", nil
	}
	return upToken, nil
}

// 根据id找评论
func (ps *PostService) SearchCommentByID(CommentID string) (models.Comment, error) {
	return ps.PostDao.SearchCommentByID(CommentID)
}

// 根据id找评论
func (ps *PostService) SearchReplyByID(ReplyID string) (models.Reply, error) {
	return ps.PostDao.SearchReplyByID(ReplyID)
}

// 查询用户已发布的帖子
func (ps *PostService) GetUserPosts(UserID string) ([]models.PostMessage, error) {
	return ps.PostDao.GetUserPosts(UserID)
}
