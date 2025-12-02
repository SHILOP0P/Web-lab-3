package handlers_cart

import (
	"backend/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ClearCart(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	db, err := database.ConnectToServer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db_connection_failed"})
		return
	}
	defer db.Close()

	_, err = db.Exec(`
		DELETE FROM cart_items
		WHERE user_id = $1
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete_failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cart_cleared"})
}
