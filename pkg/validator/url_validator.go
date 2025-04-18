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
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/yasin-yalcin-dev/go-url-shortener/pkg/errors"
)

// URLValidator provides URL validation functionality
type URLValidator struct {
	// List of blocked domain names
	blockedDomains []string

	// Maximum allowed URL length
	maxURLLength int
}

// NewURLValidator creates a new URL validator
func NewURLValidator() *URLValidator {
	return &URLValidator{
		blockedDomains: []string{
			// "example.com",
			// "test.com",
			// Other blocked domains
		},
		maxURLLength: 2048, // Standard maximum URL length
	}
}

// Validate checks the validity of a given URL
func (v *URLValidator) Validate(rawURL string) *errors.APIError {
	// Empty URL check
	if rawURL == "" {
		return errors.New(
			400,
			"URL cannot be empty",
			"Please provide a valid URL",
		)
	}

	// URL length check
	if len(rawURL) > v.maxURLLength {
		return errors.New(
			413,
			"URL is too long",
			"The entered URL exceeds the maximum allowed length",
		)
	}

	// URL format check
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return errors.New(
			400,
			"Invalid URL format",
			"Please enter a correct URL format",
		)
	}

	// Schema check (http, https)
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.New(
			400,
			"Only http and https schemes are supported",
			"Please enter a URL starting with http or https",
		)
	}

	// Domain validation
	if !govalidator.IsURL(rawURL) {
		return errors.New(
			400,
			"Invalid URL",
			"The entered URL is not a valid web address",
		)
	}

	// Blocked domain check
	if v.isDomainBlocked(parsedURL.Hostname()) {
		return errors.New(
			403,
			"This domain is blocked",
			"The specified domain is not allowed",
		)
	}

	return nil
}

func (v *URLValidator) isDomainBlocked(domain string) bool {
	domain = strings.ToLower(domain)
	for _, blockedDomain := range v.blockedDomains {
		if domain == strings.ToLower(blockedDomain) {
			return true
		}
	}
	return false
}

// AddBlockedDomain adds a new domain to the blocked domains list
func (v *URLValidator) AddBlockedDomain(domain string) {
	v.blockedDomains = append(v.blockedDomains, domain)
}
