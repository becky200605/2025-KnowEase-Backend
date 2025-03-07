package services

import (
	rag "KnowEase/RAG/rag_go/client"
	"KnowEase/dao"
	"KnowEase/models"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"
)

var LastSyncTime time.Time
var mu sync.Mutex

type QAService struct {
	QADao      dao.QADaoInterface
	PostDao    dao.PostDaoInterface
	LikeDao    dao.LikeDaoInterface
	UserDao    dao.UserDaoInterface
	RAGService *rag.RagService
}

func NewQAService(QADao dao.QADaoInterface, LikeDao dao.LikeDaoInterface, UserDao dao.UserDaoInterface, RagService *rag.RagService, PostDao dao.PostDaoInterface) *QAService {
	return &QAService{QADao: QADao, LikeDao: LikeDao, UserDao: UserDao, RAGService: RagService, PostDao: PostDao}
}

// 发布问答
func (qs *QAService) PublishQA(QA models.QAs) error {
	return qs.QADao.SyncQABodyToDB(&QA)
}

// 删除问答
func (qs *QAService) DeleteQA(PostID string) error {
	err := qs.QADao.DeleteQABody(PostID)
	if err != nil {
		return err
	}
	go qs.DeleteQAComment(PostID)
	return nil
}

// 批量删除问答
func (qs *QAService) DeleteQAs(PostIDs []string) error {
	var Err []error
	for _, PostID := range PostIDs {
		err := qs.QADao.DeleteQABody(PostID)
		if err != nil {
			Err = append(Err, err)
		}
		go qs.DeleteQAComment(PostID)
	}
	if Err != nil {
		return fmt.Errorf("failed to delete QAs")
	}
	return nil
}

// 查找用户问题同质化程度
func (qs *QAService) SearchSimilarity(Question string) (float64, string, error) {
	Answer, err := qs.RAGService.InitQuestion(Question)
	trimmedStr := strings.Trim(Answer, "`json")
	fmt.Println(trimmedStr)

	if err != nil {
		return 0, "", err
	}
	var response models.Responses
	err = json.Unmarshal([]byte(Answer), &response)
	if err != nil {
		return 0, "", fmt.Errorf("JSON parsing error:%v", err)
	}
	//var IDs []string
	//for _, PostID := range response.PostIDs {
	//	if PostID != "" {
	//		IDs = append(IDs, PostID)
	//	}
	//}
	return response.Similarity, response.PostIDs, nil
}

// 查询帖子的所有评论
func (qs *QAService) GetQAComment(PostID string) (*models.QAs, error) {
	PostMessage, err := qs.QADao.SearchQAByID(PostID)
	if err != nil {
		return nil, fmt.Errorf("failed to find this Post")
	}
	//查询所有的评论
	Comments, err := qs.PostDao.SearchAllComment(PostMessage.PostID)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(Comments); i++ {
		Reply, _ := qs.PostDao.SearchALLReply(Comments[i].CommentID)
		Comments[i].Reply = Reply
		for j := 0; j < len(Reply); j++ {
			Replys, _ := qs.PostDao.SearchALLReply(Reply[j].ReplyID)
			Reply[j].Reply = Replys
		}
	}
	PostMessage.Comment = Comments
	return &PostMessage, nil

}

// 根据标签获取问答帖子
func (qs *QAService) GetQAMessage(Tag, UserID string) ([]models.QAs, error) {
	ViewPost, err := qs.LikeDao.GetViewRecord(UserID, "问答")
	if err != nil {
		return nil, err
	}
	var PostIDs []string
	for _, Post := range ViewPost {
		PostIDs = append(PostIDs, Post.PostID)
	}
	if PostIDs == nil {
		PostIDs = append(PostIDs, "")
	}
	QAs, err := qs.QADao.SearchUnviewedQAByTag(PostIDs, Tag)
	if err != nil {
		return nil, err
	}
	sort.Slice(QAs, func(i, j int) bool {
		return QAs[i].LikeCount < QAs[j].LikeCount
	})
	return QAs, nil
}

// 删除问答相关评论
func (qs *QAService) DeleteQAComment(PostID string) {
	Comments, err := qs.PostDao.SearchAllComment(PostID)
	if err != nil {
		fmt.Printf("failed to find post %s comment", PostID)
		return
	}
	err = qs.PostDao.DeleteAllComment(PostID)
	if err != nil {
		fmt.Printf("failed to delete post %s comment", PostID)
		return
	}
	for _, Comment := range Comments {
		err := qs.PostDao.DeleteAllReply(Comment.CommentID)
		if err != nil {
			fmt.Printf("failed to delete comment %s reply", Comment.CommentID)
		}
	}
}

// 查询用户已发布的问答
func (qs *QAService) GetUserQA(UserID string) ([]models.QAs, error) {
	return qs.QADao.GetUserQA(UserID)
}

func (qs *QAService) SyncQuestion() error {
	mu.Lock()
	defer mu.Unlock()
	log.Printf("上次更新时间为%v", LastSyncTime)
	FromID, err := qs.QADao.SearchUpdateTime(LastSyncTime)
	if err != nil {
		return fmt.Errorf("Search question Error:%w", err)
	}
	if FromID != 0 {
		err := qs.RAGService.SyncQuestion(FromID)
		if err != nil {
			return fmt.Errorf("Sync question Error:%w", err)
		}
	}
	LastSyncTime = time.Now()
	return nil
}

// 获取精华帖
func (qs *QAService) GetGoodPost(QAs []models.QAs) []models.QAs {
	var GoodPost []models.QAs
	GoodPost=qs.CheckQA(GoodPost)
	for j := 0; j < len(QAs); j++ {
		CommentCount, err := qs.QADao.GetCommentCount(QAs[j].PostID)
		if err != nil {
			log.Printf("get post%s comment count error:%v", QAs[j].PostID, err)
			continue
		}
		QAs[j].Weight = (float64(QAs[j].LikeCount)*0.7 + float64(CommentCount)*0.3) * 100
	}
	sort.Slice(QAs, func(i, j int) bool {
		return int(QAs[i].Weight) < int(QAs[j].Weight)
	})
	for i := 0; i < len(QAs)/3; i++ {
		GoodPost = append(GoodPost, QAs[i])
	}
	return GoodPost
}

// 获取评论数
func (qs *QAService) GetPostCommentCount(PostID string) (int64, error) {
	return qs.QADao.GetCommentCount(PostID)
}
func (qs *QAService) GetQAByID(PostID string) (models.QAs, error) {
	return qs.QADao.SearchQAByID(PostID)
}

// 获取用户问答的响应
func (qs *QAService) GetLLMResponse(Question string) (float64, []string, string, error) {
	Answer, err := qs.RAGService.InitQuestion(Question)
	trimmedStr := strings.Trim(Answer, "`json")
	fmt.Println(trimmedStr)

	if err != nil {
		return 0, nil, "", err
	}
	var response models.Responses
	err = json.Unmarshal([]byte(Answer), &response)
	if err != nil {
		return 0, nil, "", fmt.Errorf("JSON parsing error:%v", err)
	}
	//var IDs []string
	//for _, PostID := range response.PostIDs {
	//	if PostID != "" {
	//		IDs = append(IDs, PostID)
	//	}
	//}
	return response.Similarity, nil, response.Answer, nil
}
//过滤官方文档
func(qs *QAService)CheckQA(QAs []models.QAs)[]models.QAs{
	var QAMessage []models.QAs
	for _,QA:=range QAs{
		if QA.Answer==""{
			QAMessage=append(QAMessage,QA)
		}
	}
	return QAMessage
}