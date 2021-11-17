package infogram_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/blainsmith/infogram-go"
	"github.com/frankban/quicktest"
)

type sample struct {
	Id    string `json:"id"`
	Label string `json:"label"`
}

func TestClient(t *testing.T) {
	c := quicktest.New(t)

	c.Run("NewRequest", func(c *quicktest.C) {
		client := infogram.Client{Endpoint: infogram.DefaultEndpoint}

		c.Run("simple", func(c *quicktest.C) {
			req, err := client.NewRequest(http.MethodGet, "/infographics", nil, nil)
			c.Assert(err, quicktest.IsNil)

			c.Assert(req.Method, quicktest.Equals, http.MethodGet)
			c.Assert(req.Header.Get("Content-Type"), quicktest.Equals, "")
			c.Assert(req.Body, quicktest.IsNil)
		})

		c.Run("json body", func(c *quicktest.C) {
			body := sample{
				Id:    "123",
				Label: "New Label",
			}
			req, err := client.NewRequest(http.MethodPost, "/infographics", nil, body)
			c.Assert(err, quicktest.IsNil)

			c.Assert(req.Method, quicktest.Equals, http.MethodPost)
			c.Assert(req.Header.Get("Content-Type"), quicktest.Equals, "application/json")

			reqBody, err := io.ReadAll(req.Body)
			c.Assert(err, quicktest.IsNil)
			c.Assert(reqBody, quicktest.JSONEquals, &sample{Id: "123", Label: "New Label"})
		})

		c.Run("query params", func(c *quicktest.C) {
			qs := make(url.Values)
			qs.Add("id", "1")
			qs.Add("label", "new label")

			req, err := client.NewRequest(http.MethodGet, "/infographics", qs, nil)
			c.Assert(err, quicktest.IsNil)

			c.Assert(req.Method, quicktest.Equals, http.MethodGet)
			c.Assert(req.Header.Get("Content-Type"), quicktest.Equals, "")
			c.Assert(req.Body, quicktest.IsNil)

			c.Assert(req.URL.Query().Get("id"), quicktest.Equals, qs.Get("id"))
			c.Assert(req.URL.Query().Get("label"), quicktest.Equals, qs.Get("label"))
		})
	})

	c.Run("SignRequest", func(c *quicktest.C) {
		client := infogram.Client{Endpoint: infogram.DefaultEndpoint, APIKey: "test-key", APISecret: "shh"}

		qs := make(url.Values)
		qs.Add("id", "1")
		qs.Add("label", "new label")

		req, err := client.NewRequest(http.MethodGet, "/infographics", qs, nil)
		c.Assert(err, quicktest.IsNil)

		err = client.SignRequest(req)
		c.Assert(err, quicktest.IsNil)

		for key := range qs {
			c.Assert(req.URL.Query().Get(key), quicktest.Equals, qs.Get(key))
		}

		c.Assert(req.URL.Query().Get("api_key"), quicktest.Equals, client.APIKey)
		c.Assert(req.URL.Query().Get("api_sig"), quicktest.Equals, "474554262f736572766963652f76312f696e666f6772617068696373266170695f6b6579253344746573742d6b65792532366964253344312532366c6162656c2533446e65772532426c6162656c253236a6b812acfa12a677bb2fe6b266bac6a7294d9d06")
	})

	c.Run("Do", func(c *quicktest.C) {
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
	})

}

func TestAPI(t *testing.T) {
	c := quicktest.New(t)
	c.Run("Convenience", func(c *quicktest.C) {
		c.Run("Infographics", func(c *quicktest.C) {
			infographics := []infogram.Infographic{
				{
					Id:        1,
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
				Id:        1,
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

			data, err := client.Infographic(1)
			c.Assert(err, quicktest.IsNil)
			c.Assert(data, quicktest.DeepEquals, &infographic)
		})

		c.Run("UserInfographics", func(c *quicktest.C) {
			infographics := []infogram.Infographic{
				{
					Id:        1,
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
	})
}
