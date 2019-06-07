package route

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/surplus-youyu/Youyu-se/controllers"
)

func loginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.Abort()
			c.JSON(401, gin.H{
				"status": false,
				"msg":    "you should login first",
			})
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

func Route(r *gin.Engine) {
	r.PUT("/api/user/login", controllers.LoginHandler)
	r.POST("/api/user/register", controllers.RegisterHandler)

	r.Use(loginRequired())

	r.GET("/api/user", controllers.GetUserInfo)
	r.PUT("/api/user", controllers.UpdateUserInfo)

	r.GET("/api/tasks/:tid", controllers.QueryTaskHandler)

	// TODO
	// 填写提交问卷接口 PUT or POST
	// r.PUT("/api/surveys/:sid", )

	// TODO
	// 获取一个人所有的问卷
	// r.GET("/api/:uid/surveys")

	r.GET("/api/surveys", controllers.GetAllTask)
	r.POST("/api/surveys", controllers.TaskCreateHandler)
}
