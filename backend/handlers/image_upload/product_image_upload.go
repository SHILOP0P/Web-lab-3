package handlers_image

import (
	"backend/database"
	_ "backend/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"time"
	"strconv"
)

// UploadImage - загрузка изображения для товара
func UploadImage(c *gin.Context) {
	// Получаем файл изображения
	file, _ := c.FormFile("image")
	if file == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не найден файл изображения"})
		return
	}

	// Генерация уникального имени для изображения
	fileExtension := filepath.Ext(file.Filename)
	// Уникальное имя для файла (с добавлением времени)
	fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileExtension)

	// Путь для сохранения изображения
	filePath := filepath.Join("images_db", fileName)

	// Сохраняем файл на сервер
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении файла"})
		return
	}

	// Получаем ID товара, к которому нужно привязать изображение
	productID := c.DefaultQuery("product_id", "0") // Получаем product_id из query параметра
	id, err := strconv.Atoi(productID)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ID товара"})
		return
	}

	// Название изображения теперь берем из исходного имени файла
	title := file.Filename // Название файла

	// Подключаемся к базе данных
	db, err := database.ConnectToServer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка подключения к базе данных"})
		return
	}
	defer db.Close()

	// Запрос для вставки данных в таблицу product_images
	query := `
		INSERT INTO product_images (product_id, image, title)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var imageID int
	err = db.QueryRow(query, id, fileName, title).Scan(&imageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении изображения"})
		fmt.Println("Ошибка при вставке:", err)
		return
	}

	// Возвращаем URL изображения
	c.JSON(http.StatusOK, gin.H{
		"message":   "Изображение успешно загружено",
		"image_url": fmt.Sprintf("http://localhost:8080/images_db/%s", fileName),
	})
}
