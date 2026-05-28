package main

import (
	"github.com/spf13/cobra"
)

var (
	version        = "dev"
	stripeYAMLFile string
)

var rootCmd = &cobra.Command{
	Use:     "stripe-cm",
	Short:   "Declarative Stripe configuration management",
	Version: version,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&stripeYAMLFile, "config", "stripe.yaml", "Path to Stripe configuration YAML file")
}
