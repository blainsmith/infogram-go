package infogram_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/blainsmith/infogram-go"
	"github.com/frankban/quicktest"
)

type sample struct {
	Id    string `json:"id"`
	Label string `json:"label"`
}

func TestSignRequest(t *testing.T) {
	c := quicktest.New(t)
	client := infogram.Client{Endpoint: infogram.DefaultEndpoint, APIKey: "nMECGhmHe9", APISecret: "da5xoLrCCx"}

	params := make(url.Values)
	params.Add("content", "%5B%7B%22type%22%3A%22h1%22%2C%22text%22%3A%22Hello%20infogr.am%22%7D%5D")
	params.Add("publish", "false")
	params.Add("theme_id", "45")
	params.Add("title", "Hello")

	req := httptest.NewRequest(http.MethodPost, "https://infogr.am/service/v1/infographics", strings.NewReader(params.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(params.Encode())))

	err := client.SignRequest(req)
	c.Assert(err, quicktest.IsNil)

	c.Assert(req.Form.Get("api_key"), quicktest.Equals, client.APIKey)
	c.Assert(req.Form.Get("api_sig"), quicktest.Equals, "bqwCqAk1TWDYNy3eqV0BiNuIERQ=")
}

func TestDo(t *testing.T) {
	c := quicktest.New(t)
	c.Run("> 299 error", func(c *quicktest.C) {
		server := httptest.NewServer(http.NotFoundHandler())
		defer server.Close()

		client := infogram.Client{HTTPClient: server.Client(), Endpoint: server.URL}

		req := httptest.NewRequest(http.MethodGet, server.URL, nil)

		res, err := client.Do(context.Background(), req, nil)
		c.Assert(res, quicktest.IsNil)
		c.Assert(err, quicktest.ErrorMatches, "404 page not found\n")
	})

	c.Run("encode to writer", func(c *quicktest.C) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte("non-json string"))
		}))
		defer server.Close()

		client := infogram.Client{HTTPClient: server.Client(), Endpoint: server.URL}

		req := httptest.NewRequest(http.MethodGet, server.URL, nil)

		buf := bytes.NewBuffer(nil)
		res, err := client.Do(context.Background(), req, buf)
		c.Assert(err, quicktest.IsNil)
		c.Assert(buf.String(), quicktest.Equals, "non-json string")
		c.Assert(res.StatusCode, quicktest.Equals, http.StatusOK)
	})

	c.Run("encode to struct", func(c *quicktest.C) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusOK)
			rw.Write([]byte(`{"id":"1","label":"new label"}`))
		}))
		defer server.Close()

		client := infogram.Client{HTTPClient: server.Client(), Endpoint: server.URL}

		req := httptest.NewRequest(http.MethodGet, server.URL, nil)

		var body sample
		res, err := client.Do(context.Background(), req, &body)
		c.Assert(err, quicktest.IsNil)
		c.Assert(body, quicktest.DeepEquals, sample{Id: "1", Label: "new label"})
		c.Assert(res.StatusCode, quicktest.Equals, http.StatusOK)
	})
}

func TestAPI(t *testing.T) {
	c := quicktest.New(t)
	c.Run("Infographics", func(c *quicktest.C) {
		infographics := []infogram.Infographic{
			{
				Id:        "1",
				Title:     "Number One",
				Thumbnail: &url.URL{Host: "example.com", Path: "/1.png"},
				ThemeId:   99,
				Published: false,
				URL:       &url.URL{Host: "example.com", Path: "/1"},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusOK)
			json.NewEncoder(rw).Encode(infographics)
		}))
		defer server.Close()

		client := infogram.Client{HTTPClient: server.Client(), Endpoint: server.URL}

		data, err := client.Infographics()
		c.Assert(err, quicktest.IsNil)
		c.Assert(data, quicktest.DeepEquals, infographics)
	})

	c.Run("Infographic", func(c *quicktest.C) {
		infographic := infogram.Infographic{
			Id:        "1",
			Title:     "Number One",
			Thumbnail: &url.URL{Host: "example.com", Path: "/1.png"},
			ThemeId:   99,
			Published: false,
			URL:       &url.URL{Host: "example.com", Path: "/1"},
		}

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusOK)
			json.NewEncoder(rw).Encode(&infographic)
		}))
		defer server.Close()

		client := infogram.Client{HTTPClient: server.Client(), Endpoint: server.URL}

		data, err := client.Infographic("1")
		c.Assert(err, quicktest.IsNil)
		c.Assert(data, quicktest.DeepEquals, &infographic)
	})

	c.Run("UserInfographics", func(c *quicktest.C) {
		infographics := []infogram.Infographic{
			{
				Id:        "1",
				Title:     "Number One",
				Thumbnail: &url.URL{Host: "example.com", Path: "/1.png"},
				ThemeId:   99,
				Published: false,
				URL:       &url.URL{Host: "example.com", Path: "/1"},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusOK)
			json.NewEncoder(rw).Encode(infographics)
		}))
		defer server.Close()

		client := infogram.Client{HTTPClient: server.Client(), Endpoint: server.URL}

		data, err := client.UserInfographics("12345")
		c.Assert(err, quicktest.IsNil)
		c.Assert(data, quicktest.DeepEquals, infographics)
	})

	c.Run("Themes", func(c *quicktest.C) {
		themes := []infogram.Theme{
			{
				Id:        1,
				Title:     "Number One",
				Thumbnail: &url.URL{Host: "example.com", Path: "/1.png"},
			},
		}

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(http.StatusOK)
			json.NewEncoder(rw).Encode(themes)
		}))
		defer server.Close()

		client := infogram.Client{HTTPClient: server.Client(), Endpoint: server.URL}

		data, err := client.Themes()
		c.Assert(err, quicktest.IsNil)
		c.Assert(data, quicktest.DeepEquals, themes)
	})
}
