package route

import (
	"github.com/gin-gonic/gin"
	"github.com/surplus-youyu/Youyu-se/controller"
)

func Route(r *gin.Engine) {
	r.PUT("/api/login", controller.LoginHandler)
	r.POST("/api/register", controller.RegisterHandler)

	r.GET("/api/survey/:sid", controller.QuerySurveyHandler)
	r.POST("/api/survey", controller.SurveyCreateHandler)
}

