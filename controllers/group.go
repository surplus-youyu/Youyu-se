package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/surplus-youyu/Youyu-se/models"
	"github.com/surplus-youyu/Youyu-se/utils"
)

func CreateGroup(c *gin.Context) {
	type ReqBody struct {
		Name  string `json:"name"`
		Intro string `json:"intro"`
		IsPublic int `json:"is_public"`
		Type int `json:"type"`
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
	group := models.Group{
		Owner:    user.Uid,
		Name:     req.Name,
		Intro:    req.Intro,
		IsPublic: req.IsPublic,
		Type:     req.Type,
	}
	models.CreateGroup(group)

	c.JSON(200, gin.H{
		"msg": "OK",
		"status": true,
	})
}

func GetGroupList(c *gin.Context) {
	groups := models.GetGroupList()
	c.JSON(200, gin.H{
		"msg": "OK",
		"status" : true,
		"data" :groups,
	})
}

func GetJoinedGroup(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	groups := models.JoinedGroupList(user.Uid)
	c.JSON(200, gin.H{
		"msg": "OK",
		"status" : true,
		"data" :groups,
	})
}

func RemoveMemberFromGroup(c *gin.Context) {
	type ReqBody struct {
		Member  int    `json:"member"`
		Gid     int    `json:"gid"`
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
	models.RemoveFromGroup(user.Uid, req.Gid, req.Member)
	c.JSON(200, gin.H{
		"msg": "OK",
		"status": true,
	})
}

func DeleteGroup(c *gin.Context) {
	type ReqBody struct {
		Gid     int    `json:"gid"`
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
	models.DeleteGroup(user.Uid, req.Gid)
	c.JSON(200, gin.H{
		"msg": "OK",
		"status": true,
	})
}

func JoinGroup(c *gin.Context) {
	type ReqBody struct {
		Gid     int    `json:"gid"`
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
	models.JoinGroup(user.Uid, req.Gid)

	c.JSON(200, gin.H{
		"msg": "OK",
		"status": true,
	})
}

func GetMembers(c *gin.Context)  {
	gid := utils.StringToInt(c.Param("gid"), c)
	members := models.GetGroupMembers(gid)
	c.JSON(200, gin.H{
		"msg": "OK",
		"status": true,
		"data": members,
	})
}