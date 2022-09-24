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

func (e *Error) Unwrap() error {
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
		case Code:
			e.Code = arg

		}

	}
	return e

}

func (k Kind) String() string {
	switch k {
	case Other:
		return "other_error"
	case Invalid:
		return "invalid_operation"
	case Exist:
		return "item_already_exists"
	case NotExist:
		return "item_does_not_exist"
	case Internal:
		return "internal_error"
	case Database:
		return "database_error"
	case Validation:
		return "input_validation_error"
	case InvalidRequest:
		return "invalid_request_error"
	case Unauthenticated:
		return "unauthenticated_request"
	case Unauthorized:
		return "unauthorized_request"
	}
	return "unknown_error_kind"
}

func (e *Error) isZero() bool {
	return e.Kind == 0 && e.Param == "" && e.Code == "" && e.Err == nil
}
