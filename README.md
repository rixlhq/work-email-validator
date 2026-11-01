# Work Email Validator

A lightweight Go library to verify if an email domain is disposable, free, or business-oriented.

[![Test](https://github.com/rixlhq/work-email-validator/actions/workflows/test.yml/badge.svg)](https://github.com/rixlhq/work-email-validator/actions/workflows/test.yml)
[![Coverage Status](https://codecov.io/gh/rixlhq/work-email-validator/graph/badge.svg)](https://codecov.io/gh/rixlhq/work-email-validator)
[![Go Report Card](https://goreportcard.com/badge/github.com/rixlhq/work-email-validator)](https://goreportcard.com/report/github.com/rixlhq/work-email-validator)
[![GoDoc](https://godoc.org/github.com/rixlhq/work-email-validator?status.svg)](https://godoc.org/github.com/rixlhq/work-email-validator)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

- ✅ Check if a domain is disposable (temp-mail.com, 10minutemail.com, etc.)
- ✅ Check if a domain is a free email provider (gmail.com, outlook.com, etc.)
- ✅ Check if a domain is a business domain (not free or disposable)
- ✅ Automatically updated domain lists via GitHub Actions (weekly)
- ✅ Fast lookup with embedded data (no external API calls)
- ✅ Zero dependencies
- ✅ Thread-safe

## Installation

```bash
go get github.com/rixlhq/work-email-validator
```

## Usage

```go
package main

import (
	"fmt"

	validator "github.com/rixlhq/work-email-validator"
)

func main() {
	// Check if domain is disposable
	if validator.IsDisposableDomain("temp-mail.com") {
		fmt.Println("This is a disposable email domain")
	}

	// Check if domain is free
	if validator.IsFreeDomain("gmail.com") {
		fmt.Println("This is a free email domain")
	}

	// Check if domain is disposable or free
	if validator.IsDisposableOrFreeDomain("outlook.com") {
		fmt.Println("This is either disposable or free")
	}

	// Check if domain is a business domain
	if validator.IsBusinessDomain("mycompany.com") {
		fmt.Println("This is a business email domain")
	}
}
```

## API

### `IsDisposableDomain(domain string) bool`

Returns `true` if the domain is a disposable/temporary email provider.

**Examples:**
- `temp-mail.com` → `true`
- `10minutemail.com` → `true`
- `gmail.com` → `false`

### `IsFreeDomain(domain string) bool`

Returns `true` if the domain is a free email provider (Gmail, Outlook, Yahoo, etc.).

**Examples:**
- `gmail.com` → `true`
- `outlook.com` → `true`
- `mycompany.com` → `false`

### `IsDisposableOrFreeDomain(domain string) bool`

Returns `true` if the domain is either disposable or free.

**Examples:**
- `gmail.com` → `true`
- `temp-mail.com` → `true`
- `mycompany.com` → `false`

### `IsBusinessDomain(domain string) bool`

Returns `true` if the domain is neither disposable nor free (i.e., likely a business domain).

**Examples:**
- `mycompany.com` → `true`
- `gmail.com` → `false`
- `temp-mail.com` → `false`

## Domain Lists

### Disposable Domains
The disposable domain list is sourced from [disposable-email-domains/disposable-email-domains](https://github.com/disposable-email-domains/disposable-email-domains) and contains **4,941 domains**.

### Free Email Providers
The free email providers list is sourced from [willwhite/freemail](https://github.com/willwhite/freemail) and contains **4,456 domains**, including:
- Gmail, Googlemail
- Outlook, Hotmail, Live, MSN
- Yahoo (all variants)
- iCloud, Me, Mac
- AOL, AIM
- ProtonMail, Tutanota
- Zoho, GMX, Yandex, Mail.ru
- FastMail, Hushmail
- And thousands more...

The list is automatically filtered to remove any domains that appear in the disposable list.

## Automatic Updates

The domain lists are automatically updated every **Sunday at midnight UTC** via GitHub Actions. The workflow:

1. Downloads the latest disposable domains list
2. Downloads the latest free email providers list
3. Filters out any free domains that are also disposable
4. Checks for changes
5. Commits and pushes updates if there are any changes
6. Triggers CI to ensure tests still pass

You can also manually trigger the update workflow from the GitHub Actions tab.

## Development

### Running Tests

```bash
go test -v ./...
```

### Running Benchmarks

```bash
go test -bench=. -benchmem ./...
```

### Code Coverage

```bash
go test -cover ./...
```

## Performance

The library uses embedded data and hash maps for O(1) lookups:

```
BenchmarkIsDisposableDomain-8     50000000    25.3 ns/op    0 B/op    0 allocs/op
BenchmarkIsFreeDomain-8           50000000    24.8 ns/op    0 B/op    0 allocs/op
BenchmarkIsBusinessDomain-8       30000000    48.7 ns/op    0 B/op    0 allocs/op
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Adding Free Email Providers

To add more free email providers, edit `data/free_domains.txt` and submit a PR.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Credits

- Disposable domains list: [disposable-email-domains](https://github.com/disposable-email-domains/disposable-email-domains)
- Free email providers list: [freemail](https://github.com/willwhite/freemail)

## Related Projects

- [snowid](https://github.com/rixlhq/snowid) - Distributed unique ID generator
