package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/stripe/stripe-go/v81"
	"github.com/trescale/stripe-cm/internal/config"
	scm "github.com/trescale/stripe-cm/internal/stripe"
	"gopkg.in/yaml.v3"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Apply Stripe configuration from a YAML file",
	RunE:  runSync,
}

func init() {
	rootCmd.AddCommand(syncCmd)
}

func runSync(_ *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load credentials: %w", err)
	}
	if cfg.StripeAPIKey == "" {
		return fmt.Errorf("stripe API key not configured — run 'stripe-cm init'")
	}
	stripe.Key = cfg.StripeAPIKey

	data, err := os.ReadFile(stripeYAMLFile)
	if err != nil {
		return fmt.Errorf("read stripe config: %w", err)
	}
	var sc scm.Config
	if err := yaml.Unmarshal(data, &sc); err != nil {
		return fmt.Errorf("parse stripe config: %w", err)
	}

	scm.SyncBusiness(sc.Business)
	scm.SyncTax(sc.Tax)
	scm.SyncBilling(sc.Billing)
	portalIDs := scm.SyncCustomerPortalProfiles(sc.CustomerPortal.Profiles)
	scm.SyncWebhooks(sc.Webhooks)
	scm.SyncCoupons(sc.Coupons)
	scm.SyncPromotionCodes(sc.PromotionCodes)
	scm.SyncTaxRates(sc.TaxRates)
	scm.SyncTaxRegistrations(sc.TaxRegs)
	scm.SyncShippingRates(sc.ShippingRates)
	scm.SyncBillingMeters(sc.BillingMeters)
	scm.SyncEntitlements(sc.Entitlements)
	scm.SyncApplePayDomains(sc.ApplePayDomains)
	scm.SyncBranding(sc.Branding)
	scm.SyncPayoutSchedule(sc.PayoutSchedule)
	scm.SyncPaymentMethodsConfig(sc.PaymentMethods)
	scm.SyncBillingAlerts(sc.BillingAlerts)

	if len(portalIDs) > 0 {
		out, _ := json.MarshalIndent(portalIDs, "", "  ")
		log.Printf("  portal profile IDs:\n%s", string(out))
	}

	log.Println("stripe configuration synced successfully")
	return nil
}
