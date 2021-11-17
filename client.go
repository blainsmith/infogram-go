package infogram

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
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

	req.RequestURI = ""

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
		_, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("code: %d, message: ", res.StatusCode)
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
	var data url.Values

	switch req.Method {
	case http.MethodGet, http.MethodDelete:
		data = req.URL.Query()
	default:
		req.ParseForm()
		data = req.Form
	}

	data.Set("api_key", c.APIKey)

	var sig bytes.Buffer
	sig.WriteString(req.Method)
	sig.WriteByte('&')
	sig.WriteString(url.QueryEscape(req.URL.String()))
	sig.WriteByte('&')

	var dataKeys []string
	for key := range data {
		dataKeys = append(dataKeys, key)
	}
	sort.Slice(dataKeys, func(i int, j int) bool { return dataKeys[i] < dataKeys[j] })

	var params bytes.Buffer
	for idx, key := range dataKeys {
		params.WriteString(key)
		params.WriteString("=")
		params.WriteString(data.Get(key))
		if idx < len(dataKeys)-1 {
			params.WriteString("&")
		}
	}

	sig.WriteString(url.QueryEscape(params.String()))

	h := hmac.New(sha1.New, []byte(c.APISecret))
	h.Write(sig.Bytes())
	signature := h.Sum(nil)

	data.Add("api_sig", base64.StdEncoding.EncodeToString(signature))

	switch req.Method {
	case http.MethodGet, http.MethodDelete:
		req.URL.RawQuery = data.Encode()
	default:
		req.Form = data
	}

	return nil
}

// Infographics fetches the list of infographics
func (c *Client) Infographics() ([]Infographic, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.Endpoint, "infographics"), nil)
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
func (c *Client) Infographic(id string) (*Infographic, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/%s", c.Endpoint, "infographics", id), nil)
	if err != nil {
		return nil, fmt.Errorf("new infographic request: %w", err)
	}

	err = c.SignRequest(req)
	if err != nil {
		return nil, err
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
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/%s/%s", c.Endpoint, "users", id, "infographics"), nil)
	if err != nil {
		return nil, fmt.Errorf("new user infographics request: %w", err)
	}

	err = c.SignRequest(req)
	if err != nil {
		return nil, err
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
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", c.Endpoint, "themes"), nil)
	if err != nil {
		return nil, fmt.Errorf("new themes request: %w", err)
	}

	err = c.SignRequest(req)
	if err != nil {
		return nil, err
	}

	var themes []Theme
	_, err = c.Do(context.Background(), req, &themes)
	if err != nil {
		return nil, err
	}

	return themes, nil
}
