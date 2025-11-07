// backend/handlers/reviews/reviews.go
package reviews

import (
	"backend/database"
	"backend/models"
	"net/http"
	"strconv"
	"strings"
	"time"
	"database/sql"
	"github.com/gin-gonic/gin"
	"math"
)


// ===== helpers =====

func currentUser(c *gin.Context) (id int64, username, role string, ok bool) {
	uidV, _ := c.Get("user_id")
	switch v := uidV.(type) {
	case int64:
		id = v
	case int:
		id = int64(v)
	case int32:
		id = int64(v)
	case uint:
		id = int64(v)
	case uint32:
		id = int64(v)
	case uint64:
		// не допускаем переполнения при переводе в int64
		if v <= uint64(math.MaxInt64) {
			id = int64(v)
		} else {
			id = 0
		}
	case float64: // иногда кладут float64
		id = int64(v)
	case string: // иногда кладут строку
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			id = n
		}
	default:
		id = 0
	}

	unV, _ := c.Get("username")
	roleV, _ := c.Get("role")
	username, _ = unV.(string)
	role, _ = roleV.(string)

	ok = id > 0 && username != ""
	return
}



func isAdmin(username, role string) bool {
	// Админ по логину ИЛИ по роли, если когда-нибудь заведёшь 'ADMIN'
	return username == "SHILOP0P" || strings.EqualFold(role, "ADMIN")
}

// ===== handlers =====

// GET /api/reviews  (публично)
func List(c *gin.Context) {
	db, err := database.ConnectToServer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db_connect"})
		return
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT r.id, r.user_id, u.username, r.content, r.created_at
		FROM reviews r
		JOIN users u ON u.id = r.user_id
		ORDER BY r.created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db_query"})
		return
	}
	defer rows.Close()

	out := make([]models.Review, 0, 32)
	for rows.Next() {
		var r models.Review
		if err := rows.Scan(&r.ID, &r.UserID, &r.Username, &r.Content, &r.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db_scan"})
			return
		}
		out = append(out, r)
	}

	c.JSON(http.StatusOK, gin.H{"reviews": out})
}

type createReq struct {
	Content string `json:"content"`
}

// POST /api/reviews  (только авторизованные)
func Create(c *gin.Context) {
	uid, username, _, ok := currentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req createReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_json"})
		return
	}
	req.Content = strings.TrimSpace(req.Content)
	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "content_required"})
		return
	}

	db, err := database.ConnectToServer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db_connect"})
		return
	}
	defer db.Close()

	var id int64
	var created time.Time
	if err := db.QueryRow(`
		INSERT INTO reviews (user_id, content)
		VALUES ($1, $2)
		RETURNING id, created_at
	`, uid, req.Content).Scan(&id, &created); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db_insert"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"review": models.Review{
			ID:        id,
			UserID:    uid,
			Username:  username,
			Content:   req.Content,
			CreatedAt: created,
		},
	})
}

// DELETE /api/reviews/:id  (только админ SHILOP0P)
// DELETE /api/reviews/:id  (владелец или админ SHILOP0P)
func Delete(c *gin.Context) {
	uid, username, role, ok := currentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad_id"})
		return
	}

	db, err := database.ConnectToServer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db_connect"})
		return
	}
	defer db.Close()

	// Узнаём владельца отзыва
	var ownerID int64
	if err := db.QueryRow(`SELECT user_id FROM reviews WHERE id = $1`, id).Scan(&ownerID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db_select_owner"})
		return
	}

	// Разрешаем: админ ИЛИ владелец
	if !(isAdmin(username, role) || ownerID == uid) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// Удаляем
	res, err := db.Exec(`DELETE FROM reviews WHERE id = $1`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db_delete"})
		return
	}
	if aff, _ := res.RowsAffected(); aff == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ok": true})
}

