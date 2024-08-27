package utils

import (
	"fmt"
	"net/url"
)

// AddParamsToURL thêm các tham số từ một map vào URL
func AddParamsToURL(baseURL string, params map[string]interface{}) (string, error) {
	// Parse the base URL
	u, err := url.Parse(baseURL)

	if err != nil {
		return "", err
	}

	// Prepare query parameters
	q := u.Query()

	for key, value := range params {
		q.Add(key, fmt.Sprint(value))
	}

	u.RawQuery = q.Encode()

	// Return the full URL with added parameters
	return u.String(), nil
}
