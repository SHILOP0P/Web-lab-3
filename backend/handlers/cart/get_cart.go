package handlers_cart

import (
	"backend/database"
	"github.com/gin-gonic/gin"
	"net/http"
	"backend/models"
)

func GetCart(c *gin.Context) {
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

	rows, err := db.Query(`
		SELECT
		    p.id,
		    p.name,
		    p.price,
		    COALESCE(pi.image, ''),
		    ci.quantity
		FROM cart_items AS ci
		JOIN product AS p ON p.id = ci.product_id
		LEFT JOIN product_images AS pi ON pi.product_id = p.id
		WHERE ci.user_id = $1
	`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query_failed"})
		return
	}
	defer rows.Close()

	var items []models.CartItemResponse
	for rows.Next() {
		var it models.CartItemResponse
		if err := rows.Scan(&it.ProductID, &it.Name, &it.Price, &it.Image, &it.Quantity); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan_failed"})
			return
		}
		items = append(items, it)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "rows_error"})
		return
	}

	c.JSON(http.StatusOK, items)
}
