package models

import "time"

type Review struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"userId"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}