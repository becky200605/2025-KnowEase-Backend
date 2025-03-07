package routes

import "github.com/gin-gonic/gin"

type APP struct {
	r *gin.Engine
}

func NewApp(usv *UserSvc, psv *PostSvc, upsv *UserPageSvc,qasv *QASvc,ssv *SearchSvc,acsv *AIChatSvc) *APP {
	r := gin.Default()
	usv.NewUserGroup(r)
	psv.NewPostGroup(r)
	upsv.NewUserPageGroup(r)
	qasv.NewQAGroup(r)
	ssv.NewSearchGroup(r)
	acsv.NewAIChatGroup(r)
	return &APP{
		r: r,
	}
}
func (a *APP) Run() {
	a.r.Run(":8080")
}
