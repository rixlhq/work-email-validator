package workemailvalidator_test

import (
	"strings"
	"testing"

	workemailvalidator "github.com/rixlhq/work-email-validator"
)

// Benchmark best case - exact match in map.
func BenchmarkIsDisposableDomain_ExactMatch(b *testing.B) {
	for b.Loop() {
		workemailvalidator.IsDisposableDomain(domainTempMail)
	}
}

// Benchmark worst case - needs to scan entire domain for dots.
func BenchmarkIsDisposableDomain_NoMatch(b *testing.B) {
	for b.Loop() {
		workemailvalidator.IsDisposableDomain("very.long.subdomain.that.is.not.disposable." + domainExample)
	}
}

// Benchmark subdomain match.
func BenchmarkIsDisposableDomain_SubdomainMatch(b *testing.B) {
	for b.Loop() {
		workemailvalidator.IsDisposableDomain("sub.sub.sub." + domainTempMail)
	}
}

// Benchmark with whitespace (requires trimming).
func BenchmarkIsDisposableDomain_WithWhitespace(b *testing.B) {
	for b.Loop() {
		workemailvalidator.IsDisposableDomain(domainTempMailSpaces)
	}
}

// Benchmark with case conversion.
func BenchmarkIsDisposableDomain_UpperCase(b *testing.B) {
	for b.Loop() {
		workemailvalidator.IsDisposableDomain("TEMP-MAIL.COM")
	}
}

// Benchmark free domain best case.
func BenchmarkIsFreeDomain_ExactMatch(b *testing.B) {
	for b.Loop() {
		workemailvalidator.IsFreeDomain(domainGmail)
	}
}

// Benchmark free domain worst case.
func BenchmarkIsFreeDomain_NoMatch(b *testing.B) {
	for b.Loop() {
		workemailvalidator.IsFreeDomain("a.very.long.business.domain." + domainExample)
	}
}

// Benchmark free domain subdomain.
func BenchmarkIsFreeDomain_SubdomainMatch(b *testing.B) {
	for b.Loop() {
		workemailvalidator.IsFreeDomain(emailUserMailGmail)
	}
}

// Benchmark business domain (checks both maps).
func BenchmarkIsBusinessDomain_True(b *testing.B) {
	for b.Loop() {
		workemailvalidator.IsBusinessDomain(domainExample)
	}
}

// Benchmark business domain false (disposable).
func BenchmarkIsBusinessDomain_FalseDisposable(b *testing.B) {
	for b.Loop() {
		workemailvalidator.IsBusinessDomain(domainTempMail)
	}
}

// Benchmark business domain false (free).
func BenchmarkIsBusinessDomain_FalseFree(b *testing.B) {
	for b.Loop() {
		workemailvalidator.IsBusinessDomain(domainGmail)
	}
}

// Benchmark combined check.
func BenchmarkIsDisposableOrFreeDomain_Disposable(b *testing.B) {
	for b.Loop() {
		workemailvalidator.IsDisposableOrFreeDomain(domainTempMail)
	}
}

func BenchmarkIsDisposableOrFreeDomain_Free(b *testing.B) {
	for b.Loop() {
		workemailvalidator.IsDisposableOrFreeDomain(domainGmail)
	}
}

func BenchmarkIsDisposableOrFreeDomain_Neither(b *testing.B) {
	for b.Loop() {
		workemailvalidator.IsDisposableOrFreeDomain(domainExample)
	}
}

// Benchmark work email validation.
func BenchmarkIsWorkEmail_Valid(b *testing.B) {
	for b.Loop() {
		workemailvalidator.IsWorkEmail(emailUserExample)
	}
}

func BenchmarkIsWorkEmail_InvalidFree(b *testing.B) {
	for b.Loop() {
		workemailvalidator.IsWorkEmail(emailUserGmail)
	}
}

func BenchmarkIsWorkEmail_InvalidDisposable(b *testing.B) {
	for b.Loop() {
		workemailvalidator.IsWorkEmail(emailUserTempMail)
	}
}

func BenchmarkIsWorkEmail_InvalidFormat(b *testing.B) {
	for b.Loop() {
		workemailvalidator.IsWorkEmail("not-an-email")
	}
}

// Benchmark with very long email.
func BenchmarkIsWorkEmail_LongEmail(b *testing.B) {
	email := strings.Repeat("a", 100) + "@" + strings.Repeat("subdomain.", 10) + domainExample

	b.ResetTimer()

	for b.Loop() {
		workemailvalidator.IsWorkEmail(email)
	}
}

// Benchmark with subdomain of free provider.
func BenchmarkIsWorkEmail_SubdomainFree(b *testing.B) {
	for b.Loop() {
		workemailvalidator.IsWorkEmail(emailUserMailGmail)
	}
}

// Benchmark with subdomain of business.
func BenchmarkIsWorkEmail_SubdomainBusiness(b *testing.B) {
	for b.Loop() {
		workemailvalidator.IsWorkEmail(emailUserAPICompany)
	}
}

// Benchmark parallel execution patterns.
func BenchmarkIsWorkEmail_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			workemailvalidator.IsWorkEmail(emailUserExample)
		}
	})
}

func BenchmarkIsDisposableDomain_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			workemailvalidator.IsDisposableDomain(domainTempMail)
		}
	})
}

func BenchmarkIsFreeDomain_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			workemailvalidator.IsFreeDomain(domainGmail)
		}
	})
}

// Benchmark mixed workload (realistic usage).
func BenchmarkMixedWorkload(b *testing.B) {
	domains := []string{
		emailUserExample,
		emailUserGmail,
		emailUserTempMail,
		emailAdminMyCompany,
		emailTestOutlook,
	}

	b.ResetTimer()

	for b.Loop() {
		for _, domain := range domains {
			workemailvalidator.IsWorkEmail(domain)
		}
	}
}

// Benchmark to measure allocation overhead.
func BenchmarkIsWorkEmail_Allocations(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		workemailvalidator.IsWorkEmail(emailUserExample)
	}
}

func BenchmarkIsDisposableDomain_Allocations(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		workemailvalidator.IsDisposableDomain(domainTempMail)
	}
}

func BenchmarkIsFreeDomain_Allocations(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		workemailvalidator.IsFreeDomain(domainGmail)
	}
}
