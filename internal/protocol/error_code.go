package protocol

const (
	ErrorNone                 = 0
	ErrorUnknown              = 1
	ErrorNotFound             = 2
	ErrWrongParam             = 3
	ErrInvalidParam           = 4
	ErrUserNotFound           = 5
	ErrUnauthorized           = 6
	ErrUploadFileSizeTooLarge = 7

	ErrUnknownMsgType  = 4001
	ErrSessionRequired = 4003
	ErrNotMember       = 4004
	ErrUserBanned      = 4005
	ErrInternalError   = 5000
)
