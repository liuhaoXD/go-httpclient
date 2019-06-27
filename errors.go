package httpclient

import (
	"errors"
)

var (
	ErrUrlIsEmpty    = errors.New("url is empty")
	ErrLoggerIsEmpty = errors.New("logger is empty")
)
