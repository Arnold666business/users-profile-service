package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	BlockTypeDictionaryRepository *BlockTypeDictionary
	UserRepository                *User
	UserBlockStatusRepository     *UserBlockStatus
	UserHistoryRepository         *UserHistory
	Transactor                    *Transactor
}

func Build(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		BlockTypeDictionaryRepository: NewBlockTypeDictionary(db),
		UserRepository:                NewUser(db),
		UserBlockStatusRepository:     NewUserBlockStatus(db),
		UserHistoryRepository:         NewUserHistory(db),
		Transactor:                    NewTransactor(db),
	}
}
