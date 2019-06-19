package controllers

import (
	"encoding/json"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/surplus-youyu/Youyu-se/models"
)

func GetTaskList(c *gin.Context) {
	scope := c.Query("scope")
	var taskList []models.Task
	user := c.MustGet("user").(models.User)
	switch scope {
	case "owned":
		taskList = models.GetTaskListByCreator(user.Uid)
		break
	default:
		taskList = models.GetTaskList()
		break
	}
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

func GetAssignmentByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("assgn_id"))
	if err != nil {
		c.AbortWithStatus(400)
		return
	}
	user := c.MustGet("user").(models.User)
	assign := models.GetAssignmentByID(id)
	if assign.Assignee != user.Uid {
		c.AbortWithStatusJSON(403, gin.H{
			"msg": "权限不足",
		})
		return
	}
	c.JSON(200, gin.H{
		"data": assign,
		"msg":  "OK",
	})
}

func GetAssignListByTaskID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		c.AbortWithStatus(400)
		return
	}
	user := c.MustGet("user").(models.User)
	task := models.GetTaskByID(id)
	if task.Creator != user.Uid {
		c.AbortWithStatusJSON(403, gin.H{
			"msg": "权限不足",
		})
		return
	}
	assgnList := models.GetAssignmentListByTaskID(id)
	c.JSON(200, gin.H{
		"data": assgnList,
		"msg":  "OK",
	})
}

func SubmitAssign(c *gin.Context) {
	type Req struct {
		Payload string `json:"payload"`
	}
	id, err := strconv.Atoi(c.Param("assgn_id"))
	if err != nil {
		c.AbortWithStatus(400)
		return
	}
	var req Req
	err = c.BindJSON(&req)
	if err != nil {
		c.AbortWithStatus(400)
		return
	}
	user := c.MustGet("user").(models.User)
	assgn := models.GetAssignmentByID(id)
	if assgn.Assignee != user.Uid {
		c.AbortWithStatusJSON(403, gin.H{
			"msg": "权限不足",
		})
		return
	}
	assgn.Payload = req.Payload
	assgn.Status = models.AssignmentStatusJudging
	models.UpsertAssignment(&assgn)
	c.Status(204)
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
	type Req struct {
		TaskID int `json:"task_id"`
	}
	var req Req
	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithStatus(400)
		return
	}
	user := c.MustGet("user").(models.User)
	task := models.GetTaskByID(req.TaskID)
	if task.Creator == user.Uid {
		c.AbortWithStatusJSON(403, gin.H{
			"msg": "不能认领自己发布的任务",
		})
		return
	}
	if task.Assigned == task.Limit {
		c.AbortWithStatusJSON(403, gin.H{
			"msg": "任务已被认领完",
		})
		return
	}
	if task.Status == models.TaskStatusFinished {
		c.AbortWithStatusJSON(403, gin.H{
			"msg": "任务已结束",
		})
		return
	}
	task.Assigned++
	task.Status = models.TaskStatusPending
	models.UpsertTask(&task)
	assgn := &models.Assignment{
		TaskID:   task.ID,
		Assignee: user.Uid,
		Status:   models.AssignmentStatusPending,
	}
	models.UpsertAssignment(assgn)
	c.JSON(200, gin.H{
		"data": gin.H{
			"id": assgn.ID,
		},
		"msg": "OK",
	})
}

func JudgeAssignment(c *gin.Context) {
	type Req struct {
		Pass bool `json:"pass"`
	}
	var req Req
	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithStatus(400)
		return
	}
	assgnID, err := strconv.Atoi(c.Param("assgn_id"))
	if err != nil {
		c.AbortWithStatus(400)
		return
	}
	taskID, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		c.AbortWithStatus(400)
		return
	}
	user := c.MustGet("user").(models.User)
	task := models.GetTaskByID(taskID)
	assgn := models.GetAssignmentByID(assgnID)
	if assgn.TaskID != taskID || task.Creator != user.Uid {
		c.AbortWithStatusJSON(403, gin.H{
			"msg": "权限不足",
		})
		return
	}
	if req.Pass {
		assgn.Status = models.AssignmentStatusSuccess
	} else {
		assgn.Status = models.AssignmentStatusFailed
	}
	models.UpsertAssignment(&assgn)
	c.Status(204)
}

func FinishTask(c *gin.Context) {
	taskID, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		c.AbortWithStatus(400)
		return
	}
	user := c.MustGet("user").(models.User)
	task := models.GetTaskByID(taskID)
	if task.Creator != user.Uid {
		c.AbortWithStatusJSON(403, gin.H{
			"msg": "权限不足",
		})
		return
	}
	models.IssueTaskRewards(task)
	c.Status(204)
}

func GetSurveyStatistics(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("task_id"))
	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	user := c.MustGet("user").(models.User)
	task := models.GetTaskByID(id)
	if task.Creator != user.Uid {
		c.AbortWithStatusJSON(403, gin.H{
			"msg": "权限不足",
		})
		return
	}

	if task.Type != models.TaskTypeSurvey {
		c.JSON(400, gin.H{
			"status": false,
			"msg":    "非问卷无统计",
		})
	}

	var content []interface{}
	err = json.Unmarshal([]byte(task.Content), &content)

	data := []gin.H{}

	assignList := models.GetAssignmentListByTaskID(id)

	for _, q := range content {
		question := q.(gin.H)
		if question["type"].(int) == 3 {
			continue
		}

		data = append(data, gin.H{})
		index := len(data) - 1

		for j := range assignList {
			assign := assignList[j]

			var raw interface{}
			err = json.Unmarshal([]byte(assign.Payload), &raw)
			answer := raw.(gin.H)
			options := answer["options"].([]string)

			for _, op := range options {
				if val, ok := data[index][op]; ok {
					data[index][op] = val.(int) + 1
				} else {
					data[index][op] = 0
				}
			}
		}
	}

	c.JSON(200, gin.H{
		"status": true,
		"msg":    "OK",
		"data":   data,
	})
}
