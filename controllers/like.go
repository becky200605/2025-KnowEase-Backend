package controllers

import (
	"KnowEase/models"
	"KnowEase/services"
	"fmt"
	"log"
	"net/http"
	"context"

	"github.com/gin-gonic/gin"
)

type LikeControllers struct {
	LikeService *services.LikeService
	PostService *services.PostService
	UserService *services.UserService
}

func NewLikeControllers(LikeService *services.LikeService, PostService *services.PostService, UserService *services.UserService) *LikeControllers {
	return &LikeControllers{LikeService: LikeService, PostService: PostService, UserService: UserService}
}

// @Summary 用户点赞帖子
// @Description 用户给帖子点赞，点赞数加一
// @Tags 帖子-点赞
// @Accept application/json
// @Produce application/json
// @Param postid path string true "帖子ID"
// @Param userid path string true "用户ID"
// @Param type path string true "帖子类型"
// @Success 201 {object} models.Response "用户点赞成功"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "点赞记录上传失败"
// @Router /api/{userid}/{type}/{postid}/like [post]
func (lc *LikeControllers) LikePost(c *gin.Context) {
	PostID := c.Param("postid")
	UserID := c.Param("userid")
	Type:=c.Param("type")
	if PostID == "" || UserID == "" ||Type==""{
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	Post, _ := lc.PostService.GetPostByID(PostID)
	User, _ := lc.UserService.SearchUserByID(UserID)
	if err := lc.LikeService.HandleUserAction(Type, "like", PostID, UserID, Type, "like"); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	message := fmt.Sprintf("用户%s点赞了你的帖子:\n%s", User.Username, Post.Title)
	if err := lc.LikeService.InitMessage(Post.PosterID, message, User.ImageURL, "点赞", PostID); err != nil {
		c.JSON(http.StatusMultiStatus, models.Write("点赞处理成功，消息上次失败"))
	}
	c.JSON(http.StatusCreated, models.Write("用户点赞成功！"))
}

// @Summary 用户取消点赞帖子
// @Description 用户取消点赞，点赞数减一
// @Tags 帖子-点赞
// @Accept application/json
// @Produce application/json
// @Param postid path string true "帖子ID"
// @Param userid path string true "用户D"
// @Param type path string true "帖子类型"
// @Success 201 {object} models.Response "用户取消点赞成功"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "取消点赞记录上传失败"
// @Router /api/{userid}/{type}/{postid}/cancellike [post]
func (lc *LikeControllers) CancelLike(c *gin.Context) {
	PostID := c.Param("postid")
	UserID := c.Param("userid")
	Type:=c.Param("type")
	if PostID == "" || UserID == ""||Type=="" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := lc.LikeService.HandleUserAction(Type, "like", PostID, UserID, Type, "cancel"); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, models.Write("取消点赞成功！"))
}

// @Summary 获取用户数据
// @Description 在主页获取用户获赞数，关注数，粉丝数
// @Tags 个人主页
// @Accept application/json
// @Produce application/json
// @Param userid path string true "用户ID"
// @Success 200 {object} models.Response "查询用户数据成功"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "查询用户数据失败"
// @Router /api/{userid}/userpage/count [get]
func (lc *LikeControllers) GetUserLikes(c *gin.Context) {
	UserID := c.Param("userid")
	if UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	likecount, err1 := lc.LikeService.GetUserLikes(UserID)
	followercount, err2 := lc.LikeService.GetCount("follower", UserID, "follow")
	followeecount, err3 := lc.LikeService.GetCount("followee", UserID, "follow")
	if err1 != nil || err2 != nil || err3 != nil {
		fmt.Printf("get user counts error:%v %v %v", err1, err2, err3)
		c.JSON(http.StatusInternalServerError, models.Write("查询用户数据失败"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "查询用户获赞数成功！", "LikesCount": likecount, "FollowerCount": followercount, "FolloweeCount": followeecount})
}

// @Summary 获取用户历史一个月内的历史点赞记录
// @Description 在个人主页获取用户一个月（31天）内的点赞记录
// @Tags 个人主页-记录查询
// @Accept application/json
// @Produce application/json
// @Param userid path string true "用户ID"
// @Param type path string true "分区"
// @Success 200 {object} map[string]interface{} "成功信息以及历史点赞帖子"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "获取点赞记录失败"
// @Router /api/{userid}/userpage/likerecord/{type} [get]
func (lc *LikeControllers) GetLikeRecord(c *gin.Context) {
	UserID := c.Param("userid")
	Type := c.Param("type")
	if UserID == "" || Type == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	Record, err := lc.LikeService.GetLikeRecord(UserID, Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	PostMessages, Err := lc.PostService.SearchPostByID(Record)
	if Err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("点赞记录查询失败！"))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":    "查询历史点赞记录成功！",
		"LikeRecord": PostMessages,
	})
}

// @Summary 用户收藏帖子
// @Description 用户给帖子收藏，收藏数加一
// @Tags 帖子-收藏
// @Accept application/json
// @Produce application/json
// @Param postid path string true "帖子ID"
// @Param userid path string true "用户ID"
// @Param type path string true "帖子类型"
// @Success 201 {object} models.Response "用户收藏成功"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "收藏记录上传失败"
// @Router /api/{type}/{postid}/{userid}/save [post]
func (lc *LikeControllers) SavePost(c *gin.Context) {
	PostID := c.Param("postid")
	UserID := c.Param("userid")
	Type:=c.Param("type")
	if PostID == "" || UserID == "" ||Type==""{
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := lc.LikeService.HandleUserAction(Type, "save", PostID, UserID, Type, "like"); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, models.Write("用户收藏成功！"))
}

// @Summary 用户取消收藏帖子
// @Description 用户取消收藏，收藏数减一
// @Tags 帖子-收藏
// @Accept application/json
// @Produce application/json
// @Param postid path string true "帖子ID"
// @Param userid path string true "用户ID"
// @Param type path string true "帖子类型"
// @Success 201 {object} models.Response "用户取消收藏成功"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "取消收藏记录上传失败"
// @Router /api/{userid}/{type}/{postid}/cancelsave [post]
func (lc *LikeControllers) CancelSave(c *gin.Context) {
	PostID := c.Param("postid")
	UserID := c.Param("userid")
	Type:=c.Param("type")
	if PostID == "" || UserID == "" ||Type==""{
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := lc.LikeService.HandleUserAction(Type, "save", PostID, UserID, Type, "cancel"); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, models.Write("取消收藏成功！"))
}

// @Summary 获取用户历史收藏记录
// @Description 在个人主页获取用户的历史收藏记录
// @Tags 个人主页-记录查询-生活
// @Accept application/json
// @Produce application/json
// @Param userid path string true "用户ID"
// @Param type path string true "帖子类型"
// @Param type path string true "分区"
// @Success 200 {object} map[string]interface{} "成功信息以及历史收藏帖子"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "获取收藏记录失败"
// @Router /api/{userid}/userpage/saverecord/{type} [get]
func (lc *LikeControllers) GetSaveRecord(c *gin.Context) {
	UserID := c.Param("userid")
	Type := c.Param("type")
	if UserID == "" || Type == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	Record, err := lc.LikeService.GetSaveRecord(UserID, Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	var Records []string
	for _, Post := range Record {
		Records = append(Records, Post.PostID)
	}
	PostMessages, Err := lc.PostService.SearchPostByID(Records)
	if Err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("收藏记录查询失败！"))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":    "查询历史收藏记录成功！",
		"LikeRecord": PostMessages,
	})
}

// @Summary 获取用户浏览记录
// @Description 在个人主页获取用户的历史浏览记录
// @Tags 个人主页-记录查询
// @Accept application/json
// @Produce application/json
// @Param userid path string true "用户ID"
// @Param type path string true "分区"
// @Success 200 {object} map[string]interface{} "成功信息以及历史浏览帖子"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "获取浏览记录失败"
// @Router /api/{userid}/userpage/viewrecord/{type} [get]
func (lc *LikeControllers) GetViewRecord(c *gin.Context) {
	UserID := c.Param("userid")
	Type := c.Param("type")
	if UserID == "" || Type == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	Record, err := lc.LikeService.GetViewRecord(UserID, Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	var Records []string
	for _, Post := range Record {
		Records = append(Records, Post.PostID)
	}
	PostMessages, Err := lc.PostService.SearchPostByID(Records)
	if Err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("用户浏览记录查询失败！"))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":    "查询历史浏览记录成功！",
		"LikeRecord": PostMessages,
	})
}

// @Summary 用户点赞评论
// @Description 用户给评论点赞，评论数加一
// @Tags 帖子-评论
// @Accept application/json
// @Produce application/json
// @Param commentid path string true "评论ID"
// @Param userid path string true "用户ID"
// @Param postid path string true "帖子ID"
// @Param type path string true "帖子类型"
// @Success 201 {object} models.Response "用户收藏成功"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "收藏记录上传失败"
// @Router /api/{userid}/{type}/{postid}/{commentid}/like [post]
func (lc *LikeControllers) LikeComment(c *gin.Context) {
	CommentID := c.Param("commentid")
	UserID := c.Param("userid")
	PostID := c.Param("postid")
	Type:=c.Param("type")
	if CommentID == "" || UserID == "" || PostID == "" ||Type==""{
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := lc.LikeService.HandleUserAction("comment", "like", CommentID, UserID, Type, "like"); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	User, _ := lc.UserService.SearchUserByID(UserID)
	Comment, err1 := lc.PostService.SearchCommentByID(CommentID)
	if err1 != nil {
		c.JSON(http.StatusMultiStatus, models.Write("点赞处理成功，消息上传失败"))
	}
	message := fmt.Sprintf("用户%s点赞了你的评论:\n%s", User.Username, Comment.Body)
	if err := lc.LikeService.InitMessage(Comment.CommenterID, message, User.ImageURL, "点赞", PostID); err != nil {
		c.JSON(http.StatusMultiStatus, models.Write("点赞处理成功，消息上传失败"))
	}
	c.JSON(http.StatusCreated, models.Write("点赞评论成功！"))
}

// @Summary 用户取消点赞评论
// @Description 用户取消点赞评论，评论数减一
// @Tags 帖子-评论
// @Accept application/json
// @Produce application/json
// @Param commentid path string true "评论ID"
// @Param userid path string true "用户ID"
// @Param type path string true "帖子类型"
// @Success 201 {object} models.Response "用户取消点赞成功"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "取消点赞记录上传失败"
// @Router /api/{userid}/{type}/{postid}/{commentid}/cancellike [post]
func (lc *LikeControllers) CancelCommentLike(c *gin.Context) {
	CommentID := c.Param("commentid")
	UserID := c.Param("userid")
	Type:=c.Param("type")
	if CommentID == "" || UserID == ""||Type=="" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := lc.LikeService.HandleUserAction("comment", "like", CommentID, UserID, Type, "cancel"); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, models.Write("取消点赞成功！"))
}

// @Summary 用户点赞回复
// @Description 用户给回复点赞，评论数加一
// @Tags 帖子-评论
// @Accept application/json
// @Produce application/json
// @Param replyid path string true "回复ID"
// @Param userid path string true "用户ID"
// @Param postid path string true "帖子ID"
// @Param type path string true "帖子类型"
// @Success 201 {object} models.Response "用户点赞成功"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "点赞记录上传失败"
// @Router /api/{userid}/{type}/{postid}/{commentid}/{replyid}/like [post]
func (lc *LikeControllers) LikeReply(c *gin.Context) {
	ReplyID := c.Param("replyid")
	UserID := c.Param("userid")
	PostID := c.Param("postid")
	Type:=c.Param("type")
	if ReplyID == "" || UserID == "" || PostID == "" ||Type==""{
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := lc.LikeService.HandleUserAction("reply", "like", ReplyID, UserID, "life", "like"); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	User, _ := lc.UserService.SearchUserByID(UserID)
	Reply, err1 := lc.PostService.SearchReplyByID(ReplyID)
	if err1 != nil {
		c.JSON(http.StatusMultiStatus, models.Write("点赞处理成功，消息上传失败"))
	}
	message := fmt.Sprintf("用户%s点赞了你的评论:\n%s", User.Username, Reply.Body)
	if err := lc.LikeService.InitMessage(Reply.ReplyerID, message, User.ImageURL, "life", PostID); err != nil {
		c.JSON(http.StatusMultiStatus, models.Write("点赞处理成功，消息上传失败"))
	}
	c.JSON(http.StatusCreated, models.Write("点赞评论成功！"))
}

// @Summary 用户取消点赞回复
// @Description 用户取消点赞回复，点赞数减一
// @Tags 帖子-评论
// @Accept application/json
// @Produce application/json
// @Param replyid path string true "回复ID"
// @Param userid path string true "用户ID"
// @Param type path string true "帖子类型"
// @Success 201 {object} models.Response "用户取消点赞成功"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "取消点赞记录上传失败"
// @Router /api/{userid}/{type}/{postid}/{commentid}/{replyid}/cancellike [post]
func (lc *LikeControllers) CancelReplyLike(c *gin.Context) {
	ReplyID := c.Param("replyid")
	UserID := c.Param("userid")
	Type:=c.Param("type")
	if ReplyID == "" || UserID == ""||Type=="" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := lc.LikeService.HandleUserAction("reply", "like", ReplyID, UserID, Type, "cancel"); err != nil {
		c.JSON(http.StatusInternalServerError, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, models.Write("取消点赞成功！"))
}

// @Summary 获取帖子相关数值
// @Description 获取帖子的点赞数，收藏数，浏览量
// @Tags 帖子
// @Accept application/json
// @Produce application/json
// @Param postid path string true "帖子ID"
// @Success 200 {object} map[string]interface{} "相关数值"
// @Failure 400 {object} models.Response "输入无效"
// @Router /api/{userid}/post/{postid}/getcounts [get]
func (lc *LikeControllers) GetPostCounts(c *gin.Context) {
	PostID := c.Param("postid")
	if PostID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	LikeCounts, err1 := lc.LikeService.GetCount("post", PostID, "like")
	SaveCounts, err2 := lc.LikeService.GetCount("post", PostID, "save")
	ViewCounts, err3 := lc.LikeService.GetCount("post", PostID, "view")
	if err1 != nil || err2 != nil || err3 != nil {
		fmt.Print(err1)
		log.Printf("get post %s count error", PostID)
	}
	c.JSON(http.StatusOK, gin.H{"likecount": LikeCounts, "savecount": SaveCounts, "ViewCount": ViewCounts})
}

// @Summary 获取评论点赞数
// @Description 获取评论的点赞数
// @Tags 帖子-评论
// @Accept application/json
// @Produce application/json
// @Param commentid path string true "评论ID"
// @Success 200 {object} map[string]interface{} "相关数值"
// @Failure 400 {object} models.Response "输入无效"
// @Router /api/{userid}/post/{postid}/{commentid}/getcounts [get]
func (lc *LikeControllers) GetCommentCounts(c *gin.Context) {
	CommentID := c.Param("commentid")
	if CommentID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	LikeCount, err := lc.LikeService.GetCount("comment", CommentID, "like")
	if err != nil {
		log.Printf("get comment %s count error", CommentID)
	}
	c.JSON(http.StatusOK, gin.H{"likecount": LikeCount})
}

// @Summary 获取回复点赞数
// @Description 获取回复的点赞数
// @Tags 帖子-评论
// @Accept application/json
// @Produce application/json
// @Param replyid path string true "回复ID"
// @Success 200 {object} map[string]interface{} "相关数值"
// @Failure 400 {object} models.Response "输入无效"
// @Router /api/{userid}/post/{postid}/{commentid}/{replyid}/getcounts [get]
func (lc *LikeControllers) GetReplyCounts(c *gin.Context) {
	ReplyID := c.Param("replyid")
	if ReplyID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	LikeCount, err := lc.LikeService.GetCount("reply", ReplyID, "like")
	if err != nil {
		log.Printf("get reply %s count error", ReplyID)
	}
	c.JSON(http.StatusOK, gin.H{"likecount": LikeCount})
}

// 定期更新数据至数据库
func (lc *LikeControllers) UpdateAllCount(ctx context.Context) {
	lc.LikeService.StartUpdateTicker(ctx)
}

// @Summary 关注用户
// @Description 关注用户并上传关注消息
// @Tags 关注
// @Accept application/json
// @Produce application/json
// @Param userid path string true "用户ID"
// @Param followid path string true "被关注用户ID"
// @Success 201 {object} models.Response "关注成功"
// @Failure 409 {object} models.Response "用户已关注"
// @Failure 400 {object} models.Response "输入无效"
// @Success 207 {object} models.Response "消息上传失败，关注成功"
// @Router /api/{userid}/userpage/{followid}/follow [post]
func (lc *LikeControllers) FollowUser(c *gin.Context) {
	UserID := c.Param("userid")
	FolloweeID := c.Param("followid")
	if UserID == "" || FolloweeID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := lc.LikeService.Follow(UserID, FolloweeID); err != nil {
		c.JSON(http.StatusConflict, models.Write(err.Error()))
		return
	}
	User, _ := lc.UserService.SearchUserByID(UserID)
	message := fmt.Sprintf("用户%s关注了你", User.Username)
	if err := lc.LikeService.InitMessage(FolloweeID, message, User.ImageURL, "关注", UserID); err != nil {
		c.JSON(http.StatusMultiStatus, models.Write("关注处理成功，消息上传失败"))
	}
	c.JSON(http.StatusCreated, models.Write("关注成功！"))
}

// @Summary 取消关注用户
// @Description 取消关注用户
// @Tags 关注
// @Accept application/json
// @Produce application/json
// @Param userid path string true "用户ID"
// @Param followid path string true "被关注用户ID"
// @Success 201 {object} models.Response "取消关注成功"
// @Failure 409 {object} models.Response "用户未关注"
// @Failure 400 {object} models.Response "输入无效"
// @Router /api/{userid}/userpage/{followid}/cancelfollow [post]
func (lc *LikeControllers) CancelFollowUser(c *gin.Context) {
	UserID := c.Param("userid")
	FolloweeID := c.Param("followid")
	if UserID == "" || FolloweeID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	if err := lc.LikeService.CancelFollow(UserID, FolloweeID); err != nil {
		c.JSON(http.StatusConflict, models.Write(err.Error()))
		return
	}
	c.JSON(http.StatusCreated, models.Write("取消关注成功！"))
}

// @Summary 关注列表
// @Description 获取用户关注列表
// @Tags 个人主页-关注列表
// @Accept application/json
// @Produce application/json
// @Param userid path string true "用户ID"
// @Success 200 {object} map[string]interface{} "响应信息和关注列表"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "获取关注列表失败"
// @Router /api/{userid}/userpage/followeelist  [get]
func (ls *LikeControllers) GetFolloweeList(c *gin.Context) {
	UserID := c.Param("userid")
	if UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试"))
		return
	}
	FolloweeMessages, err := ls.LikeService.GetUserFoloweeList(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("获取用户列表失败"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "获取用户列表成功！", "followees": FolloweeMessages})
}

// @Summary 粉丝列表
// @Description 获取用户粉丝列表
// @Tags 个人主页-关注列表
// @Accept application/json
// @Produce application/json
// @Param userid path string true "用户ID"
// @Success 200 {object} map[string]interface{} "响应信息和粉丝列表"
// @Failure 400 {object} models.Response "输入无效"
// @Failure 500 {object} models.Response "获取粉丝列表失败"
// @Router /api/{userid}/userpage/followerlist  [get]
func (ls *LikeControllers) GetFollowerList(c *gin.Context) {
	UserID := c.Param("userid")
	if UserID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试"))
		return
	}
	FollowerMessages, err := ls.LikeService.GetUserFolowerList(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Write("获取用户列表失败"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "获取用户列表成功！", "followers": FollowerMessages})
}

// @Summary 获取帖子点赞收藏状态
// @Description 获取帖子点赞收藏状态
// @Tags 生活区
// @Accept  json
// @Produce  json
// @Param postid path string true "帖子ID"
// @Param userid path string true "用户ID"
// @Success 200 {object} map[string]interface{} "帖子状态信息"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Router /api/{userid}/post/{postid}/getstatus [get]
func (lc *LikeControllers) GetPostStatus(c *gin.Context) {
	UserID := c.Param("userid")
	PostID := c.Param("postid")
	if UserID == "" || PostID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	LikeStatus := lc.LikeService.GetEntityStatus("post_like", UserID, PostID)
	SaveStatus := lc.LikeService.GetEntityStatus("post_save", UserID, PostID)
	c.JSON(http.StatusOK, gin.H{"LikeStatus": LikeStatus, "SaveStatus": SaveStatus})
}

// @Summary 获取评论点赞状态
// @Description 获取评论点赞状态
// @Tags 生活区
// @Accept  json
// @Produce  json
// @Param commentid path string true "评论ID"
// @Param userid path string true "用户ID"
// @Success 200 {object} map[string]interface{} "评论状态信息"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Router /api/{userid}/post/{postid}/{commentid}/getstatus [get]
func (lc *LikeControllers) GetCommentStatus(c *gin.Context) {
	UserID := c.Param("userid")
	CommentID := c.Param("commentid")
	if UserID == "" || CommentID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	LikeStatus := lc.LikeService.GetEntityStatus("comment_like", UserID, CommentID)
	c.JSON(http.StatusOK, gin.H{"LikeStatus": LikeStatus})
}

// @Summary 获取回复点赞状态
// @Description 获取回复点赞状态
// @Tags 生活区
// @Accept  json
// @Produce  json
// @Param replyid path string true "回复ID"
// @Param userid path string true "用户ID"
// @Success 200 {object} map[string]interface{} "回复状态信息"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Router /api/{userid}/post/{postid}/{commentid}/{replyid}/getstatus [get]
func (lc *LikeControllers) GetReplyStatus(c *gin.Context) {
	UserID := c.Param("userid")
	ReplyID := c.Param("replyid")
	if UserID == "" || ReplyID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	LikeStatus := lc.LikeService.GetEntityStatus("reply_like", UserID, ReplyID)
	c.JSON(http.StatusOK, gin.H{"LikeStatus": LikeStatus})
}

// @Summary 获取关注状态
// @Description 获取关注状态
// @Tags 生活区
// @Accept  json
// @Produce  json
// @Param followid path string true "关注人ID"
// @Param userid path string true "用户ID"
// @Success 200 {object} map[string]interface{} "关注状态信息"
// @Failure 400 {object} models.Response "输入无效，请重试!"
// @Router /api/{userid}/userpage/{followid}/getstatus [get]
func (lc *LikeControllers) GetFollowStatus(c *gin.Context) {
	UserID := c.Param("userid")
	FolloweeID := c.Param("followid")
	if UserID == "" || FolloweeID == "" {
		c.JSON(http.StatusBadRequest, models.Write("输入无效，请重试！"))
		return
	}
	FollowStatus := lc.LikeService.GetEntityStatus("follow", UserID, FolloweeID)
	c.JSON(http.StatusOK, gin.H{"FollowStatus": FollowStatus})
}
