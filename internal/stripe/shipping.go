package stripe

import (
	"log"
	"sort"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/shippingrate"
)

func SyncShippingRates(rates []ShippingRateConfig) {
	existing := listManagedShippingRates()

	desired := make(map[string]bool, len(rates))
	for _, r := range rates {
		desired[r.DisplayName] = true

		if existingRate, ok := existing[r.DisplayName]; ok {
			params := &stripe.ShippingRateParams{
				Active:      stripe.Bool(r.Active),
				DisplayName: stripe.String(r.DisplayName),
			}
			params.AddMetadata("stripe_cm_managed", "true")
			if r.TaxBehavior != "" {
				params.TaxBehavior = stripe.String(r.TaxBehavior)
			}
			_, err := shippingrate.Update(existingRate.ID, params)
			if err != nil {
				log.Printf("  shipping [%s]: update warning — %v", r.DisplayName, err)
				continue
			}
			log.Printf("  shipping [%s]: updated", r.DisplayName)
		} else {
			params := &stripe.ShippingRateParams{
				DisplayName: stripe.String(r.DisplayName),
				Type:        stripe.String("fixed_amount"),
				FixedAmount: &stripe.ShippingRateFixedAmountParams{
					Amount:   stripe.Int64(r.Amount),
					Currency: stripe.String(r.Currency),
				},
			}
			params.AddMetadata("stripe_cm_managed", "true")
			if r.TaxBehavior != "" {
				params.TaxBehavior = stripe.String(r.TaxBehavior)
			}
			created, err := shippingrate.New(params)
			if err != nil {
				log.Printf("  shipping [%s]: create warning — %v", r.DisplayName, err)
				continue
			}
			log.Printf("  shipping [%s]: created (%s)", r.DisplayName, created.ID)
		}
	}

	for name, rate := range existing {
		if desired[name] {
			continue
		}
		_, err := shippingrate.Update(rate.ID, &stripe.ShippingRateParams{
			Active: stripe.Bool(false),
		})
		if err != nil {
			log.Printf("  shipping [%s]: deactivate warning — %v", name, err)
			continue
		}
		log.Printf("  shipping [%s]: deactivated", name)
	}
}

func ImportShippingRates() []ShippingRateConfig {
	params := &stripe.ShippingRateListParams{
		Active: stripe.Bool(true),
	}
	params.Filters.AddFilter("limit", "", "100")

	var rates []ShippingRateConfig
	iter := shippingrate.List(params)
	for iter.Next() {
		sr := iter.ShippingRate()
		if sr.Metadata["stripe_cm_managed"] != "true" {
			continue
		}
		taxBehavior := string(sr.TaxBehavior)
		if taxBehavior == "unspecified" {
			taxBehavior = ""
		}
		rc := ShippingRateConfig{
			DisplayName: sr.DisplayName,
			Active:      sr.Active,
			TaxBehavior: taxBehavior,
		}
		if sr.FixedAmount != nil {
			rc.Amount = sr.FixedAmount.Amount
			rc.Currency = string(sr.FixedAmount.Currency)
		}
		rates = append(rates, rc)
	}
	if err := iter.Err(); err != nil {
		log.Printf("  shipping: warning listing rates — %v", err)
	}

	sort.Slice(rates, func(i, j int) bool {
		return rates[i].DisplayName < rates[j].DisplayName
	})

	log.Printf("  shipping: imported %d managed rate(s)", len(rates))
	return rates
}

func listManagedShippingRates() map[string]*stripe.ShippingRate {
	result := make(map[string]*stripe.ShippingRate)
	params := &stripe.ShippingRateListParams{}
	params.Filters.AddFilter("limit", "", "100")

	iter := shippingrate.List(params)
	for iter.Next() {
		sr := iter.ShippingRate()
		if sr.Metadata["stripe_cm_managed"] != "true" {
			continue
		}
		result[sr.DisplayName] = sr
	}
	if err := iter.Err(); err != nil {
		log.Printf("  shipping: warning listing existing rates — %v", err)
	}
	return result
}
