package controllers

import (
	"KnowEase/models"
	"KnowEase/services"
	"fmt"
	"log"
	"net/http"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type QAControllers struct {
	QAService    *services.QAService
	EmailService *services.EmailService
	LikeService  *services.LikeService
	UserService  *services.UserService
	PostService  *services.PostService
}

func NewQAControllers(QAService *services.QAService, EmailService *services.EmailService, UserService *services.UserService, LikeService *services.LikeService, PostService *services.PostService) *QAControllers {
	return &QAControllers{QAService: QAService, EmailService: EmailService, UserService: UserService, LikeService: LikeService, PostService: PostService}
}

// @Summary 搜索相似问答
// @Description 根据用户输入的问题，搜索其同质化程度
// @Tags 问答
// @Accept json
// @Produce json
// @Param QA body models.QAs true "用户输入的问题"
// @Success 200 {object} map[string]interface{} "未检测到类似问答"
// @Success 302 {object} map[string]interface{} "检测到有类似问答帖子"
// @Failure 400 {object} models.Response "输入无效，请重试"
// @Failure 500 {object} models.Response "服务器内部错误"
// @Router /api/{userid}/QA/similarity [post]
func (qc *QAControllers) SearchSimilarity(c *gin.Context) {
	var QA models.QAs
	if err := c.BindJSON(&QA); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试!"))
		return
	}
	fmt.Printf("question:%s",QA.Question)
	Similarity, PostIDs, err := qc.QAService.SearchSimilarity(QA.Question)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	if Similarity >= 0.6 {
		c.JSON(http.StatusFound, gin.H{"message": "检测到有类似问答帖子", "PostIDs": PostIDs, "similarity": Similarity})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "未检测到类似问答", "similarity": Similarity})
}

// @Summary 发布问答
// @Description 用户发布一个新的问答帖子
// @Tags 问答
// @Accept json
// @Produce json
// @Param userid path string true "用户ID"
// @Param QA body models.QAs true "问答内容"
// @Success 201 {object} map[string]interface{} "发布问答成功"
// @Failure 400 {object} models.Response "输入无效，请重试"
// @Failure 500 {object} models.Response "服务器内部错误"
// @Router /api/{userid}/QA/publish [post]
func (qc *QAControllers) PublishQA(c *gin.Context) {
	var QA models.QAs
	UserID := c.Param("userid")
	if UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := c.BindJSON(&QA); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试!"))
		return
	}
	QA.AuthorID = UserID
	for {
		//生成问答帖子id
		QA.PostID = qc.EmailService.RandomCode(7)
		_, err := qc.QAService.GetQAByID(QA.PostID)
		if err == gorm.ErrRecordNotFound {
			break
		}
	}
	PosterName, PosterURL, err := qc.UserService.SearchPosterMessage(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("查询发帖人信息失败！"))
		return
	}
	QA.AuthorName = PosterName
	QA.AuthorURL = PosterURL
	if err := qc.QAService.PublishQA(QA); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("发布帖子失败"))
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "发布问答成功！", "QAMessage": QA})
}

// @Summary 删除问答
// @Description 用户删除已发布的问答
// @Tags 个人主页-我的发布
// @Accept json
// @Produce json
// @Param postid path string true "帖子ID"
// @Success 201 {object} models.Response "删除失败"
// @Failure 400 {object} models.Response "输入无效，请重试"
// @Failure 500 {object} models.Response "删除失败"
// @Router /api/mypost/{postid}/deleteqa [delete]
func (qc *QAControllers) DeleteQA(c *gin.Context) {
	PostID := c.Param("postid")
	if PostID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := qc.QAService.DeleteQA(PostID); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("删除问答帖子失败！"))
		return
	}
	c.JSON(http.StatusCreated, models.Write("删除问答帖子成功！"))
}

// @Summary 批量删除问答
// @Description 用户删除已发布的问答
// @Tags 个人主页-我的发布
// @Accept json
// @Produce json
// @Param post body models.PostIDs true "帖子id们"
// @Success 201 {object} models.Response "删除失败"
// @Failure 400 {object} models.Response "输入无效，请重试"
// @Failure 500 {object} models.Response "删除失败"
// @Router /api/mypost/deleteqas [delete]
func (qc *QAControllers) DeleteQAs(c *gin.Context) {
	var Post models.PostIDs
	if err := c.BindJSON(&Post); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效,请重试！"))
		return
	}
	var PostIDS []string
	if err := json.Unmarshal([]byte(Post.PostID), &PostIDS); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("ID输入无效，请重试！"))
		return
	}
	if err := qc.QAService.DeleteQAs(PostIDS); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("批量删除帖子失败！"))
		return
	}
	c.JSON(http.StatusCreated, models.Write("批量删除帖子成功！"))
}

// @Summary 获取问答帖子详细信息
// @Description 获取问答帖子详细信息
// @Tags 问答
// @Accept json
// @Produce json
// @Param postid path string true "帖子ID"
// @Success 200 {object} map[string]interface{} "获取成功以及问答信息"
// @Failure 400 {object} models.Response "输入无效，请重试"
// @Failure 500 {object} models.Response "查询失败"
// @Router /api/QA/{postid} [get]
func (qc *QAControllers) GetQAMessage(c *gin.Context) {
	PostID := c.Param("postid")
	if PostID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	QAMessage, err := qc.QAService.GetQAComment(PostID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("查询问答详情失败！"))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":   "获取问答帖详情成功",
		"QAMessage": QAMessage,
	})
}

// @Summary 根据标签获取问答帖子
// @Description 根据标签获取问答帖子
// @Tags 问答
// @Accept json
// @Produce json
// @Param post body models.Type true "标签名以及用户ID"
// @Success 200 {object} map[string]interface{} "获取成功以及问答信息"
// @Failure 400 {object} models.Response "输入无效，请重试"
// @Failure 500 {object} models.Response "查询失败"
// @Router /api/QA/getbytag [post]
func (qc *QAControllers) GetQAPostByTag(c *gin.Context) {
	var Tag models.Type
	if err := c.BindJSON(&Tag); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	QAs, err := qc.QAService.GetQAMessage(Tag.Tag, Tag.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("查询帖子失败"))
		//c.Error(&models.Error{Code: 500, Message: "查询帖子失败！"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "获取问答帖子成功！", "QAPosts": QAs})
}

// @Summary 获取用户问答帖子
// @Description 获取该用户发布的问答帖子
// @Tags 他人主页
// @Accept json
// @Produce json
// @Param userid path string true "用户ID"
// @Success 200 {object} map[string]interface{} "获取成功以及问答信息"
// @Failure 400 {object} models.Response "输入无效，请重试"
// @Failure 500 {object} models.Response "查询失败"
// @Router /api/userpage/{userid}/getuserqa [get]
func (qc *QAControllers) GetUserQA(c *gin.Context) {
	UserID := c.Param("userid")
	if UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	QAs, err := qc.QAService.GetUserQA(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("查询用户发布的问答帖子出错！"))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "获取用户发布问答成功！",
		"QAs":     QAs,
	})
}

// @Summary 获取精华帖
// @Description 获取精华帖(如果不用标签筛选的话 tag里传一个空字符串)
// @Tags 问答
// @Accept json
// @Produce json
// @Param tag body models.Type true "用户ID以及标签信息"
// @Success 200 {object} map[string]interface{} "获取成功以及问答信息"
// @Failure 400 {object} models.Response "输入无效，请重试"
// @Failure 500 {object} models.Response "查询失败"
// @Router /api/QA/getgoodqa [post]
func (qc *QAControllers) GetGoodPosts(c *gin.Context) {
	var Tag models.Type
	if err := c.BindJSON(&Tag); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	QAs, err := qc.QAService.GetQAMessage(Tag.Tag, Tag.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("查询帖子失败"))
		//c.Error(&models.Error{Code: 500, Message: "查询帖子失败！"})
		return
	}
	QAs = qc.QAService.GetGoodPost(QAs)
	c.JSON(http.StatusOK, gin.H{"message": "获取精华帖成功！", "QAs": QAs})
}

// @Summary 获取问答帖子相关数值
// @Description 获取帖子的点赞数，收藏数，评论量
// @Tags 帖子
// @Accept application/json
// @Produce application/json
// @Param postid path string true "帖子ID"
// @Success 200 {object} map[string]interface{} "相关数值"
// @Failure 400 {object} models.Response "输入无效"
// @Router /api/{userid}/QA/{postid}/getcounts [get]
func (qc *QAControllers) GetPostCounts(c *gin.Context) {
	PostID := c.Param("postid")
	if PostID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	LikeCounts, err1 := qc.LikeService.GetCount("qa", PostID, "like")
	SaveCounts, err2 := qc.LikeService.GetCount("qa", PostID, "save")
	CommentCounts, err3 := qc.QAService.GetPostCommentCount(PostID)
	if err1 != nil || err2 != nil || err3 != nil {
		fmt.Print(err1)
		log.Printf("get post %s count error", PostID)
	}
	c.JSON(http.StatusOK, gin.H{"likecount": LikeCounts, "savecount": SaveCounts, "CommentCount": CommentCounts})
}

// 异步检测问答数据更新
func (qc *QAControllers) SyncQuestion() error {
	return qc.QAService.SyncQuestion()
}
