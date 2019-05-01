package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/surplus-youyu/Youyu-se/models"
	"strconv"
)

func QuerySurveyHandler(c *gin.Context) {
	param := c.Param("sid")
	sid, err := strconv.ParseInt(param, 10, 32)
	statusCode := 200
	msg := ""
	data := models.Survey{}
	if err != nil {
		c.JSON(400, gin.H{
			"status": "fail",
			"msg": "invalid param",
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
	c.JSON(statusCode, gin.H{
		"status": "OK",
		"msg": msg,
		"data": data,
	})
}