package models

import "time"

type UserBlockStatus struct {
	Id          int64
	UserId      int64
	BlockTypeId int
	ForeverFlag bool
	UnblockDate time.Time
}

func (u *UserBlockStatus) SetUnblockDate(hours int) {
	unblockDate := time.Now().Add(time.Duration(hours) * time.Hour)
	u.UnblockDate = unblockDate
}
