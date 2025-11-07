package auth

import (
	"backend/database"
	"backend/middleware"
	"backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type loginReq struct {
	UsernameOrEmail string `json:"usernameOrEmail"`
	Password        string `json:"password"`
}

func Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
	 c.JSON(http.StatusBadRequest, gin.H{"error":"bad_json"}); return
	}
	req.UsernameOrEmail = strings.TrimSpace(req.UsernameOrEmail)
	if req.UsernameOrEmail=="" || req.Password=="" {
	 c.JSON(http.StatusBadRequest, gin.H{"error":"missing_fields"}); return
	}

	db, err := database.ConnectToServer()
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"db_connect"}); return }
	defer db.Close()

	var id int64
	var username, email, role, hash string
	err = db.QueryRow(`
		select id, username, email, role, password_hash
		from users where username=$1 or email=$1 limit 1
	`, req.UsernameOrEmail).Scan(&id, &username, &email, &role, &hash)
	if err != nil { c.JSON(http.StatusUnauthorized, gin.H{"error":"invalid_credentials"}); return }

	if err := utils.CheckPassword(req.Password, hash, passwordPepper); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error":"invalid_credentials"}); return
	}

	_ = middleware.SetSession(c, id, username, role)

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{"id": id, "username": username, "email": email, "role": role},
	})
}
