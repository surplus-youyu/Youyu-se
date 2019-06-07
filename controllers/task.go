package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/surplus-youyu/Youyu-se/models"
	"strconv"
)

func QueryTaskHandler(c *gin.Context) {
	param := c.Param("tid")
	tid, err := strconv.ParseInt(param, 10, 32)
	msg := ""
	data := models.Task{}
	if err != nil {
		c.JSON(400, gin.H{
			"status": false,
			"msg":    "invalid param",
		})
		return
	}
	result := models.GetTaskById(int32(tid))
	if len(result) == 0 {
		msg = "task not found"
	} else {
		msg = "success"
		data = result[0]
	}
	c.JSON(200, gin.H{
		"status": true,
		"msg":    msg,
		"data":   data,
	})
}

func TaskCreateHandler(c *gin.Context) {
	type ReqBody struct {
		Title     string `json:"title"`
		Summary   string `json:"content"`
		Type      string `json:"type"`
		Bounty    int    `json:"bounty"`
		Extra     string `json:"extra"`
		Enclosure string `json:"enclosure"`
	}
	var body ReqBody
	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(400, gin.H{
			"status": false,
			"msg":    "invalid param",
		})
		return
	}

	user := c.MustGet("user").(models.User)

	newTask := models.Task{
		Owner:     user.Uid,
		Title:     body.Title,
		Summary:   body.Summary,
		Type:      body.Type,
		Bounty:    body.Bounty,
		Extra:     body.Extra,
		Enclosure: body.Enclosure,
	}
	models.CreateNewTask(newTask)
	c.JSON(200, gin.H{
		"status": true,
		"msg":    "success",
	})
}

func GetAllTask(c *gin.Context) {
	result := models.GetAllTasks()
	c.JSON(200, gin.H{
		"status": true,
		"msg":    "success",
		"data":   result,
	})
}

// TODO
//

func AnswerSubmit(c *gin.Context) {

}
