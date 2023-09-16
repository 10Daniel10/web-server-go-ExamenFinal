package internal

import (
	"errors"
	"gorm.io/gorm"
)

type Error struct {
	inter   error
	Code    int
	Message string
}

func (i Error) Error() string {

	//TODO implement me
	panic("implement me")
}

var (
	ErrDuplicate = Error{
		inter:   gorm.ErrDuplicatedKey,
		Code:    0,
		Message: "",
	}
	ErrNotFound = errors.New("record not found")
	ErrBadReq   = errors.New("bad request")
)
