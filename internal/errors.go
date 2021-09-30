package internal

var (
	ErrNoVersionSpecified   = NewUserError("no version specified")
	ErrInvalidVersionFormat = NewUserError("invalid version format")
)

type UserError struct {
	msg string
}

func NewUserError(msg string) error {
	return UserError{msg}
}

func (e UserError) Error() string {
	return e.msg
}
