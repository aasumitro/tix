package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aasumitro/tix/common"
	"github.com/aasumitro/tix/internal/domain"
	"github.com/aasumitro/tix/internal/domain/response"
	"io"
	"net/http"
	"time"
)

type authRESTRepository struct {
	supabaseAPIURL, supabaseAPIKey,
	supabaseRootAPIKey string
}

func (repository *authRESTRepository) SendMagicLink(
	ctx context.Context,
	email string,
) (data *response.SupabaseRespond, err error) {
	reqBody, err := json.Marshal(map[string]string{"email": email})
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(reqBody)
	endpoint := fmt.Sprintf("%s/%s/magiclink",
		repository.supabaseAPIURL, common.SupabaseAuthEndpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"apiKey":       []string{repository.supabaseAPIKey},
		"Content-Type": []string{"application/json"},
	}
	c := &http.Client{Timeout: 10 * time.Second}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return nil, err
	}

	var res interface{}
	if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
		return nil, err
	}
	if rsp, ok := res.(map[string]interface{}); ok {
		data = &response.SupabaseRespond{
			Code: func() int {
				if code, ok := rsp["code"].(float64); ok {
					return int(code)
				}
				return http.StatusOK
			}(),
			Message: func() string {
				if message, ok := rsp["msg"].(string); ok {
					return message
				}
				return "OK"
			}(),
		}
	}

	return data, nil
}

func (repository *authRESTRepository) InviteUserByEmail(
	ctx context.Context,
	email string,
) (data *response.SupabaseRespond, err error) {
	reqBody, err := json.Marshal(map[string]string{"email": email})
	if err != nil {
		return nil, err
	}
	body := bytes.NewBuffer(reqBody)
	endpoint := fmt.Sprintf("%s/%s/invite",
		repository.supabaseAPIURL, common.SupabaseAuthEndpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", repository.supabaseRootAPIKey)},
		"apiKey":        []string{repository.supabaseAPIKey},
		"Content-Type":  []string{"application/json"},
	}
	c := &http.Client{Timeout: 10 * time.Second}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return nil, err
	}

	var res interface{}
	if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
		return nil, err
	}
	if rsp, ok := res.(map[string]interface{}); ok {
		data = &response.SupabaseRespond{
			Code: func() int {
				if code, ok := rsp["code"].(float64); ok {
					return int(code)
				}
				return http.StatusOK
			}(),
			Message: func() string {
				if message, ok := rsp["msg"].(string); ok {
					return message
				}
				return "OK"
			}(),
		}
	}

	return data, nil
}

func (repository *authRESTRepository) DeleteUser(
	ctx context.Context,
	uuid string,
) (data *response.SupabaseRespond, err error) {
	endpoint := fmt.Sprintf("%s/%s/admin/users/%s",
		repository.supabaseAPIURL, common.SupabaseAuthEndpoint, uuid)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Authorization": []string{fmt.Sprintf("Bearer %s", repository.supabaseRootAPIKey)},
		"apiKey":        []string{repository.supabaseAPIKey},
		"Content-Type":  []string{"application/json"},
	}
	c := &http.Client{Timeout: 10 * time.Second}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) { _ = Body.Close() }(resp.Body)

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return nil, err
	}

	var res interface{}
	if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
		return nil, err
	}
	if rsp, ok := res.(map[string]interface{}); ok {
		data = &response.SupabaseRespond{
			Code: func() int {
				if code, ok := rsp["code"].(float64); ok {
					return int(code)
				}
				return http.StatusOK
			}(),
			Message: func() string {
				if message, ok := rsp["msg"].(string); ok {
					return message
				}
				return "OK"
			}(),
		}
	}

	return data, nil
}

func NewAuthRESTRepository(
	supabaseAPIURL, supabaseAPIKey,
	supabaseRootAPIKey string,
) domain.IAuthRESTRepository {
	return &authRESTRepository{
		supabaseAPIURL:     supabaseAPIURL,
		supabaseAPIKey:     supabaseAPIKey,
		supabaseRootAPIKey: supabaseRootAPIKey,
	}
}
