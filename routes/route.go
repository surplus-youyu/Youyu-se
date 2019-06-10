package route

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/surplus-youyu/Youyu-se/controllers"
	"github.com/surplus-youyu/Youyu-se/models"
	"github.com/surplus-youyu/Youyu-se/utils"
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
	api.Use(utils.HandleError)
	{
		// auth api
		api.PUT("/login", controllers.LoginHandler)
		api.POST("/register", controllers.RegisterHandler)

		// login middleware
		api.Use(loginRequired())

		// tasks apis
		api.GET("/tasks", controllers.GetTaskList)
		api.POST("/tasks", controllers.CreateTask)
		api.GET("/tasks/:task_id", controllers.GetTaskByID)
		api.POST("/tasks/:task_id/assign", controllers.AssignTask)

		// assignments apis
		api.GET("/assignments", controllers.GetAssignList)

		// user apis
		user := api.Group("/user")
		{
			user.GET("/", controllers.GetUserInfo)
			user.PUT("/", controllers.UpdateUserInfo)
			user.GET("/avatar", controllers.GetAvatar)
			user.PUT("/avatar", controllers.UpdateAvatar)
		}

		// tasks api
		// task := api.Group("/tasks")
		// {
		// 	task.GET("/", controllers.GetAllSurvey)
		// 	task.POST("/", controllers.SurveyCreateHandler)
		// 	task.GET("/:tid", controllers.QuerySurveyHandler)
		// }
	}
}
