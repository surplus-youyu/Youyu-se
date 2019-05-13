package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/surplus-youyu/Youyu-se/models"
	"time"
)

func LoginHandler(c *gin.Context) {
	type ReqBody struct {
		Email string `json:"email"`
		Pwd   string `json:"password"`
	}
	var req ReqBody
	var msg string
	statusCode := 200
	err := c.BindJSON(&req)
	if err != nil {
		msg = "invalid param"
		statusCode = 400
		c.JSON(statusCode, gin.H{
			"status": false,
			"msg":    msg,
		})
		return
	}
	user := models.GetUserByEmail(req.Email)
	if len(user) > 0 && user[0].Password == req.Pwd {
		msg = "login successfully"
	} else {
		msg = "username not exists or password mismatch"
		statusCode = 403
	}
	c.JSON(statusCode, gin.H{
		"status": true,
		"msg":    msg,
	})
}

func RegisterHandler(c *gin.Context) {
	type ReqBody struct {
		Password string `json:"password"`
		RealName string `json:"real_name"`
		NickName string `json:"nick_name"`
		Age      int    `json:"age"`
		Gender   string `json:"gender"`
		Major    string `json:"major"`
		Grade    int    `json:"grade"`
		Phone    string `json:"phone"`
		Email    string `json:"email"`
	}
	var req ReqBody
	var msg string
	statusCode := 200
	err := c.BindJSON(&req)
	if err != nil {
		msg = "invalid param"
		statusCode = 400
		c.JSON(statusCode, gin.H{
			"status": false,
			"msg":    msg,
		})
		return
	}
	user := models.GetUserByEmail(req.Email)
	if len(user) > 0 {
		statusCode = 400
		msg = "email address has been registered"
	} else {
		msg = "success"
		newUser := models.User{
			Uid:      int(time.Now().Unix()),
			Password: req.Password,
			RealName: req.RealName,
			NickName: req.NickName,
			Age:      req.Age,
			Gender:   req.Gender,
			Balance:  0.0,
			Major:    req.Major,
			Grade:    req.Grade,
			Phone:    req.Phone,
			Email:    req.Email,
		}
		models.CreateNewUser(newUser)
	}
	c.JSON(statusCode, gin.H{
		"status": true,
		"msg":    msg,
	})
}
