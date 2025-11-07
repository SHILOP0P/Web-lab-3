package handlers_catalog

import (
	"backend/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CardItem struct {
	ID               int     `json:"id"`
	Name             string  `json:"name"`
	Price            float64 `json:"price"`
	ShortDescription string  `json:"short_description"`
	Image            string  `json:"image"`
}

// GET /catalog/cards?category=<alias>   (пример: "Строительные материалы")
func GetCatalogCards(c *gin.Context) {
	alias := c.Query("category")
	if alias == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category is required"})
		return
	}

	db, err := database.ConnectToServer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db connection failed"})
		return
	}
	defer db.Close()

	// берём по одному (первому) изображению на товар
	query := `
SELECT p.id,
       p.name,
       COALESCE(p.price, 0) AS price,
       COALESCE(p.short_description, '') AS short_description,
       COALESCE(pi.image, '') AS image
FROM product p
LEFT JOIN LATERAL (
    SELECT image
    FROM product_images
    WHERE product_id = p.id
    ORDER BY id ASC
    LIMIT 1
) pi ON true
WHERE p.alias = $1
ORDER BY p.id DESC;`

	rows, err := db.Query(query, alias)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "query failed"})
		return
	}
	defer rows.Close()

	items := make([]CardItem, 0, 32)
	for rows.Next() {
		var it CardItem
		if err := rows.Scan(&it.ID, &it.Name, &it.Price, &it.ShortDescription, &it.Image); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "scan failed"})
			return
		}
		items = append(items, it)
	}
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "rows error"})
		return
	}

	c.JSON(http.StatusOK, items)
}
