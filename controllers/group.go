package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/surplus-youyu/Youyu-se/models"
	"github.com/surplus-youyu/Youyu-se/utils"
	"net/http"
)

func CreateGroup(c *gin.Context) {
	type ReqBody struct {
		Name     string `json:"name"`
		Intro    string `json:"intro"`
		IsPublic int    `json:"is_public"`
		Type     int    `json:"type"`
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
	group := models.Group{
		Owner:    user.Uid,
		Name:     req.Name,
		Intro:    req.Intro,
		IsPublic: req.IsPublic,
		Type:     req.Type,
	}
	models.CreateGroup(group)

	// 200
	c.JSON(http.StatusOK, gin.H{
		"msg":    "OK",
		"status": true,
	})
}

func GetGroupList(c *gin.Context) {
	groups := models.GetGroupList()

	// 200
	c.JSON(http.StatusOK, gin.H{
		"msg":    "OK",
		"status": true,
		"data":   groups,
	})
}

func GetJoinedGroup(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	groups := models.JoinedGroupList(user.Uid)

	// 200
	c.JSON(http.StatusOK, gin.H{
		"msg":    "OK",
		"status": true,
		"data":   groups,
	})
}

func RemoveMemberFromGroup(c *gin.Context) {
	type ReqBody struct {
		Member int `json:"member"`
		Gid    int `json:"gid"`
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
	models.RemoveFromGroup(user.Uid, req.Gid, req.Member)
	// 200
	c.JSON(http.StatusOK, gin.H{
		"msg":    "OK",
		"status": true,
	})
}

func DeleteGroup(c *gin.Context) {
	type ReqBody struct {
		Gid int `json:"gid"`
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
	models.DeleteGroup(user.Uid, req.Gid)
	// 200
	c.JSON(http.StatusOK, gin.H{
		"msg":    "OK",
		"status": true,
	})
}

func JoinGroup(c *gin.Context) {
	type ReqBody struct {
		Gid int `json:"gid"`
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
	models.JoinGroup(user.Uid, req.Gid)

	// 200
	c.JSON(http.StatusOK, gin.H{
		"msg":    "OK",
		"status": true,
	})
}

func GetMembers(c *gin.Context) {
	gid := utils.StringToInt(c.Param("gid"), c)
	members := models.GetGroupMembers(gid)
	// 200
	c.JSON(http.StatusOK, gin.H{
		"msg":    "OK",
		"status": true,
		"data":   members,
	})
}
