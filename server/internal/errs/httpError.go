package errs

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/VrMolodyakov/stock-market/pkg/logging"
	"github.com/gin-gonic/gin"
)

type ErrResponse struct {
	Error ServiceError `json:"error"`
}

type ServiceError struct {
	Kind    string `json:"kind,omitempty"`
	Code    string `json:"code,omitempty"`
	Param   string `json:"param,omitempty"`
	Message string `json:"message,omitempty"`
}

func HTTPErrorResponse(ctx *gin.Context, logger *logging.Logger, err error) {
	var e *Error
	if errors.As(err, &e) {
		switch e.Kind {
		case Unauthenticated:
			badRequesteResponse(ctx, logger, e)
			return
		case Validation:
			validationRequesteResponse(ctx, logger, e)
			return
		case Unauthorized:
			unauthorizedErrorResponse(ctx, logger, e)
			return
		default:
			commonErrorResponse(ctx, logger, e)
			return
		}
	}
	unknownErrorResponse(ctx, logger, err)
}

func badRequesteResponse(c *gin.Context, logger *logging.Logger, err *Error) {
	logger.Errorf("http status code %v\n error = %v", http.StatusUnauthorized, err)
	c.JSON(http.StatusBadRequest, err.Error())
}

func validationRequesteResponse(c *gin.Context, logger *logging.Logger, err *Error) {
	logger.Errorf("http status code %v\n error = %v", http.StatusUnauthorized, err)
	c.JSON(http.StatusBadRequest, err.Param)
}

func unauthorizedErrorResponse(c *gin.Context, logger *logging.Logger, err *Error) {
	logger.Errorf("http status code %v\n error = %v", http.StatusForbidden, err)
	c.JSON(http.StatusForbidden, err.Error())
}

func unknownErrorResponse(c *gin.Context, logger *logging.Logger, err error) {
	logger.Errorf("http status code %v\n error = %v", http.StatusInternalServerError, err)
	c.Header("Content-Type", "application/json")
	c.Header("X-Content-Type-Options", "nosniff")
	c.JSON(http.StatusInternalServerError, err.Error())
}

func commonErrorResponse(c *gin.Context, logger *logging.Logger, err *Error) {

	if err.isZero() {
		c.JSON(http.StatusInternalServerError, "internal server error - please contact support")
		return
	}

	logger.Errorf("kind =  %v\n parameter = %v\n code =	%v\n error = %v", err.Kind.String(), string(err.Param), string(err.Code), err)
	e := newErrResponse(err)
	errJSON, _ := json.Marshal(e)
	errAsStr := string(errJSON)
	c.Header("Content-Type", "application/json")
	c.Header("X-Content-Type-Options", "nosniff")
	c.JSON(http.StatusInternalServerError, errAsStr)
}

func newErrResponse(err *Error) ErrResponse {
	const msg string = "internal server error - please contact support"

	switch err.Kind {
	case Internal, Database:
		return ErrResponse{
			Error: ServiceError{
				Kind:    Internal.String(),
				Message: msg,
			},
		}
	default:
		return ErrResponse{
			Error: ServiceError{
				Kind:    err.Kind.String(),
				Code:    string(err.Code),
				Param:   string(err.Param),
				Message: err.Error(),
			},
		}
	}
}
