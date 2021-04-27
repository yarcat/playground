package redis

import "errors"

const EOL = "\r\n"

var (
	ErrInternal  = errors.New("internal error")
	ErrResponse  = errors.New("response error")
	ErrFirstByte = errors.New("unexpected first byte")
)
