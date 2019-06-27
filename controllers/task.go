package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/surplus-youyu/Youyu-se/models"
	"github.com/surplus-youyu/Youyu-se/utils"
	"io"
	"net/http"
	"os"
)

const filePath = "static/files/"

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
	// 200
	c.JSON(http.StatusOK, gin.H{
		"data": taskList,
		"msg":  "OK",
	})
}

func CreateTask(c *gin.Context) {
	user := c.MustGet("user").(models.User)

	form, err := c.MultipartForm()
	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	values := form.Value
	reward := utils.StringToFloat32(values["reward"][0], c)
	if user.Balance < reward {
		// 403
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"msg": "余额不足",
		})
		return
	}
	insertID := models.CreateTask(models.Task{
		Title:       values["title"][0],
		Type:        values["type"][0],
		Limit:       utils.StringToInt(values["limit"][0], c),
		Description: values["description"][0],
		Content:     values["content"][0],
		Reward:      reward,
		Creator:     user.Uid,
		Status:      models.TaskStatusCreated,
		Files:       "",
	}, user)

	files := form.File

	if len(files) != 0 {
		task := models.GetTaskByID(insertID)

		path := filePath + utils.IntToString(insertID)
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			// 500
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"msg":    "fail to upload files",
				"status": false,
			})
			return
		}

		for _, file := range files["files"] {
			err := c.SaveUploadedFile(file, path+"/"+file.Filename)
			if err != nil {
				// 500
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"msg":    "fail to upload files",
					"status": false,
				})
				return
			}
			task.Files += file.Filename + "/"
		}
		task.Files = task.Files[0 : len(task.Files)-1]
		models.DB.Save(&task)
	}

	// 200
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id": insertID,
		},
		"msg": "OK",
	})
}

func GetTaskByID(c *gin.Context) {
	id := utils.StringToInt(c.Param("task_id"), c)
	task := models.GetTaskByID(id)

	// 200
	c.JSON(http.StatusOK, gin.H{
		"data": task,
		"msg":  "OK",
	})
}

func GetTaskFiles(c *gin.Context) {
	filename := c.Param("filename")

	path := filePath + c.Param("task_id") + "/" + filename
	f, err := os.Open(path)
	if err != nil {
		// 404
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"msg":    "cannot find file",
			"status": false,
		})
		return
	}
	defer f.Close()

	w := c.Writer
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", c.Request.Header.Get("Content-Type"))

	_, _ = io.Copy(w, f)
}

func GetAssignmentByID(c *gin.Context) {
	id := utils.StringToInt(c.Param("assign_id"), c)
	user := c.MustGet("user").(models.User)
	assign := models.GetAssignmentByID(id)
	if assign.Assignee != user.Uid {
		// 403
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"msg": "权限不足",
		})
		return
	}
	// 200
	c.JSON(http.StatusOK, gin.H{
		"data": assign,
		"msg":  "OK",
	})
}

func GetAssignListByTaskID(c *gin.Context) {
	id := utils.StringToInt(c.Param("task_id"), c)
	user := c.MustGet("user").(models.User)
	task := models.GetTaskByID(id)
	if task.Creator != user.Uid {
		// 403
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"msg": "权限不足",
		})
		return
	}
	assgnList := models.GetAssignmentListByTaskID(id)
	// 200
	c.JSON(http.StatusOK, gin.H{
		"data": assgnList,
		"msg":  "OK",
	})
}

func SubmitAssign(c *gin.Context) {
	type Req struct {
		Payload string `json:"payload"`
	}
	id := utils.StringToInt(c.Param("assign_id"), c)

	var req Req
	err := c.BindJSON(&req)
	if err != nil {
		// 400
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	user := c.MustGet("user").(models.User)
	assign := models.GetAssignmentByID(id)
	if assign.Assignee != user.Uid {
		// 403
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"msg": "权限不足",
		})
		return
	}
	assign.Payload = req.Payload
	assign.Status = models.AssignmentStatusJudging
	models.UpsertAssignment(&assign)
	// 204
	c.Status(http.StatusNoContent)
}

func GetAssignList(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	assignlist := models.GetAssignmentListByUid(user.Uid)
	// 200
	c.JSON(http.StatusOK, gin.H{
		"data": assignlist,
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
		// 400
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	user := c.MustGet("user").(models.User)
	task := models.GetTaskByID(req.TaskID)
	if task.Creator == user.Uid {
		// 403
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"msg": "不能认领自己发布的任务",
		})
		return
	}
	if task.Assigned == task.Limit {
		// 403
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"msg": "任务已被认领完",
		})
		return
	}
	if task.Status == models.TaskStatusFinished {
		// 403
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"msg": "任务已结束",
		})
		return
	}
	task.Assigned++
	task.Status = models.TaskStatusPending
	models.UpsertTask(&task)
	assign := &models.Assignment{
		TaskID:   task.ID,
		Assignee: user.Uid,
		Status:   models.AssignmentStatusPending,
	}
	models.UpsertAssignment(assign)
	// 200
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id": assign.ID,
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
		// 400
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	assignID := utils.StringToInt(c.Param("assign_id"), c)
	taskID := utils.StringToInt(c.Param("task_id"), c)

	user := c.MustGet("user").(models.User)
	task := models.GetTaskByID(taskID)
	assign := models.GetAssignmentByID(assignID)
	if assign.TaskID != taskID || task.Creator != user.Uid {
		c.AbortWithStatusJSON(403, gin.H{
			"msg": "权限不足",
		})
		return
	}
	if req.Pass {
		assign.Status = models.AssignmentStatusSuccess
	} else {
		assign.Status = models.AssignmentStatusFailed
	}
	models.UpsertAssignment(&assign)
	// 204
	c.JSON(http.StatusNoContent, gin.H{
		"msg":    "OK",
		"status": true,
	})
}

func FinishTask(c *gin.Context) {
	taskID := utils.StringToInt(c.Param("task_id"), c)
	user := c.MustGet("user").(models.User)
	task := models.GetTaskByID(taskID)
	if task.Creator != user.Uid {
		// 403
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"msg": "权限不足",
		})
		return
	}
	models.IssueTaskRewards(task)
	// 204
	c.Status(http.StatusNoContent)
}

func GetSurveyStatistics(c *gin.Context) {
	id := utils.StringToInt(c.Param("task_id"), c)
	user := c.MustGet("user").(models.User)
	task := models.GetTaskByID(id)
	if task.Creator != user.Uid {
		// 403
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"msg": "权限不足",
		})
		return
	}

	if task.Type != models.TaskTypeSurvey {
		// 400
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"msg":    "非问卷无统计",
		})
	}

	type Item struct {
		Title    string   `json:"title"`
		Type     int      `json:"type"`
		Options  []string `json:"options"`
		Optional bool     `json:"optional"`
		Limit    int      `json:"limit"`
		Answer   []string `json:"answer"`
	}
	var content []Item
	_ = json.Unmarshal([]byte(task.Content), &content)

	var data []map[string]interface{}

	assignList := models.GetAssignmentListByTaskID(id)

	for _, question := range content {
		if question.Type == 3 {
			continue
		}

		data = append(data, map[string]interface{}{})
		index := len(data) - 1

		data[index]["title"] = question.Title
		for _, opt := range question.Options {
			data[index][opt] = 0
		}

		for j := range assignList {
			assign := assignList[j]

			var items []Item
			_ = json.Unmarshal([]byte(assign.Payload), &items)
			item := items[index]
			for _, answer := range item.Answer {
				val := data[index][answer]
				data[index][answer] = val.(int) + 1
			}
		}
	}

	// 200
	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"msg":    "OK",
		"data":   data,
	})
}
