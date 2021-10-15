package internal

var (
	ErrNoVersionSpecified   = NewUserError("no version specified")
	ErrInvalidVersionFormat = NewUserError("invalid version format, correct format should be like: 1.17.2, v1.17.2 or go1.17.2, etc.")
	ErrVersionNotInstalled  = NewUserError("version not installed")
	ErrVersionIsInUse       = NewUserError("version is in use")
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
