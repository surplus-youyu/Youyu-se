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
		api.GET("/tasks", controllers.GetTaskList)                           // all tasks
		api.POST("/tasks", controllers.CreateTask)                           // create task
		api.GET("/tasks/:task_id", controllers.GetTaskByID)                  // get task
		api.GET("/tasks/:task_id/files/:filename", controllers.GetTaskFiles) // get task
		api.PUT("/tasks/:task_id", controllers.FinishTask)                   // finish the task and close it, only can be accessed by owner

		// assignments apis
		api.GET("/assignments", controllers.GetAssignList)                             // get current user's assignments
		api.POST("/assignments", controllers.AssignTask)                               // create assignment with task id
		api.GET("/assignments/:assign_id", controllers.GetAssignmentByID)              // get assignment detail, only can be accessed by assignee
		api.PUT("/assignments/:assign_id", controllers.SubmitAssign)                   // submit assginment content
		api.GET("/tasks/:task_id/assignments", controllers.GetAssignListByTaskID)      // get assignments with task id, only can be accessed by owner
		api.PUT("/tasks/:task_id/assignments/:assign_id", controllers.JudgeAssignment) // judge the assignment
		api.GET("/tasks/:task_id/statistics", controllers.GetSurveyStatistics)         // get survey statistics

		// user apis
		api.GET("/user", controllers.GetUserInfo)
		api.PUT("/user", controllers.UpdateUserInfo)
		api.GET("/user/:uid", controllers.GetUserInfoById) // fuck httprouter
		api.GET("/user/:uid/avatar", controllers.GetAvatar)
		api.POST("/user/:uid/avatar", controllers.UpdateAvatar)

	}
}
