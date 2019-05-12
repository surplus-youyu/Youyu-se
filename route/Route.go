package route

import (
	"github.com/gin-gonic/gin"
	"github.com/surplus-youyu/Youyu-se/controller"
)

func Route(r *gin.Engine) {
	r.PUT("/api/login", controller.LoginHandler)
	r.POST("/api/register", controller.RegisterHandler)

	r.GET("/api/surveys/:sid", controller.QuerySurveyHandler)
	// TODO
	// 填写提交问卷接口 PUT or POST
	// r.PUT("/api/surveys/:sid", )
	
	// TODO
	// 获取一个人所有的问卷
	// r.GET("/api/:uid/surveys")

	r.GET("/api/surveys", controller.GetAllSurvey)
	r.POST("/api/survey", controller.SurveyCreateHandler)
}
