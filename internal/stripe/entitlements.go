package stripe

import (
	"log"
	"sort"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/entitlements/feature"
)

func SyncEntitlements(features []EntitlementConfig) {
	existing := make(map[string]*stripe.EntitlementsFeature)
	params := &stripe.EntitlementsFeatureListParams{}
	params.Filters.AddFilter("limit", "", "100")
	iter := feature.List(params)
	for iter.Next() {
		f := iter.EntitlementsFeature()
		if f.Metadata["stripe_cm_managed"] != "true" {
			continue
		}
		existing[f.LookupKey] = f
	}
	if err := iter.Err(); err != nil {
		log.Printf("  entitlements: warning listing existing — %v (entitlements API may require opt-in)", err)
		return
	}

	desired := make(map[string]bool, len(features))
	for _, cfg := range features {
		lookupKey := cfg.Lookup
		if lookupKey == "" {
			lookupKey = cfg.Name
		}
		desired[lookupKey] = true

		if ex, ok := existing[lookupKey]; ok {
			if ex.Name != cfg.Name {
				updateParams := &stripe.EntitlementsFeatureParams{
					Name: stripe.String(cfg.Name),
				}
				_, err := feature.Update(ex.ID, updateParams)
				if err != nil {
					log.Printf("  entitlements [%s]: update warning — %v", lookupKey, err)
					continue
				}
				log.Printf("  entitlements [%s]: updated", lookupKey)
			} else {
				log.Printf("  entitlements [%s]: unchanged", lookupKey)
			}
		} else {
			createParams := &stripe.EntitlementsFeatureParams{
				Name:      stripe.String(cfg.Name),
				LookupKey: stripe.String(lookupKey),
				Metadata: map[string]string{
					"stripe_cm_managed": "true",
				},
			}
			_, err := feature.New(createParams)
			if err != nil {
				log.Printf("  entitlements [%s]: create warning — %v", lookupKey, err)
				continue
			}
			log.Printf("  entitlements [%s]: created", lookupKey)
		}
	}

	for lookupKey, f := range existing {
		if desired[lookupKey] {
			continue
		}
		if !f.Active {
			continue
		}
		archiveParams := &stripe.EntitlementsFeatureParams{
			Active: stripe.Bool(false),
		}
		_, err := feature.Update(f.ID, archiveParams)
		if err != nil {
			log.Printf("  entitlements [%s]: archive warning — %v", lookupKey, err)
			continue
		}
		log.Printf("  entitlements [%s]: archived", lookupKey)
	}
}

func ImportEntitlements() []EntitlementConfig {
	params := &stripe.EntitlementsFeatureListParams{}
	params.Filters.AddFilter("limit", "", "100")

	var configs []EntitlementConfig
	iter := feature.List(params)
	for iter.Next() {
		f := iter.EntitlementsFeature()
		if f.Metadata["stripe_cm_managed"] != "true" {
			continue
		}
		if !f.Active {
			continue
		}
		cfg := EntitlementConfig{
			Name:   f.Name,
			Lookup: f.LookupKey,
		}
		configs = append(configs, cfg)
	}
	if err := iter.Err(); err != nil {
		log.Printf("  entitlements: warning — %v (entitlements API may require opt-in)", err)
		return nil
	}

	sort.Slice(configs, func(i, j int) bool {
		return configs[i].Name < configs[j].Name
	})

	log.Printf("  entitlements: imported %d active feature(s)", len(configs))
	return configs
}
