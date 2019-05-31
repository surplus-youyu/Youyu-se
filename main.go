package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/surplus-youyu/Youyu-se/configs"
	"github.com/surplus-youyu/Youyu-se/routes"
)

func main() {
	r := gin.Default()
	store := cookie.NewStore([]byte(configs.SessionKey))
	r.Use(sessions.Sessions("youyu-session", store))

	route.Route(r)
	err := r.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		panic(err)
	}
}
