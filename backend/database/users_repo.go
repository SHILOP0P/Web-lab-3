package database

import (
	"context"
	"time"
	"backend/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct{ pool *pgxpool.Pool }

func NewUserRepo(p *pgxpool.Pool) *UserRepo { return &UserRepo{pool: p} }

func (r *UserRepo) ExistsByUsernameOrEmail(ctx context.Context, username, email string) (bool, error) {
	var x int
	err := r.pool.QueryRow(ctx, `select 1 from users where username=$1 or email=$2 limit 1`, username, email).Scan(&x)
	if err == nil {
		return true, nil
	}
	return false, nil
}

type CreateUserParams struct {
	Username, Email, PasswordHash, Role string
	FirstName, LastName, Phone, Gender, Region *string
	Birthdate *time.Time
}

func (r *UserRepo) Create(ctx context.Context, p CreateUserParams) (models.User, error) {
	if p.Role == "" { p.Role = "USER" }
	row := r.pool.QueryRow(ctx, `
		insert into users
		  (username,email,password_hash,role,first_name,last_name,phone,gender,birthdate,region)
		values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		returning id, username, email, password_hash, role,
		          first_name, last_name, phone, gender, birthdate, region,
		          created_at, updated_at
	`, p.Username, p.Email, p.PasswordHash, p.Role,
		p.FirstName, p.LastName, p.Phone, p.Gender, p.Birthdate, p.Region)

	var u models.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role,
		&u.FirstName, &u.LastName, &u.Phone, &u.Gender, &u.Birthdate, &u.Region,
		&u.CreatedAt, &u.UpdatedAt)
	return u, err
}

func (r *UserRepo) GetByUsernameOrEmail(ctx context.Context, login string) (models.User, error) {
	row := r.pool.QueryRow(ctx, `
		select id, username, email, password_hash, role,
		       first_name, last_name, phone, gender, birthdate, region,
		       created_at, updated_at
		from users
		where username=$1 or email=$1
		limit 1
	`, login)

	var u models.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role,
		&u.FirstName, &u.LastName, &u.Phone, &u.Gender, &u.Birthdate, &u.Region,
		&u.CreatedAt, &u.UpdatedAt)
	return u, err
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (models.User, error) {
	row := r.pool.QueryRow(ctx, `
		select id, username, email, password_hash, role,
		       first_name, last_name, phone, gender, birthdate, region,
		       created_at, updated_at
		from users where id=$1
	`, id)

	var u models.User
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.PasswordHash, &u.Role,
		&u.FirstName, &u.LastName, &u.Phone, &u.Gender, &u.Birthdate, &u.Region,
		&u.CreatedAt, &u.UpdatedAt)
	return u, err
}
