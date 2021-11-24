package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johan-ag/wishlist/pkg/store"
	"gorm.io/gorm"
)

var (
	db *gorm.DB = store.DBConnection()
)

func main() {
	app := gin.Default()
	app.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "Pong")
	})

	app.Run()
}
