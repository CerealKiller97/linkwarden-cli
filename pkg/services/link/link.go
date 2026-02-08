package link

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/CerealKiller97/linkwarden-cli/pkg/dto"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Interface interface {
	Get(any) []any
	Find(int) any
	Create(ctx context.Context, request dto.CreateLinkRequest) (any, error)
	UpdateBulk([]any) []any
	Update(any) any
	Archive(int) string
	Delete(int) int
	DeleteBulk([]int) any
}

type Service struct {
	instanceURL string
	httpClient  *http.Client
	logger      zerolog.Logger
	token       string
}

func New(instanceURL string, token string, logger zerolog.Logger) *Service {
	return &Service{
		instanceURL: instanceURL,
		token:       fmt.Sprintf("Bearer %s", token),
		httpClient:  &http.Client{},
		logger:      logger,
	}
}

func (s *Service) Create(ctx context.Context, request dto.CreateLinkRequest) (any, error) {
	url := s.instanceURL + "/api/v1/links"
	body, err := json.Marshal(request)
	if err != nil {
		log.Err(err).Msg("failed to marshal request")
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(body),
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

	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			log.Err(err).Msg("Failed to close body")
		}
	}(res.Body)

	if res.StatusCode == http.StatusConflict {
		return nil, errors.New("link already exists")
	}

	body, err = io.ReadAll(res.Body)
	if err != nil {
		log.Err(err).Msg("Failed to read response body")
		return nil, err
	}

	log.Debug().Str("body", string(body)).Msg("Got response body")

	if res.StatusCode != http.StatusOK {
		log.Error().Int("statusCode", res.StatusCode).Msg("Failed to create a link")
		return nil, errors.New(res.Status)
	}

	return body, nil
}
