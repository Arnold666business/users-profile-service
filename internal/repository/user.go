package repository

import (
	"context"
	"errors"
	"fmt"
	"users-profile-service/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	NotFoundUserError = errors.New("users not found")
)

type User struct {
	db *pgxpool.Pool
}

func NewUser(db *pgxpool.Pool) *User {
	return &User{db: db}
}

func (user *User) GetById(ctx context.Context, id int64) (models.User, error) {
	var u models.User
	err := user.db.QueryRow(ctx, `SELECT id, login, email, role, is_deleted, delete_at
          FROM users_profile.users WHERE id = $1;`, id).Scan(
		&u.Id, &u.Login, &u.Email, &u.Role, &u.IsDeleted, &u.DeletedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NotFoundUserError
		} else {
			return nil, fmt.Errorf("user %v: %w", id, err)
		}
	}
	return &u, nil
}
