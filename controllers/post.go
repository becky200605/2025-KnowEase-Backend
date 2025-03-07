package controllers

import (
	"KnowEase/models"
	"KnowEase/services"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PostControllers struct {
	PostService  *services.PostService
	LikeService  *services.LikeService
	EmailService *services.EmailService
	UserService  *services.UserService
}

func NewPostControllers(PostService *services.PostService, LikeService *services.LikeService, EmailService *services.EmailService, UserService *services.UserService) *PostControllers {
	return &PostControllers{PostService: PostService, LikeService: LikeService, EmailService: EmailService, UserService: UserService}
}

// @Summary 用户发帖
// @Description 用户发帖，上传帖子正文、图片等信息。
// @Tags 生活区-发帖
// @Accept  json
// @Produce  json
// @Param userid path string true "用户ID"
// @Param post body models.PostMessage true "发帖内容"
// @Success 201 {object} map[string]interface{} "成功的响应信息以及帖子信息"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Failure 500 {object} models.Response "查询发帖人信息失败！"
// @Failure 500 {object} models.Response "帖子发布失败！"
// @Router /api/{userid}/post/publish [post]
func (pc *PostControllers) PublishPostBody(c *gin.Context) {
	var Post models.PostMessage
	UserID := c.Param("userid")
	if UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := c.BindJSON(&Post); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试!"))
		return
	}
	Post.PosterID = UserID
	for {
		//生成帖子id
		Post.PostID = pc.EmailService.RandomCode(6)
		_, err := pc.PostService.GetPostByID(Post.PostID)
		if err == gorm.ErrRecordNotFound {
			break
		}
	}
	PosterName, PosterURL, err := pc.UserService.SearchPosterMessage(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("查询发帖人信息失败！"))
		return
	}
	Post.PosterName = PosterName
	Post.PosterURL = PosterURL
	if err := pc.PostService.PublishPost(Post); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("帖子发布失败！"))
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "发帖成功！", "PostMessage": Post})
}

// @Summary 删除帖子
// @Description 用户在个人主页选择删除帖子
// @Tags 个人主页-我的发布
// @Accept  json
// @Produce  json
// @Param postid path string true "帖子ID"
// @Success 201 {object} models.Response "删帖成功！"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Failure 500 {object} models.Response "删除帖子及其相关信息失败！"
// @Router /api/{userid}/userpage/mypost/delete/{postid} [delete]
func (pc *PostControllers) DeletePost(c *gin.Context) {
	PostID := c.Param("postid")
	if PostID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := pc.PostService.DeletePost(PostID); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("删除帖子失败！"))
		return
	}
	c.JSON(http.StatusCreated, models.Write("删帖成功！"))
}

// @Summary 批量删除帖子
// @Description 用户在个人主页选择批量删除帖子
// @Tags 个人主页-我的发布
// @Accept  json
// @Produce  json
// @Param userid path string true "用户ID"
// @Success 201 {object} models.Response "删帖成功！"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "删除帖子及其相关信息失败！"
// @Router /api/{userid}/userpage/mypost/deleteposts [delete]
func (pc *PostControllers) DeletePosts(c *gin.Context) {
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
	if err := pc.PostService.DeletePosts(PostIDS); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("批量删除帖子失败！"))
		return
	}
	c.JSON(http.StatusCreated, models.Write("批量删除帖子成功！"))
}

// @Summary 推荐标签页面
// @Description 获取生活区的推荐标签页面的帖子信息
// @Tags 生活区-分区
// @Accept  json
// @Produce  json
// @Param userid path string true "用户ID"
// @Success 200 {object} map[string]interface{} "成功响应信息以及帖子信息"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Failure 500 {object} models.Response "推荐帖子查询失败"
// @Router /api/{userid}/post/recommend [get]
func (pc *PostControllers) RecommendationPost(c *gin.Context) {
	UserID := c.Param("userid")
	if UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	WeightedRecord, err := pc.PostService.WeightedRecommendation(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	var PostIDs []string
	for _, Posts := range WeightedRecord {
		PostIDs = append(PostIDs, Posts.PostID)
	}
	PostMessages, Err := pc.PostService.SearchPostByID(PostIDs)
	if Err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("推荐帖子信息查询失败！"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "请求推荐帖成功！", "posts": PostMessages})
}

// @Summary 校园标签页面
// @Description 获取生活区的校园标签页面的帖子信息
// @Tags 生活区-分区
// @Accept  json
// @Produce  json
// @Param userid path string true "用户ID"
// @Success 200 {object} map[string]interface{} "成功响应信息以及帖子信息"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Failure 500 {object} models.Response "校园帖子查询失败"
// @Router /api/{userid}/post/campus [get]
func (pc *PostControllers) CampusPost(c *gin.Context) {
	UserID := c.Param("userid")
	if UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	Posts, err := pc.PostService.SearchUnviewedPostsByTag(UserID, "校园")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "请求校园标签帖成功！", "posts": Posts})
}

// @Summary 生活标签页面
// @Description 获取生活区的生活标签页面的帖子信息
// @Tags 生活区-分区
// @Accept  json
// @Produce  json
// @Param userid path string true "用户ID"
// @Success 200 {object} map[string]interface{} "成功响应信息以及帖子信息"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Failure 500 {object} models.Response "生活帖子查询失败"
// @Router /api/{userid}/post/life [get]
func (pc *PostControllers) LifePost(c *gin.Context) {
	UserID := c.Param("userid")
	if UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	Posts, err := pc.PostService.SearchUnviewedPostsByTag(UserID, "生活")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "请求生活标签帖成功！", "posts": Posts})
}

// @Summary 美食标签页面
// @Description 获取生活区的美食标签页面的帖子信息
// @Tags 生活区-分区
// @Accept  json
// @Produce  json
// @Param userid path string true "用户ID"
// @Success 200 {object} map[string]interface{} "成功响应信息以及帖子信息"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Failure 500 {object} models.Response "美食帖子查询失败"
// @Router /api/{userid}/post/food [get]
func (pc *PostControllers) FoodPost(c *gin.Context) {
	UserID := c.Param("userid")
	if UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	Posts, err := pc.PostService.SearchUnviewedPostsByTag(UserID, "美食")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "请求美食标签帖成功！", "posts": Posts})
}

// @Summary 绘画标签页面
// @Description 获取生活区的绘画标签页面的帖子信息
// @Tags 生活区-分区
// @Accept  json
// @Produce  json
// @Param userid path string true "用户ID"
// @Success 200 {object} map[string]interface{} "成功响应信息以及帖子信息"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Failure 500 {object} models.Response "绘画帖子查询失败"
// @Router /api/{userid}/post/paint [get]
func (pc *PostControllers) PaintPost(c *gin.Context) {
	UserID := c.Param("userid")
	if UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	Posts, err := pc.PostService.SearchUnviewedPostsByTag(UserID, "绘画")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "请求绘画标签帖成功！", "posts": Posts})
}

// @Summary 发布评论
// @Description 用户在帖子下面发布评论
// @Tags 生活区-评论
// @Accept  json
// @Produce  json
// @Param userid path string true "用户ID"
// @Param postid path string true "帖子ID"
// @Success 201 {object} map[string]interface{} "成功响应信息以及评论信息"
// @Failure 500 {object} models.Response "发布失败"
// @Failure 400 {object} models.Response "输入无效"
// @Router /api/{userid}/post/{postid}/publishcomment [post]
func (pc *PostControllers) PublishComment(c *gin.Context) {
	var Post models.Comment
	PostID := c.Param("postid")
	UserID := c.Param("userid")
	if UserID == "" || PostID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := c.BindJSON(&Post); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试!"))
		return
	}
	Post.PostID = PostID
	Post.CommenterID = UserID
	for {
		//生成评论id
		Post.CommentID = pc.EmailService.RandomCode(6)
		_, err := pc.PostService.SearchCommentByID(Post.CommentID)
		_, err1 := pc.PostService.SearchReplyByID(Post.CommentID)
		if err == gorm.ErrRecordNotFound && err1 == gorm.ErrRecordNotFound {
			break
		}
	}
	PosterName, PosterURL, err := pc.UserService.SearchPosterMessage(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("查询回复人信息失败！"))
		return
	}
	Post.CommenterName = PosterName
	Post.CommenterURL = PosterURL
	if err := pc.PostService.PublishComment(Post); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("评论发布失败！"))
		return
	}
	PostMessage, _ := pc.PostService.GetPostByID(PostID)
	User, _ := pc.UserService.SearchUserByID(UserID)
	message := fmt.Sprintf("用户%s评论了帖子:\n%s", User.Username, PostMessage.Title)
	if err := pc.LikeService.InitMessage(PostMessage.PosterID, message, User.ImageURL, "评论", PostID); err != nil {
		c.JSON(http.StatusMultiStatus, models.Write("评论处理成功，消息上传失败"))
	}
	c.JSON(http.StatusCreated, gin.H{"message": "发布评论成功！", "CommentMessage": Post})
}

// @Summary 删除评论
// @Description 用户在帖子下面发布评论
// @Tags 生活区-评论
// @Accept  json
// @Produce  json
// @Param commentid path string true "评论ID"
// @Success 201 {object} models.Response "删除成功"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Failure 500 {object} models.Response "查询失败"
// @Router /api/{userid}/post/{postid}/{commentid}/deletecomment [delete]
func (pc *PostControllers) DeleteComment(c *gin.Context) {
	CommentID := c.Param("commentid")
	if CommentID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := pc.PostService.DeleteComment(CommentID); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("删除评论失败！"))
		return
	}
	if err := pc.PostService.DeleteAllReply(CommentID); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("删除评论回复失败！"))
		return
	}
	c.JSON(http.StatusCreated, models.Write("删除评论成功！"))
}

// @Summary 发布回复
// @Description 用户在帖子下面发布回复
// @Tags 生活区-评论
// @Accept  json
// @Produce  json
// @Param userid path string true "用户ID"
// @Param postid path string true "帖子ID"
// @Success 207 {object} models.Response "发布成功但消息发送失败"
// @Success 201 {object} map[string]interface{} "成功响应信息以及评论信息"
// @Failure 500 {object} models.Response "查询失败"
// @Failure 400 {object} models.Response "输入无效"
// @Router /api/{userid}/post/{postid}/{commentid}/publishreply [post]
func (pc *PostControllers) PublishReply(c *gin.Context) {
	var Post models.Reply
	CommentID, UserID, PostID := c.Param("commentid"), c.Param("userid"), c.Param("postid")
	if CommentID == "" || UserID == "" || PostID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := c.BindJSON(&Post); err != nil {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试!"))
		return
	}
	Post.PostID = PostID
	Post.CommentID = CommentID
	Post.ReplyerID = UserID
	for {
		//生成回复id
		Post.ReplyID = pc.EmailService.RandomCode(6)
		_, err := pc.PostService.SearchCommentByID(Post.ReplyID)
		_, err1 := pc.PostService.SearchReplyByID(Post.ReplyID)
		if err == gorm.ErrRecordNotFound && err1 == gorm.ErrRecordNotFound {
			break
		}
	}
	PosterName, PosterURL, err := pc.UserService.SearchPosterMessage(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("查询回复人信息失败！"))
		return
	}
	Post.ReplyerName = PosterName
	Post.ReplyURL = PosterURL
	if err := pc.PostService.PublishReply(Post); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("回复发布失败！"))
		return
	}
	PostMessage, _ := pc.PostService.GetPostByID(PostID)
	Comment, _ := pc.PostService.SearchCommentByID(CommentID)
	User, _ := pc.UserService.SearchUserByID(UserID)
	message := fmt.Sprintf("用户%s回复了你在帖子%s的评论:\n%s", User.Username, PostMessage.Title, Comment.Body)
	if err := pc.LikeService.InitMessage(Comment.CommenterID, message, User.ImageURL, "评论", PostID); err != nil {
		c.JSON(http.StatusMultiStatus, models.Write("回复处理成功，消息上次失败"))
	}
	c.JSON(http.StatusCreated, gin.H{"message": "发布回复成功！", "ReplyMessage": Post})
}

// @Summary 删除回复
// @Description 用户删除自己的回复
// @Tags 生活区-评论
// @Accept  json
// @Produce  json
// @Param replyid path string true "回复ID"
// @Success 201 {object} models.Response "删除成功"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Failure 500 {object} models.Response "删除失败"
// @Router /api/{userid}/post/{postid}/{commentid}/{replyid} [delete]
func (pc *PostControllers) DeleteReply(c *gin.Context) {
	ReplyID := c.Param("replyid")
	if ReplyID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := pc.PostService.DeleteReply(ReplyID); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("删除评论失败！"))
		return
	}
	//删除回复的回复
	if err := pc.PostService.DeleteAllReply(ReplyID); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("删除评论回复失败！"))
		return
	}
	c.JSON(http.StatusCreated, models.Write("删除评论成功！"))
}

// @Summary 获取帖子详细信息
// @Description 用户点开帖子详情页，获取帖子详情信息,帖子浏览量加一
// @Tags 生活区-分区
// @Accept  json
// @Produce  json
// @Param postid path string true "帖子ID"
// @Success 200 {object} map[string]interface{} "成功响应信息以及帖子信息"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Failure 500 {object} models.Response "查询失败"
// @Router /api/{userid}/post/{postid} [get]
func (pc *PostControllers) GetPostMessage(c *gin.Context) {
	PostID := c.Param("postid")
	UserID := c.Param("userid")
	if PostID == "" || UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	pc.LikeService.HandleUserAction("post", "view", PostID, UserID, "生活", "like")
	PostMessage, err := pc.PostService.GetAllComment(PostID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("查询帖子具体信息出错！"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "查询帖子信息成功！", "postMessage": *PostMessage})
}

// @Summary 获取未读消息通知-点赞
// @Description 用户点开消息通知，获取未读的消息(帖子，评论被点赞），并把消息状态更新为已读
// @Tags 消息-点赞
// @Accept  json
// @Produce  json
// @Param userid path string true "用户ID"
// @Success 200 {object} map[string]interface{} "成功响应信息以及消息信息"
// @Success 207 {object} map[string]interface{} "状态更新错误以及消息信息"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Failure 500 {object} models.Response "查询失败"
// @Router /api/{userid}/likemessage [get]
func (pc *PostControllers) GetUserUnreadLikeMessage(c *gin.Context) {
	UserID := c.Param("userid")
	if UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	Message, err := pc.PostService.SearchAllUnreadMessage(UserID, "点赞")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("查询未读消息出错！"))
		return
	}
	//将所有消息更新为已读
	go pc.PostService.UpdateMessageStatus(UserID, "点赞")
	c.JSON(http.StatusOK, gin.H{"message": "获取并更新消息成功！", "messages": Message})
}

// @Summary 获取未读消息通知-评论
// @Description 用户点开消息通知，获取未读的消息(帖子，评论被评论），并把消息状态更新为已读
// @Tags 消息-评论
// @Accept  json
// @Produce  json
// @Param userid path string true "用户ID"
// @Success 200 {object} map[string]interface{} "成功响应信息以及消息信息"
// @Success 207 {object} map[string]interface{} "状态更新错误以及消息信息"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Failure 500 {object} models.Response "查询失败"
// @Router /api/{userid}/commentmessage [get]
func (pc *PostControllers) GetUserUnreadCommentMessage(c *gin.Context) {
	UserID := c.Param("userid")
	if UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	Message, err := pc.PostService.SearchAllUnreadMessage(UserID, "评论")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("查询未读消息出错！"))
		return
	}
	//将所有消息更新为已读
	go pc.PostService.UpdateMessageStatus(UserID, "评论")
	c.JSON(http.StatusOK, gin.H{"message": "获取并更新消息成功！", "messages": Message})
}

// @Summary 获取未读消息通知-关注
// @Description 用户点开消息通知，获取未读的消息（被关注），并把消息状态更新为已读
// @Tags 消息-关注
// @Accept  json
// @Produce  json
// @Param userid path string true "用户ID"
// @Success 200 {object} map[string]interface{} "成功响应信息以及消息信息"
// @Success 207 {object} map[string]interface{} "状态更新错误以及消息信息"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Failure 500 {object} models.Response "查询失败"
// @Router /api/{userid}/followmessage [get]
func (pc *PostControllers) GetUserUnreadFollowMessage(c *gin.Context) {
	UserID := c.Param("userid")
	if UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	Message, err := pc.PostService.SearchAllUnreadMessage(UserID, "关注")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("查询未读消息出错！"))
		return
	}
	//将所有消息更新为已读
	go pc.PostService.UpdateMessageStatus(UserID, "关注")
	c.JSON(http.StatusOK, gin.H{"message": "获取并更新消息成功！", "messages": Message})
}

// @Summary 获取图床token
// @Description 获取token
// @Tags 图床
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{} "成功响应信息以及token信息"
// @Failure 500 {object} models.Response "获取配置文件失败或生成token失败"
// @Router /api/getToken [get]
func (ps *PostControllers) GetToken(c *gin.Context) {
	Config, err := ps.PostService.ReadConfig("config.yaml")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("获取配置文件出错"))
		return
	}
	token, err := ps.PostService.GetToken(Config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// @Summary 获取用户生活帖子
// @Description 获取该用户发布的生活帖子
// @Tags 他人主页
// @Accept json
// @Produce json
// @Param userid path string true "用户ID"
// @Success 200 {object} map[string]interface{} "获取成功以及帖子信息"
// @Failure 400 {object} models.Response "输入无效，请重试"
// @Failure 500 {object} models.Response "查询失败"
// @Router /api/userpage/{userid}/getuserpost [get]
func (pc *PostControllers) GetUserPosts(c *gin.Context) {
	UserID := c.Param("userid")
	if UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	Posts, err := pc.PostService.GetUserPosts(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("查询用户发布的帖子出错！"))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "获取用户发布帖子成功！",
		"Posts":   Posts,
	})
}
