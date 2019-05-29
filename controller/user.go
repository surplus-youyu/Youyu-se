package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/surplus-youyu/Youyu-se/models"
)

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
		},
	})
}

func UpdateUserInfo(c *gin.Context) {
	type ReqBody struct {
		Pwd      string `json:"password"`
		NickName string `json:"nickname"`
		Age      int    `json:"age"`
		Gender   string `json:"gender"`
		Phone    string `json:"Phone"`
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

	models.UpdateUser(user)

	c.JSON(200, gin.H{
		"status": true,
		"msg":    "success",
	})
}
