package collection

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog"
)

type (
	Collection struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	Interface interface {
		Get(context.Context) ([]Collection, error)
	}

	Service struct {
		url        string
		httpClient *http.Client
		token      string
		logger     zerolog.Logger
	}
)

func New(
	url string,
	token string,
	logger zerolog.Logger,
) *Service {
	return &Service{
		url:        url,
		httpClient: &http.Client{},
		token:      fmt.Sprintf("Bearer %s", token),
		logger:     logger,
	}
}

func (s *Service) Get(ctx context.Context) ([]Collection, error) {
	url := fmt.Sprintf("%s/api/v1/collections", s.url)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		s.logger.Error().Err(err).Msg("error creating request")
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", s.token)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Error().Err(err).Msg("error making request")
		return nil, err
	}

	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			s.logger.Err(err).Msg("error closing body")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.Err(err).Msg("error reading body")
		return nil, err
	}

	type response struct {
		Response []Collection `json:"response"`
	}

	var responseBody response
	if err := json.Unmarshal(body, &responseBody); err != nil {
		s.logger.Err(err).Msg("error unmarshalling collections")
		return nil, err
	}

	return responseBody.Response, nil
}
