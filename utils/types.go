package utils

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func StringToFloat32(str string, c *gin.Context) float32 {
	f, err := strconv.ParseFloat(str, 32)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"mag":    "invalid param",
			"status": false,
		})
		panic(err)
	}
	return float32(f)
}

func StringToInt(str string, c *gin.Context) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{
			"mag":    "invalid param",
			"status": false,
		})
		panic(err)
	}
	return i
}

func IntToString(i int) string {
	str := strconv.Itoa(i)
	return str
}
