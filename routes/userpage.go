package routes

import (
	"KnowEase/controllers"
	"KnowEase/middleware"

	"github.com/gin-gonic/gin"
)

type UserPageSvc struct {
	lc *controllers.LikeControllers
	m  *middleware.Middleware
	pc *controllers.PostControllers
	uc *controllers.UserControllers
	qc *controllers.QAControllers
}

func NewUserPageSvc(lc *controllers.LikeControllers, m *middleware.Middleware, pc *controllers.PostControllers, uc *controllers.UserControllers, qc *controllers.QAControllers) *UserPageSvc {
	return &UserPageSvc{
		lc: lc,
		m:  m,
		pc: pc,
		uc: uc,
		qc: qc,
	}
}
func (up *UserPageSvc) NewUserPageGroup(r *gin.Engine) {
	r.Use(up.m.Cors())
	r.Use(up.m.Verifytoken())
	userpage := r.Group("/api")
	{
		userpage.GET("/:userid/userpage/count", up.lc.GetUserLikes)
		userpage.GET("/:userid/userpage/likerecord/:type", up.lc.GetLikeRecord)
		userpage.GET("/:userid/userpage/viewrecord/:type", up.lc.GetViewRecord)
		userpage.GET("/:userid/userpage/saverecord/:type", up.lc.GetSaveRecord)
		userpage.DELETE("/:userid/userpage/mypost/delete/:postid", up.pc.DeletePost)
		userpage.DELETE("/:userid/userpage/mypost/deleteposts", up.pc.DeletePosts)
		userpage.POST("/logout", up.uc.Logout)
		userpage.POST("/:userid/userpage/alterbackground", up.uc.ChangeUserBackground)
		userpage.POST("/:userid/userpage/alterimage", up.uc.ChangeUserPicture)
		userpage.POST("/:userid/userpage/alterpassword", up.uc.ChangeUserPassword)
		userpage.POST("/:userid/userpage/sendemail", up.uc.AlterSendemail)
		userpage.POST("/:userid/userpage/alteremail", up.uc.ChangeUserEmail)
		userpage.POST("/:userid/userpage/alterusername", up.uc.ChangeUsername)
		userpage.POST("/:userid/userpage/:followid/follow", up.lc.FollowUser)
		userpage.POST("/:userid/userpage/:followid/cancelfollow", up.lc.CancelFollowUser)
		userpage.GET("/:userid/userpage/followeelist", up.lc.GetFolloweeList)
		userpage.GET("/:userid/userpage/followerlist", up.lc.GetFollowerList)
		userpage.GET("/:userid/userpage/:followid/getstatus", up.lc.GetFollowStatus)
		userpage.GET("/:userid/followmessage", up.pc.GetUserUnreadFollowMessage)
		userpage.GET("/:userid/likemessage", up.pc.GetUserUnreadLikeMessage)
		userpage.GET("/:userid/commentmessage", up.pc.GetUserUnreadCommentMessage)
		userpage.DELETE("/mypost/:postid/deleteqas", up.qc.DeleteQAs)
		userpage.DELETE("/mypost/deleteqa", up.qc.DeleteQA)
		userpage.GET("/userpage/:userid/getuserpost", up.pc.GetUserPosts)
		userpage.GET("/userpage/:userid/getuserqa", up.qc.GetUserQA)
		userpage.GET("/userpage/:userid", up.uc.GetUserMessage)
	}
}
