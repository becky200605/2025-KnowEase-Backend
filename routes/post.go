package routes

import (
	"KnowEase/controllers"
	"KnowEase/middleware"
	"context"

	"github.com/gin-gonic/gin"
)

type PostSvc struct {
	lc *controllers.LikeControllers
	pc *controllers.PostControllers
	m  *middleware.Middleware
}

func NewPostSvc(lc *controllers.LikeControllers, m *middleware.Middleware, pc *controllers.PostControllers) *PostSvc {
	return &PostSvc{
		lc: lc,
		m:  m,
		pc: pc,
	}
}
func (p *PostSvc) NewPostGroup(r *gin.Engine) {
	ctx,cancel:=context.WithCancel(context.Background())
	defer cancel()
	go p.lc.UpdateAllCount(ctx)
	r.Use(p.m.Cors())
	r.Use(p.m.Verifytoken())
	posts := r.Group("/api")
	{
		posts.POST("/:userid/post/publish", p.pc.PublishPostBody)
		posts.GET("/:userid/post/recommend", p.pc.RecommendationPost)
		posts.GET("/:userid/post/campus", p.pc.CampusPost)
		posts.GET("/:userid/post/life", p.pc.LifePost)
		posts.GET("/:userid/post/food", p.pc.FoodPost)
		posts.GET("/:userid/post/paint", p.pc.PaintPost)
		posts.GET("/:userid/post/:postid", p.pc.GetPostMessage)
		posts.POST("/:userid/post/:postid/publishcomment", p.pc.PublishComment)
		posts.DELETE("/:userid/post/:postid/:commentid/deletecomment", p.pc.DeleteComment)
		posts.POST("/:userid/post/:postid/:commentid/publishreply", p.pc.PublishReply)
		posts.DELETE("/:userid/post/:postid/:commentid/:replyid", p.pc.DeleteReply)
		posts.POST("/:userid/:type/:postid/like", p.lc.LikePost)
		posts.POST("/:userid/:type/:postid/cancellike", p.lc.CancelLike)
		posts.POST("/:userid/:type/:postid/save", p.lc.SavePost)
		posts.POST("/:userid/:type/:postid/cancelsave", p.lc.CancelSave)
		posts.POST("/:userid/:type/:postid/:commentid/like", p.lc.LikeComment)
		posts.POST("/:userid/:type/:postid/:commentid/cancellike", p.lc.CancelCommentLike)
		posts.POST("/:userid/:type/:postid/:commentid/:replyid/like", p.lc.LikeReply)
		posts.POST("/:userid/:type/:postid/:commentid/:replyid/cancellike", p.lc.CancelReplyLike)
		posts.GET("/:userid/post/:postid/getcounts", p.lc.GetPostCounts)
		posts.GET("/:userid/post/:postid/:commentid/getcounts", p.lc.GetCommentCounts)
		posts.GET("/:userid/post/:postid/:commentid/:replyid/getcounts", p.lc.GetReplyCounts)
		posts.GET("/getToken", p.pc.GetToken)
		posts.GET("/:userid/post/:postid/getstatus", p.lc.GetPostStatus)
		posts.GET("/:userid/post/:postid/:commentid/getstatus", p.lc.GetCommentStatus)
		posts.GET("/:userid/post/:postid/:commentid/:replyid/getstatus", p.lc.GetReplyStatus)
	}
}
