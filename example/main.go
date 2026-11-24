// Command example demonstrates basic usage of the work-email-validator package.
package main

import (
	"fmt"
	"os"
	"strings"

	validator "github.com/rixlhq/work-email-validator/v2"
)

const headerWidth = 49

func outf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stdout, format, args...)
}

func main() {
	outf("Work Email Validator - Example Usage\n")
	outf("%s\n", "="+strings.Repeat("=", headerWidth))
	outf("\n")

	testDomains := []string{
		"gmail.com",
		"outlook.com",
		"mycompany.com",
		"temp-mail.com",
		"10minutemail.com",
		"protonmail.com",
		"example.com",
		"business-domain.io",
	}

	for _, domain := range testDomains {
		outf("Domain: %-25s\n", domain)
		outf("  ├─ Disposable: %-5t\n", validator.IsDisposableDomain(domain))
		outf("  ├─ Free:       %-5t\n", validator.IsFreeDomain(domain))
		outf("  └─ Business:   %-5t\n", validator.IsBusinessDomain(domain))
		outf("\n")
	}

	outf("\nExtract domain from email and validate:\n")

	emails := []string{
		"user@gmail.com",
		"contact@mycompany.com",
		"test@temp-mail.com",
	}

	for _, email := range emails {
		if _, domain, ok := strings.Cut(email, "@"); ok {
			outf("Email: %s → Domain: %s (Business: %t)\n",
				email, domain, validator.IsBusinessDomain(domain))
		}
	}
}
