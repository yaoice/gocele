package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Ping(c *gin.Context) {
	s := sessions.Default(c)
	c.JSON(200, gin.H{
		"app_id":       s.Get("app_id"),
		"account_name": s.Get("account_name"),
		"display_name": s.Get("display_name"),
		"token":        s.Get("token"),
		"message":      "pong",
	})
}
