package handlers_product

import (
	"backend/database"
	"backend/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// UpdateProduct - обновление данных о продукте
func UpdateProduct(c *gin.Context) {
	// Извлекаем ID продукта из URL
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil || productID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID продукта"})
		return
	}

	var product models.Product

	// Получаем данные из тела запроса
	if err := c.BindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные"})
		return
	}

	// Подключаемся к базе данных
	db, err := database.ConnectToServer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка подключения к базе данных"})
		return
	}
	defer db.Close()

	// Формируем динамический SQL запрос для обновления
	query := "UPDATE product SET "
	var args []interface{}
	i := 1

	// Проверяем и добавляем только те поля, которые передаются в запрос
	if product.ManufacturerID != 0 {
		query += fmt.Sprintf("manufacturer_id = $%d, ", i)
		args = append(args, product.ManufacturerID)
		i++
	}
	if product.Name != "" {
		query += fmt.Sprintf("name = $%d, ", i)
		args = append(args, product.Name)
		i++
	}
	if product.Alias != "" {
		query += fmt.Sprintf("alias = $%d, ", i)
		args = append(args, product.Alias)
		i++
	}
	if product.Price != 0 {
		query += fmt.Sprintf("price = $%d, ", i)
		args = append(args, product.Price)
		i++
	}
	if product.Description != "" {
		query += fmt.Sprintf("description = $%d, ", i)
		args = append(args, product.Description)
		i++
	}
	if product.Available != 0 {
		query += fmt.Sprintf("available = $%d, ", i)
		args = append(args, product.Available)
		i++
	}
	if product.MetaKeywords != "" {
		query += fmt.Sprintf("meta_keywords = $%d, ", i)
		args = append(args, product.MetaKeywords)
		i++
	}
	if product.MetaDescription != "" {
		query += fmt.Sprintf("meta_description = $%d, ", i)
		args = append(args, product.MetaDescription)
		i++
	}
	if product.MetaTitle != "" {
		query += fmt.Sprintf("meta_title = $%d, ", i)
		args = append(args, product.MetaTitle)
		i++
	}
	if product.ShortDescription != "" {
		query += fmt.Sprintf("short_description = $%d, ", i)
		args = append(args, product.ShortDescription)
		i++
	}

	// Удаляем последнюю запятую
	query = query[:len(query)-2]

	// Добавляем условие WHERE для обновления по ID
	query += fmt.Sprintf(" WHERE id = $%d RETURNING id", i)
	args = append(args, productID)

	// Выполняем запрос на обновление
	var id int
	err = db.QueryRow(query, args...).Scan(&id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении продукта"})
		fmt.Println("Ошибка при обновлении:", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Продукт успешно обновлен",
		"id":      id,
	})
}
