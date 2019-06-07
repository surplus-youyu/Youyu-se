package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/surplus-youyu/Youyu-se/models"
)

func GetTaskList(c *gin.Context) {
	taskList := models.GetTaskList()
	c.JSON(200, gin.H{
		"data": taskList,
		"msg":  "OK",
	})
}

func CreateTask(c *gin.Context) {
	type Req struct {
		Title       string  `json:"title"`
		Type        string  `json:"type"`
		Description string  `json:"description"`
		Content     string  `json:"content"`
		Reward      float32 `json:"reward"`
	}
	user := c.MustGet("user").(models.User)
	var body Req
	err := c.BindJSON(&body)
	if err != nil {
		c.AbortWithStatus(400)
		return
	}
	if user.Balance < body.Reward {
		c.AbortWithStatusJSON(403, gin.H{
			"msg": "余额不足",
		})
		return
	}
	insertID := models.CreateTask(models.Task{
		Title:       body.Title,
		Type:        body.Type,
		Description: body.Description,
		Content:     body.Content,
		Reward:      body.Reward,
		Creator:     user.Uid,
		State:       models.TaskStateCreated,
	}, user)
	c.JSON(200, gin.H{
		"data": gin.H{
			"id": insertID,
		},
		"msg": "OK",
	})
}

func GetTaskByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		c.AbortWithStatus(400)
		return
	}
	task := models.GetTaskByID(id)
	c.JSON(200, gin.H{
		"data": task,
		"msg":  "OK",
	})
}

// func AssginTask(c *gin.Context) {
// 	user := c.MustGet("user").(models.User)
// 	id, err := strconv.Atoi(c.Param("task_id"))
// 	if err != nil {
// 		c.AbortWithStatus(400)
// 	}
// 	// tricky serialization isolation level, should only be used with single-process-server
// 	models.DB.Lock()
// 	defer models.DB.Unlock()
// 	task := models.GetTaskByID(id)
// 	if task.AssignedTo != 0 {
// 		c.AbortWithStatusJSON(403, gin.H{
// 			"msg": "任务已被认领",
// 		})
// 		return
// 	}
// 	task.AssignedTo = user.Uid
// 	task.State = models.TaskStatePending
// 	models.UpsertTask(task)
// 	c.JSON(200, gin.H{
// 		"msg": "OK",
// 	})
// }
