package routes

import (
	"KnowEase/controllers"
	"KnowEase/middleware"

	"github.com/gin-gonic/gin"
)

type AIChatSvc struct {
	acc *controllers.AIChatControllers
	m   *middleware.Middleware
}

func NewAIChatSvc(acc *controllers.AIChatControllers, m *middleware.Middleware) *AIChatSvc {
	return &AIChatSvc{
		acc: acc,
		m:   m,
	}
}
func (acs *AIChatSvc) NewAIChatGroup(r *gin.Engine) {
	r.Use(acs.m.Cors())
	r.Use(acs.m.Verifytoken())
	aichats := r.Group("/api/aichat")
	{
		aichats.POST("/chat/:userid/:chatid/postchat", acs.acc.UserInitQuestion)
		aichats.POST("/:userid/:chatid/sseconnect", acs.acc.SSEConnect)
		aichats.GET("/:userid/search/:chatid/gethistory",acs.acc.GetHistoryChatByID)
		aichats.GET("/:userid/getlist",acs.acc.GetUserAIChatList)
		aichats.POST("/create",acs.acc.CreateNewChat)
	}
}
