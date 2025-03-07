package services

import (
	rag "KnowEase/RAG/rag_go/client"
	"KnowEase/dao"
	"KnowEase/models"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type AIChatService struct {
	AIChatDao  dao.AIChatInterface
	RAGService *rag.RagService
}

func NewAIChatService(AIChatDao dao.AIChatInterface, RagService *rag.RagService) *AIChatService {
	return &AIChatService{AIChatDao: AIChatDao, RAGService: RagService}
}

// 聊天记录写入数据库
func (acs *AIChatService) SyncAIChatMessage(Body *models.AIChatMessage) error {
	Body.CreatedAt = time.Now()
	return acs.AIChatDao.SyncChatBodyToDB(Body)
}

// 获取回答
func (acs *AIChatService) GetResponse(Message string) (string, error) {
	Answer, err := acs.RAGService.InitQuestion(Message)
	fmt.Printf("%s\n",Answer)
	trimmedStr := strings.Trim(Answer, "`json")
	fmt.Println(trimmedStr)

	if err != nil {
		return "", err
	}
	var response models.Responses
	err = json.Unmarshal([]byte(Answer), &response)
	if err != nil {
		return "", fmt.Errorf("JSON parsing error:%v", err)
	}
	//var Chat models.AIChatMessage
	return response.Answer, nil
}

// 查询用户聊天记录
func (acs *AIChatService) SearchUserChatMessage(ChatID, UserID string) ([]models.AIChatMessage, error) {
	AIChatHistory, err := acs.AIChatDao.SearchChatsByID(ChatID, UserID)
	if err != nil {
		return nil, err
	}
	if AIChatHistory == nil {
		return nil, fmt.Errorf("empty chat history")
	}
	return AIChatHistory, nil
}

// 生成聊天历史记录
func (acs *AIChatService) CreateHistory(History *models.AIChatHistory) error {
	History.CreatedAt = time.Now()
	return acs.AIChatDao.SyncChatHistory(History)
}

// 查询用户所有记录
func (acs *AIChatService) GetUserChatHistory(UserID string) ([]models.AIChatHistory, error) {
	return acs.AIChatDao.GetUserChatHistory(UserID)
}
