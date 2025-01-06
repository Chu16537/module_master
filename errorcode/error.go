package errorcode

type Error struct {
	Code int   `json:"code"`
	Err  error `json:"err"`
}

func New(code int, err error) *Error {
	return &Error{
		Code: code,
		Err:  err,
	}
}

func (e *Error) GetCode() int {
	return e.Code
}

func (e *Error) GetErr() error {
	return e.Err
}

func Success() *Error {
	return &Error{
		Code: Code_Success,
		Err:  nil,
	}
}

func (e *Error) IsNotSuccess() bool {
	return e.Code != Code_Success
}
