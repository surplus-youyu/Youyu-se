package route

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/surplus-youyu/Youyu-se/controllers"
	"github.com/surplus-youyu/Youyu-se/models"
)

func loginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		email := session.Get("userEmail")
		if email == nil {
			c.Abort()
			c.JSON(401, gin.H{
				"status": false,
				"msg":    "you should login first",
			})
			return
		}
		user := models.GetUserByEmail(email.(string))[0]
		c.Set("user", user)
		c.Next()
	}
}

func Route(r *gin.Engine) {

	api := r.Group("/api")
	{
		// auth api
		api.PUT("/login", controllers.LoginHandler)
		api.POST("/register", controllers.RegisterHandler)

		// login middleware
		api.Use(loginRequired())

		// user apis
		user := api.Group("/user")
		{
			user.GET("/", controllers.GetUserInfo)
			user.PUT("/", controllers.UpdateUserInfo)
			user.GET("/avatar", controllers.GetAvatar)
			user.PUT("/avatar", controllers.UpdateAvatar)
		}

		// tasks api
		task := api.Group("/tasks")
		{
			task.GET("/", controllers.GetAllSurvey)
			task.POST("/", controllers.SurveyCreateHandler)
			task.GET("/:tid", controllers.QuerySurveyHandler)
		}
	}
}
