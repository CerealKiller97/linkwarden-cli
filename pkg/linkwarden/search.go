package linkwarden

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

// SearchOptions filters/paginates SearchService.Links.
type SearchOptions struct {
	SearchQueryString string
	// Sort is left nil to omit the query parameter and let the API pick its default.
	Sort   *int
	Cursor int
}

func (o SearchOptions) toQuery() url.Values {
	q := url.Values{}
	setIfNonEmpty(q, "searchQueryString", o.SearchQueryString)
	if o.Sort != nil {
		q.Set("sort", strconv.Itoa(*o.Sort))
	}
	setIfPositive(q, "cursor", o.Cursor)

	return q
}

type searchData struct {
	NextCursor *int   `json:"nextCursor"`
	Links      []Link `json:"links"`
}

// UnmarshalJSON also accepts a bare array, which the API sends as "data"
// instead of an object when nothing matches (`{"data":[],"message":"Nothing found."}`).
func (d *searchData) UnmarshalJSON(b []byte) error {
	if trimmed := bytes.TrimSpace(b); len(trimmed) > 0 && trimmed[0] == '[' {
		*d = searchData{}
		return json.Unmarshal(trimmed, &d.Links)
	}

	type plain searchData
	return json.Unmarshal(b, (*plain)(d))
}

// SearchService accesses the /api/v1/search endpoint.
type SearchService struct {
	c *Client
}

func (s *SearchService) Links(ctx context.Context, opts SearchOptions) ([]Link, error) {
	_, body, err := s.c.request(ctx, requestOptions{
		method: http.MethodGet,
		path:   "/api/v1/search",
		query:  opts.toQuery(),
	})
	if err != nil {
		return nil, err
	}

	data, err := decodeDataEnvelope[searchData](body)
	if err != nil {
		return nil, err
	}

	return data.Links, nil
}
