package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version        = "dev"
	commit         = "none"
	stripeYAMLFile string
)

var rootCmd = &cobra.Command{
	Use:   "stripe-cm",
	Short: "Declarative Stripe configuration management",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("stripe-cm %s (commit: %s)\n", version, commit)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&stripeYAMLFile, "config", "stripe.yaml", "Path to Stripe configuration YAML file")
	rootCmd.AddCommand(versionCmd)
	rootCmd.Version = version
}
