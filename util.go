package workemailvalidator

import "golang.org/x/net/idna"

var lookupProfile = idna.New(
	idna.MapForLookup(),
	idna.BidiRule(),
	idna.StrictDomainName(true),
	idna.VerifyDNSLength(true),
)

// domainToASCII converts an internationalized domain name to the ASCII form used for DNS lookups
// and reports whether the name is valid.
// Reference: https://en.wikipedia.org/wiki/Punycode
func domainToASCII(domain string) (string, bool) {
	ascii, err := lookupProfile.ToASCII(domain)
	return ascii, err == nil
}
