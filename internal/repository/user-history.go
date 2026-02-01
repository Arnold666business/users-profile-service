package repository

import (
	"context"
	"encoding/json"
	"time"
	"users-profile-service/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserHistory struct {
	db *pgxpool.Pool
}

func NewUserHistory(db *pgxpool.Pool) *UserHistory {
	return &UserHistory{db: db}
}

func (uHistory *UserHistory) Save(ctx context.Context, u models.User, action models.ACTION) (int64, error) {
	db := GetQuerier(ctx, uHistory.db)
	userHistory := models.UserHistory{
		UserId:   u.Id,
		Action:   action,
		CreateAt: time.Now(),
		Email:    u.Email,
		Login:    u.Login,
		OldFields: map[string]interface{}{
			"id":                  u.Id,
			"login":               u.Login,
			"email":               u.Email,
			"access_email_status": u.AccessEmailStatus,
		},
	}

	oldFields, _ := json.Marshal(userHistory.OldFields)
	query := `
		INSERT INTO users_profile.users_h (user_id, action, create_at, email, login, old_fields) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err := db.QueryRow(ctx, query,
		&userHistory.UserId, userHistory.Action.String(), &userHistory.CreateAt, &userHistory.Email, &userHistory.Login, &oldFields,
	).Scan(&userHistory.Id)
	if err != nil {
		return 0, err
	}
	return userHistory.Id, nil
}
