package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/trescale/stripe-cm/internal/config"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize stripe-cm configuration",
	Long:  "Sets up credentials and generates a sample stripe.yaml in the current directory.",
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(_ *cobra.Command, _ []string) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Stripe API key: ")
	scanner.Scan()
	apiKey := strings.TrimSpace(scanner.Text())
	if apiKey == "" {
		return fmt.Errorf("API key is required")
	}

	fmt.Print("Storage mode (file/keychain) [file]: ")
	scanner.Scan()
	storage := strings.TrimSpace(scanner.Text())
	if storage == "" {
		storage = "file"
	}
	if storage != "file" && storage != "keychain" {
		return fmt.Errorf("storage must be 'file' or 'keychain'")
	}

	cfg := &config.Config{
		StripeAPIKey: apiKey,
		Storage:      storage,
	}
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("save config: %w", err)
	}
	fmt.Printf("Configuration saved to %s/config.yaml\n", config.ConfigDir())

	samplePath := "stripe.yaml"
	if _, err := os.Stat(samplePath); err == nil {
		fmt.Printf("%s already exists, skipping sample generation\n", samplePath)
		return nil
	}

	if err := os.WriteFile(samplePath, []byte(sampleYAML), 0644); err != nil {
		return fmt.Errorf("write sample config: %w", err)
	}
	fmt.Printf("Sample configuration written to %s\n", samplePath)

	return nil
}

const sampleYAML = `business:
  name: "My Company"
  support_email: "support@example.com"
  support_url: "https://example.com/support"
  support_phone: ""

tax:
  automatic: true
  default_tax_code: ""

billing:
  subscriptions:
    default_payment_method: "card"
    proration_behavior: "create_prorations"
    missing_payment_method: "default_incomplete"
  invoices:
    days_until_due: 30
    footer: ""

customer_portal:
  profiles:
    - name: default
      default_return_url: "https://example.com/dashboard"
      features:
        subscription_cancel:
          enabled: true
          mode: "at_period_end"
        subscription_update:
          enabled: true
          proration_behavior: "create_prorations"
        payment_method_update:
          enabled: true
        invoice_history:
          enabled: true

checkout:
  payment_methods:
    - card
  allow_promotion_codes: false
  collect_phone_number: false
  tax_id_collection: true

webhooks: []
coupons: []
promotion_codes: []
tax_rates: []
tax_registrations: []
shipping_rates: []
billing_meters: []
entitlements: []
apple_pay_domains: []
`
