package utils

import (
	"github.com/gin-gonic/gin"
)

type Error struct {
	Status int
	Msg    string
	Err    error
}

func HandleError(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			if e, ok := err.(Error); ok {
				c.AbortWithStatusJSON(e.Status, gin.H{
					"msg": e.Msg,
				})
				panic(e.Err)
			} else {
				panic(e)
			}
		}
	}()
	c.Next()
}
