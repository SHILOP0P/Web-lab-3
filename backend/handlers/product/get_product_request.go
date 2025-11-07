package handlers_product

import (
	"database/sql"
	"net/http"
	"strings"

	"backend/database"
	"backend/models"
	"github.com/gin-gonic/gin"

)

// GetProducts — вернуть список продуктов с (опциональными) главной картинкой и свойствами.
// Все чтения из БД сделаны NULL-safe: если какая-то колонка = NULL, сервер не падает.
func GetProducts(c *gin.Context) {
	db, err := database.ConnectToServer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка подключения к базе данных"})
		return
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT
			id,
			manufacturer_id,
			name,
			alias,
			short_description,
			description,
			price,
			available,
			meta_keywords,
			meta_description,
			meta_title
		FROM product
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении данных о продуктах"})
		return
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		// Читаем все потенциально NULL поля через sql.Null*
		var (
			id              sql.NullInt64
			manufacturerID  sql.NullInt64
			name            sql.NullString
			alias           sql.NullString
			shortDesc       sql.NullString
			desc            sql.NullString
			price           sql.NullFloat64
			available       sql.NullInt64
			metaKeywords    sql.NullString
			metaDescription sql.NullString
			metaTitle       sql.NullString
		)

		if err := rows.Scan(
			&id,
			&manufacturerID,
			&name,
			&alias,
			&shortDesc,
			&desc,
			&price,
			&available,
			&metaKeywords,
			&metaDescription,
			&metaTitle,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при чтении данных о продукте"})
			return
		}

		// Собираем продукт в модель (для NULL берём «пустые» значения)
		p := models.Product{
			ID:               int(id.Int64),
			ManufacturerID:   int(manufacturerID.Int64),
			Name:             name.String,
			Alias:            alias.String,
			ShortDescription: shortDesc.String,
			Description:      desc.String,
			Price:            price.Float64,
			Available:        int(available.Int64),
			MetaKeywords:     metaKeywords.String,
			MetaDescription:  metaDescription.String,
			MetaTitle:        metaTitle.String,
		}

		// --- Главная картинка (NULL-safe) ---
		{
			var (
				imgID    sql.NullInt64
				imgPID   sql.NullInt64
				imgPath  sql.NullString
				imgTitle sql.NullString
			)
			err = db.QueryRow(`
				SELECT id, product_id, image, title
				FROM product_images
				WHERE product_id = $1
				ORDER BY id ASC
				LIMIT 1
			`, p.ID).Scan(&imgID, &imgPID, &imgPath, &imgTitle)

			if err == nil {
				p.ProductImage = models.ProductImage{
					ID:        int(imgID.Int64),
					ProductID: int(imgPID.Int64),
					Image:     strings.TrimPrefix(imgPath.String, "images_db/"),
					Title:     imgTitle.String,
				}
			} else if err != sql.ErrNoRows {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении изображения"})
				return
			}
		}

		// --- Свойства (NULL-safe) ---
		{
			var (
				propID   sql.NullInt64
				propPID  sql.NullInt64
				propChar sql.NullString
			)
			err = db.QueryRow(`
				SELECT id, product_id, characteristics
				FROM product_properties
				WHERE product_id = $1
				ORDER BY id ASC
				LIMIT 1
			`, p.ID).Scan(&propID, &propPID, &propChar)

			if err == nil {
				p.ProductProperty = models.ProductProperties{
					ID:              int(propID.Int64),
					ProductID:       int(propPID.Int64),
					Characteristics: propChar.String,
				}
			} else if err != sql.ErrNoRows {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении характеристик"})
				return
			}
		}

		products = append(products, p)
	}

	// На всякий случай обработаем возможную ошибку итерации
	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при чтении списка продуктов"})
		return
	}

	c.JSON(http.StatusOK, products)
}
