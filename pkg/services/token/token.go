package token

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog"
)

type (
	Token struct {
		ID        int       `json:"id"`
		Name      string    `json:"name"`
		IsSession bool      `json:"isSession"`
		Expires   time.Time `json:"expires"`
		CreatedAt time.Time `json:"createdAt"`
	}

	Interface interface {
		Get(context.Context) ([]Token, error)
		Create(context.Context, string, string) (Token, error)
		Revoke(context.Context, int) error
	}

	Service struct {
		http   *http.Client
		token  string
		url    string
		logger zerolog.Logger
	}
)

func New(
	url string,
	token string,
	logger zerolog.Logger,
) *Service {
	return &Service{
		http:   &http.Client{},
		token:  fmt.Sprintf("Bearer %s", token),
		url:    url,
		logger: logger,
	}
}

func (s *Service) Get(ctx context.Context) ([]Token, error) {
	url := fmt.Sprintf("%s/api/v1/tokens", s.url)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		s.logger.Err(err).Msg("error creating request")
		return nil, err
	}
	req.Header.Add("Authorization", s.token)
	req.Header.Add("Accept", "application/json")

	resp, err := s.http.Do(req)
	if err != nil {
		s.logger.Err(err).Msg("error executing request")
		return nil, err
	}
	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			s.logger.Err(err).Msg("error closing body")
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Err(err).Msg("error reading response body")
		return nil, err
	}

	var responseBody struct {
		Response []Token `json:"response"`
	}

	if err := json.Unmarshal(body, &responseBody); err != nil {
		s.logger.Err(err).Msg("error unmarshalling response body")
		return nil, err
	}

	return responseBody.Response, nil
}

func (s *Service) Create(ctx context.Context, name string, expires int) (string, error) {
	url := fmt.Sprintf("%s/api/v1/tokens", s.url)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBufferString(fmt.Sprintf(`{"name":"%s","expires":%d}`, name, expires)),
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", s.token)
	resp, err := s.http.Do(req)
	if err != nil {
		s.logger.Err(err).Msg("error creating tokens")
		return "", err
	}

	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			s.logger.Err(err).Msg("error closing body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Err(err).Msg("error reading body")
		return "", err
	}

	type createTokenResponse struct {
		Response struct {
			SecretKey string `json:"secretKey"`
			Token     struct {
				ID         int         `json:"id"`
				Name       string      `json:"name"`
				UserID     int         `json:"userId"`
				Token      string      `json:"token"`
				Revoked    bool        `json:"revoked"`
				IsSession  bool        `json:"isSession"`
				Expires    time.Time   `json:"expires"`
				LastUsedAt interface{} `json:"lastUsedAt"`
				CreatedAt  time.Time   `json:"createdAt"`
				UpdatedAt  time.Time   `json:"updatedAt"`
			} `json:"token"`
		} `json:"response"`
	}

	var responseBody createTokenResponse
	if err := json.Unmarshal(body, &responseBody); err != nil {
		s.logger.Err(err).Msg("error unmarshalling body")
		return "", err
	}

	return responseBody.Response.SecretKey, nil
}

func (s *Service) Revoke(ctx context.Context, id int) error {
	url := fmt.Sprintf("%s/api/v1/tokens/%d", s.url, id)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		url,
		nil,
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", s.token)
	resp, err := s.http.Do(req)
	if err != nil {
		s.logger.Err(err).Msg("error revoking token")
		return err
	}

	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			s.logger.Err(err).Msg("error closing body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	return nil
}
