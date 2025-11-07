package handlers_product

import (
	"backend/database"
	"backend/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"encoding/json"
	"path/filepath"
	"time"
)

// Добавление нового продукта с изображениями и характеристиками
func AddProduct(c *gin.Context) {
	// Получаем данные из формы
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при получении формы"})
		return
	}

	// Извлекаем JSON данные о продукте
	productData := form.Value["product"][0]
	var product models.Product
	err = json.Unmarshal([]byte(productData), &product)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при обработке данных о продукте"})
		return
	}

	// Логируем данные формы для отладки
	fmt.Println(form.Value)

	// Подключаемся к базе данных
	db, err := database.ConnectToServer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка подключения к базе данных"})
		return
	}
	defer db.Close()

	// Вставляем данные о продукте в таблицу product
	query := `
		INSERT INTO product (manufacturer_id, name, alias, price, description, available, meta_keywords, meta_description, meta_title, short_description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`
	var productID int
	err = db.QueryRow(query, product.ManufacturerID, product.Name, product.Alias, product.Price, product.Description, product.Available, product.MetaKeywords, product.MetaDescription, product.MetaTitle, product.ShortDescription).Scan(&productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении продукта"})
		return
	}

	// Загружаем изображения
	files := form.File["images[]"]
	for _, file := range files {
		// Генерация имени файла
		fileExtension := filepath.Ext(file.Filename)
		fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileExtension)

		// Путь для сохранения изображения
		filePath := filepath.Join("images_db", fileName)

		// Сохраняем файл на сервере
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении файла"})
			return
		}

		// Вставляем запись в таблицу product_images
		imageQuery := `
			INSERT INTO product_images (product_id, image, title)
			VALUES ($1, $2, $3)
		`
		_, err = db.Exec(imageQuery, productID, fileName, file.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении изображения"})
			return
		}
	}

	// Получаем характеристики товара из формы
	characteristics := ""
	if len(form.Value["characteristics"]) > 0 {
		characteristics = form.Value["characteristics"][0]
	}

	// Вставляем характеристики товара в таблицу product_properties
	_, err = db.Exec(`
		INSERT INTO product_properties (product_id, characteristics)
		VALUES ($1, $2)`,
		productID, characteristics)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении характеристик"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Продукт и изображения успешно добавлены",
		"id":      productID,
	})
}
