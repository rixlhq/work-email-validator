package workemailvalidator_test

import (
	"strings"
	"testing"

	workemailvalidator "github.com/rixlhq/work-email-validator"
)

// testCase represents a generic test case for domain/email validation.
type testCase struct {
	name     string
	input    string
	expected bool
}

// runDomainTests is a DRY helper that runs domain validation tests.
func runDomainTests(t *testing.T, tests []testCase, validatorFunc func(string) bool) {
	t.Helper()

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			result := validatorFunc(testCase.input)
			if result != testCase.expected {
				t.Errorf("got %v, want %v for input %q", result, testCase.expected, testCase.input)
			}
		})
	}
}

// Edge case tests for IsDisposableDomain.
func TestIsDisposableDomain(t *testing.T) {
	t.Parallel()

	tests := []testCase{
		// Known disposable domains
		{nameKnownDisposable, domainTempMail, true},
		{nameKnownDisposable2, domain10MinuteMail, true},
		{nameKnownDisposable3, "guerrillamail.com", true},

		// Case insensitivity
		{"uppercase", "TEMP-MAIL.COM", true},
		{"mixedcase", "TeMp-MaIl.CoM", true},

		// Whitespace handling
		{"leading_whitespace", "  temp-mail.org", true},
		{"trailing_whitespace", "temp-mail.org  ", true},
		{"both_whitespace", "  temp-mail.org  ", true},
		{"tabs_and_spaces", "\t temp-mail.org \t", true},

		// Non-disposable domains
		{"gmail_not_disposable", domainGmail, false},
		{"business_domain", domainExample, false},
		{"corporate_domain", domainBusiness, false},

		// Edge cases
		{nameEmptyString, "", false},
		{"single_char", "a", false},
		{"no_tld", "domain", false},
		{"just_dot", ".", false},
		{"multiple_dots", "...", false},
		{"leading_dot", ".domain.com", false},
		{"trailing_dot", "domain.com.", false},

		// Subdomain tests (if parent is disposable)
		{"subdomain_of_disposable", domainTempMailSub, true},
		{"deep_subdomain_disposable", domainTempMailDeepSub, true},

		// Special characters
		{"hyphen_in_domain", "my-domain.com", false},
		{"numbers_in_domain", "123domain.com", false},
		{"mixed_alphanumeric", "test123-domain456.com", false},

		// Very long inputs
		{"long_domain", strings.Repeat("a", 100) + ".com", false},
		{"long_subdomain", strings.Repeat("sub.", 20) + "domain.com", false},

		// Unicode domains (IDN)
		{"unicode_domain", domainMuenchen, false},
		{"emoji_domain", domainEmoji, false},
	}

	runDomainTests(t, tests, workemailvalidator.IsDisposableDomain)
}

// Edge case tests for IsFreeDomain.
func TestIsFreeDomain(t *testing.T) {
	t.Parallel()

	tests := []testCase{
		// Known free domains
		{"gmail", domainGmail, true},
		{"outlook", domainOutlook, true},
		{"yahoo", "yahoo.com", true},
		{"hotmail", "hotmail.com", true},
		{"icloud", "icloud.com", true},
		{"protonmail", "protonmail.com", true},

		// Case insensitivity
		{"uppercase_gmail", domainGmailUpper, true},
		{"mixedcase_outlook", "OuTlOoK.cOm", true},

		// Whitespace handling
		{"whitespace_yahoo", "  yahoo.com  ", true},
		{"tabs_hotmail", "\t\thotmail.com\t\t", true},

		// Non-free domains
		{"business_domain", domainExample, false},
		{"disposable_not_free", domainTempMail, false},
		{"corporate", domainBusiness, false},

		// Edge cases
		{nameEmptyString, "", false},
		{"single_letter", "g", false},
		{"incomplete_domain", "gma", false},
		{"typo_gmail", "gmial.com", false},

		// Subdomains (if parent is free)
		{"subdomain_gmail", emailUserMailGmail, true},
		{"subdomain_outlook", domainOutlookAccounts, true},
		{"deep_subdomain_yahoo", "a.b.yahoo.com", true},

		// Edge case subdomains
		{"empty_subdomain", domainGmailEmptySub, true},
		{"numeric_subdomain", domainGmailNumericSub, true},
	}

	runDomainTests(t, tests, workemailvalidator.IsFreeDomain)
}

// Edge case tests for IsDisposableOrFreeDomain.
func TestIsDisposableOrFreeDomain(t *testing.T) {
	t.Parallel()

	tests := []testCase{
		// Free domains
		{nameFreeGmail, domainGmail, true},
		{nameFreeOutlook, domainOutlook, true},

		// Disposable domains
		{"disposable_tempmail", domainTempMail, true},
		{"disposable_10min", domain10MinuteMail, true},

		// Business domains (neither)
		{"business_example", domainExample, false},
		{"business_corporate", domainBusiness, false},
		{"business_custom", "custom-business.io", false},

		// Edge cases
		{"empty", "", false},
		{"whitespace_only", "   ", false},
		{"invalid_format", "not a domain", false},

		// Mixed cases
		{"uppercase_free", domainGmailUpper, true},
		{"whitespace_disposable", domainTempMailSpaces, true},

		// Subdomains
		{"subdomain_free", emailUserMailGmail, true},
		{"subdomain_disposable", domainTempMailXSub, true},
		{"subdomain_business", domainBusinessSub, false},
	}

	runDomainTests(t, tests, workemailvalidator.IsDisposableOrFreeDomain)
}

// Edge case tests for IsBusinessDomain.
func TestIsBusinessDomain(t *testing.T) {
	t.Parallel()

	tests := []testCase{
		// Business domains
		{"simple_business", domainExample, true},
		{"corporate", domainBusiness, true},
		{"startup", "startupname.io", true},
		{"enterprise", "bigcorp.net", true},

		// Non-business (free)
		{nameFreeGmail, domainGmail, false},
		{nameFreeOutlook, domainOutlook, false},

		// Non-business (disposable)
		{"disposable_tempmail", domainTempMail, false},
		{"disposable_guerrilla", "guerrillamail.com", false},

		// Edge cases - invalid domain syntax should return false
		{"empty", "", false},
		{"single_char", "x", false},
		{"tld_only", ".com", false},
		{"single_char_tld", "domain.a", false}, // TLD must be at least 2 chars

		// Subdomains of business domains
		{"business_subdomain", domainBusinessSub, true},
		{"business_deep_sub", "v1." + domainBusinessSub, true},

		// Subdomains of free/disposable (should be non-business)
		{"free_subdomain", domainGmailCustomSub, false},
		{"disposable_subdomain", domainTempMailTestSub, false},

		// Case and whitespace
		{"uppercase_business", "EXAMPLE.COM", true},
		{"whitespace_business", "  " + domainExample + "  ", true},
	}

	runDomainTests(t, tests, workemailvalidator.IsBusinessDomain)
}

// Edge case tests for IsWorkEmail.
func TestIsWorkEmail(t *testing.T) {
	t.Parallel()

	tests := []testCase{
		// Valid work emails
		{"valid_work", emailUserMailCompany, true},
		{"valid_work_2", emailContactExample, true},
		{"valid_work_subdomain", "admin@mail.company.com", true},

		// Non-work (free)
		{"free_gmail", emailUserGmail, false},
		{"free_outlook", emailUserOutlook, false},
		{"free_yahoo", emailUserYahoo, false},

		// Non-work (disposable)
		{"disposable", emailUserTempMail, false},
		{"disposable_2", emailUserTempMailAt, false},

		// Invalid email formats
		{"no_at_sign", "invalid-email", false},
		{nameEmptyString, "", false},
		{"only_at", "@", false},
		{"at_at_start", emailAtDomain, false},
		{"at_at_end", emailUserAtEnd, false},
		// Multiple @ signs - domain extraction uses LastIndexByte
		// "user@@domain.com" -> "domain.com" (valid and business)
		// "user@domain@com" -> "com" (invalid - TLD only, fails validation)
		{"multiple_at", "user@@domain.com", true},
		{"multiple_at_2", "user@domain@com", false},

		// Edge cases with @ symbol
		{"just_domain", "domain.com", false},
		{"no_local_part", emailAtDomain, false},
		{"no_domain_part", emailUserAtEnd, false},

		// Whitespace
		{"whitespace_domain", emailUserWhitespace, true},
		{"whitespace_full", emailUserWhitespaceFull, true},

		// Case insensitivity
		{"uppercase_domain", emailUserExampleUpper, true},
		{"uppercase_free", emailUserUpperGmail, false},

		// Special characters in local part (should still validate domain)
		{"plus_addressing", emailPlusExample, true},
		{"dots_in_local", emailDotsExample, true},
		{"underscore", emailUnderscoreExample, true},
		{"hyphen", emailHyphenExample, true},

		// Subdomain edge cases
		{"subdomain_business", emailUserMailCompany, true},
		{"subdomain_free", emailUserMailGmail, false},
		{"subdomain_disposable", emailUserXTempMail, false},

		// Multiple @ signs - uses LastIndexByte, so extracts domain after last @
		// "user@host@company.com" extracts "company.com" which is valid and business
		{"email_with_at_in_local", emailUserHostCompany, true},

		// Very long emails
		{"long_local_part", strings.Repeat("a", 64) + emailAtExample, true},
		{"long_domain", "user@" + strings.Repeat("sub.", 10) + domainExample, true},

		// Unicode
		{"unicode_local", emailUnicodeExample, true},
		{"unicode_domain", emailUserMuenchen, true},
	}

	runDomainTests(t, tests, workemailvalidator.IsWorkEmail)
}

// Test subdomain hierarchy handling.
func TestSubdomainHierarchy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		domain     string
		disposable bool
		free       bool
		business   bool
	}{
		{"root_disposable", domainTempMail, true, false, false},
		{"sub_disposable", domainTempMailTestSub, true, false, false},
		{"deep_sub_disposable", domainTempMailDeepSub, true, false, false},

		{"root_free", domainGmail, false, true, false},
		{"sub_free", emailUserMailGmail, false, true, false},
		{"deep_sub_free", domainGmailDeepSub, false, true, false},

		{"root_business", domainExample, false, false, true},
		{"sub_business", domainExampleSub, false, false, true},
		{"deep_sub_business", domainExampleDeepSub, false, false, true},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			disposable := workemailvalidator.IsDisposableDomain(testCase.domain)
			free := workemailvalidator.IsFreeDomain(testCase.domain)
			business := workemailvalidator.IsBusinessDomain(testCase.domain)

			if disposable != testCase.disposable {
				t.Errorf("IsDisposableDomain(%q) = %v, want %v", testCase.domain, disposable, testCase.disposable)
			}

			if free != testCase.free {
				t.Errorf("IsFreeDomain(%q) = %v, want %v", testCase.domain, free, testCase.free)
			}

			if business != testCase.business {
				t.Errorf("IsBusinessDomain(%q) = %v, want %v", testCase.domain, business, testCase.business)
			}
		})
	}
}

// Test consistency between functions.
func TestFunctionConsistency(t *testing.T) {
	t.Parallel()

	domains := []string{
		domainGmail,
		domainTempMail,
		domainExample,
		domainOutlook,
		domainBusiness,
		"",
		"invalid",
	}

	for _, domain := range domains {
		t.Run(domain, func(t *testing.T) {
			t.Parallel()

			disposable := workemailvalidator.IsDisposableDomain(domain)
			free := workemailvalidator.IsFreeDomain(domain)
			disposableOrFree := workemailvalidator.IsDisposableOrFreeDomain(domain)
			business := workemailvalidator.IsBusinessDomain(domain)

			// IsDisposableOrFreeDomain should be true if either disposable or free
			expectedDisposableOrFree := disposable || free
			if disposableOrFree != expectedDisposableOrFree {
				t.Errorf("IsDisposableOrFreeDomain(%q) = %v, but disposable=%v, free=%v",
					domain, disposableOrFree, disposable, free)
			}

			// IsBusinessDomain should be opposite of IsDisposableOrFreeDomain
			// BUT only for VALID domains. Invalid domains (empty, invalid syntax) return false for both

			// Check if domain would pass basic validation
			normalized := strings.ToLower(strings.TrimSpace(domain))
			isLikelyValid := len(normalized) >= 4 && strings.Contains(normalized, ".")

			if isLikelyValid {
				// For valid-looking domains, business should be opposite of disposable/free
				if business == disposableOrFree {
					t.Errorf("IsBusinessDomain(%q) = %v, but IsDisposableOrFreeDomain = %v (should be opposite for valid domains)",
						domain, business, disposableOrFree)
				}
			} else {
				// For invalid domains, business should be false
				// disposableOrFree can be false (if not in lists) or could theoretically be true (if matched)
				if business {
					t.Errorf("IsBusinessDomain(%q) = true, but domain appears invalid", domain)
				}
			}

			// A domain cannot be both disposable and free
			// (this is a logical constraint, not enforced by code but should be true for data)
			if disposable && free {
				t.Logf("Warning: %q is marked as both disposable and free", domain)
			}
		})
	}
}
