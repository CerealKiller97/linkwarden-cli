package migration

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"net/http"
)

type (
	Interface interface {
		Export(context.Context) ([]byte, string, error)
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

func (s *Service) Export(ctx context.Context) ([]byte, string, error) {
	url := fmt.Sprintf("%s/api/v1/migration", s.url)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		s.logger.Error().Err(err).Msg("error creating request")
		return nil, "", err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", s.token)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		s.logger.Error().Err(err).Msg("error making request")
		return nil, "", err
	}

	contentDispositionHeader := resp.Header.Get("Content-Disposition")

	backupFile := contentDispositionHeader[len(contentDispositionHeader)-11:]
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return data, backupFile, nil
}
