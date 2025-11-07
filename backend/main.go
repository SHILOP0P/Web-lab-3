package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	handlers_catalog "backend/handlers/catalog"
	handlers_image "backend/handlers/image_upload"
	handlers_product "backend/handlers/product"
	handlers_reviews "backend/handlers/reviews"

	"backend/handlers/auth"
	"backend/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func pickDir(dir string) string {
    if fi, err := os.Stat(dir); err == nil && fi.IsDir() {
        return dir
    }
    up := filepath.Join("..", dir)
    if fi, err := os.Stat(up); err == nil && fi.IsDir() {
        return up
    }
    // последняя попытка — вернуть исходный (пусть упадёт заметно)
    return dir
}

func main() {
	gin.SetMode(gin.DebugMode)
	r := gin.Default()
	r.SetTrustedProxies(nil)

	// ============================================================
	// CORS (можно выключить, если фронт всегда открывается с :8080)
	// ============================================================
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))


	frontDir := pickDir("front")
	imagesDir := pickDir("images_db")

	// ============================================================
	// СТАТИКА ФРОНТА (всё с того же origin :8080)
	// Директория фронта — ./front (HTML/CSS/JS/картинки)
	// Картинки товаров — ./images_db
	// ============================================================
	// --- Статика ---
	r.Static("/front", frontDir)
	r.Static("/images_db", imagesDir)

	// удобные алиасы (если нужны)
	r.Static("/js", filepath.Join(frontDir, "js"))
	r.Static("/images", filepath.Join(frontDir, "images"))
	r.Static("/admin", filepath.Join(frontDir, "admin"))
	r.Static("/products", filepath.Join(frontDir, "products"))

	// отдать конкретные файлы
	r.StaticFile("/index.html", filepath.Join(frontDir, "index.html"))
	r.StaticFile("/style.css", filepath.Join(frontDir, "style.css"))
	
	r.StaticFile("/catalog.html",  filepath.Join(frontDir, "catalog.html"))
	r.StaticFile("/register.html", filepath.Join(frontDir, "register.html"))
	// при желании в таком же стиле:
	r.StaticFile("/about.html",    filepath.Join(frontDir, "about.html"))
	r.StaticFile("/contacts.html", filepath.Join(frontDir, "contacts.html"))

	// ЯВНО: корень сайта → index.html
	r.GET("/", func(c *gin.Context) {
		c.File(filepath.Join(frontDir, "index.html"))
	})

	

	// ============================================================
	// BACKEND API
	// ============================================================
	api := r.Group("/api")
	{
		// catalog
		api.GET("/catalog/cards", handlers_catalog.GetCatalogCards)

		// products CRUD + by id
		api.GET("/products", handlers_product.GetProducts)
		api.POST("/products", handlers_product.AddProduct)
		api.PUT("/products/:id", handlers_product.UpdateProduct)
		api.DELETE("/products/:id", handlers_product.DeleteProduct)
		api.GET("/products/:id", handlers_product.GetProductByID)

		api.GET("/reviews", handlers_reviews.List)
    	api.POST("/reviews", middleware.AuthRequired(), handlers_reviews.Create)
    	api.DELETE("/reviews/:id", middleware.AuthRequired(), handlers_reviews.Delete)

		// upload
		api.POST("/upload-image", handlers_image.UploadImage)

		// auth
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", auth.Register)
			authGroup.POST("/login", auth.Login)
			authGroup.POST("/logout", auth.Logout)
			authGroup.GET("/me", middleware.AuthRequired(), auth.Me)
		}
	}

	// ============================================================
	// Fallback для прямых ссылок на фронт
	// - Не перехватываем /api/*
	// - Если файл существует в ./front — отдаем его
	// - Иначе показываем index.html
	// ============================================================
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path

		// /api/* не трогаем
		if strings.HasPrefix(path, "/api/") {
			c.Status(http.StatusNotFound)
			return
		}

		// пробуем отдать реальный файл из frontDir
		rel := strings.TrimPrefix(filepath.Clean(path), "/")
		if rel != "" && !strings.HasSuffix(rel, "/") {
			tryFile := filepath.Join(frontDir, rel)
			if fi, err := os.Stat(tryFile); err == nil && !fi.IsDir() {
				c.File(tryFile)
				return
			}
    }

    // иначе главная
    c.File(filepath.Join(frontDir, "index.html"))
	})





	// ============================================================
	// Стартуем на :8080
	// ============================================================
	if err := r.Run(":8080"); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
