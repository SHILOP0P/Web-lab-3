// handlers/auth/me.go
package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Me(c *gin.Context) {
	// сюда попадаем только если токен валиден (AuthRequired уже прошёл)
	uid, _ := c.Get("user_id")
	username, _ := c.Get("username")
	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       uid,
			"username": username,
			"role":     role,
		},
	})
}
