package wrapper_test

import (
	"github.com/aasumitro/tix/pkg/http/wrapper"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewHttpRespond(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		data     interface{}
		args     []any
		expected interface{}
	}{
		{
			name:     "success with no pagination",
			code:     http.StatusOK,
			data:     []string{"foo", "bar"},
			args:     nil,
			expected: wrapper.CommonRespond{Code: http.StatusOK, Status: "OK", Data: []string{"foo", "bar"}},
		},
		{
			name: "success with pagination",
			code: http.StatusOK,
			data: []string{"foo", "bar"},
			args: []any{
				2,
				1,
				wrapper.Paging{URL: "http://example.com/next", Path: "/next"},
				wrapper.Paging{URL: "http://example.com/prev", Path: "/prev"},
				100,
			},
			expected: wrapper.SuccessWithPaginationRespond{
				Code:     http.StatusOK,
				Status:   "OK",
				Total:    2,
				Current:  1,
				Next:     wrapper.Paging{URL: "http://example.com/next", Path: "/next"},
				Previous: wrapper.Paging{URL: "http://example.com/prev", Path: "/prev"},
				Data:     []string{"foo", "bar"},
			},
		},
		{
			name:     "error with data",
			code:     http.StatusBadRequest,
			data:     "invalid request",
			expected: wrapper.ErrorRespond{Code: http.StatusBadRequest, Status: "Bad Request", Data: "invalid request"},
		},
		{
			name:     "error with no data",
			code:     http.StatusBadRequest,
			expected: wrapper.ErrorRespond{Code: http.StatusBadRequest, Status: "Bad Request", Data: "something went wrong with the request"},
		},
		{
			name:     "error with no data and server error code",
			code:     http.StatusInternalServerError,
			expected: wrapper.ErrorRespond{Code: http.StatusInternalServerError, Status: "Internal Server Error", Data: "something went wrong with the server"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			writer := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(writer)
			wrapper.NewHTTPRespondWrapper(c, test.code, test.data, test.args...)
		})
	}
}
