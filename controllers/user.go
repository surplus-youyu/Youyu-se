package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/surplus-youyu/Youyu-se/models"
)

const avatarPath = "static/avatars/"

func GetUserInfo(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	c.JSON(200, gin.H{
		"status": true,
		"msg":    "OK",
		"data": gin.H{
			"nickname": user.NickName,
			"email":    user.Email,
			"age":      user.Age,
			"gender":   user.Gender,
			"phone":    user.Phone,
			"balance":  user.Balance,
			"grade":    user.Grade,
			"major":    user.Major,
		},
	})
}

func UpdateUserInfo(c *gin.Context) {
	type ReqBody struct {
		Pwd      string `json:"password"`
		NickName string `json:"nickname"`
		Age      int    `json:"age"`
		Gender   string `json:"gender"`
		Phone    string `json:"phone"`
		Grade    string `json:"grade"`
		major    string `json:"major"`
	}
	var req ReqBody
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"status": false,
			"msg":    "invalid param",
		})
		return
	}

	user := c.MustGet("user").(models.User)
	user.Password = req.Pwd
	user.NickName = req.NickName
	user.Age = req.Age
	user.Gender = req.Gender
	user.Phone = req.Phone
	user.Major = req.major
	user.Grade = req.Grade

	models.UpdateUser(user)

	c.JSON(200, gin.H{
		"status": true,
		"msg":    "success",
	})
}

func GetAvatar(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	c.File(avatarPath + user.Avatar)
}

func UpdateAvatar(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	newAvatar, err := c.FormFile("avatar")

	if err != nil {
		c.JSON(400, gin.H{
			"status":  false,
			"message": "invalid avatar",
		})
	}

	if user.Avatar == "default" {
		user.Avatar = string(user.Uid)
		models.UpdateUser(user)
	}

	if err := c.SaveUploadedFile(newAvatar, avatarPath+string(user.Uid)); err != nil {
		c.JSON(500, gin.H{
			"status":  false,
			"message": "Cannot upload avatar",
		})
	}

	c.JSON(200, gin.H{
		"status": true,
		"msg":    "success",
	})
}
