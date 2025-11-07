package models

import "time"

type User struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Role         string     `json:"role"`
	FirstName    *string    `json:"firstName,omitempty"`
	LastName     *string    `json:"lastName,omitempty"`
	Phone        *string    `json:"phone,omitempty"`
	Gender       *string    `json:"gender,omitempty"`
	Birthdate    *time.Time `json:"birthdate,omitempty"`
	Region       *string    `json:"region,omitempty"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
}

type PublicUser struct {
	ID        string     `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	FirstName *string    `json:"firstName,omitempty"`
	LastName  *string    `json:"lastName,omitempty"`
	Phone     *string    `json:"phone,omitempty"`
	Gender    *string    `json:"gender,omitempty"`
	Birthdate *time.Time `json:"birthdate,omitempty"`
	Region    *string    `json:"region,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

func (u User) Public() PublicUser {
	return PublicUser{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Phone:     u.Phone,
		Gender:    u.Gender,
		Birthdate: u.Birthdate,
		Region:    u.Region,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
