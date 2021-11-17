package infogram

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Infographic defines the type returned by the Infogram API
type Infographic struct {
	Id        int
	Title     string
	Thumbnail *url.URL
	ThemeId   int
	Published bool
	Modified  time.Time
	URL       *url.URL
}

func (i *Infographic) reader(client *Client, format string) (io.Reader, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/%d?api_key=%s&format=%s", client.Endpoint, "infographics", i.Id, client.APIKey, format), nil)
	if err != nil {
		return nil, fmt.Errorf("new infographic PDF reader request: %w", err)
	}

	err = client.SignRequest(req)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(context.Background(), req, nil)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

// PDFReader returns an io.Reader of the Infographic in PDF format
func (i *Infographic) PDFReader(client *Client) (io.Reader, error) {
	return i.reader(client, "pdf")
}

// PNGReader returns an io.Reader of the Infographic in PNG format
func (i *Infographic) PNGReader(client *Client) (io.Reader, error) {
	return i.reader(client, "png")
}

// HTMLReader returns an io.Reader of the Infographic in HTML format
func (i *Infographic) HTMLReader(client *Client) (io.Reader, error) {
	return i.reader(client, "html")
}

// MarshalJSON implements json.Marshaler
func (i *Infographic) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})

	data["id"] = i.Id
	data["title"] = i.Title
	if i.Thumbnail != nil {
		data["thumbnail_url"] = i.Thumbnail.String()
	}
	data["theme_id"] = i.ThemeId
	data["published"] = i.Published
	if i.URL != nil {
		data["url"] = i.URL.String()
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshalling infographic: %w", err)
	}

	return bytes, nil
}

// UnmarshalJSON implements json.Unarshaler
func (i *Infographic) UnmarshalJSON(bytes []byte) error {
	data := make(map[string]interface{})

	if err := json.Unmarshal(bytes, &data); err != nil {
		return fmt.Errorf("unmarshalling infographic: %w", err)
	}

	if val, found := data["id"]; found {
		v, ok := val.(float64)
		if !ok {
			return errors.New("id needs to be an int")
		}
		i.Id = int(v)
	}
	if val, found := data["title"]; found {
		v, ok := val.(string)
		if !ok {
			return errors.New("title needs to be an string")
		}
		i.Title = v
	}
	if val, found := data["thumbnail_url"]; found {
		v, ok := val.(string)
		if !ok {
			return errors.New("thumbnail_url needs to be an string")
		}
		var err error
		i.Thumbnail, err = url.Parse(v)
		if err != nil {
			return errors.New("thumbnail_url needs to be a parsable URL")
		}
	}
	if val, found := data["theme_id"]; found {
		v, ok := val.(float64)
		if !ok {
			return errors.New("theme_id needs to be an int")
		}
		i.ThemeId = int(v)
	}
	if val, found := data["published"]; found {
		v, ok := val.(bool)
		if !ok {
			return errors.New("published needs to be an boolean")
		}
		i.Published = v
	}
	if val, found := data["date_modified"]; found {
		v, ok := val.(string)
		if !ok {
			return errors.New("date_modified needs to be an string")
		}
		var err error
		i.Modified, err = time.Parse(time.RFC3339, v)
		if err != nil {
			return errors.New("date_modified needs to be a parsable RFC 3339 time")
		}
	}
	if val, found := data["url"]; found {
		v, ok := val.(string)
		if !ok {
			return errors.New("url needs to be an string")
		}
		var err error
		i.URL, err = url.Parse(v)
		if err != nil {
			return errors.New("url needs to be a parsable URL")
		}
	}

	return nil
}

// Theme defines the type returned by the Infogram API
type Theme struct {
	Id        int
	Title     string
	Thumbnail *url.URL
}

// MarshalJSON implements json.Marshaler
func (t *Theme) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})

	data["id"] = t.Id
	data["title"] = t.Title
	data["thumbnail_url"] = t.Thumbnail.String()

	return json.Marshal(data)
}

// UnmarshalJSON implements json.Unmarshaler
func (t *Theme) UnmarshalJSON(bytes []byte) error {
	data := make(map[string]interface{})

	if err := json.Unmarshal(bytes, &data); err != nil {
		return fmt.Errorf("unmarshalling theme: %w", err)
	}

	if val, found := data["id"]; found {
		v, ok := val.(float64)
		if !ok {
			return errors.New("id needs to be an int")
		}
		t.Id = int(v)
	}
	if val, found := data["title"]; found {
		v, ok := val.(string)
		if !ok {
			return errors.New("title needs to be an string")
		}
		t.Title = v
	}
	if val, found := data["thumbnail_url"]; found {
		v, ok := val.(string)
		if !ok {
			return errors.New("thumbnail_url needs to be an string")
		}
		var err error
		t.Thumbnail, err = url.Parse(v)
		if err != nil {
			return errors.New("thumbnail_url needs to be a parsable URL")
		}
	}

	return nil
}
