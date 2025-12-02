package handlers_cart

import (
	"backend/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func DeleteCartItem(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	productIDStr := c.Param("productID")
	productID, err := strconv.Atoi(productIDStr)
	if err != nil || productID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_productID"})
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
		WHERE user_id = $1 AND product_id = $2
	`, userID, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete_failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item_deleted"})
}
