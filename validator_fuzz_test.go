package workemailvalidator_test

import (
	"strings"
	"testing"

	workemailvalidator "github.com/rixlhq/work-email-validator"
)

// FuzzIsDisposableDomain tests IsDisposableDomain with random inputs.
func FuzzIsDisposableDomain(f *testing.F) {
	// Seed with interesting test cases
	seeds := []string{
		domainTempMail,
		domainTempMailUpper,
		domainTempMailSpaces,
		domainTempMailSub,
		domainGmail,
		domainExample,
		"",
		".",
		"..",
		"...",
		"a",
		"a.b",
		"a.b.c",
		strings.Repeat("a", 1000),
		"domain.com.",
		".domain.com",
		"domain..com",
		"@",
		"domain@com",
		"\x00domain.com",
		"domain.com\x00",
		domainMuenchen,
		domainEmoji,
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, domain string) {
		// Should never panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("IsDisposableDomain panicked with input %q: %v", domain, r)
			}
		}()

		result := workemailvalidator.IsDisposableDomain(domain)

		// Result should be deterministic
		result2 := workemailvalidator.IsDisposableDomain(domain)
		if result != result2 {
			t.Errorf("IsDisposableDomain not deterministic for %q: %v != %v", domain, result, result2)
		}

		// Case insensitive - upper and lower should give same result
		if domain != "" {
			upperResult := workemailvalidator.IsDisposableDomain(strings.ToUpper(domain))

			lowerResult := workemailvalidator.IsDisposableDomain(strings.ToLower(domain))
			if upperResult != lowerResult {
				t.Errorf("Case sensitivity issue for %q: upper=%v, lower=%v", domain, upperResult, lowerResult)
			}
		}
	})
}

// FuzzIsFreeDomain tests IsFreeDomain with random inputs.
func FuzzIsFreeDomain(f *testing.F) {
	seeds := []string{
		domainGmail,
		domainGmailUpper,
		domainGmailSpaces,
		emailUserMailGmail,
		domainOutlook,
		domainYahoo,
		domainExample,
		"",
		".",
		strings.Repeat("x", 500),
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, domain string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("IsFreeDomain panicked with input %q: %v", domain, r)
			}
		}()

		result := workemailvalidator.IsFreeDomain(domain)

		// Deterministic check
		result2 := workemailvalidator.IsFreeDomain(domain)
		if result != result2 {
			t.Errorf("IsFreeDomain not deterministic for %q", domain)
		}

		// Case insensitivity check
		if domain != "" {
			upperResult := workemailvalidator.IsFreeDomain(strings.ToUpper(domain))

			lowerResult := workemailvalidator.IsFreeDomain(strings.ToLower(domain))
			if upperResult != lowerResult {
				t.Errorf("Case sensitivity issue for %q", domain)
			}
		}
	})
}

// FuzzIsBusinessDomain tests IsBusinessDomain with random inputs.
func FuzzIsBusinessDomain(f *testing.F) {
	seeds := []string{
		domainExample,
		domainBusiness,
		domainGmail,
		domainTempMail,
		"",
		"x",
		strings.Repeat("business", 100),
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, domain string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("IsBusinessDomain panicked with input %q: %v", domain, r)
			}
		}()

		business := workemailvalidator.IsBusinessDomain(domain)
		disposableOrFree := workemailvalidator.IsDisposableOrFreeDomain(domain)

		// Business domain should be opposite of disposable or free, but only for VALID domains
		// Invalid domains (empty, invalid syntax) should return false for business
		normalized := strings.ToLower(strings.TrimSpace(domain))
		isLikelyValid := len(normalized) >= 4 && strings.Contains(normalized, ".")

		if isLikelyValid {
			// For valid-looking domains, business should be opposite of disposableOrFree
			if business == disposableOrFree {
				t.Errorf("Inconsistency for %q: business=%v, disposableOrFree=%v (should be opposite for valid domains)",
					domain, business, disposableOrFree)
			}
		} else {
			// For invalid domains, business should be false
			if business {
				t.Errorf("IsBusinessDomain returned true for invalid domain %q", domain)
			}
		}

		// Deterministic check
		business2 := workemailvalidator.IsBusinessDomain(domain)
		if business != business2 {
			t.Errorf("IsBusinessDomain not deterministic for %q", domain)
		}
	})
}

// FuzzIsDisposableOrFreeDomain tests the combined function.
func FuzzIsDisposableOrFreeDomain(f *testing.F) {
	seeds := []string{
		domainGmail,
		domainTempMail,
		domainExample,
		"",
		"test",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, domain string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("IsDisposableOrFreeDomain panicked with input %q: %v", domain, r)
			}
		}()

		disposableOrFree := workemailvalidator.IsDisposableOrFreeDomain(domain)
		disposable := workemailvalidator.IsDisposableDomain(domain)
		free := workemailvalidator.IsFreeDomain(domain)

		// Should be true if either disposable or free
		expected := disposable || free
		if disposableOrFree != expected {
			t.Errorf("Inconsistency for %q: IsDisposableOrFreeDomain=%v, but disposable=%v, free=%v",
				domain, disposableOrFree, disposable, free)
		}
	})
}

// FuzzIsWorkEmail tests IsWorkEmail with random inputs
//
//nolint:cyclop,funlen
func FuzzIsWorkEmail(f *testing.F) {
	seeds := []string{
		emailUserExample,
		emailUserGmail,
		emailUserTempMail,
		"@",
		emailAtDomain,
		emailUserAtEnd,
		"",
		"no-at-sign",
		"multiple@@at.com",
		"user@domain@com",
		strings.Repeat("a", 100) + emailAtExample,
		"user@" + strings.Repeat("sub.", 20) + domainExample,
		emailUnicodeExample,
		emailUserMuenchen,
		emailPlusExample,
		emailDotsExample,
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, email string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("IsWorkEmail panicked with input %q: %v", email, r)
			}
		}()

		result := workemailvalidator.IsWorkEmail(email)

		// Deterministic check
		result2 := workemailvalidator.IsWorkEmail(email)
		if result != result2 {
			t.Errorf("IsWorkEmail not deterministic for %q", email)
		}

		// If email is valid work email, extract domain and verify consistency
		if result && strings.Contains(email, "@") {
			atIndex := strings.LastIndexByte(email, '@')
			if atIndex > 0 && atIndex < len(email)-1 {
				domain := email[atIndex+1:]

				business := workemailvalidator.IsBusinessDomain(domain)
				if !business {
					t.Errorf("IsWorkEmail(%q)=true but domain %q is not business", email, domain)
				}
			}
		}

		// If email is not work email due to invalid format, check edge cases
		if !result && strings.Contains(email, "@") {
			atIndex := strings.LastIndexByte(email, '@')
			if atIndex > 0 && atIndex < len(email)-1 {
				// Valid format, so must be non-business domain
				domain := email[atIndex+1:]

				business := workemailvalidator.IsBusinessDomain(domain)
				if business {
					t.Errorf("IsWorkEmail(%q)=false but domain %q is business", email, domain)
				}
			}
		}
	})
}

// FuzzConsistency ensures all functions remain consistent with each other.
//
//nolint:cyclop
func FuzzConsistency(f *testing.F) {
	seeds := []string{
		domainGmail,
		domainTempMail,
		domainExample,
		domainBusiness,
		"",
		"test",
		domainGmailSub,
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, domain string) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Consistency test panicked with input %q: %v", domain, r)
			}
		}()

		disposable := workemailvalidator.IsDisposableDomain(domain)
		free := workemailvalidator.IsFreeDomain(domain)
		disposableOrFree := workemailvalidator.IsDisposableOrFreeDomain(domain)
		business := workemailvalidator.IsBusinessDomain(domain)

		// disposableOrFree should equal disposable OR free
		if disposableOrFree != (disposable || free) {
			t.Errorf("IsDisposableOrFreeDomain inconsistent for %q", domain)
		}

		// business should be opposite of disposableOrFree, but only for valid domains
		// Invalid domains return false for business regardless
		normalized := strings.ToLower(strings.TrimSpace(domain))
		isLikelyValid := len(normalized) >= 4 && strings.Contains(normalized, ".")

		if isLikelyValid {
			if business == disposableOrFree {
				t.Errorf("IsBusinessDomain inconsistent for valid domain %q", domain)
			}
		} else {
			// Invalid domains should return false for business
			if business {
				t.Errorf("IsBusinessDomain returned true for invalid domain %q", domain)
			}
		}

		// A domain should not be both disposable and free (data integrity check)
		if disposable && free {
			t.Logf("Warning: %q marked as both disposable and free", domain)
		}
	})
}
