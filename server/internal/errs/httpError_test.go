package errs

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/VrMolodyakov/stock-market/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_httpErrorStatusCode(t *testing.T) {
	tests := []struct {
		title string
		err   Error
		want  string
	}{
		{"Exist", Error{Kind: Database}, "internal_error"},
		{"Exist", Error{Kind: Internal}, "internal_error"},
		{"Exist", Error{Kind: Other, Err: errors.New("some error")}, "other_error"},
	}
	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			got := newErrResponse(&test.err)
			assert.Equal(t, got.Error.Kind, test.want)

		})
	}
}

func TestHTTPErrorResponse_StatusCode(t *testing.T) {

	type args struct {
		w   *httptest.ResponseRecorder
		l   *logging.Logger
		err error
	}

	l := logging.GetLogger("debug")

	unauthenticatedErr := New(Unauthenticated, "some error from Google")
	unauthorizedErr := New(Unauthorized, "some authorization error")

	tests := []struct {
		name string
		args args
		want int
	}{
		{"empty *Error", args{httptest.NewRecorder(), l, &Error{Err: errors.New("")}}, http.StatusInternalServerError},
		{"unauthenticated", args{httptest.NewRecorder(), l, unauthenticatedErr}, http.StatusBadRequest},
		{"unauthorized", args{httptest.NewRecorder(), l, unauthorizedErr}, http.StatusForbidden},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(test.args.w)
			HTTPErrorResponse(c, l, test.args.err)
			got := test.args.w.Result().StatusCode
			assert.Equal(t, got, test.want)
		})
	}
}

func TestHTTPErrorResponse_Body(t *testing.T) {

	type args struct {
		w   *httptest.ResponseRecorder
		l   *logging.Logger
		err error
	}

	lgr := logging.GetLogger("debug")

	tests := []struct {
		name string
		args args
		want string
	}{
		{"empty Error", args{httptest.NewRecorder(), lgr, &Error{}}, "\"internal server error - please contact support\""},
		{"unauthenticated", args{httptest.NewRecorder(), lgr, New(Unauthenticated, "some unauthenticated error")}, "\"some unauthenticated error\""},
		{"unauthorized", args{httptest.NewRecorder(), lgr, New(Unauthorized, "some authorization error")}, "\"some authorization error\""},
		{"normal", args{httptest.NewRecorder(), lgr, New(Exist, Parameter("some_param"), Code("some_code"), errors.New("some error"))}, "\"{\\\"error\\\":{\\\"kind\\\":\\\"item_already_exists\\\",\\\"code\\\":\\\"some_code\\\",\\\"param\\\":\\\"some_param\\\",\\\"message\\\":\\\"some error\\\"}}\""},
		{"not via New", args{httptest.NewRecorder(), lgr, errors.New("some error")}, "\"some error\""},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(test.args.w)
			HTTPErrorResponse(c, lgr, test.args.err)
			got := strings.TrimSpace(test.args.w.Body.String())
			assert.Equal(t, got, test.want)
		})
	}
}
