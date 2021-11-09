package infogram

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
)

const (
	// DefaultEndpoint is the default API endpoint for Infogram
	DefaultEndpoint = "https://infogr.am/service/v1"
)

// ClientOpts defines optional configuration settings when creating a Client
type ClientOpts func(*Client)

// Client is used to interact with the Infogram API
type Client struct {
	httpClient *http.Client
	endpoint   string
	apiKey     string
	apiSecret  string
}

// NewClient creates a Client with the specified API key and secret and any ClientOpts provided
func NewClient(apiKey string, apiSecret string, options ...ClientOpts) *Client {
	c := Client{
		httpClient: http.DefaultClient,
		endpoint:   DefaultEndpoint,
		apiKey:     apiKey,
		apiSecret:  apiSecret,
	}

	for _, opts := range options {
		opts(&c)
	}

	return &c
}

// ClientOptHTTPClient overrides the http.DefaultClient with the one specified
func ClientOptHTTPClient(httpClient *http.Client) func(*Client) {
	return func(client *Client) {
		client.httpClient = httpClient
	}
}

// ClientOptEndpoint overrides the DefaultEndpoint with the one specified
func ClientOptEndpoint(endpoint string) func(*Client) {
	return func(client *Client) {
		client.endpoint = endpoint
	}
}

// SignRequest adds the `api_sig` query parameter in accordance with https://developers.infogr.am/rest/request-signing.html
func (c *Client) SignRequest(req *http.Request) error {
	query := req.URL.Query()

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

	h := hmac.New(sha1.New, []byte(c.apiSecret))
	signature := h.Sum(sig.Bytes())

	query.Set("api_sig", fmt.Sprintf("%x", signature))
	req.URL.RawQuery = query.Encode()

	return nil
}

func (c *Client) signAndDo(req *http.Request) (*http.Response, error) {
	if err := c.SignRequest(req); err != nil {
		return nil, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Infographics fetches the list of infographics
func (c *Client) Infographics() ([]Infographic, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s?api_key=%s", c.endpoint, "infographics", c.apiKey), nil)
	if err != nil {
		return nil, fmt.Errorf("new infographics request: %w", err)
	}

	res, err := c.signAndDo(req)
	if err != nil {
		return nil, err
	}

	var infographics []Infographic
	if err := json.NewDecoder(res.Body).Decode(&infographics); err != nil {
		return nil, err
	}

	return infographics, nil
}

// Infographics fetches a single infographic by identification number
func (c *Client) Infographic(id int) (*Infographic, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/%d?api_key=%s", c.endpoint, "infographics", id, c.apiKey), nil)
	if err != nil {
		return nil, fmt.Errorf("new infographic request: %w", err)
	}

	res, err := c.signAndDo(req)
	if err != nil {
		return nil, nil
	}

	var infographic Infographic
	if err := json.NewDecoder(res.Body).Decode(&infographic); err != nil {
		return nil, err
	}

	return &infographic, nil
}

// UserInfographics fetches the list of infographics for the user's identification number
func (c *Client) UserInfographics(id string) ([]Infographic, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/%s/%s?api_key=%s", c.endpoint, "users", id, "infographics", c.apiKey), nil)
	if err != nil {
		return nil, fmt.Errorf("new infographics request: %w", err)
	}

	res, err := c.signAndDo(req)
	if err != nil {
		return nil, err
	}

	var infographics []Infographic
	if err := json.NewDecoder(res.Body).Decode(&infographics); err != nil {
		return nil, err
	}

	return infographics, nil
}

// Infographics fetches a available themes to use for infographics
func (c *Client) Themes() ([]Theme, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s?api_key=%s", c.endpoint, "themes", c.apiKey), nil)
	if err != nil {
		return nil, fmt.Errorf("new themes request: %w", err)
	}

	res, err := c.signAndDo(req)
	if err != nil {
		return nil, err
	}

	var themes []Theme
	if err := json.NewDecoder(res.Body).Decode(&themes); err != nil {
		return nil, err
	}

	return themes, nil
}
