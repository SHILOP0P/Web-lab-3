package models

// CartItem — одна позиция в корзине (как лежит в таблице cart_items).
type CartItem struct {
	UserID    int64 `json:"userId"`
	ProductID int   `json:"productId"`
	Quantity  int   `json:"quantity"`
}

// CartItemResponse — то, что отдаём на фронт в GET /api/cart.
type CartItemResponse struct {
	ProductID int     `json:"productId"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Image     string  `json:"image"`
	Quantity  int     `json:"quantity"`
}
