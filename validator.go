// Package workemailvalidator provides email validation utilities to determine
// whether an email address is from a disposable, free, or business domain.
package workemailvalidator

import (
	"net/mail"
	"strings"
)

// normalize prepares the domain for lookup and reports whether it is a valid DNS name.
// Invalid names are still lowercased so they can be matched against the lists.
func normalize(domain string) (string, bool) {
	domain = strings.TrimSuffix(strings.TrimSpace(domain), ".")

	ascii, valid := domainToASCII(domain)
	if !valid {
		return strings.ToLower(domain), false
	}

	return ascii, true
}

// contains checks if the domain or any of its parents exist in the map.
// It expects the domain to be already normalized (lowercase, trimmed).
func contains(domain string, domainMap map[string]struct{}) bool {
	if _, ok := domainMap[domain]; ok {
		return true
	}

	for i := range len(domain) {
		if domain[i] == '.' {
			if _, ok := domainMap[domain[i+1:]]; ok {
				return true
			}
		}
	}

	return false
}

// isValidDomainSyntax checks the shape of a name that already passed IDNA validation.
func isValidDomainSyntax(domain string) bool {
	if len(domain) < 4 || len(domain) > 253 {
		return false
	}

	lastDotIndex := strings.LastIndexByte(domain, '.')
	if lastDotIndex <= 0 {
		return false
	}

	const TLDMinLength = 2

	tld := domain[lastDotIndex+1:]

	return len(tld) >= TLDMinLength && strings.Trim(tld, "0123456789") != ""
}

// IsDisposableDomain checks if the given domain is a disposable/temporary email domain.
func IsDisposableDomain(domain string) bool {
	normalized, _ := normalize(domain)
	return contains(normalized, disposableDomains)
}

// IsFreeDomain checks if the given domain is a free email provider domain.
func IsFreeDomain(domain string) bool {
	normalized, _ := normalize(domain)
	return contains(normalized, freeDomains)
}

// IsDisposableOrFreeDomain checks if the given domain is either disposable or free.
func IsDisposableOrFreeDomain(domain string) bool {
	normalized, _ := normalize(domain)
	return contains(normalized, disposableDomains) || contains(normalized, freeDomains)
}

// IsBusinessDomain checks if the given domain is neither disposable nor free.
func IsBusinessDomain(domain string) bool {
	normalized, valid := normalize(domain)
	if !valid || !isValidDomainSyntax(normalized) {
		return false
	}

	return !contains(normalized, disposableDomains) && !contains(normalized, freeDomains)
}

// IsWorkEmail checks if the given email address is from a business domain.
func IsWorkEmail(email string) bool {
	email = strings.TrimSpace(email)

	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email {
		return false
	}

	return IsBusinessDomain(email[strings.LastIndexByte(email, '@')+1:])
}
