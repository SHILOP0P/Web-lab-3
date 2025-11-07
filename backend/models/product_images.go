package models

// ProductImage — изображение товара.
type ProductImage struct {
	ID        int    `json:"id"`
	ProductID int    `json:"product_id"`
	Image     string `json:"image"` // имя файла или относительный путь (например, "cement.png")
	Title     string `json:"title"`
}
