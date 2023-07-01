package token

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

type IJSONWebToken interface {
	Claim(payload interface{}) (string, error)
	ExtractAndValidateJWT(secret string, cookie *http.Cookie) (claim *JSONWebTokenClaim, err error)
}

type JSONWebTokenClaim struct {
	jwt.RegisteredClaims
	Email     string      `json:"email"`
	SessionID string      `json:"session_id"`
	Payload   interface{} `json:"payload"`
}

type JSONWebToken struct {
	Issuer    string
	SecretKey []byte
	IssuedAt  time.Time
	ExpiredAt time.Time
}

func (j *JSONWebToken) Claim(payload interface{}) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JSONWebTokenClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.Issuer,
			IssuedAt:  &jwt.NumericDate{Time: j.IssuedAt},
			ExpiresAt: &jwt.NumericDate{Time: j.ExpiredAt},
		},
		Payload: payload,
	})

	return token.SignedString(j.SecretKey)
}

func ExtractAndValidateJWT(
	secret string, token string,
) (claim *JSONWebTokenClaim, err error) {
	var parseToken *jwt.Token
	var ok bool

	if parseToken, err = jwt.ParseWithClaims(
		token,
		&JSONWebTokenClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	); err != nil && !parseToken.Valid {
		return nil, err
	}

	if claim, ok = parseToken.Claims.(*JSONWebTokenClaim); !ok {
		return nil, errors.New("invalid claim token")
	}

	return claim, nil
}
