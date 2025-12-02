package handlers_cart

import "github.com/gin-gonic/gin"

// getUserID достаёт user_id, который положил middleware.AuthRequired().
func getUserID(c *gin.Context) (int64, bool) {
	v, ok := c.Get("user_id")
	if !ok {
		return 0, false
	}
	id, ok := v.(int64)
	if !ok || id <= 0 {
		return 0, false
	}
	return id, true
}
