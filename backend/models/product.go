package models

// Product — товар.
type Product struct {
	ID               int                `json:"id"`
	ManufacturerID   int                `json:"manufacturer_id"`
	Name             string             `json:"name"`
	Alias            string             `json:"alias"`
	ShortDescription string             `json:"short_description"`
	Description      string             `json:"description"`
	Price            float64            `json:"price"`
	Available        int                `json:"available"`
	MetaKeywords     string             `json:"meta_keywords"`
	MetaDescription  string             `json:"meta_description"`
	MetaTitle        string             `json:"meta_title"`

	// Связанные коллекции (для удобной отдачи на фронт)
	ProductImage     ProductImage      `json:"product_images"`
	ProductProperty  ProductProperties `json:"product_properties"`
}
