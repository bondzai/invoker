package util

import (
	"errors"
)

func NewError(msg string) error {
	return errors.New(msg)
}

var ErrCommon = NewError("common error")
var ErrInvalidTaskInterval = NewError("invalid task interval")
var ErrInvalidTaskCronExpr = NewError("invalid task cron expression")
