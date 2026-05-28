package stripe

import (
	"log"
	"sort"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/applepaydomain"
)

func SyncApplePayDomains(domains []string) {
	existing := make(map[string]string) // domain name -> ID
	iter := applepaydomain.List(&stripe.ApplePayDomainListParams{})
	for iter.Next() {
		d := iter.ApplePayDomain()
		existing[d.DomainName] = d.ID
	}
	if err := iter.Err(); err != nil {
		log.Printf("  apple_pay: warning listing domains — %v", err)
		return
	}

	desired := make(map[string]bool, len(domains))
	for _, domain := range domains {
		desired[domain] = true
		if _, ok := existing[domain]; ok {
			continue
		}
		_, err := applepaydomain.New(&stripe.ApplePayDomainParams{
			DomainName: stripe.String(domain),
		})
		if err != nil {
			log.Printf("  apple_pay [%s]: warning — %v", domain, err)
			continue
		}
		log.Printf("  apple_pay [%s]: registered", domain)
	}

	for domain, id := range existing {
		if desired[domain] {
			continue
		}
		_, err := applepaydomain.Del(id, nil)
		if err != nil {
			log.Printf("  apple_pay [%s]: warning removing — %v", domain, err)
			continue
		}
		log.Printf("  apple_pay [%s]: removed", domain)
	}
}

func ImportApplePayDomains() []string {
	iter := applepaydomain.List(&stripe.ApplePayDomainListParams{})
	var domains []string
	for iter.Next() {
		domains = append(domains, iter.ApplePayDomain().DomainName)
	}
	if err := iter.Err(); err != nil {
		log.Printf("  apple_pay: warning — %v", err)
		return []string{}
	}

	sort.Strings(domains)
	log.Printf("  apple_pay: imported %d domain(s)", len(domains))

	if len(domains) == 0 {
		return []string{}
	}
	return domains
}
