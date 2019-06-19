package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/surplus-youyu/Youyu-se/models"
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

		session := sessions.Default(c)
		session.Set("userEmail", user[0].Email)
		err = session.Save()

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
		NickName string `json:"nickname"`
		Email    string `json:"email"`
		Age      int    `gorm:"column:age"`
		Gender   string `gorm:"column:gender"`
		Phone    string `gorm:"column:phone"`
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
			Password: req.Password,
			NickName: req.NickName,
			Balance:  0.0,
			Email:    req.Email,
			Age:      req.Age,
			Gender:   req.Gender,
			Phone:    req.Phone,
		}
		models.CreateNewUser(newUser)
	}
	c.JSON(statusCode, gin.H{
		"status": true,
		"msg":    msg,
	})
}
