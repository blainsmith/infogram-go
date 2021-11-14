package infogram

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"sync"
)

const (
	// DefaultEndpoint is the default API endpoint for Infogram
	DefaultEndpoint = "https://infogr.am/service/v1"
)

// Client is used to interact with the Infogram API
type Client struct {
	setup sync.Once

	HTTPClient *http.Client
	Endpoint   string
	APIKey     string
	APISecret  string
}

// NewRequest creates a proper http.Request for the HTTP call with the correct body data and encoding headers
func (c *Client) NewRequest(method string, path string, params url.Values, body interface{}) (*http.Request, error) {
	path = c.Endpoint + path

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, path, buf)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = params.Encode()

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil

}

// Do performs the *http.Request and decodes the http.Response.Body into v and return the *http.Response. If v is an io.Writer it will copy the body to the writer.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	c.setup.Do(func() {
		if c.HTTPClient == nil {
			c.HTTPClient = http.DefaultClient
		}

		if c.Endpoint == "" {
			c.Endpoint = DefaultEndpoint
		}
	})

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
	}
	defer res.Body.Close()

	if res.StatusCode > 299 {
		errBody, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(string(errBody))
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, res.Body)
		} else {
			decErr := json.NewDecoder(res.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				return nil, decErr
			}
		}
	}

	return res, nil
}

// SignRequest adds the `api_key` and `api_sig` query parameter in accordance with https://developers.infogr.am/rest/request-signing.html
func (c *Client) SignRequest(req *http.Request) error {
	query := req.URL.Query()
	query.Set("api_key", c.APIKey)

	var sig bytes.Buffer
	sig.WriteString(req.Method)
	sig.WriteByte('&')
	sig.WriteString(req.URL.EscapedPath())
	sig.WriteByte('&')

	var queryKeys []string
	for key := range query {
		queryKeys = append(queryKeys, key)
	}
	sort.Slice(queryKeys, func(i int, j int) bool { return queryKeys[i] < queryKeys[j] })

	var params bytes.Buffer
	for idx, key := range queryKeys {
		params.WriteString(url.QueryEscape(key))
		params.WriteByte('=')
		params.WriteString(url.QueryEscape(query.Get(key)))
		if idx < len(queryKeys) {
			params.WriteByte('&')
		}
	}

	sig.WriteString(url.QueryEscape(params.String()))

	h := hmac.New(sha1.New, []byte(c.APISecret))
	signature := h.Sum(sig.Bytes())

	query.Set("api_sig", fmt.Sprintf("%x", signature))
	req.URL.RawQuery = query.Encode()

	return nil
}

// Infographics fetches the list of infographics
func (c *Client) Infographics() ([]Infographic, error) {
	req, err := c.NewRequest(http.MethodGet, "/infographics", nil, nil)
	if err != nil {
		return nil, fmt.Errorf("new infographics request: %w", err)
	}

	err = c.SignRequest(req)
	if err != nil {
		return nil, err
	}

	var infographics []Infographic
	_, err = c.Do(context.Background(), req, &infographics)
	if err != nil {
		return nil, err
	}

	return infographics, nil
}

// Infographics fetches a single infographic by identification number
func (c *Client) Infographic(id int) (*Infographic, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/%d?api_key=%s", c.Endpoint, "infographics", id, c.APIKey), nil)
	if err != nil {
		return nil, fmt.Errorf("new infographic request: %w", err)
	}

	var infographic Infographic
	_, err = c.Do(context.Background(), req, &infographic)
	if err != nil {
		return nil, nil
	}

	return &infographic, nil
}

// UserInfographics fetches the list of infographics for the user's identification number
func (c *Client) UserInfographics(id string) ([]Infographic, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/%s/%s?api_key=%s", c.Endpoint, "users", id, "infographics", c.APIKey), nil)
	if err != nil {
		return nil, fmt.Errorf("new infographics request: %w", err)
	}

	var infographics []Infographic
	_, err = c.Do(context.Background(), req, &infographics)
	if err != nil {
		return nil, nil
	}

	return infographics, nil
}

// Infographics fetches a available themes to use for infographics
func (c *Client) Themes() ([]Theme, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s?api_key=%s", c.Endpoint, "themes", c.APIKey), nil)
	if err != nil {
		return nil, fmt.Errorf("new themes request: %w", err)
	}

	var themes []Theme
	_, err = c.Do(context.Background(), req, &themes)
	if err != nil {
		return nil, nil
	}

	return themes, nil
}
