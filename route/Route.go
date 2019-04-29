package route

import (
	"github.com/gin-gonic/gin"
	"github.com/surplus-youyu/Youyu-se/controller"
)

func Route(r *gin.Engine) {
	r.PUT("/api/login", controller.LoginHandler)
}

