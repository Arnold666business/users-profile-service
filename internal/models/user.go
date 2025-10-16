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

type UserAggregate struct {
	Id               int64
	Login            string
	Email            string
	EmailAccess      bool
	Role             int
	IsDeleted        bool
	DeletedAt        time.Time
	BlockTypeId      int
	ForeverFlag      bool
	UnBlockDate      time.Time
	BlockTitle       string
	BlockDescription string
}
