package main

//import "github.com/gin-gonic/gin"
//
//func main() {
//	r := gin.Default()
//	r.GET("/ping", func(c *gin.Context) {
//		c.JSON(200, gin.H{
//			"message": "pong",
//		})
//	})
//	r.Run() // listen and serve on 0.0.0.0:8080
//}
import (
	//"fmt"
	"github.com/surplus-youyu/Youyu-se/models"
)

func main() {
	models.GetUserByUid("123")
	//fmt.Println(result)
}