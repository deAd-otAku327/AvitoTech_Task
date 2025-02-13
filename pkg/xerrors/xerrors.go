package xerrors

type Xerror interface {
	Error() string
	Code() int
}

type xerror struct {
	Err        error
	StatusCode int
}

func New(err error, code int) Xerror {
	return &xerror{
		Err:        err,
		StatusCode: code,
	}
}

func (xe *xerror) Error() string {
	return xe.Err.Error()
}

func (xe *xerror) Code() int {
	return xe.StatusCode
}
