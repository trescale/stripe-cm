package stripe

import (
	"log"
	"sort"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/billing/meter"
)

func SyncBillingMeters(meters []BillingMeterConfig) {
	existing := make(map[string]*stripe.BillingMeter)
	params := &stripe.BillingMeterListParams{}
	params.Filters.AddFilter("limit", "", "100")
	iter := meter.List(params)
	for iter.Next() {
		m := iter.BillingMeter()
		existing[m.EventName] = m
	}
	if err := iter.Err(); err != nil {
		log.Printf("  meters: warning listing existing — %v", err)
		return
	}

	desired := make(map[string]bool, len(meters))
	for _, cfg := range meters {
		desired[cfg.EventName] = true

		if ex, ok := existing[cfg.EventName]; ok {
			updateParams := &stripe.BillingMeterParams{
				DisplayName: stripe.String(cfg.DisplayName),
			}
			_, err := meter.Update(ex.ID, updateParams)
			if err != nil {
				log.Printf("  meters [%s]: update warning — %v", cfg.EventName, err)
				continue
			}
			log.Printf("  meters [%s]: updated", cfg.EventName)
		} else {
			createParams := &stripe.BillingMeterParams{
				EventName:   stripe.String(cfg.EventName),
				DisplayName: stripe.String(cfg.DisplayName),
				DefaultAggregation: &stripe.BillingMeterDefaultAggregationParams{
					Formula: stripe.String(cfg.Aggregation),
				},
			}
			if cfg.Aggregation == "sum" && cfg.ValueKey != "" {
				createParams.ValueSettings = &stripe.BillingMeterValueSettingsParams{
					EventPayloadKey: stripe.String(cfg.ValueKey),
				}
			}
			if cfg.EventTimeWindow != "" {
				createParams.EventTimeWindow = stripe.String(cfg.EventTimeWindow)
			}
			if cfg.CustomerMapping != "" {
				createParams.CustomerMapping = &stripe.BillingMeterCustomerMappingParams{
					Type:            stripe.String(cfg.CustomerMapping),
					EventPayloadKey: stripe.String("stripe_customer_id"),
				}
			}
			_, err := meter.New(createParams)
			if err != nil {
				log.Printf("  meters [%s]: create warning — %v", cfg.EventName, err)
				continue
			}
			log.Printf("  meters [%s]: created", cfg.EventName)
		}
	}

	for eventName, m := range existing {
		if desired[eventName] {
			continue
		}
		if m.Status != stripe.BillingMeterStatusActive {
			continue
		}
		_, err := meter.Deactivate(m.ID, &stripe.BillingMeterDeactivateParams{})
		if err != nil {
			log.Printf("  meters [%s]: deactivate warning — %v", eventName, err)
			continue
		}
		log.Printf("  meters [%s]: deactivated", eventName)
	}
}

func ImportBillingMeters() []BillingMeterConfig {
	params := &stripe.BillingMeterListParams{}
	params.Filters.AddFilter("limit", "", "100")
	params.Status = stripe.String(string(stripe.BillingMeterStatusActive))

	var configs []BillingMeterConfig
	iter := meter.List(params)
	for iter.Next() {
		m := iter.BillingMeter()
		cfg := BillingMeterConfig{
			EventName:   m.EventName,
			DisplayName: m.DisplayName,
		}
		if m.DefaultAggregation != nil {
			cfg.Aggregation = string(m.DefaultAggregation.Formula)
		}
		if m.ValueSettings != nil && m.ValueSettings.EventPayloadKey != "" {
			cfg.ValueKey = m.ValueSettings.EventPayloadKey
		}
		if m.EventTimeWindow != "" {
			cfg.EventTimeWindow = string(m.EventTimeWindow)
		}
		if m.CustomerMapping != nil && m.CustomerMapping.Type != "" {
			cfg.CustomerMapping = string(m.CustomerMapping.Type)
		}
		configs = append(configs, cfg)
	}
	if err := iter.Err(); err != nil {
		log.Printf("  meters: warning — %v", err)
	}

	sort.Slice(configs, func(i, j int) bool {
		return configs[i].EventName < configs[j].EventName
	})

	log.Printf("  meters: imported %d active meter(s)", len(configs))
	return configs
}
