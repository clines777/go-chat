package protocol

type DispatchError struct {
	Code    int
	Message string
	Fatal   bool
}

func (e *DispatchError) Error() string {
	return e.Message
}
