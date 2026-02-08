package tag

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Item struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type (
	Interface interface {
		Get(context.Context) ([]Item, error)
	}

	Service struct {
		httpClient *http.Client
		token      string
		url        string
		logger     zerolog.Logger
	}
)

var _ Interface = &Service{}

func New(url, token string, logger zerolog.Logger) *Service {
	return &Service{
		httpClient: &http.Client{},
		token:      fmt.Sprintf("Bearer %s", token),
		url:        url,
		logger:     logger,
	}
}

func (s *Service) Get(ctx context.Context) ([]Item, error) {
	url := fmt.Sprintf("%s/api/v1/tags", s.url)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		log.Err(err).Msg("Failed to create request")
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", s.token)

	res, err := s.httpClient.Do(req)
	if err != nil {
		log.Err(err).Msg("Failed to create link")
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			s.logger.Err(err).Msg("Failed to close body")
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		s.logger.Err(err).Msg("Failed to read response body")
	}

	var items []Item
	if err := json.Unmarshal(body, &items); err != nil {
		s.logger.Err(err).Msg("Failed to unmarshal response body")
		return nil, err
	}

	return items, nil
}
