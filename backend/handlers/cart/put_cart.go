package handlers_cart

import (
	"backend/database"
	"backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func UpdateCartItem(c *gin.Context) {
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

	var item models.CartItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_json"})
		return
	}

	item.UserID = userID
	item.ProductID = productID // берём id из URL, а не из тела

	db, err := database.ConnectToServer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db_connection_failed"})
		return
	}
	defer db.Close()

	if item.Quantity <= 0 {
		// quantity <= 0 — считаем, что нужно удалить позицию
		_, err = db.Exec(`
			DELETE FROM cart_items
			WHERE user_id = $1 AND product_id = $2
		`, item.UserID, item.ProductID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "delete_failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "item_removed"})
		return
	}

	_, err = db.Exec(`
		UPDATE cart_items
		SET quantity = $3
		WHERE user_id = $1 AND product_id = $2
	`, item.UserID, item.ProductID, item.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update_failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item_updated"})
}
