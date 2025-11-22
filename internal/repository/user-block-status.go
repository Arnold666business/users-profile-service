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
	NotFoundUserBlockStatusError = errors.New("user_block_status not found")
)

type UserBlockStatus struct {
	db *pgxpool.Pool
}

func NewUserBlockStatus(db *pgxpool.Pool) *UserBlockStatus {
	return &UserBlockStatus{db: db}
}

func (ubsr *UserBlockStatus) GetByUserId(ctx context.Context, id int64) (*models.UserBlockStatus, error) {
	db := GetQuerier(ctx, ubsr.db)
	var ubs models.UserBlockStatus
	query := `SELECT id, user_id, block_type_id, forever_flag, unblock_date, is_active, unblock_event_sent
          FROM users_profile.users_block_status WHERE user_id = $1;`
	err := db.QueryRow(ctx, query, id).Scan(
		&ubs.Id, &ubs.UserId, &ubs.BlockTypeId, &ubs.ForeverFlag, &ubs.UnblockDate, &ubs.IsActive, &ubs.UnBlockEventSent)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NotFoundUserBlockStatusError
		} else {
			return nil, fmt.Errorf("user_block_status, user_id: %v: %w", id, err)
		}
	}
	return &ubs, nil
}

func (ubsr *UserBlockStatus) Save(ctx context.Context, ubs *models.UserBlockStatus) (int64, error) {
	db := GetQuerier(ctx, ubsr.db)
	query := `
		INSERT INTO users_profile.users_block_status (user_id, block_type_id, forever_flag, unblock_date, is_active, unblock_event_sent)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;
	`

	err := db.QueryRow(ctx, query,
		ubs.UserId,
		ubs.BlockTypeId,
		ubs.ForeverFlag,
		ubs.UnblockDate,
		ubs.IsActive,
		ubs.UnBlockEventSent,
	).Scan(&ubs.Id)

	if err != nil {
		return 0, fmt.Errorf("failed to insert user_block_status for user_id=%d: %w", ubs.UserId, err)
	}
	return ubs.Id, nil
}

func (ubsr *UserBlockStatus) Update(ctx context.Context, ubs *models.UserBlockStatus) error {
	db := GetQuerier(ctx, ubsr.db)
	query := `
		UPDATE users_profile.users_block_status  SET block_type_id=$1, forever_flag=$2, unblock_date=$3, is_active=$4, unblock_event_sent=$5
		 WHERE id=$6;
	`

	_, err := db.Exec(ctx, query,
		ubs.BlockTypeId,
		ubs.ForeverFlag,
		ubs.UnblockDate,
		ubs.IsActive,
		ubs.UnBlockEventSent,
		ubs.Id,
	)
	if err != nil {
		return fmt.Errorf("failed to update user_block_status for user_id=%d: %w", ubs.UserId, err)
	}
	return nil
}

func (ubsr *UserBlockStatus) UpsertByUserId(ctx context.Context, ubs *models.UserBlockStatus) error {
	db := GetQuerier(ctx, ubsr.db)
	ubs.IsActive = true // на всякий и наче constraint не сработает
	query := `
		INSERT INTO users_profile.users_block_status (user_id, block_type_id, forever_flag, unblock_date, is_active, unblock_event_sent)
		VALUES ($1, $2, $3, $4, $5) ON CONFLICT(user_id) WHERE is_active=true DO UPDATE SET 
		block_type_id = EXCLUDED.block_type_id,
        forever_flag = EXCLUDED.forever_flag,
        unblock_date = EXCLUDED.unblock_date
	`

	_, err := db.Exec(ctx, query,
		ubs.UserId,
		ubs.BlockTypeId,
		ubs.ForeverFlag,
		ubs.UnblockDate,
		ubs.IsActive,
	)
	if err != nil {
		return fmt.Errorf("failed to upesrt user_block_status for user_id=%d: %w", ubs.UserId, err)
	}
	return nil
}

func (ubsr *UserBlockStatus) FindForUnblockProcessing(ctx context.Context, limit int) ([]*models.UserBlockStatus, error) {
	db := GetQuerier(ctx, ubsr.db)
	query := `SELECT id, user_id, block_type_id, forever_flag, unblock_date, is_active, unblock_event_sent
		FROM users_profile.users_block_status 
		WHERE unblock_event_sent = false AND unblock_date < now()
		LIMIT $1 FOR UPDATE SKIP LOCKED`
	rows, err := db.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ubs []*models.UserBlockStatus
	for rows.Next() {
		ub := &models.UserBlockStatus{}
		err := rows.Scan(&ub.Id,
			&ub.UserId,
			&ub.BlockTypeId,
			&ub.ForeverFlag,
			&ub.UnblockDate,
			&ub.IsActive,
			&ub.UnBlockEventSent)
		if err != nil {
			return nil, err
		}
		ubs = append(ubs, ub)
	}
	return ubs, nil
}
