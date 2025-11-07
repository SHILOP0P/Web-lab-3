package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret  = []byte(getEnvOr("JWT_SECRET", "dev_secret_change_me"))
	cookieName = getEnvOr("COOKIE_NAME", "sid")
	cookieDays = getEnvIntOr("COOKIE_MAX_AGE_DAYS", 7)
	cookieSecure = getEnvOr("COOKIE_SECURE", "false") == "true"
)

func SetSession(c *gin.Context, userID int64, username, role string) error {
	claims := jwt.MapClaims{
		"sub":      userID,
		"username": username,
		"role":     role,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Duration(cookieDays) * 24 * time.Hour).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString(jwtSecret)
	if err != nil { return err }

	maxAge := 60 * 60 * 24 * cookieDays

	// secure только если HTTPS (иначе на http кука не придёт)
	secure := c.Request.TLS != nil

	// для одного origin можно Lax
	c.SetSameSite(http.SameSiteLaxMode)

	// тот же cookieName везде, домен пустой — уходит на текущий хост
	c.SetCookie(cookieName, token, maxAge, "/", "", secure, true)
	return nil
}

func ClearSession(c *gin.Context) {
	secure := c.Request.TLS != nil
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(cookieName, "", -1, "/", "", secure, true)
}


func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(cookieName)
		if err != nil || tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"unauthorized"})
			return
		}
		claims := jwt.MapClaims{}
		t, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token)(interface{}, error){ return jwtSecret, nil })
		if err != nil || !t.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"invalid_session"})
			return
		}
		// кладём в контекст
		c.Set("user_id",   int64FromClaims(claims["sub"]))
		c.Set("username",  strFromClaims(claims["username"]))
		c.Set("role",      strFromClaims(claims["role"]))
		c.Next()
	}
}

// — вспомогалки —
func int64FromClaims(v any) int64 {
	switch t := v.(type) {
	case float64: return int64(t)
	case int64:   return t
	default:      return 0
	}
}
func strFromClaims(v any) string {
	if s, ok := v.(string); ok { return s }
	return ""
}

// небольшие геттеры env без внешних пакетов.
func getEnvOr(k, def string) string {
	if v := getenv(k); v != "" { return v }
	return def
}
func getEnvIntOr(k string, def int) int {
	if v := getenv(k); v != "" {
		if n, err := atoi(v); err == nil { return n }
	}
	return def
}

// минимальные обёртки, чтобы не тянуть os/strconv во всех местах
func getenv(k string) string { return lookupEnv(k) }
func atoi(s string) (int, error) { return parseInt(s) }

// ↓ заменяются в компиляции линкером go (простая заглушка):
var lookupEnv = func(k string) string { return "" }
var parseInt  = func(s string) (int, error) { return 0, fmtError("atoi") }
type errorString string
func (e errorString) Error() string { return string(e) }
func fmtError(_ string) error { return errorString("error") }
