package middleware

import (
	"KnowEase/services"
	"fmt"
	"net/http"

	"log"
	"time"

	"KnowEase/utils"

	"github.com/go-co-op/gocron"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	TokenService *services.TokenService
}

func NewMiddleWare(TokenService *services.TokenService) *Middleware {
	return &Middleware{TokenService: TokenService}
}

// 检验token
func (m *Middleware) Verifytoken() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求头中的 token
		tokenString := c.GetHeader("Authorization")
		role, _ := m.TokenService.VerifyToken(tokenString)
		if role != "user" {
			c.JSON(http.StatusForbidden, gin.H{"message": "身份检验失败！"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// 解决跨域问题
func (m *Middleware) Cors() gin.HandlerFunc {
	c := cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:    []string{"Content-Type", "Access-Token", "Authorization"},
		MaxAge:          6 * time.Hour,
	}

	return cors.New(c)
}

// 捕获错误并写入日志
func (m *Middleware) ErrorAndLogger(Task func() error) func() {
	return func() {
		err := Task()
		if err != nil {
			log.Printf("执行定时任务时出错：%v", err)
		}
	}
}

// 定时任务中间件
func (m *Middleware) CronJobMiddleware(task func() error) {
	loc := time.Now().Location()
	s := gocron.NewScheduler(loc)
	Task := m.ErrorAndLogger(task)
	s.Every(10).Minutes().Do(Task)
	go s.StartBlocking()
}

// 数据埋点中间件
func (m *Middleware) Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.RequestCounter.With(prometheus.Labels{
			"Method":     c.Request.Method,
			"Path":       c.Request.URL.Path,
			"StatusCode": "200", //预设状态码
		})
		c.Next()
		//获取请求的状态码并更新
		StatusCode := fmt.Sprintf("%d", c.Writer.Status())
		utils.RequestCounter.With(prometheus.Labels{
			"Method":     c.Request.Method,
			"Path":       c.Request.URL.Path,
			"StatusCode": StatusCode, //预设状态码
		}).Inc()

	}
}
