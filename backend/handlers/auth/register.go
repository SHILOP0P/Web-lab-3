package auth

import (
	"backend/database"
	"backend/utils"
	"backend/middleware"
	_ "database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var passwordPepper = getEnvOr("PASSWORD_PEPPER", "+Pepper#2025!")

type registerReq struct {
	Username  string  `json:"username"`
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	FirstName *string `json:"firstName"`
	LastName  *string `json:"lastName"`
	Phone     *string `json:"phone"`
	Gender    *string `json:"gender"`
	Birthdate *string `json:"birthdate"` // "YYYY-MM-DD"
	Region    *string `json:"region"`
}

func Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error":"bad_json"}); return
	}
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	if req.Username=="" || req.Email=="" || req.Password=="" {
		c.JSON(http.StatusBadRequest, gin.H{"error":"missing_fields"}); return
	}

	db, err := database.ConnectToServer()
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"db_connect"}); return }
	defer db.Close()

	// проверка уникальности
	var exists int
	if err := db.QueryRow(`select 1 from users where username=$1 or email=$2 limit 1`, req.Username, req.Email).Scan(&exists); err == nil {
		c.JSON(http.StatusConflict, gin.H{"error":"user_exists"}); return
	}

	hash, err := utils.HashPassword(req.Password, passwordPepper)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"hash_error"}); return }

	var birth *time.Time
	if req.Birthdate != nil && *req.Birthdate != "" {
		if t, err := time.Parse("2006-01-02", *req.Birthdate); err == nil {
			birth = &t
		}
	}

	// вставка
	var id int64
	err = db.QueryRow(`
		insert into users (username,email,password_hash,role,first_name,last_name,phone,gender,birthdate,region)
		values ($1,$2,$3,'USER',$4,$5,$6,$7,$8,$9)
		returning id
	`, req.Username, req.Email, hash, req.FirstName, req.LastName, req.Phone, req.Gender, birth, req.Region).Scan(&id)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"db_insert"}); return }

	_ = middleware.SetSession(c, id, req.Username, "USER")

	// публичный ответ без password_hash
	c.JSON(http.StatusCreated, gin.H{
		"user": gin.H{
			"id": id, "username": req.Username, "email": req.Email, "role": "USER",
			"firstName": req.FirstName, "lastName": req.LastName, "phone": req.Phone,
			"gender": req.Gender, "birthdate": birth, "region": req.Region,
		},
	})
}

// простые env-геттеры (как в middleware) — можно заменить на os.Getenv
func getEnvOr(k, def string) string { return def }
