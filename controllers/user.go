package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/surplus-youyu/Youyu-se/models"
	"strconv"
)

const avatarPath = "static/avatars/"

func GetUserInfo(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	c.JSON(200, gin.H{
		"status": true,
		"msg":    "OK",
		"data": gin.H{
			"uid":      user.Uid,
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

func GetUserInfoById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		c.AbortWithStatus(400)
		return
	}
	user := models.GetUserById(id)
	c.JSON(200, gin.H{
		"status": true,
		"msg":    "OK",
		"data": gin.H{
			"nickname": user.NickName,
			"email":    user.Email,
			"age":      user.Age,
			"gender":   user.Gender,
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
		Major    string `json:"major"`
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

	if len(req.Pwd) != 0 {
		user.Password = req.Pwd
	}

	user.NickName = req.NickName
	user.Age = req.Age
	user.Gender = req.Gender
	user.Phone = req.Phone
	user.Major = req.Major
	user.Grade = req.Grade

	models.UpdateUser(user)

	c.JSON(200, gin.H{
		"status": true,
		"msg":    "success",
	})
}

func GetAvatar(c *gin.Context) {
	uid, err := strconv.Atoi(c.Param("uid"))
	if err != nil {
		c.AbortWithStatus(400)
	}
	user := models.GetUserById(uid)
	c.File(avatarPath + user.Avatar)
}

func UpdateAvatar(c *gin.Context) {
	uid, err := strconv.Atoi(c.Param("uid"))

	if err != nil {
		c.AbortWithStatus(400)
	}

	user := c.MustGet("user").(models.User)

	if uid != user.Uid {
		c.AbortWithStatusJSON(403, gin.H{
			"msg":    "只能修改自己的头像",
			"status": true,
		})
	}

	newAvatar, err := c.FormFile("avatar")

	if err != nil {
		c.JSON(400, gin.H{
			"status":  false,
			"message": "invalid avatar",
		})
	}

	if user.Avatar == "default" {
		user.Avatar = strconv.Itoa(user.Uid)
		models.UpdateUser(user)
	}

	if err := c.SaveUploadedFile(newAvatar, avatarPath+user.Avatar); err != nil {
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
