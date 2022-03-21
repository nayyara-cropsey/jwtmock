package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/mitchellh/mapstructure"
)

func formUnmarshal(r *http.Request, v interface{}) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}

	rawMap, err := url.ParseQuery(string(body))
	if err != nil {
		return fmt.Errorf("form unmarshal: %w", err)
	}

	rawMapSingle := make(map[string]interface{})
	for k, v := range rawMap {
		rawMapSingle[k] = v[0]
	}

	return mapstructure.Decode(rawMapSingle, &v)
}
