package routes

import (
	"KnowEase/controllers"
	"KnowEase/middleware"

	"github.com/gin-gonic/gin"
)

type SearchSvc struct {
	sc *controllers.SearchControllers
	m  *middleware.Middleware
}

func NewSearchSvc(sc *controllers.SearchControllers, m *middleware.Middleware) *SearchSvc {
	return &SearchSvc{
		sc: sc,
		m:  m,
	}
}
func (s *SearchSvc) NewSearchGroup(r *gin.Engine) {
	r.Use(s.m.Cors())
	r.Use(s.m.Verifytoken())
	searchs := r.Group("/api")
	{
		searchs.POST("/search", s.sc.Search)
		searchs.GET("/search/:userid/getrecord", s.sc.GetUserSearchRecord)
		searchs.POST("/search/:userid/recommend", s.sc.GetSearchRecommend)
		searchs.POST("/search/delete", s.sc.DeleteUserRecord)
		searchs.DELETE("/search/:userid/deleteall", s.sc.DeleteUserALLRecord)
		searchs.GET("/search/gethotrank", s.sc.GetHotRank)
		searchs.POST("/search/syncRecord", s.sc.SyncSearchRecord )

	}
}
