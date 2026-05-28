package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/stripe/stripe-go/v81"
	"github.com/trescale/stripe-cm/internal/config"
	scm "github.com/trescale/stripe-cm/internal/stripe"
	"gopkg.in/yaml.v3"
)

var (
	importOutput string
	importMerge  string
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import current Stripe configuration into a YAML file",
	Long: `Reads the current Stripe account configuration and writes it as a YAML file.

Sections that are stored on Stripe (customer_portal profiles) are fetched live.
Sections that are local-only (business, billing, checkout, tax.automatic) are
preserved from an existing file when --merge is set.`,
	RunE: runImport,
}

func init() {
	importCmd.Flags().StringVarP(&importOutput, "output", "o", "stripe.yaml", "Output YAML file path")
	importCmd.Flags().StringVarP(&importMerge, "merge", "m", "", "Existing YAML file to merge local-only sections from")
	rootCmd.AddCommand(importCmd)
}

func runImport(_ *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load credentials: %w", err)
	}
	if cfg.StripeAPIKey == "" {
		return fmt.Errorf("stripe API key not configured — run 'stripe-cm init'")
	}
	stripe.Key = cfg.StripeAPIKey

	var existing scm.Config
	if importMerge != "" {
		data, err := os.ReadFile(importMerge)
		if err != nil {
			return fmt.Errorf("read merge file: %w", err)
		}
		if err := yaml.Unmarshal(data, &existing); err != nil {
			return fmt.Errorf("parse merge file: %w", err)
		}
		log.Printf("merging local-only sections from %s", importMerge)
	}

	sc := scm.Config{}
	sc.Business = scm.ImportBusiness(existing.Business)
	sc.Tax = scm.ImportTax(existing.Tax)
	sc.Billing = existing.Billing
	if sc.Billing.Subscriptions.DefaultPaymentMethod == "" {
		sc.Billing = scm.ImportBillingDefaults()
	}
	sc.CustomerPortal.Profiles = scm.ImportPortalProfiles()
	sc.Checkout = existing.Checkout
	if len(sc.Checkout.PaymentMethods) == 0 {
		sc.Checkout = scm.ImportCheckoutDefaults()
	}
	sc.Webhooks = scm.ImportWebhooks()
	sc.Coupons = scm.ImportCoupons()
	sc.PromotionCodes = scm.ImportPromotionCodes()
	sc.TaxRates = scm.ImportTaxRates()
	sc.TaxRegs = scm.ImportTaxRegistrations()
	sc.ShippingRates = scm.ImportShippingRates()
	sc.BillingMeters = scm.ImportBillingMeters()
	sc.Entitlements = scm.ImportEntitlements()
	sc.ApplePayDomains = scm.ImportApplePayDomains()
	sc.Branding = scm.ImportBranding(existing.Branding)
	sc.PayoutSchedule = scm.ImportPayoutSchedule(existing.PayoutSchedule)
	sc.PaymentMethods = scm.ImportPaymentMethodsConfig()
	sc.BillingAlerts = scm.ImportBillingAlerts()

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	if err := enc.Encode(&sc); err != nil {
		return fmt.Errorf("marshal yaml: %w", err)
	}
	if err := enc.Close(); err != nil {
		return fmt.Errorf("close yaml encoder: %w", err)
	}
	if err := os.WriteFile(importOutput, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	log.Printf("stripe configuration imported to %s", importOutput)
	return nil
}
