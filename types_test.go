package infogram_test

import (
	"bytes"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/blainsmith/infogram-go"
	"github.com/frankban/quicktest"
)

func TestInfographic(t *testing.T) {
	c := quicktest.New(t)

	c.Run("json", func(c *quicktest.C) {
		thumbnailURL, err := url.Parse("https://example.com/thumbnail.png")
		c.Assert(err, quicktest.IsNil)

		infogramURL, err := url.Parse("https://example.com/100")
		c.Assert(err, quicktest.IsNil)

		infographic := infogram.Infographic{
			Id:        100,
			Title:     "One Hundred",
			Thumbnail: thumbnailURL,
			ThemeId:   200,
			Published: false,
			URL:       infogramURL,
		}

		var buf bytes.Buffer
		err = json.NewEncoder(&buf).Encode(&infographic)
		c.Assert(err, quicktest.IsNil)

		var newInfographic infogram.Infographic
		err = json.NewDecoder(&buf).Decode(&newInfographic)
		c.Assert(err, quicktest.IsNil)

		c.Assert(newInfographic, quicktest.DeepEquals, infographic)
	})
}

func TestTheme(t *testing.T) {
	c := quicktest.New(t)

	c.Run("json", func(c *quicktest.C) {
		thumbnailURL, err := url.Parse("https://example.com/thumbnail.png")
		c.Assert(err, quicktest.IsNil)

		theme := infogram.Theme{
			Id:        100,
			Title:     "One Hundred",
			Thumbnail: thumbnailURL,
		}

		var buf bytes.Buffer
		err = json.NewEncoder(&buf).Encode(&theme)
		c.Assert(err, quicktest.IsNil)

		var newTheme infogram.Theme
		err = json.NewDecoder(&buf).Decode(&newTheme)
		c.Assert(err, quicktest.IsNil)

		c.Assert(newTheme, quicktest.DeepEquals, theme)
	})
}
