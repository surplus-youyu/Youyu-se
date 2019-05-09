package main

import (
	"github.com/gin-gonic/gin"
	"github.com/surplus-youyu/Youyu-se/route"
)

func main() {
	r := gin.Default()
	route.Route(r)
	err := r.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		panic(err)
	}
}
