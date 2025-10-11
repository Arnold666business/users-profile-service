package models

import "time"

type UserHistory struct {
	Id        int64
	UserId    int64
	Action    ACTION
	CreateAt  time.Time
	Email     string
	Login     string
	OldFields map[string]interface{}
}

type ACTION int

const (
	CREATE ACTION = iota
	UPDATE
	DELETE
	BLOCKED
	UNBLOCKED
)

func (a ACTION) String() string {
	switch a {
	case CREATE:
		return "CREATE"
	case UPDATE:
		return "UPDATE"
	case DELETE:
		return "DELETE"
	case BLOCKED:
		return "BLOCKED"
	case UNBLOCKED:
		return "UNBLOCKED"
	default:
		return "UNKNOWN"
	}
}
