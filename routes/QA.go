package routes

import (
	"KnowEase/controllers"
	"KnowEase/middleware"

	"github.com/gin-gonic/gin"
)

type QASvc struct {
	qc *controllers.QAControllers
	m  *middleware.Middleware
}

func NewQASvc(qc *controllers.QAControllers, m *middleware.Middleware) *QASvc {
	return &QASvc{
		qc: qc,
		m:  m,
	}
}
func (qa *QASvc) NewQAGroup(r *gin.Engine) {
	go qa.m.CronJobMiddleware(qa.qc.SyncQuestion)
	r.Use(qa.m.Cors())
	r.Use(qa.m.Verifytoken())
	qas := r.Group("/api")
	{
		qas.POST("/:userid/QA/similarity", qa.qc.SearchSimilarity)
		qas.POST("/:userid/QA/publish", qa.qc.PublishQA)
		qas.GET("/QA/:postid", qa.qc.GetQAMessage)
		qas.POST("/QA/getbytag", qa.qc.GetQAPostByTag)
		qas.POST("/QA/getgoodqa", qa.qc.GetGoodPosts)
		qas.GET("/:userid/QA/:postid/getcounts", qa.qc.GetPostCounts)
	}
}
