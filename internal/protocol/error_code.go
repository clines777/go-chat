package protocol

var (
	ErrorNone       = 0
	ErrorUnknown    = 1
	ErrorNotFound   = 2
	ErrWrongParam   = 3
	ErrInvalidParam = 4
	ErrUserNotFound = 5
	ErrUnauthorized = 401

	ErrUnknownMsgType  = 4001
	ErrSessionRequired = 4003
	ErrInternalError   = 5000
)
