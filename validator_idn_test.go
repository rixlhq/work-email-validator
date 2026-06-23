package workemailvalidator_test

import (
	"testing"

	workemailvalidator "github.com/rixlhq/work-email-validator"
)

// TestInternationalizedDomains tests support for IDN (Internationalized Domain Names).
func TestInternationalizedDomains(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		domain   string
		expected bool
	}{
		// German domains with umlauts
		{"german_umlaut", domainMuenchen, true},
		{"german_business", "bücher.com", true},
		{"german_ae", "äpfel.de", true},

		// Cyrillic domains
		{"cyrillic_russian", domainMoscow, true},
		{"cyrillic_ukraine", "київ.ua", true},
		{"cyrillic_example", "пример.com", true},

		// Chinese domains
		{"chinese_simplified", domainChinese, true},
		{"chinese_traditional", "範例.台灣", true}, //nolint:gosmopolitan
		{"chinese_mixed", "测试.com", true},      //nolint:gosmopolitan

		// Arabic domains
		{"arabic_example", "مثال.com", true},
		{"arabic_saudi", "السعودية.sa", true},

		// Japanese domains
		{"japanese_hiragana", "にほん.jp", true},
		{"japanese_katakana", "テスト.jp", true},
		{"japanese_kanji", domainJapanese, true},

		// Greek domains
		{"greek_example", "παράδειγμα.gr", true},

		// Hebrew domains
		{"hebrew_example", "דוגמה.com", true},

		// Mixed script (should still work)
		{"mixed_english_cyrillic", "test-тест.com", true},

		// Emoji domains (these exist!)
		{"emoji_heart", "❤️.com", true},
		{"emoji_smile", domainEmoji, true},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// Test IsBusinessDomain - all these should be business domains
			// (none should be in free/disposable lists)
			result := workemailvalidator.IsBusinessDomain(testCase.domain)
			if result != testCase.expected {
				t.Errorf("IsBusinessDomain(%q) = %v, want %v", testCase.domain, result, testCase.expected)
			}

			// They should NOT be free or disposable
			if workemailvalidator.IsFreeDomain(testCase.domain) {
				t.Errorf("IsFreeDomain(%q) = true, but internationalized domain shouldn't be in free list", testCase.domain)
			}

			if workemailvalidator.IsDisposableDomain(testCase.domain) {
				t.Errorf("IsDisposableDomain(%q) = true, but internationalized domain shouldn't be in disposable list",
					testCase.domain)
			}
		})
	}
}

// TestInternationalizedEmails tests email validation with IDN.
func TestInternationalizedEmails(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		// German
		{"german_umlaut_domain", emailUserMuenchen, true},
		{"german_umlaut_local", "müller@" + domainExample, true},

		// Cyrillic
		{"cyrillic_domain", emailUserMoscow, true},
		{"cyrillic_local", "пользователь@" + domainExample, true},

		// Chinese
		{"chinese_domain", emailUserChinese, true},
		{"chinese_local", "用户@" + domainExample, true}, //nolint:gosmopolitan

		// Japanese
		{"japanese_domain", emailUserJapanese, true},
		{"japanese_local", "ユーザー@" + domainExample, true},

		// Mixed
		{"mixed_local_domain", "user@тест.com", true},
		{"both_internationalized", emailUserMoscowFull, true},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			result := workemailvalidator.IsWorkEmail(testCase.email)
			if result != testCase.expected {
				t.Errorf("IsWorkEmail(%q) = %v, want %v", testCase.email, result, testCase.expected)
			}
		})
	}
}

// TestIDNEquivalence tests that IDN domains are normalized correctly.
func TestIDNEquivalence(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		domain1            string
		domain2            string
		shouldBeEquivalent bool
	}{
		// German umlaut normalization
		{"german_ae_vs_a", "müller.com", "muller.com", false}, // ü ≠ u
		{"german_same", domainMuenchen, "MÜNCHEN.DE", true},

		// Cyrillic normalization
		{"cyrillic_case", domainMoscow, "МОСКВА.РФ", true},

		// Mixed case should normalize the same
		{"mixed_case", "Test.COM", "test.com", true},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			result1 := workemailvalidator.IsBusinessDomain(testCase.domain1)
			result2 := workemailvalidator.IsBusinessDomain(testCase.domain2)

			if testCase.shouldBeEquivalent {
				if result1 != result2 {
					t.Errorf("Domains %q and %q should be equivalent, but got %v != %v",
						testCase.domain1, testCase.domain2, result1, result2)
				}
			}
		})
	}
}

// TestIDNPunycodeConversion tests that Punycode conversion works.
func TestIDNPunycodeConversion(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		unicodeDomain  string
		punycodeDomain string
	}{
		// These should be treated as equivalent
		{"german", domainMuenchen, "xn--mnchen-3ya.de"},
		{"russian", domainMoscow, "xn--80adxhks.xn--p1ai"},
		{"chinese", domainChinese, "xn--fsqu00a.xn--fiqs8s"},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// Both unicode and punycode versions should give same result
			unicodeResult := workemailvalidator.IsBusinessDomain(testCase.unicodeDomain)
			punycodeResult := workemailvalidator.IsBusinessDomain(testCase.punycodeDomain)

			if unicodeResult != punycodeResult {
				t.Errorf("Unicode %q and Punycode %q should give same result, but got %v != %v",
					testCase.unicodeDomain, testCase.punycodeDomain, unicodeResult, punycodeResult)
			}
		})
	}
}

// TestIDNInvalidDomains tests that invalid IDN domains are handled gracefully.
func TestIDNInvalidDomains(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		domain string
	}{
		{"only_emoji", "😀"},
		{"emoji_no_tld", "❤️"},
		{"invalid_unicode", "\x00\x01\x02.com"},
		// Test cases that trigger idna.ToASCII errors to cover error path in domainToASCII
		{"invalid_idna_long_label", "xn--" + string(make([]byte, 64)) + ".com"}, // Label too long for IDNA
		{"control_characters_null", "test\x00domain.com"},                       // Null byte in domain
		{"control_characters_bell", "test\x07.com"},                             // Bell character
		{"invalid_punycode_prefix", "xn--invalid.com"},                          // Invalid punycode
		{"disallowed_bidi_chars", "test\u200E.com"},                             // Left-to-right mark
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// These should not panic - the main goal is to cover the error path in domainToASCII
			// Some of these domains may or may not trigger IDNA errors, but we're testing robustness
			_ = workemailvalidator.IsBusinessDomain(testCase.domain)
			_ = workemailvalidator.IsDisposableDomain(testCase.domain)
			_ = workemailvalidator.IsFreeDomain(testCase.domain)

			// Test passes if no panic occurs
		})
	}
}
