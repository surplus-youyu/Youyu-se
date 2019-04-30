package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/surplus-youyu/Youyu-se/models"
)

func LoginHandler(c *gin.Context) {
	type ReqBody struct {
		Email string `json:"email"`
		Pwd string `json:"password"`
	}
	var req ReqBody
	var msg string
	statusCode := 200
	err := c.BindJSON(&req)
	if err != nil {
		msg = "invalid param"
		statusCode = 400
		c.JSON(statusCode, gin.H{
			"status": "fail",
			"msg": msg,
		})
		return
	}
	user := models.GetUserByEmail(req.Email)
	if user.Password == req.Pwd {
		msg = "login successfully"
	} else {
		msg = "username not exists or password mismatch"
		statusCode = 403
	}
	c.JSON(statusCode, gin.H{
		"status": "OK",
		"msg": msg,
	})
}
