package models

// ProductProperties — свойства товара, теперь одно поле для всех характеристик
type ProductProperties struct {
    ID            int    `json:"id"`
    ProductID     int    `json:"product_id"`    // ID продукта, ссылается на таблицу product
    Characteristics string `json:"characteristics,omitempty"` // Характеристики товара в одном поле
}
