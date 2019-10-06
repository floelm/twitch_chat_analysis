package main

import (
	"github.com/gin-gonic/gin"
	"twitch_chat_analysis/cmd/consumer/domain"
)

func main() {
	r := gin.Default()

	cache := domain.NewTermsCache()

	r.GET("/terms", func(c *gin.Context) {
		terms := cache.GetAll()
		println(terms)
		c.JSON(200, terms)
	})
	r.Run()
}
