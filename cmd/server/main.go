package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johan-ag/wishlist/cmd/server/handler"
	"github.com/johan-ag/wishlist/internal/authentication"
	"github.com/johan-ag/wishlist/internal/users"
	"github.com/johan-ag/wishlist/pkg/store"
	"gorm.io/gorm"
)

var (
	db *gorm.DB = store.DBConnection()

	usersRepository users.Repository = users.NewUsersRepository(db)

	usersService users.Service              = users.NewUsersService(usersRepository)
	authService  authentication.AuthService = authentication.NewAuthService(usersService)

	authHandler handler.Auth = handler.NewAuthHandler(authService)
)

func main() {
	app := gin.Default()
	app.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "Pong")
	})
	auth := app.Group("/api/auth")
	{
		auth.POST("/sign-in", func(c *gin.Context) {

		}, 
		authHandler.Login())
		auth.POST("/sign-up", authHandler.Register())
	}
	app.Run()
}
