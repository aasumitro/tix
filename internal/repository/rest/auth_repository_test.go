package rest_test

import (
	"context"
	"github.com/aasumitro/tix/internal/domain/response"
	"github.com/aasumitro/tix/internal/repository/rest"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type testcase struct {
	name        string
	response    string
	endpoint    string
	status      int
	expected    interface{}
	method      string
	wantErr     bool
	setupServer func(*httptest.Server, string)
}

type authRESTRepositoryTestSuite struct {
	suite.Suite
}

func (s *authRESTRepositoryTestSuite) Test_SendMagicLink() {
	testcases := []testcase{
		{
			name:     "valid request",
			response: `{"code": 200, "msg":"lorem"}`,
			status:   200,
			method:   http.MethodPost,
			expected: &response.SupabaseRespond{Code: 200, Message: "lorem"},
			wantErr:  false,
		},
		{
			name:     "valid request",
			response: `{}`,
			status:   200,
			method:   http.MethodPost,
			expected: &response.SupabaseRespond{Code: 200, Message: "OK"},
			wantErr:  false,
		},
		{
			method:   http.MethodPost,
			name:     "invalid endpoint",
			endpoint: "invalid endpoint",
			wantErr:  true,
		},
		{
			method:   http.MethodPost,
			name:     "nil endpoint",
			endpoint: "",
			wantErr:  true,
		},
		{
			method:   http.MethodPost,
			name:     "invalid JSON response",
			response: `{"foo": "bar`,
			wantErr:  true,
		},
		{
			method:   "lorem_ipsum",
			name:     "invalid method",
			response: `{"foo": "bar`,
			wantErr:  true,
		},
		{
			name:     "error in http.NewRequestWithContext",
			response: `{"code": 200, "msg":"lorem"}`,
			wantErr:  true,
		},
		{
			name:     "invalid JSON response",
			response: `{"foo": "bar`,
			wantErr:  true,
			setupServer: func(server *httptest.Server, response string) {
				server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(response))
				})
			},
		},
	}

	for _, ts := range testcases {
		s.T().Run(ts.name, func(t *testing.T) {
			var server *httptest.Server
			if ts.setupServer != nil {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					ts.setupServer(server, ts.response)
				}))
			} else {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(ts.status)
					_, _ = w.Write([]byte(ts.response))
				}))
			}
			defer server.Close()
			api := rest.NewAuthRESTRepository(
				server.URL, "", "")
			data, err := api.SendMagicLink(context.Background(), "hello@tix.id")
			if (err != nil) != ts.wantErr {
				t.Errorf("unexpected error: %v", err)
			}

			if err == nil && !reflect.DeepEqual(data, ts.expected) {
				t.Errorf("expected %v, got %v", ts.expected, data)
			}
		})
	}
}

func (s *authRESTRepositoryTestSuite) Test_InviteUserByEmail() {
	testcases := []testcase{
		{
			name:     "valid request",
			response: `{"code": 200, "msg":"lorem"}`,
			status:   200,
			method:   http.MethodPost,
			expected: &response.SupabaseRespond{Code: 200, Message: "lorem"},
			wantErr:  false,
		},
		{
			name:     "valid request",
			response: `{}`,
			status:   200,
			method:   http.MethodPost,
			expected: &response.SupabaseRespond{Code: 200, Message: "OK"},
			wantErr:  false,
		},
		{
			method:   http.MethodPost,
			name:     "invalid endpoint",
			endpoint: "invalid endpoint",
			wantErr:  true,
		},
		{
			method:   http.MethodPost,
			name:     "nil endpoint",
			endpoint: "",
			wantErr:  true,
		},
		{
			method:   http.MethodPost,
			name:     "invalid JSON response",
			response: `{"foo": "bar`,
			wantErr:  true,
		},
		{
			method:   "lorem_ipsum",
			name:     "invalid method",
			response: `{"foo": "bar`,
			wantErr:  true,
		},
		{
			name:     "invalid JSON response",
			response: `{"foo": "bar`,
			wantErr:  true,
			setupServer: func(server *httptest.Server, response string) {
				server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(response))
				})
			},
		},
	}

	for _, ts := range testcases {
		s.T().Run(ts.name, func(t *testing.T) {
			var server *httptest.Server
			if ts.setupServer != nil {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					ts.setupServer(server, ts.response)
				}))
			} else {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(ts.status)
					_, _ = w.Write([]byte(ts.response))
				}))
			}
			defer server.Close()
			api := rest.NewAuthRESTRepository(
				server.URL, "", "")
			data, err := api.InviteUserByEmail(context.Background(), "hello@tix.id")
			if (err != nil) != ts.wantErr {
				t.Errorf("unexpected error: %v", err)
			}

			if err == nil && !reflect.DeepEqual(data, ts.expected) {
				t.Errorf("expected %v, got %v", ts.expected, data)
			}
		})
	}
}

func (s *authRESTRepositoryTestSuite) Test_DeleteUser() {
	testcases := []testcase{
		{
			name:     "valid request",
			response: `{"code": 200, "msg":"lorem"}`,
			status:   200,
			method:   http.MethodDelete,
			expected: &response.SupabaseRespond{Code: 200, Message: "lorem"},
			wantErr:  false,
		},
		{
			name:     "valid request",
			response: `{}`,
			status:   200,
			method:   http.MethodDelete,
			expected: &response.SupabaseRespond{Code: 200, Message: "OK"},
			wantErr:  false,
		},
		{
			method:   http.MethodDelete,
			name:     "invalid endpoint",
			endpoint: "invalid endpoint",
			wantErr:  true,
		},
		{
			method:   http.MethodDelete,
			name:     "nil endpoint",
			endpoint: "",
			wantErr:  true,
		},
		{
			method:   http.MethodDelete,
			name:     "invalid JSON response",
			response: `{"foo": "bar`,
			wantErr:  true,
		},
		{
			method:   "lorem_ipsum",
			name:     "invalid method",
			response: `{"foo": "bar`,
			wantErr:  true,
		},
		{
			name:     "invalid JSON response",
			response: `{"foo": "bar`,
			wantErr:  true,
			setupServer: func(server *httptest.Server, response string) {
				server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(response))
				})
			},
		},
	}

	for _, ts := range testcases {
		s.T().Run(ts.name, func(t *testing.T) {
			var server *httptest.Server
			if ts.setupServer != nil {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					ts.setupServer(server, ts.response)
				}))
			} else {
				server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(ts.status)
					_, _ = w.Write([]byte(ts.response))
				}))
			}
			defer server.Close()
			api := rest.NewAuthRESTRepository(
				server.URL, "", "")
			data, err := api.DeleteUser(context.Background(), "12345")
			if (err != nil) != ts.wantErr {
				t.Errorf("unexpected error: %v", err)
			}

			if err == nil && !reflect.DeepEqual(data, ts.expected) {
				t.Errorf("expected %v, got %v", ts.expected, data)
			}
		})
	}
}

func TestAuthRESTRepository(t *testing.T) {
	suite.Run(t, new(authRESTRepositoryTestSuite))
}
