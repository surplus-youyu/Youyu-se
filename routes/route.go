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

	r.Group("/api")
	{
		// auth api
		r.PUT("/login", controllers.LoginHandler)
		r.POST("/register", controllers.RegisterHandler)

		// login middleware
		r.Use(loginRequired())

		// user apis
		r.Group("/user")
		{
			r.GET("/", controllers.GetUserInfo)
			r.PUT("/", controllers.UpdateUserInfo)
			r.GET("/avatar", controllers.GetAvatar)
			r.PUT("/avatar", controllers.UpdateAvatar)
		}

		// tasks api
		r.Group("/tasks")
		{
			r.GET("/", controllers.GetAllSurvey)
			r.POST("/", controllers.SurveyCreateHandler)
			r.GET("/:tid", controllers.QuerySurveyHandler)
		}
	}
}
