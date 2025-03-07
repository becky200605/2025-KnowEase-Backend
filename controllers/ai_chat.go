package controllers

import (
	"KnowEase/models"
	"KnowEase/services"
	"fmt"
	"net/http"

	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// 房间号为key，client数组为value
var clients = make(map[string][]chan string)

var clientsMutex sync.RWMutex

type AIChatControllers struct {
	AIChatService *services.AIChatService
}

func NewAIChatControllers(AIChatService *services.AIChatService) *AIChatControllers {
	return &AIChatControllers{AIChatService: AIChatService}
}

// 生成回答
// @Summary 生成用户回答
// @Description 用户发送ai请求，调用api生成回答
// @Tags AI
// @Accept json
// @Produce json
// @Param UserMessage body models.AIChatMessage true "用户输入的问题"
// @Success 200 {object} map[string]interface{} "生成回答成功"
// @Failure 400 {object} models.Response "输入无效，请重试"
// @Failure 500 {object} models.Response "服务器内部错误"
// @Router /api/aichat/chat/{userid}/{chatid}/ [post]
func (acc *AIChatControllers) UserInitQuestion(c *gin.Context) {
	var UserMessage models.AIChatMessage
	if err := c.BindJSON(&UserMessage); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	ChatID := c.Param("chatid")
	UserID:=c.Param("userid")
	roomID:=UserID+ChatID
	data, err := acc.AIChatService.GetResponse(UserMessage.Request)
	if err!=nil{
		broadcast(roomID, "服务器繁忙，我们换个话题聊吧")
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Data sent to clients"})
		return
	}
	UserMessage.Response = data
	go acc.AIChatService.SyncAIChatMessage(&UserMessage)
	println("Data received from client :", data)
	broadcast(roomID, data)
	c.JSON(http.StatusOK, gin.H{"message": "Data sent to clients"})
}

// 广播房间内的所有用户
func broadcast(roomID string, data string) {
	for _, client := range clients[roomID] {
        //将data里的内容一一写入
        for _, char := range data {
            client <- string(char)
        }
    }
}

// 流式响应连接
// @Summary 流式响应连接
// @Description 发送请求，建立sse连接 
// @Tags 问答
// @Accept json
// @Produce json
// @Param chatid path string true "聊天id"
// @Success 200 {object} map[string]interface{} "历史聊天记录"
// @Router /api/{userid}/{chatid}/sseconnect [post]
func (acc *AIChatControllers) SSEConnect(c *gin.Context) {
	ChatID := c.Param("chatid")
	UserID:=c.Param("userid")
	roomID:=ChatID+UserID
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	println("Client connected")

	eventChan := make(chan string)

	clientsMutex.Lock()
	if clients[roomID] == nil {
		clients[roomID] = []chan string{}
	}
	clients[roomID] = append(clients[roomID], eventChan)
	clientsMutex.Unlock()

	defer func() {
		clientsMutex.Lock()
		for idx, v := range clients[roomID] {
			if v == eventChan {
				clients[roomID] = append(clients[roomID][:idx], clients[roomID][idx+1:]...)
				break
			}
		}
		clientsMutex.Unlock()
		close(eventChan)
	}()

	// 心跳机制
	go func() {
		for {
			time.Sleep(15 * time.Second)
			if _, err := fmt.Fprintf(c.Writer, "event: ping\ndata: {}\n\n"); err != nil {
				return
			}
		}
	}()

	// 监听客户端断开
	notify := c.Writer.CloseNotify()
	go func() {
		<-notify
		fmt.Println("Client disconnected")
	}()

	// 推送消息
	for {
		select {
		case data := <-eventChan:
			if _, err := fmt.Fprintf(c.Writer, "data: %s\n\n", data); err != nil {
				fmt.Println("Failed to send data:", err)
				return
			}
			println("Sending data to client", data)
			c.Writer.Flush()
			//超时了
		case <-time.After(30 * time.Second):
			if _, err := fmt.Fprintf(c.Writer, "event: ping\ndata: {}\n\n"); err != nil {
				return
			}
			//断开连接
		case <-c.Request.Context().Done():
			fmt.Println("Client disconnected")
			return
		}
	}
}

//获取历史聊天
// @Summary 获取历史聊天记录
// @Description 获取用户某一个聊天历史里的历史聊天记录
// @Tags AI
// @Accept json
// @Produce json
// @Param UserID path string true "用户ID"
// @Param ChatID path string true "聊天ID"
// @Success 200 {object} map[string]interface{} "历史聊天记录"
// @Failure 400 {object} models.Response "输入无效，请重试"
// @Failure 500 {object} models.Response "搜索记录失败"
// @Router /api/aichat/{userid}/search/{chatid}/gethistory [get]
func(acc *AIChatControllers)GetHistoryChatByID(c *gin.Context){
	UserID:=c.Param("userid")
	ChatID:=c.Param("chatid")
	if UserID==""||ChatID==""{
		c.JSON(http.StatusBadRequest,models.Write("输入无效，请重试！"))
		return
	}
	Record,err:=acc.AIChatService.SearchUserChatMessage(ChatID,UserID)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,models.Write("搜索记录失败！"))
		return
	}
	c.JSON(http.StatusOK,gin.H{"message":"获取历史聊天记录成功！","chathistory":Record})
}

//获取聊天记录列表
// @Summary 获取历史聊天记录列表
// @Description 获取用户创建的所有聊天记录列表
// @Tags AI
// @Accept json
// @Produce json
// @Param UserID path string true "用户ID"
// @Success 200 {object} map[string]interface{} "历史聊天记录列表"
// @Failure 400 {object} models.Response "输入无效，请重试"
// @Failure 500 {object} models.Response "搜索记录失败"
// @Router /api/aichat/{userid}/getlist [get]
func(acc *AIChatControllers)GetUserAIChatList(c *gin.Context){
	UserID:=c.Param("userid")
	if UserID==""{
		c.JSON(http.StatusBadRequest,models.Write("输入无效，请重试！"))
		return
	}
	ChatList,err:=acc.AIChatService.GetUserChatHistory(UserID)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,models.Write("获取聊天记录列表失败！"))
		return
	}
	c.JSON(http.StatusOK,gin.H{"message":"获取聊天记录列表成功！","chatList":ChatList})
}

//创建新的聊天
// @Summary 创建新的聊天
// @Description 用户创建新的聊天记录
// @Tags AI
// @Accept json
// @Produce json
// @Param AIChatHistory body models.AIChatHistory true "用户聊天信息"
// @Success 201 {object} map[string]interface{} "聊天信息"
// @Failure 400 {object} models.Response "输入无效，请重试"
// @Failure 500 {object} models.Response "创建聊天失败"
// @Router /api/aichat/create [post]
func(acc *AIChatControllers)CreateNewChat(c *gin.Context){
	var AIChatHistory models.AIChatHistory
	if err:=c.BindJSON(&AIChatHistory);err!=nil{
		c.JSON(http.StatusBadRequest,models.Write("输入无效，请重试！"))
		return
	}
	err:=acc.AIChatService.CreateHistory(&AIChatHistory)
	if err!=nil{
		c.JSON(http.StatusInternalServerError,models.Write("创建新聊天失败！"))
		return
	}
	c.JSON(http.StatusCreated,gin.H{"message":"创建新聊天成功！","chathistory":AIChatHistory})

}
