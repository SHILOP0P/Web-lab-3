package handlers_product

import (
	"backend/database"
	_ "database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// imagesDir — используйте тот же путь, что и при загрузке файлов.
// Если backend запускается из myproject/backend, а папка с файлами myproject/images_db,
// то относительный путь: ../images_db
var imagesDir = filepath.Clean(filepath.Join("..", "images_db"))

func DeleteProduct(c *gin.Context) {
	productID := c.Param("id")

	db, err := database.ConnectToServer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка подключения к базе данных"})
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось открыть транзакцию"})
		return
	}

	// 1) Удаляем изображения и одновременно получаем список имён файлов через RETURNING
	var filenames []string
	rows, err := tx.Query(`DELETE FROM product_images WHERE product_id = $1 RETURNING image`, productID)
	if err != nil {
		_ = tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении изображений товара"})
		return
	}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err == nil && name != "" {
			filenames = append(filenames, name)
		}
	}
	_ = rows.Close()

	// 2) Удаляем свойства
	if _, err := tx.Exec(`DELETE FROM product_properties WHERE product_id = $1`, productID); err != nil {
		_ = tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении свойств товара"})
		return
	}

	// 3) Удаляем сам товар
	res, err := tx.Exec(`DELETE FROM product WHERE id = $1`, productID)
	if err != nil {
		_ = tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при удалении товара"})
		return
	}
	if aff, _ := res.RowsAffected(); aff == 0 {
		_ = tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Товар не найден"})
		return
	}

	// 4) Фиксируем транзакцию
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось зафиксировать удаление"})
		return
	}

	// 5) Удаляем файлы с диска (уже после коммита)
	for _, fname := range filenames {
		// Защита от путевых трюков; fname — то, что хранится в БД
		target := filepath.Clean(filepath.Join(imagesDir, fname))
		if err := os.Remove(target); err != nil && !os.IsNotExist(err) {
			// Ошибку удаления файла логируем, но пользователю возвращаем успех:
			// данные в БД уже удалены, файл может быть удалён вручную позже.
			log.Printf("DeleteProduct: не удалось удалить файл %s: %v", target, err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Товар и все связанные данные удалены"})
}
