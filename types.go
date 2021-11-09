package infogram

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"
)

type Infoghaphic struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Thumbnail *url.URL  `json:"thumbnail_url"`
	ThemeId   int       `json:"theme_id"`
	Published bool      `json:"published"`
	Modified  time.Time `json:"date_modified"`
	URL       *url.URL  `json:"url"`
}

func (i *Infoghaphic) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})

	data["id"] = i.Id
	data["title"] = i.Title
	data["thumbnail_url"] = i.Thumbnail.String()
	data["theme_id"] = i.ThemeId
	data["published"] = i.Published
	data["url"] = i.URL.String()

	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshalling infographic: %w", err)
	}

	return bytes, nil
}

func (i *Infoghaphic) UnmarshalJSON(bytes []byte) error {
	data := make(map[string]interface{})

	if err := json.Unmarshal(bytes, &data); err != nil {
		return fmt.Errorf("unmarshalling infographic: %w", err)
	}

	if val, found := data["id"]; found {
		v, ok := val.(int)
		if !ok {
			return errors.New("id needs to be an int")
		}
		i.Id = v
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
		v, ok := val.(int)
		if !ok {
			return errors.New("theme_id needs to be an int")
		}
		i.ThemeId = v
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

type Theme struct {
	Id        int      `json:"id"`
	Title     string   `json:"title"`
	Thumbnail *url.URL `json:"thumbnail_url"`
}

func (t *Theme) MarshalJSON() ([]byte, error) {
	data := make(map[string]interface{})

	data["id"] = t.Id
	data["title"] = t.Title
	data["thumbnail_url"] = t.Thumbnail.String()

	return json.Marshal(data)
}

func (t *Theme) UnmarshalJSON(bytes []byte) error {
	data := make(map[string]interface{})

	if val, found := data["id"]; found {
		v, ok := val.(int)
		if !ok {
			return errors.New("id needs to be an int")
		}
		t.Id = v
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
