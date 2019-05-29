package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/surplus-youyu/Youyu-se/models"
	"strconv"
)

func QuerySurveyHandler(c *gin.Context) {
	param := c.Param("sid")
	sid, err := strconv.ParseInt(param, 10, 32)
	msg := ""
	data := models.Survey{}
	if err != nil {
		c.JSON(400, gin.H{
			"status": false,
			"msg":    "invalid param",
		})
		return
	}
	result := models.GetSurveyById(int32(sid))
	if len(result) == 0 {
		msg = "survey not found"
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

func SurveyCreateHandler(c *gin.Context) {
	type ReqBody struct {
		Title   string `json:"title"`
		Content string `json:"content"`
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

	newSurvey := models.Survey{
		PublisherId: user.Uid,
		Title:       body.Title,
		Content:     body.Content,
	}
	models.CreateNewSurvey(newSurvey)
	c.JSON(200, gin.H{
		"status": true,
		"msg":    "success",
	})
}

func GetAllSurvey(c *gin.Context) {
	result := models.GetAllSurvey()
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
