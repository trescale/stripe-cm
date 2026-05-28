package stripe

import (
	"log"
	"sort"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhookendpoint"
)

func SyncWebhooks(webhooks []WebhookConfig) {
	existing := listManagedWebhookEndpoints()

	desired := make(map[string]bool, len(webhooks))
	for _, wh := range webhooks {
		desired[wh.URL] = true

		if existingEP, ok := existing[wh.URL]; ok {
			params := &stripe.WebhookEndpointParams{
				EnabledEvents: toStringPtrs(wh.EnabledEvents),
				Description:   stripe.String(wh.Description),
				Disabled:      stripe.Bool(wh.Disabled),
			}
			_, err := webhookendpoint.Update(existingEP.ID, params)
			if err != nil {
				log.Printf("  webhook [%s]: update warning — %v", wh.URL, err)
				continue
			}
			log.Printf("  webhook [%s]: updated (%s)", wh.URL, existingEP.ID)
		} else {
			params := &stripe.WebhookEndpointParams{
				URL:           stripe.String(wh.URL),
				EnabledEvents: toStringPtrs(wh.EnabledEvents),
				Description:   stripe.String(wh.Description),
			}
			if wh.APIVersion != "" {
				params.APIVersion = stripe.String(wh.APIVersion)
			}
			params.AddMetadata("stripe_cm_managed", "true")
			created, err := webhookendpoint.New(params)
			if err != nil {
				log.Printf("  webhook [%s]: create warning — %v", wh.URL, err)
				continue
			}
			log.Printf("  webhook [%s]: created (%s)", wh.URL, created.ID)
		}
	}

	for url, ep := range existing {
		if desired[url] {
			continue
		}
		_, err := webhookendpoint.Del(ep.ID, nil)
		if err != nil {
			log.Printf("  webhook [%s]: delete warning — %v", url, err)
			continue
		}
		log.Printf("  webhook [%s]: deleted (%s)", url, ep.ID)
	}
}

func ImportWebhooks() []WebhookConfig {
	params := &stripe.WebhookEndpointListParams{}
	params.Filters.AddFilter("limit", "", "100")

	var configs []WebhookConfig
	iter := webhookendpoint.List(params)
	for iter.Next() {
		ep := iter.WebhookEndpoint()
		if ep.Metadata["stripe_cm_managed"] != "true" {
			continue
		}
		wc := WebhookConfig{
			URL:           ep.URL,
			Description:   ep.Description,
			EnabledEvents: ep.EnabledEvents,
			APIVersion:    ep.APIVersion,
			Disabled:      ep.Status == "disabled",
		}
		configs = append(configs, wc)
		log.Printf("  webhook [%s]: imported (%s)", ep.URL, ep.ID)
	}
	if err := iter.Err(); err != nil {
		log.Printf("  webhook: warning — %v", err)
	}
	sort.Slice(configs, func(i, j int) bool {
		return configs[i].URL < configs[j].URL
	})
	if len(configs) == 0 {
		log.Println("  webhook: no managed endpoints found")
	}
	return configs
}

func listManagedWebhookEndpoints() map[string]*stripe.WebhookEndpoint {
	result := make(map[string]*stripe.WebhookEndpoint)
	params := &stripe.WebhookEndpointListParams{}
	params.Filters.AddFilter("limit", "", "100")
	iter := webhookendpoint.List(params)
	for iter.Next() {
		ep := iter.WebhookEndpoint()
		if ep.Metadata["stripe_cm_managed"] == "true" {
			result[ep.URL] = ep
		}
	}
	if err := iter.Err(); err != nil {
		log.Printf("  webhook: warning listing existing endpoints — %v", err)
	}
	return result
}

func toStringPtrs(ss []string) []*string {
	ptrs := make([]*string, len(ss))
	for i, s := range ss {
		ptrs[i] = stripe.String(s)
	}
	return ptrs
}
