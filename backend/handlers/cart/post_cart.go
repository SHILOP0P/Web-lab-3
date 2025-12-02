package handlers_cart

import (
	"backend/database"
	"backend/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddToCart(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var item models.CartItem
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_json"})
		return
	}

	// userId никогда не берём из тела запроса — только из авторизации
	item.UserID = userID

	if item.ProductID <= 0 || item.Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "productId_and_quantity_must_be_positive"})
		return
	}

	db, err := database.ConnectToServer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db_connection_failed"})
		return
	}
	defer db.Close()

	_, err = db.Exec(`
		INSERT INTO cart_items (user_id, product_id, quantity)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, product_id)
		DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity
	`, item.UserID, item.ProductID, item.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "insert_or_update_failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cart_updated"})
}
