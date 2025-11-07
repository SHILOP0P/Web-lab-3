package auth

import (
	"backend/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	middleware.ClearSession(c)
	c.JSON(http.StatusOK, gin.H{"ok": true})
}
