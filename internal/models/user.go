package models

import "time"

type User struct {
	Id                int64
	Login             string
	Email             string
	AccessEmailStatus bool
	Role              int
	IsDeleted         bool
	DeletedAt         time.Time
}
