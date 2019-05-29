package route

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/surplus-youyu/Youyu-se/controller"
)

func loginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		email := session.Get("email")
		if email == nil {
			c.Abort()
			c.JSON(401, gin.H{
				"status": false,
				"msg":    "you should login first",
			})
			return
		}
		c.Set("email", email)
		c.Next()
	}
}

func Route(r *gin.Engine) {
	r.PUT("/api/login", controller.LoginHandler)
	r.POST("/api/register", controller.RegisterHandler)

	r.Use(loginRequired())

	r.GET("/api/surveys/:sid", controller.QuerySurveyHandler)
	// TODO
	// 填写提交问卷接口 PUT or POST
	// r.PUT("/api/surveys/:sid", )

	// TODO
	// 获取一个人所有的问卷
	// r.GET("/api/:uid/surveys")

	r.GET("/api/surveys", controller.GetAllSurvey)
	r.POST("/api/surveys", controller.SurveyCreateHandler)
}
