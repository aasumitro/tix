package token_test

import (
	"github.com/aasumitro/tix/pkg/token"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestJSONWebToken_ClaimJWTToken(t *testing.T) {
	type fields struct {
		Issuer    string
		SecretKey []byte
		Payload   interface{}
		IssuedAt  time.Time
		ExpiredAt time.Time
	}

	issuedAt, _ := time.Parse("",
		"2022-12-05 17:57:44.321843 +0800 WITA m=+25.737606459")
	expiredAt, _ := time.Parse("",
		"2022-12-06 17:57:44.321851 +0800 WITA m=+86425.737614876")

	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "NEW JWT TEST SHOULD SUCCESS",
			fields: fields{
				Issuer: "BECOOP_TEST",
				Payload: map[string]string{
					"data": "hello world",
				},
				SecretKey: []byte("123"),
				IssuedAt:  issuedAt,
				ExpiredAt: expiredAt,
			},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJCRUNPT1BfVEVTVCIsImV4cCI6LTYyMTM1NTk2ODAwLCJpYXQiOi02MjEzNTU5NjgwMCwiZW1haWwiOiIiLCJzZXNzaW9uX2lkIjoiIiwicGF5bG9hZCI6eyJkYXRhIjoiaGVsbG8gd29ybGQifX0.ZKpZgkuL0A9QrLOTSZ7oaZoHMIM96o2NWbBoAXWZkgg",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jwt := &token.JSONWebToken{
				Issuer:    tt.fields.Issuer,
				SecretKey: tt.fields.SecretKey,
				IssuedAt:  tt.fields.IssuedAt,
				ExpiredAt: tt.fields.ExpiredAt,
			}
			got, err := jwt.Claim(tt.fields.Payload)
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.AddCookie(&http.Cookie{Name: "access_token", Value: got})
			if data, err := req.Cookie("access_token"); err == nil {
				_, _ = token.ExtractAndValidateJWT(string(tt.fields.SecretKey), data.String())
			}
			if !tt.wantErr(t, err, "ClaimJWTToken()") {
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
