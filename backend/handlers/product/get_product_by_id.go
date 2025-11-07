package handlers_product

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"backend/database"
	"backend/models"
	"github.com/gin-gonic/gin"
)

func GetProductByID(c *gin.Context) {
	id := c.Param("id")

	db, err := database.ConnectToServer()
	if err != nil {
		log.Println("db_connect:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db_connect"})
		return
	}
	// Пул НЕ закрываем здесь

	var (
		p                               models.Product
		imgID, imgProductID             sql.NullInt64
		imgPath, imgTitle               sql.NullString
		propID, propProductID           sql.NullInt64
		propChars                       sql.NullString
	)

	// Один запрос: товар + первая картинка + первые свойства
	err = db.QueryRow(`
		SELECT
			p.id,
			p.manufacturer_id,
			COALESCE(p.name, ''),
			COALESCE(p.alias, ''),
			COALESCE(p.short_description, ''),
			COALESCE(p.description, ''),
			COALESCE(p.price, 0),
			COALESCE(p.available, 0),
			COALESCE(p.meta_keywords, ''),
			COALESCE(p.meta_description, ''),
			COALESCE(p.meta_title, ''),

			pi.id,
			pi.product_id,
			COALESCE(pi.image, ''),
			COALESCE(pi.title, ''),

			pp.id,
			pp.product_id,
			COALESCE(pp.characteristics, '')

		FROM public.product p
		LEFT JOIN LATERAL (
			SELECT id, product_id, image, title
			FROM public.product_images
			WHERE product_id = p.id
			ORDER BY id ASC
			LIMIT 1
		) AS pi ON TRUE
		LEFT JOIN LATERAL (
			SELECT id, product_id, characteristics
			FROM public.product_properties
			WHERE product_id = p.id
			ORDER BY id ASC
			LIMIT 1
		) AS pp ON TRUE
		WHERE p.id = $1
		LIMIT 1
	`, id).Scan(
		&p.ID,
		&p.ManufacturerID,
		&p.Name,
		&p.Alias,
		&p.ShortDescription,
		&p.Description,
		&p.Price,
		&p.Available,
		&p.MetaKeywords,
		&p.MetaDescription,
		&p.MetaTitle,

		&imgID,
		&imgProductID,
		&imgPath,
		&imgTitle,

		&propID,
		&propProductID,
		&propChars,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	if err != nil {
		log.Println("scan_product_row:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "scan_product"})
		return
	}

	// Вложенные объекты (с приведение типов int64 -> int)
	if imgID.Valid {
		img := models.ProductImage{
			ID:        int(imgID.Int64),
			ProductID: int(imgProductID.Int64),
			Image:     strings.TrimPrefix(imgPath.String, "images_db/"),
			Title:     imgTitle.String,
		}
		p.ProductImage = img
	}
	if propID.Valid {
		prop := models.ProductProperties{
			ID:              int(propID.Int64),
			ProductID:       int(propProductID.Int64),
			Characteristics: propChars.String,
		}
		p.ProductProperty = prop
	}

	c.JSON(http.StatusOK, p)
}
