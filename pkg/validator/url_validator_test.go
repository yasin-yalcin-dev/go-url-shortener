/*
 ** ** ** ** ** **
  \ \ / / \ \ / /
   \ V /   \ V /
    | |     | |
    |_|     |_|
   Yasin   Yalcin
*/

package validator

import (
	"testing"
)

func TestURLValidator_Validate(t *testing.T) {
	// Test cases
	testCases := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "Valid HTTP URL",
			url:     "http://google.com",
			wantErr: false,
		},
		{
			name:    "Valid HTTPS URL",
			url:     "https://www.google.com",
			wantErr: false,
		},
		{
			name:    "Empty URL",
			url:     "",
			wantErr: true,
		},
		{
			name:    "Invalid URL Format",
			url:     "not-a-url",
			wantErr: true,
		},
		{
			name:    "URL Without Scheme",
			url:     "example.com",
			wantErr: true,
		},
		{
			name:    "Very Long URL",
			url:     "http://" + string(make([]byte, 2500)),
			wantErr: true,
		},
	}

	// Create a new URLValidator instance
	validator := NewURLValidator()

	// Run through each test case
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.Validate(tc.url)

			// Control if the error matches the expected result
			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

func TestURLValidator_BlockedDomains(t *testing.T) {
	// Define test cases for blocked domains
	validator := NewURLValidator()

	blockedTestCases := []struct {
		name string
		url  string
	}{
		{
			name: "Blocked Example Domain",
			url:  "http://specific-blocked-domain.com/test",
		},
	}

	for _, tc := range blockedTestCases {
		t.Run(tc.name, func(t *testing.T) {
			validator.blockedDomains = []string{"specific-blocked-domain.com"}
			err := validator.Validate(tc.url)

			if err == nil {
				t.Errorf("Expected blocked domain error for URL: %s", tc.url)
			}
		})
	}
}

func TestURLValidator_AddBlockedDomain(t *testing.T) {
	validator := NewURLValidator()

	// Add a blocked domain
	validator.AddBlockedDomain("blocked-example.com")

	testCases := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "Newly Blocked Domain",
			url:     "http://blocked-example.com/test",
			wantErr: true,
		},
		{
			name:    "Non-Blocked Domain",
			url:     "http://example.org",
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.Validate(tc.url)

			if (err != nil) != tc.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

// Performance test for URL validation
func BenchmarkValidate(b *testing.B) {
	validator := NewURLValidator()
	validURL := "https://example.com"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		validator.Validate(validURL)
	}
}
