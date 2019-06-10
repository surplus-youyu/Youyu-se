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
		Title string `json:"title"`
		Type  string `json:"type"`
		// Demand      string  `json:"demand"`
		Limit       int     `json:"limit"`
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
		Limit:       body.Limit,
		Description: body.Description,
		Content:     body.Content,
		Reward:      body.Reward,
		Creator:     user.Uid,
		Status:      models.TaskStatusCreated,
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

func GetAssignList(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	assgnList := models.GetAssignmentListByUid(user.Uid)
	c.JSON(200, gin.H{
		"data": assgnList,
		"msg":  "OK",
	})
}

func AssignTask(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	id, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		c.AbortWithStatus(400)
	}
	task := models.GetTaskByID(id)
	if task.Assigned == task.Limit {
		c.AbortWithStatusJSON(403, gin.H{
			"msg": "任务已被认领完",
		})
		return
	}
	task.Assigned++
	task.Status = models.TaskStatusPending
	models.UpsertTask(task)
	models.UpsertAssignment(models.Assignment{
		TaskID:   task.ID,
		Assignee: user.Uid,
		Status:   models.AssignmentStatusPending,
	})
	c.JSON(200, gin.H{
		"data": gin.H{},
		"msg":  "OK",
	})
}
