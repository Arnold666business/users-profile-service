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

// govno nujno?
func (ubsr *UserBlockStatus) GetByUserId(ctx context.Context, id int64) (*models.UserBlockStatus, error) {
	db := GetQuerier(ctx, ubsr.db)
	var ubs models.UserBlockStatus
	query := `SELECT id, user_id, block_type_id, forever_flag, unblock_date
          FROM users_profile.users_block_status WHERE user_id = $1;`
	err := db.QueryRow(ctx, query, id).Scan(
		&ubs.Id, &ubs.UserId, &ubs.BlockTypeId, &ubs.ForeverFlag, &ubs.UnblockDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NotFoundUserBlockStatusError
		} else {
			return nil, fmt.Errorf("user_block_status, user_id: %v: %w", id, err)
		}
	}
	return &ubs, nil
}

// govno nujno?
func (ubsr *UserBlockStatus) Save(ctx context.Context, ubs *models.UserBlockStatus) (int64, error) {
	db := GetQuerier(ctx, ubsr.db)
	query := `
		INSERT INTO users_profile.users_block_status (user_id, block_type_id, forever_flag, unblock_date)
		VALUES ($1, $2, $3, $4) RETURNING id;
	`

	err := db.QueryRow(ctx, query,
		ubs.UserId,
		ubs.BlockTypeId,
		ubs.ForeverFlag,
		ubs.UnblockDate,
	).Scan(&ubs.Id)

	if err != nil {
		return 0, fmt.Errorf("failed to insert user_block_status for user_id=%d: %w", ubs.UserId, err)
	}
	return ubs.Id, nil
}

// govno nujno?
func (ubsr *UserBlockStatus) Update(ctx context.Context, ubs *models.UserBlockStatus) error {
	db := GetQuerier(ctx, ubsr.db)
	query := `
		UPDATE users_profile.users_block_status  SET block_type_id=$1, forever_flag=$2, unblock_date=$3 WHERE id=$4;
	`

	_, err := db.Exec(ctx, query,
		ubs.BlockTypeId,
		ubs.ForeverFlag,
		ubs.UnblockDate,
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
		INSERT INTO users_profile.users_block_status (user_id, block_type_id, forever_flag, unblock_date, is_active)
		VALUES ($1, $2, $3, $4, $5) ON CONFLICT ON CONSTRAINT users_block_status_user_active_idx DO UPDATE SET 
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

//func (ubsr *UserBlockStatus) UpsertByUserId(ctx context.Context, ubs *models.UserBlockStatus) error {
//	db := GetQuerier(ctx, ubsr.db)
//
//
//	_, err := db.Exec(ctx,
//		"UPDATE users_profile.users_block_status SET is_active = false WHERE user_id = $1 AND is_active = true",
//		ubs.UserId,
//	)
//	if err != nil {
//		return fmt.Errorf("failed to deactivate old blocks for user_id=%d: %w", ubs.UserId, err)
//	}
//
//
//	query := `
//      INSERT INTO users_profile.users_block_status
//    (user_id, block_type_id, forever_flag, unblock_date, is_active)
//       VALUES ($1, $2, $3, $4, true)
// `
//
//	_, err = db.Exec(ctx, query,
//		ubs.UserId,
//		ubs.BlockTypeId,
//		ubs.ForeverFlag,
///		ubs.UnblockDate,
//	)
//	if err != nil {
//		return fmt.Errorf("failed to insert user_block_status for user_id=%d: %w", ubs.UserId, err)
//	}
//
//	return nil
//}
