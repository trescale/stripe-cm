package stripe

import (
	"log"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/account"
	portalconfig "github.com/stripe/stripe-go/v81/billingportal/configuration"
	"github.com/stripe/stripe-go/v81/price"
	"github.com/stripe/stripe-go/v81/product"
	taxsettings "github.com/stripe/stripe-go/v81/tax/settings"
)

func SyncBusiness(cfg BusinessConfig) {
	params := &stripe.AccountParams{}
	needsUpdate := false

	if cfg.Name != "" || cfg.SupportEmail != "" || cfg.SupportURL != "" || cfg.SupportPhone != "" {
		params.BusinessProfile = &stripe.AccountBusinessProfileParams{}
		needsUpdate = true
		if cfg.Name != "" {
			params.BusinessProfile.Name = stripe.String(cfg.Name)
		}
		if cfg.SupportEmail != "" {
			params.BusinessProfile.SupportEmail = stripe.String(cfg.SupportEmail)
		}
		if cfg.SupportURL != "" {
			params.BusinessProfile.SupportURL = stripe.String(cfg.SupportURL)
		}
		if cfg.SupportPhone != "" {
			params.BusinessProfile.SupportPhone = stripe.String(cfg.SupportPhone)
		}
	}

	if !needsUpdate {
		log.Println("  business: no changes")
		return
	}

	acct, err := account.Get()
	if err != nil {
		log.Printf("  business: warning — %v", err)
		return
	}
	_, err = account.Update(acct.ID, params)
	if err != nil {
		log.Printf("  business: warning — %v", err)
		return
	}
	log.Println("  business: updated")
}

func SyncTax(cfg TaxConfig) {
	params := &stripe.TaxSettingsParams{
		Defaults: &stripe.TaxSettingsDefaultsParams{},
	}

	if cfg.DefaultTaxCode != "" {
		params.Defaults.TaxCode = stripe.String(cfg.DefaultTaxCode)
	}

	_, err := taxsettings.Update(params)
	if err != nil {
		log.Printf("  tax: warning — %v (may require dashboard setup first)", err)
		return
	}
	log.Printf("  tax: automatic=%v default_code=%s", cfg.Automatic, cfg.DefaultTaxCode)
}

func SyncBilling(cfg BillingConfig) {
	log.Printf("  billing: invoices_due_in=%dd proration=%s missing_payment=%s",
		cfg.Invoices.DaysUntilDue, cfg.Subscriptions.ProrationBehavior, cfg.Subscriptions.MissingPaymentMethod)
	log.Printf("  billing: invoice_footer=%q", cfg.Invoices.Footer)
	log.Println("  billing: stored (applied per-request via checkout/subscription params)")
}

func ImportBusiness(merged BusinessConfig) BusinessConfig {
	acct, err := account.Get()
	if err != nil {
		log.Printf("  business: warning — %v", err)
		return merged
	}
	bc := merged
	if acct.BusinessProfile != nil {
		if acct.BusinessProfile.Name != "" {
			bc.Name = acct.BusinessProfile.Name
		}
		if acct.BusinessProfile.SupportEmail != "" {
			bc.SupportEmail = acct.BusinessProfile.SupportEmail
		}
		if acct.BusinessProfile.SupportURL != "" {
			bc.SupportURL = acct.BusinessProfile.SupportURL
		}
		if acct.BusinessProfile.SupportPhone != "" {
			bc.SupportPhone = acct.BusinessProfile.SupportPhone
		}
	}
	log.Printf("  business: imported (name=%q)", bc.Name)
	return bc
}

func ImportTax(merged TaxConfig) TaxConfig {
	settings, err := taxsettings.Get(&stripe.TaxSettingsParams{})
	if err != nil {
		log.Printf("  tax: warning — %v", err)
		return merged
	}
	tc := merged
	if settings.Defaults != nil && string(settings.Defaults.TaxCode) != "" {
		tc.DefaultTaxCode = string(settings.Defaults.TaxCode)
	}
	log.Printf("  tax: imported (automatic=%v code=%s)", tc.Automatic, tc.DefaultTaxCode)
	return tc
}

func ImportBillingDefaults() BillingConfig {
	log.Println("  billing: using defaults (billing settings are local-only)")
	return BillingConfig{
		Subscriptions: SubscriptionConfig{
			DefaultPaymentMethod: "card",
			ProrationBehavior:    "create_prorations",
			MissingPaymentMethod: "default_incomplete",
		},
		Invoices: InvoiceConfig{
			DaysUntilDue: 30,
		},
	}
}

func ImportCheckoutDefaults() CheckoutConfig {
	log.Println("  checkout: using defaults (checkout settings are local-only)")
	return CheckoutConfig{
		PaymentMethods:  []string{"card"},
		TaxIDCollection: true,
	}
}

// --- Portal sync/import ---

func SyncCustomerPortalProfiles(profiles []PortalProfile) map[string]string {
	existing := listExistingPortalConfigs()
	result := make(map[string]string, len(profiles))

	desired := make(map[string]bool, len(profiles))
	for _, p := range profiles {
		desired[p.Name] = true
		params := buildPortalParams(p)
		params.Active = stripe.Bool(true)

		if existingID, ok := existing[p.Name]; ok {
			_, err := portalconfig.Update(existingID, params)
			if err != nil {
				log.Printf("  portal [%s]: update warning — %v", p.Name, err)
				result[p.Name] = existingID
				continue
			}
			log.Printf("  portal [%s]: updated (%s)", p.Name, existingID)
			result[p.Name] = existingID
		} else {
			params.AddMetadata("stripe_cm_profile", p.Name)
			created, err := portalconfig.New(params)
			if err != nil {
				log.Printf("  portal [%s]: create warning — %v", p.Name, err)
				continue
			}
			log.Printf("  portal [%s]: created (%s)", p.Name, created.ID)
			result[p.Name] = created.ID
		}
	}

	for name, id := range existing {
		if desired[name] {
			continue
		}
		_, err := portalconfig.Update(id, &stripe.BillingPortalConfigurationParams{
			Active: stripe.Bool(false),
		})
		if err != nil {
			log.Printf("  portal [%s]: deactivate warning — %v", name, err)
			continue
		}
		log.Printf("  portal [%s]: deactivated (%s)", name, id)
	}

	return result
}

func buildPortalParams(p PortalProfile) *stripe.BillingPortalConfigurationParams {
	params := &stripe.BillingPortalConfigurationParams{
		BusinessProfile: &stripe.BillingPortalConfigurationBusinessProfileParams{},
		Features: &stripe.BillingPortalConfigurationFeaturesParams{
			CustomerUpdate: &stripe.BillingPortalConfigurationFeaturesCustomerUpdateParams{
				Enabled: stripe.Bool(true),
				AllowedUpdates: []*string{
					stripe.String("email"),
					stripe.String("address"),
					stripe.String("phone"),
					stripe.String("tax_id"),
				},
			},
			InvoiceHistory: &stripe.BillingPortalConfigurationFeaturesInvoiceHistoryParams{
				Enabled: stripe.Bool(p.Features.InvoiceHistory.Enabled),
			},
			PaymentMethodUpdate: &stripe.BillingPortalConfigurationFeaturesPaymentMethodUpdateParams{
				Enabled: stripe.Bool(p.Features.PaymentMethodUpdate.Enabled),
			},
			SubscriptionCancel: &stripe.BillingPortalConfigurationFeaturesSubscriptionCancelParams{
				Enabled: stripe.Bool(p.Features.SubscriptionCancel.Enabled),
			},
			SubscriptionUpdate: &stripe.BillingPortalConfigurationFeaturesSubscriptionUpdateParams{
				Enabled: stripe.Bool(p.Features.SubscriptionUpdate.Enabled),
			},
		},
	}

	if p.DefaultReturnURL != "" {
		params.DefaultReturnURL = stripe.String(p.DefaultReturnURL)
	}

	if p.Features.SubscriptionCancel.Enabled {
		if p.Features.SubscriptionCancel.Mode != "" {
			params.Features.SubscriptionCancel.Mode = stripe.String(p.Features.SubscriptionCancel.Mode)
		}
		if p.Features.SubscriptionCancel.ProrationBehavior != "" {
			params.Features.SubscriptionCancel.ProrationBehavior = stripe.String(p.Features.SubscriptionCancel.ProrationBehavior)
		}
	}

	if p.Features.SubscriptionUpdate.Enabled {
		params.Features.SubscriptionUpdate.DefaultAllowedUpdates = []*string{
			stripe.String("price"),
			stripe.String("quantity"),
			stripe.String("promotion_code"),
		}
		if p.Features.SubscriptionUpdate.ProrationBehavior != "" {
			params.Features.SubscriptionUpdate.ProrationBehavior = stripe.String(p.Features.SubscriptionUpdate.ProrationBehavior)
		}
		params.Features.SubscriptionUpdate.Products = buildPortalProducts()
	}

	return params
}

func buildPortalProducts() []*stripe.BillingPortalConfigurationFeaturesSubscriptionUpdateProductParams {
	listParams := &stripe.ProductListParams{}
	listParams.Filters.AddFilter("active", "", "true")
	listParams.Filters.AddFilter("limit", "", "100")

	var products []*stripe.BillingPortalConfigurationFeaturesSubscriptionUpdateProductParams
	iter := product.List(listParams)
	for iter.Next() {
		p := iter.Product()
		priceListParams := &stripe.PriceListParams{Product: stripe.String(p.ID)}
		priceListParams.Filters.AddFilter("active", "", "true")
		var priceIDs []*string
		priceIter := price.List(priceListParams)
		for priceIter.Next() {
			priceIDs = append(priceIDs, stripe.String(priceIter.Price().ID))
		}
		if len(priceIDs) > 0 {
			products = append(products, &stripe.BillingPortalConfigurationFeaturesSubscriptionUpdateProductParams{
				Product: stripe.String(p.ID),
				Prices:  priceIDs,
			})
		}
		if len(products) >= 10 {
			break
		}
	}
	return products
}

func listExistingPortalConfigs() map[string]string {
	result := make(map[string]string)
	params := &stripe.BillingPortalConfigurationListParams{}
	params.Filters.AddFilter("limit", "", "100")
	iter := portalconfig.List(params)
	for iter.Next() {
		cfg := iter.BillingPortalConfiguration()
		if name, ok := cfg.Metadata["stripe_cm_profile"]; ok && name != "" {
			result[name] = cfg.ID
		}
	}
	if err := iter.Err(); err != nil {
		log.Printf("  portal: warning listing existing configs — %v", err)
	}
	return result
}

func ImportPortalProfiles() []PortalProfile {
	params := &stripe.BillingPortalConfigurationListParams{}
	params.Filters.AddFilter("limit", "", "100")
	params.Filters.AddFilter("active", "", "true")

	var profiles []PortalProfile
	iter := portalconfig.List(params)
	for iter.Next() {
		cfg := iter.BillingPortalConfiguration()
		name := cfg.Metadata["stripe_cm_profile"]
		if name == "" {
			log.Printf("  portal [%s]: skipped (no stripe_cm_profile metadata)", cfg.ID)
			continue
		}
		p := PortalProfile{
			Name:             name,
			DefaultReturnURL: cfg.DefaultReturnURL,
			Features:         importPortalFeatures(cfg.Features),
		}
		profiles = append(profiles, p)
		log.Printf("  portal [%s]: imported (%s)", name, cfg.ID)
	}
	if err := iter.Err(); err != nil {
		log.Printf("  portal: warning — %v", err)
	}
	sortPortalProfiles(profiles)
	if len(profiles) == 0 {
		log.Println("  portal: no managed configurations found")
	}
	return profiles
}

func sortPortalProfiles(profiles []PortalProfile) {
	for i := 1; i < len(profiles); i++ {
		for j := i; j > 0 && profiles[j].Name < profiles[j-1].Name; j-- {
			profiles[j], profiles[j-1] = profiles[j-1], profiles[j]
		}
	}
}

func importPortalFeatures(f *stripe.BillingPortalConfigurationFeatures) PortalFeatures {
	if f == nil {
		return PortalFeatures{}
	}
	pf := PortalFeatures{}
	if f.SubscriptionCancel != nil {
		pf.SubscriptionCancel = PortalSubscriptionCancel{
			Enabled: f.SubscriptionCancel.Enabled,
		}
		if f.SubscriptionCancel.Enabled {
			pf.SubscriptionCancel.Mode = string(f.SubscriptionCancel.Mode)
			pf.SubscriptionCancel.ProrationBehavior = string(f.SubscriptionCancel.ProrationBehavior)
		}
	}
	if f.SubscriptionUpdate != nil {
		pf.SubscriptionUpdate = PortalSubscriptionUpdate{
			Enabled: f.SubscriptionUpdate.Enabled,
		}
		if f.SubscriptionUpdate.Enabled {
			pf.SubscriptionUpdate.ProrationBehavior = string(f.SubscriptionUpdate.ProrationBehavior)
		}
	}
	if f.PaymentMethodUpdate != nil {
		pf.PaymentMethodUpdate = PortalToggle{Enabled: f.PaymentMethodUpdate.Enabled}
	}
	if f.InvoiceHistory != nil {
		pf.InvoiceHistory = PortalToggle{Enabled: f.InvoiceHistory.Enabled}
	}
	return pf
}
