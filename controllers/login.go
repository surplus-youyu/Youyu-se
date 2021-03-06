package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/surplus-youyu/Youyu-se/models"
	"net/http"
)

func LoginHandler(c *gin.Context) {
	type ReqBody struct {
		Email string `json:"email"`
		Pwd   string `json:"password"`
	}
	var req ReqBody
	var msg string
	err := c.BindJSON(&req)
	if err != nil {
		msg = "invalid param"
		// 400
		c.JSON(http.StatusBadRequest, gin.H{
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
		// 403
		c.JSON(http.StatusForbidden, gin.H{
			"status": false,
			"msg":    msg,
		})
		return
	}
	// 200
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"msg":    msg,
	})
}

func LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	// 删除用户登陆状态
	session.Delete("userEmail")
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"msg":    "login out succeed",
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
		Grade    string `json:"grade"`
		Major    string `json:"major"`
	}
	var req ReqBody
	var msg string

	err := c.BindJSON(&req)
	if err != nil {
		msg = "invalid param"
		// 400
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    msg,
		})
		return
	}
	user := models.GetUserByEmail(req.Email)
	if len(user) > 0 {
		msg = "email address has been registered"
		// 400
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    msg,
		})
		return
	}

	msg = "success"
	newUser := models.User{
		Password: req.Password,
		NickName: req.NickName,
		Balance:  100.0,
		Email:    req.Email,
		Age:      req.Age,
		Gender:   req.Gender,
		Phone:    req.Phone,
		Grade:    req.Grade,
		Major:    req.Major,
		Avatar:   "default",
	}
	models.CreateNewUser(newUser)
	// 200
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"msg":    msg,
	})
}
