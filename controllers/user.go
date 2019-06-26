package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/surplus-youyu/Youyu-se/models"
	"github.com/surplus-youyu/Youyu-se/utils"
	"net/http"
)

const avatarPath = "static/avatars/"

func GetUserInfo(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	// 200
	c.JSON(http.StatusOK, gin.H{
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
	id := utils.StringToInt(c.Param("uid"), c)
	user := models.GetUserById(id)
	// 200
	c.JSON(http.StatusOK, gin.H{
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
		// 400
		c.JSON(http.StatusBadRequest, gin.H{
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

	// 200
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"msg":    "success",
	})
}

func GetAvatar(c *gin.Context) {
	uid := utils.StringToInt(c.Param("uid"), c)
	user := models.GetUserById(uid)
	c.File(avatarPath + user.Avatar)
}

func UpdateAvatar(c *gin.Context) {
	uid := utils.StringToInt(c.Param("uid"), c)

	user := c.MustGet("user").(models.User)

	if uid != user.Uid {
		// 403
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"msg":    "只能修改自己的头像",
			"status": true,
		})
	}

	newAvatar, err := c.FormFile("avatar")

	if err != nil {
		// 400
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "invalid avatar",
		})
	}

	if user.Avatar == "default" {
		user.Avatar = utils.IntToString(user.Uid)
		models.UpdateUser(user)
	}

	if err := c.SaveUploadedFile(newAvatar, avatarPath+user.Avatar); err != nil {
		// 500
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Cannot upload avatar",
		})
	}

	// 200
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"msg":    "success",
	})
}
