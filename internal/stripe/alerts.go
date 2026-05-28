package stripe

import (
	"log"
	"sort"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/billing/alert"
)

func SyncBillingAlerts(alerts []BillingAlertConfig) {
	existing := make(map[string]*stripe.BillingAlert)
	params := &stripe.BillingAlertListParams{}
	params.Filters.AddFilter("limit", "", "100")
	iter := alert.List(params)
	for iter.Next() {
		a := iter.BillingAlert()
		existing[a.Title] = a
	}
	if err := iter.Err(); err != nil {
		log.Printf("  alerts: warning listing existing — %v", err)
		return
	}

	// Note: billing alerts do not support custom metadata, matching by title only
	desired := make(map[string]bool, len(alerts))
	for _, cfg := range alerts {
		desired[cfg.Title] = true

		if _, ok := existing[cfg.Title]; ok {
			log.Printf("  alerts [%s]: already exists, skipping", cfg.Title)
		} else {
			createParams := &stripe.BillingAlertParams{
				Title:     stripe.String(cfg.Title),
				AlertType: stripe.String(cfg.AlertType),
			}
			if cfg.AlertType == "usage_threshold" && cfg.Threshold > 0 {
				createParams.UsageThreshold = &stripe.BillingAlertUsageThresholdParams{
					GTE:        stripe.Int64(cfg.Threshold),
					Recurrence: stripe.String(string(stripe.BillingAlertUsageThresholdRecurrenceOneTime)),
				}
			}
			_, err := alert.New(createParams)
			if err != nil {
				log.Printf("  alerts [%s]: create warning — %v", cfg.Title, err)
				continue
			}
			log.Printf("  alerts [%s]: created", cfg.Title)
		}
	}

	for title, a := range existing {
		if desired[title] {
			continue
		}
		if a.Status != stripe.BillingAlertStatusActive {
			continue
		}
		_, err := alert.Deactivate(a.ID, &stripe.BillingAlertDeactivateParams{})
		if err != nil {
			log.Printf("  alerts [%s]: deactivate warning — %v", title, err)
			continue
		}
		log.Printf("  alerts [%s]: deactivated", title)
	}
}

func ImportBillingAlerts() []BillingAlertConfig {
	params := &stripe.BillingAlertListParams{}
	params.Filters.AddFilter("limit", "", "100")

	var configs []BillingAlertConfig
	iter := alert.List(params)
	for iter.Next() {
		a := iter.BillingAlert()
		if a.Status != stripe.BillingAlertStatusActive {
			continue
		}
		cfg := BillingAlertConfig{
			Title:     a.Title,
			AlertType: string(a.AlertType),
		}
		if a.UsageThreshold != nil {
			cfg.Threshold = a.UsageThreshold.GTE
		}
		configs = append(configs, cfg)
	}
	if err := iter.Err(); err != nil {
		log.Printf("  alerts: warning — %v", err)
		return nil
	}

	sort.Slice(configs, func(i, j int) bool {
		return configs[i].Title < configs[j].Title
	})

	log.Printf("  alerts: imported %d active alert(s)", len(configs))
	return configs
}
