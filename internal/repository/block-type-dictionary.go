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
	NotFoundBlockTypeDictionaryError = errors.New("block_type_dictionary not found")
)

type BlockTypeDictionary struct {
	db *pgxpool.Pool
}

func NewBlockTypeDictionary(db *pgxpool.Pool) *BlockTypeDictionary {
	return &BlockTypeDictionary{db: db}
}

func (btdr *BlockTypeDictionary) GetByBlockType(ctx context.Context, typeId int) (*models.BlockTypeDictionary, error) {
	db := GetQuerier(ctx, btdr.db)
	var btd models.BlockTypeDictionary
	err := db.QueryRow(ctx, `SELECT type_id, title, description, hour
          FROM users_profile.block_type_dictionary WHERE type_id = $1;`, typeId).Scan(
		&btd.BlockType, &btd.Title, &btd.Description, &btd.Hour)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NotFoundBlockTypeDictionaryError
		} else {
			return nil, fmt.Errorf("block_type_dictionary, typeId: %v: %w", typeId, err)
		}
	}
	return &btd, nil
}
