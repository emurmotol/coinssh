package external

import (
	"net/http"
	"time"
	"encoding/json"
	"net/url"
)

func GetJson(url string, target interface{}) error {
	var httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}

	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func PostJson(url string, data url.Values, target interface{}) error {
	var httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}

	r, err := httpClient.PostForm(url, data)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
