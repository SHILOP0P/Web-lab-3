package config

import (
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL    string
	JWTSecret      []byte
	PasswordPepper string
	CookieName     string
	CookieSecure   bool
	CookieMaxAge   time.Duration
	CORSOrigins    []string
}

func Load() (Config, error) {
	_ = godotenv.Load()

	secure := strings.EqualFold(os.Getenv("COOKIE_SECURE"), "true")
	maxAgeDays, _ := strconv.Atoi(get("COOKIE_MAX_AGE_DAYS", "7"))

	var origins []string
	if s := os.Getenv("CORS_ORIGINS"); s != "" {
		for _, o := range strings.Split(s, ",") {
			origins = append(origins, strings.TrimSpace(o))
		}
	}
	return Config{
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		JWTSecret:      []byte(get("JWT_SECRET", "dev_secret")),
		PasswordPepper: get("PASSWORD_PEPPER", ""),
		CookieName:     get("COOKIE_NAME", "sid"),
		CookieSecure:   secure,
		CookieMaxAge:   time.Duration(maxAgeDays) * 24 * time.Hour,
		CORSOrigins:    origins,
	}, nil
}

func get(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
