package errs

import "errors"

type Kind uint8
type Parameter string
type Code string

const (
	Other           Kind = iota // Unclassified error. This value is not printed in the error message.
	Invalid                     // Invalid operation for this type of item.
	Exist                       // Item already exists.
	NotExist                    // Item does not exist.
	Private                     // Information withheld.
	Internal                    // Internal error or inconsistency.
	Database                    // Error from database.
	Validation                  // Input validation error.
	InvalidRequest              // Invalid Request
	Unauthenticated             // Unauthenticated Request

	Unauthorized
)

type Error struct {
	Kind  Kind
	Param Parameter
	Code  Code
	Err   error
}

func (e *Error) Unweap() error {
	return e.Err
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func New(args ...interface{}) error {
	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case string:
			e.Err = errors.New(arg)
		case Kind:
			e.Kind = arg
		case error:
			e.Err = arg
		case Parameter:
			e.Param = arg
		}

	}
	return e

}
