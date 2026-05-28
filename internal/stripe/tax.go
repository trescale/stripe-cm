package stripe

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/tax/registration"
	"github.com/stripe/stripe-go/v81/taxrate"
)

func taxRateKey(displayName, country, state string) string {
	return fmt.Sprintf("%s|%s|%s", displayName, country, state)
}

func SyncTaxRates(rates []TaxRateConfig) {
	existing := listManagedTaxRates()

	desired := make(map[string]bool, len(rates))
	for _, r := range rates {
		key := taxRateKey(r.DisplayName, r.Country, r.State)
		desired[key] = true

		if existingRate, ok := existing[key]; ok {
			params := &stripe.TaxRateParams{
				Active:      stripe.Bool(r.Active),
				Description: stripe.String(r.Description),
				TaxType:     stripe.String(r.TaxType),
			}
			params.AddMetadata("stripe_cm_managed", "true")
			_, err := taxrate.Update(existingRate.ID, params)
			if err != nil {
				log.Printf("  tax_rates [%s]: update warning — %v", r.DisplayName, err)
				continue
			}
			log.Printf("  tax_rates [%s]: updated (%s)", r.DisplayName, existingRate.ID)
		} else {
			params := &stripe.TaxRateParams{
				DisplayName: stripe.String(r.DisplayName),
				Description: stripe.String(r.Description),
				Percentage:  stripe.Float64(r.Percentage),
				Country:     stripe.String(r.Country),
				State:       stripe.String(r.State),
				TaxType:     stripe.String(r.TaxType),
				Inclusive:    stripe.Bool(r.Inclusive),
				Active:      stripe.Bool(r.Active),
			}
			params.AddMetadata("stripe_cm_managed", "true")
			created, err := taxrate.New(params)
			if err != nil {
				log.Printf("  tax_rates [%s]: create warning — %v", r.DisplayName, err)
				continue
			}
			log.Printf("  tax_rates [%s]: created (%s)", r.DisplayName, created.ID)
		}
	}

	for key, rate := range existing {
		if desired[key] {
			continue
		}
		_, err := taxrate.Update(rate.ID, &stripe.TaxRateParams{
			Active: stripe.Bool(false),
		})
		if err != nil {
			log.Printf("  tax_rates [%s]: deactivate warning — %v", rate.DisplayName, err)
			continue
		}
		log.Printf("  tax_rates [%s]: deactivated (%s)", rate.DisplayName, rate.ID)
	}
}

func listManagedTaxRates() map[string]*stripe.TaxRate {
	result := make(map[string]*stripe.TaxRate)
	params := &stripe.TaxRateListParams{}
	params.Filters.AddFilter("limit", "", "100")
	iter := taxrate.List(params)
	for iter.Next() {
		r := iter.TaxRate()
		if r.Metadata["stripe_cm_managed"] != "true" {
			continue
		}
		key := taxRateKey(r.DisplayName, r.Country, r.State)
		result[key] = r
	}
	if err := iter.Err(); err != nil {
		log.Printf("  tax_rates: warning listing existing — %v", err)
	}
	return result
}

func SyncTaxRegistrations(regs []TaxRegConfig) {
	existing := listActiveRegistrations()

	desired := make(map[string]bool, len(regs))
	for _, r := range regs {
		key := regKey(r.Country, r.State)
		desired[key] = true

		if _, ok := existing[key]; ok {
			log.Printf("  tax_registrations [%s/%s]: already exists", r.Country, r.State)
			continue
		}

		params := buildRegistrationParams(r)
		created, err := registration.New(params)
		if err != nil {
			log.Printf("  tax_registrations [%s/%s]: create warning — %v", r.Country, r.State, err)
			continue
		}
		log.Printf("  tax_registrations [%s/%s]: created (%s)", r.Country, r.State, created.ID)
	}

	for key, reg := range existing {
		if desired[key] {
			continue
		}
		expiresAt := time.Now().Unix()
		_, err := registration.Update(reg.ID, &stripe.TaxRegistrationParams{
			ExpiresAt: stripe.Int64(expiresAt),
		})
		if err != nil {
			log.Printf("  tax_registrations [%s]: expire warning — %v", key, err)
			continue
		}
		log.Printf("  tax_registrations [%s]: expired (%s)", key, reg.ID)
	}
}

func regKey(country, state string) string {
	if state != "" {
		return fmt.Sprintf("%s|%s", country, state)
	}
	return country
}

func buildRegistrationParams(r TaxRegConfig) *stripe.TaxRegistrationParams {
	params := &stripe.TaxRegistrationParams{
		Country:       stripe.String(r.Country),
		ActiveFromNow: stripe.Bool(true),
		CountryOptions: &stripe.TaxRegistrationCountryOptionsParams{
			US: &stripe.TaxRegistrationCountryOptionsUSParams{
				Type: stripe.String(r.Type),
			},
		},
	}

	if r.Country == "US" && r.State != "" {
		params.CountryOptions.US.State = stripe.String(r.State)
	}

	if r.Country == "CA" {
		params.CountryOptions = &stripe.TaxRegistrationCountryOptionsParams{
			Ca: &stripe.TaxRegistrationCountryOptionsCaParams{
				Type: stripe.String(r.Type),
			},
		}
		if r.State != "" {
			params.CountryOptions.Ca.ProvinceStandard = &stripe.TaxRegistrationCountryOptionsCaProvinceStandardParams{
				Province: stripe.String(r.State),
			}
		}
	}

	return params
}

func listActiveRegistrations() map[string]*stripe.TaxRegistration {
	result := make(map[string]*stripe.TaxRegistration)
	params := &stripe.TaxRegistrationListParams{
		Status: stripe.String("active"),
	}
	params.Filters.AddFilter("limit", "", "100")
	iter := registration.List(params)
	for iter.Next() {
		r := iter.TaxRegistration()
		state := extractRegistrationState(r)
		key := regKey(r.Country, state)
		result[key] = r
	}
	if err := iter.Err(); err != nil {
		log.Printf("  tax_registrations: warning listing existing — %v", err)
	}
	return result
}

func extractRegistrationState(r *stripe.TaxRegistration) string {
	if r.CountryOptions == nil {
		return ""
	}
	if r.Country == "US" && r.CountryOptions.US != nil {
		return r.CountryOptions.US.State
	}
	if r.Country == "CA" && r.CountryOptions.Ca != nil && r.CountryOptions.Ca.ProvinceStandard != nil {
		return r.CountryOptions.Ca.ProvinceStandard.Province
	}
	return ""
}

func extractRegistrationType(r *stripe.TaxRegistration) string {
	if r.CountryOptions == nil {
		return ""
	}
	if r.Country == "US" && r.CountryOptions.US != nil {
		return string(r.CountryOptions.US.Type)
	}
	if r.Country == "CA" && r.CountryOptions.Ca != nil {
		return string(r.CountryOptions.Ca.Type)
	}
	if r.Country == "GB" && r.CountryOptions.GB != nil {
		return string(r.CountryOptions.GB.Type)
	}
	if r.Country == "AU" && r.CountryOptions.Au != nil {
		return string(r.CountryOptions.Au.Type)
	}
	if r.Country == "DE" && r.CountryOptions.DE != nil {
		return string(r.CountryOptions.DE.Type)
	}
	if r.Country == "FR" && r.CountryOptions.FR != nil {
		return string(r.CountryOptions.FR.Type)
	}
	return ""
}

func ImportTaxRates() []TaxRateConfig {
	params := &stripe.TaxRateListParams{
		Active: stripe.Bool(true),
	}
	params.Filters.AddFilter("limit", "", "100")

	var rates []TaxRateConfig
	iter := taxrate.List(params)
	for iter.Next() {
		r := iter.TaxRate()
		if r.Metadata["stripe_cm_managed"] != "true" {
			continue
		}
		rates = append(rates, TaxRateConfig{
			DisplayName: r.DisplayName,
			Description: r.Description,
			Percentage:  r.Percentage,
			Country:     r.Country,
			State:       r.State,
			TaxType:     string(r.TaxType),
			Inclusive:    r.Inclusive,
			Active:      r.Active,
		})
	}
	if err := iter.Err(); err != nil {
		log.Printf("  tax_rates: warning — %v", err)
	}

	sort.Slice(rates, func(i, j int) bool {
		return rates[i].DisplayName < rates[j].DisplayName
	})

	log.Printf("  tax_rates: imported %d managed rate(s)", len(rates))
	return rates
}

func ImportTaxRegistrations() []TaxRegConfig {
	params := &stripe.TaxRegistrationListParams{
		Status: stripe.String("active"),
	}
	params.Filters.AddFilter("limit", "", "100")

	var regs []TaxRegConfig
	iter := registration.List(params)
	for iter.Next() {
		r := iter.TaxRegistration()
		state := extractRegistrationState(r)
		regType := extractRegistrationType(r)
		regs = append(regs, TaxRegConfig{
			Country: r.Country,
			State:   state,
			Type:    regType,
		})
	}
	if err := iter.Err(); err != nil {
		log.Printf("  tax_registrations: warning — %v", err)
	}

	sort.Slice(regs, func(i, j int) bool {
		if regs[i].Country != regs[j].Country {
			return regs[i].Country < regs[j].Country
		}
		return regs[i].State < regs[j].State
	})

	log.Printf("  tax_registrations: imported %d active registration(s)", len(regs))
	return regs
}
