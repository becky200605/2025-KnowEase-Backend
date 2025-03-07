package controllers

import (
	"KnowEase/models"
	"KnowEase/services"
	"KnowEase/utils"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SearchControllers struct {
	SearchService *services.SearchService
}

func NewSearchController(SearchService *services.SearchService) *SearchControllers {
	return &SearchControllers{SearchService: SearchService}
}

// @Summary 搜索
// @Description 搜索查询
// @Tags 搜索
// @Accept json
// @Produce json
// @Param SearchMessage body models.SearchRecord true "用户输入的问题及id"
// @Success 200 {object} map[string]interface{} "获取成功以及搜索的信息"
// @Failure 400 {object} models.Response "输入无效，请重试"
// @Failure 500 {object} models.Response "查询失败"
// @Router /api/search [post]
func (sc *SearchControllers) Search(c *gin.Context) {
	var Record models.SearchRecord
	if err := c.BindJSON(&Record); err != nil {
		//fmt.Print(err)
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	Posts, err1 := sc.SearchService.SearchPostByMessage(Record.SearchMessage)
	QAs, err2 := sc.SearchService.SearchQAByMessage(Record.SearchMessage)
	if err1 != nil || err2 != nil {
		log.Printf("search %s error:%v %v", Record.SearchMessage, err1, err2)
		c.JSON(http.StatusInternalServerError, models.Write("获取搜索结果失败！"))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "获取搜索结果成功！",
		"posts":   Posts,
		"QAs":     QAs,
	})

}

// @Summary 获取用户搜索记录
// @Description 获取用户的最近几条搜索记录
// @Tags 搜索
// @Accept json
// @Produce json
// @Param userid path string true "用户ID"
// @Success 200 {object} map[string]interface{} "获取成功以及搜索记录"
// @Failure 400 {object} models.Response "输入无效，请重试"
// @Failure 500 {object} models.Response "查询失败"
// @Router /api/search/{userid}/getrecord [get]
func (sc *SearchControllers) GetUserSearchRecord(c *gin.Context) {
	UserID := c.Param("userid")
	if UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	Record, err := sc.SearchService.GetUserSearchRecord(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("获取搜索进度失败！"))
		return
	}
	go sc.SearchService.UpdateSearchRecord(UserID)
	c.JSON(http.StatusOK, gin.H{"message": "获取搜索记录成功！", "records": Record})
}

// @Summary 猜你想搜
// @Description 根据搜索记录获取用户可能想要搜索的内容
// @Tags 搜索
// @Accept json
// @Produce json
// @Param userid path string true "用户ID"
// @Success 200 {object} map[string]interface{} "获取成功以及搜索推荐内容"
// @Failure 400 {object} models.Response "输入无效，请重试"
// @Failure 500 {object} models.Response "查询失败"
// @Router /api/search/:userid/recommend [post]
func (sc *SearchControllers) GetSearchRecommend(c *gin.Context) {
	UserID := c.Param("userid")
	if UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	Record, err := sc.SearchService.GetUserSearchRecord(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("获取搜索记录失败！"))
		return
	}
	var Messages []string
	for _, message := range Record {
		Messages = append(Messages, message.SearchMessage)
	}
	if Messages==nil{
		Messages=append(Messages, "")
	}
	jsonArray, err := json.Marshal(Messages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("解析失败！"))
		return
	}
	response, err := utils.InitRecommend(string(jsonArray), utils.SeaarchRecommendation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "获取推荐搜索成功！", "recommendSearch": response})
}

// @Summary 删除用户搜索记录
// @Description 用户删除一个搜索记录
// @Tags 搜索
// @Accept application/json
// @Produce application/json
// @Param message body models.SearchRecord true "搜索记录和用户ID"
// @Success 201 {object} models.Response "删除用户搜索记录成功"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "删除用户搜索记录失败"
// @Router /api/search/delete [delete]
func (sc *SearchControllers) DeleteUserRecord(c *gin.Context) {
	var Message models.SearchRecord
	if err := c.BindJSON(&Message); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	err := sc.SearchService.DeleteUserSearchRecord(Message.UserID, Message.SearchMessage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("删除用户搜索记录失败！"))
		return
	}
	c.JSON(http.StatusCreated, models.Write("删除用户搜索记录成功！"))
}

// @Summary 清空用户搜索记录
// @Description 用户清空搜索记录
// @Tags 搜索
// @Accept application/json
// @Produce application/json
// @Param userid path string true "用户ID"
// @Success 201 {object} models.Response "删除用户搜索记录成功"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "删除用户搜索记录失败"
// @Router /api/search/:userid/deleteall [delete]
func (sc *SearchControllers) DeleteUserALLRecord(c *gin.Context) {
	UserID := c.Param("userid")
	if UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	err := sc.SearchService.DeleteUserAllSearchRecord(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("清空用户搜索记录失败！"))
		return
	}
	c.JSON(http.StatusCreated, models.Write("清空用户搜索记录成功！"))
}

// @Summary 获取热搜榜
// @Description 获取热搜榜
// @Tags 搜索
// @Accept application/json
// @Produce application/json
// @Success 200 {object} map[string]interface{} "获取热搜榜成功并返回具体信息"
// @Failure 500 {object} models.Response "获取热搜榜失败"
// @Router /api/search/gethotrank [post]
func (sc *SearchControllers) GetHotRank(c *gin.Context) {
	Response, err := sc.SearchService.GetHotRank()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "获取热搜榜成功！", "hot": Response})
}

// 搜索记录存入库
// @Summary 搜索记录储存
// @Description 搜索记录储存入库
// @Tags 搜索
// @Accept application/json
// @Produce application/json
// @Param record body models.SearchRecord true "用户搜索记录"
// @Success 200 {object} map[string]interface{} "储存成功并返回具体信息"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "储存失败"
// @Router /api/search/syncRecord [post]
func (sc *SearchControllers) SyncSearchRecord(c *gin.Context) {
	var Record models.SearchRecord
	if err := c.BindJSON(&Record); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	err := sc.SearchService.SyncSearchRecord(Record)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("搜索记录储存失败！"))
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "搜索记录储存成功！", "record": Record})
}
