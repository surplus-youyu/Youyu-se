package controllers

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/surplus-youyu/Youyu-se/models"
	"github.com/surplus-youyu/Youyu-se/utils"
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
	c.JSON(200, gin.H{
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
		c.AbortWithStatusJSON(403, gin.H{
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
	}, user)

	files := form.File

	if len(files) != 0 {
		path := filePath + utils.IntToString(insertID)
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"msg":    "fail to upload files",
				"status": false,
			})
			return
		}

		for _, file := range files["files"] {
			err := c.SaveUploadedFile(file, path+"/"+file.Filename)
			if err != nil {
				c.AbortWithStatusJSON(500, gin.H{
					"msg":    "fail to upload files",
					"status": false,
				})
				return
			}
		}
	}

	c.JSON(200, gin.H{
		"data": gin.H{
			"id": insertID,
		},
		"msg": "OK",
	})
}

func GetTaskByID(c *gin.Context) {
	id := utils.StringToInt(c.Param("task_id"), c)
	task := models.GetTaskByID(id)

	files, err := ioutil.ReadDir(filePath + c.Param("task_id"))
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"msg":    "fail to get files name",
			"status": false,
		})
	}

	var filenames []string
	for _, file := range files {
		filenames = append(filenames, file.Name())
	}

	c.JSON(200, gin.H{
		"data": gin.H{
			"task":  task,
			"files": filenames,
		},
		"msg": "OK",
	})
}

func GetTaskFiles(c *gin.Context) {
	filename := c.Param("filename")

	path := filePath + c.Param("task_id") + "/" + filename
	//if  _, err := os.Stat(path); os.IsNotExist(err) {
	//	c.AbortWithStatusJSON(404, gin.H{
	//		"msg":    "cannot find file",
	//		"status": false,
	//	})
	//	return
	//}

	f, err := os.Open(path)
	if err != nil {
		c.AbortWithStatusJSON(404, gin.H{
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
	id := utils.StringToInt(c.Param("task_id"), c)
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
	id := utils.StringToInt(c.Param("assign_id"), c)

	var req Req
	err := c.BindJSON(&req)
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
	assign.Payload = req.Payload
	assign.Status = models.AssignmentStatusJudging
	models.UpsertAssignment(&assign)
	c.Status(204)
}

func GetAssignList(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	assignlist := models.GetAssignmentListByUid(user.Uid)
	c.JSON(200, gin.H{
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
	assign := &models.Assignment{
		TaskID:   task.ID,
		Assignee: user.Uid,
		Status:   models.AssignmentStatusPending,
	}
	models.UpsertAssignment(assign)
	c.JSON(200, gin.H{
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
		c.AbortWithStatus(400)
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
	c.JSON(204, gin.H{
		"msg":    "OK",
		"status": true,
	})
}

func FinishTask(c *gin.Context) {
	taskID := utils.StringToInt(c.Param("task_id"), c)
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
	id := utils.StringToInt(c.Param("task_id"), c)
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
	_ = json.Unmarshal([]byte(task.Content), &content)

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
			_ = json.Unmarshal([]byte(assign.Payload), &raw)
			answer := raw.(gin.H)
			options := answer["answer"].([]string)

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
