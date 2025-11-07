package handlers_image

import (
	"backend/database"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

// DeleteImage - удаление изображения товара
func DeleteImage(c *gin.Context) {
	// Получаем ID изображения из URL
	imageID := c.Param("id")

	// Подключаемся к базе данных
	db, err := database.ConnectToServer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка подключения к базе данных"})
		return
	}
	defer db.Close()

	// Запрос для получения пути изображения из базы данных
	var imagePath string
	err = db.QueryRow("SELECT image FROM product_images WHERE id = $1", imageID).Scan(&imagePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении пути к изображению"})
		fmt.Println("Ошибка при получении пути изображения:", err)
		return
	}

	// Формируем полный путь к файлу
	filePath := filepath.Join("images_db", imagePath)

	// Удаляем изображение с сервера
	if err := os.Remove(filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении изображения с сервера"})
		return
	}

	// Удаляем запись из базы данных
	_, err = db.Exec("DELETE FROM product_images WHERE id = $1", imageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении записи изображения из базы данных"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Изображение успешно удалено",
	})
}
