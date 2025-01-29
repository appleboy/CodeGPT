package proxy

import (
	"net/http"
	"reflect"
	"testing"
)

func TestConvertHeaders(t *testing.T) {
	testCases := []struct {
		headers  []string
		expected http.Header
	}{
		{
			headers: []string{"Content-Type=application/json", "Authorization=Bearer token"},
			expected: http.Header{
				"Content-Type":  []string{"application/json"},
				"Authorization": []string{"Bearer token"},
			},
		},
		{
			headers: []string{"X-Custom-Header=custom_value", "InvalidHeader"},
			expected: http.Header{
				"X-Custom-Header": []string{"custom_value"},
			},
		},
		{
			headers:  []string{"KeyOnly=", "=ValueOnly"},
			expected: http.Header{},
		},
		{
			headers:  []string{},
			expected: http.Header{},
		},
	}

	for _, tc := range testCases {
		result := convertHeaders(tc.headers)
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("convertHeaders(%v) = %v, expected %v", tc.headers, result, tc.expected)
		}
	}
}
