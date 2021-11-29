package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/johan-ag/wishlist/internal/authentication"
	"github.com/johan-ag/wishlist/internal/users"
)

type Auth interface {
	Login() gin.HandlerFunc
	Register() gin.HandlerFunc
}

func NewAuthHandler(authService authentication.AuthService) Auth {
	return &auth{
		authService,
	}
}

type auth struct {
	authService authentication.AuthService
}

func (h *auth) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user users.User
		err := c.Bind(&user)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		token, err := h.authService.Login(user)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	}
}

func (h *auth) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user users.User
		err := c.Bind(&user)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		_, err = h.authService.Register(user)
		if err != nil {
			if err.Error() == "This user already existed" {
				c.AbortWithStatus(http.StatusConflict)
				return
			}
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		c.Status(http.StatusCreated)
	}
}
