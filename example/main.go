package main

import (
	"fmt"
	"strings"

	validator "github.com/rixlhq/work-email-validator/v2"
)

func main() {
	fmt.Println("Work Email Validator - Example Usage")
	fmt.Println("=" + strings.Repeat("=", 49))
	fmt.Println()

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
		fmt.Printf("Domain: %-25s\n", domain)
		fmt.Printf("  ├─ Disposable: %-5t\n", validator.IsDisposableDomain(domain))
		fmt.Printf("  ├─ Free:       %-5t\n", validator.IsFreeDomain(domain))
		fmt.Printf("  └─ Business:   %-5t\n", validator.IsBusinessDomain(domain))
		fmt.Println()
	}

	fmt.Println("\nExtract domain from email and validate:")
	emails := []string{
		"user@gmail.com",
		"contact@mycompany.com",
		"test@temp-mail.com",
	}

	for _, email := range emails {
		parts := strings.Split(email, "@")
		if len(parts) == 2 {
			domain := parts[1]
			fmt.Printf("Email: %s → Domain: %s (Business: %t)\n",
				email, domain, validator.IsBusinessDomain(domain))
		}
	}
}
