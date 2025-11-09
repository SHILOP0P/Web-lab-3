package handlers_product

import (
	"backend/database"
	"backend/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// PUT /api/products/:id
// Обработка PUT-запроса для обновления товара
func UpdateProduct(c *gin.Context) {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil || productID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_id"})
		return
	}

	var payload models.Product
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_json"})
		return
	}

	// Проверим, что есть что обновлять
	if payload.Name == "" && payload.Price == 0 && payload.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no_fields"})
		return
	}

	db, err := database.ConnectToServer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db_connect"})
		return
	}
	defer db.Close()

	// Получаем текущие данные о товаре
	var currentProduct models.Product
	err = db.QueryRow(`
		SELECT id, name, alias, price, description, available, meta_title, meta_description, short_description, image
		FROM public.product WHERE id = $1`, productID).Scan(&currentProduct.ID, &currentProduct.Name, &currentProduct.Alias,
		&currentProduct.Price, &currentProduct.Description, &currentProduct.Available, &currentProduct.MetaTitle, &currentProduct.MetaDescription,
		&currentProduct.ShortDescription, &currentProduct.ProductImage.Image)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product_not_found"})
		return
	}

	// Если изображение не передано — оставляем текущее
	if payload.ProductImage.Image == "" {
		payload.ProductImage.Image = currentProduct.ProductImage.Image
	}

	// --- динамический UPDATE public.product ---
	set := make([]string, 0, 10)
	args := make([]interface{}, 0, 10)
	i := 1

	// Обновляем все остальные поля
	if payload.Name != "" {
		set = append(set, fmt.Sprintf("name = $%d", i))
		args = append(args, payload.Name)
		i++
	}
	if payload.Alias != "" {
		set = append(set, fmt.Sprintf("alias = $%d", i))
		args = append(args, payload.Alias)
		i++
	}
	if payload.Price != 0 {
		set = append(set, fmt.Sprintf("price = $%d", i))
		args = append(args, payload.Price)
		i++
	}
	// ... (т.е. аналогично для всех других полей)

	// Обновляем изображение, если оно было передано
	if payload.ProductImage.Image != "" {
		set = append(set, fmt.Sprintf("image = $%d", i))
		args = append(args, payload.ProductImage.Image)
		i++
	}

	// Строим запрос
	if len(set) > 0 {
		query := fmt.Sprintf("UPDATE public.product SET %s WHERE id = $%d", strings.Join(set, ", "), i)
		args = append(args, productID)

		// Выполняем запрос на обновление
		res, err := db.Exec(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "update_product"})
			return
		}
		if rows, _ := res.RowsAffected(); rows == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
			return
		}
	}

	// Обновляем характеристики товара, если они были переданы
	if payload.ProductProperty.Characteristics != "" {
		// Логика для обновления характеристик
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "updated",
		"id":      productID,
	})
}

