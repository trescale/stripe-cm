package main

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var (
	version        = "dev"
	commit         = "none"
	stripeYAMLFile string
)

func resolveVersion() {
	if version != "dev" {
		return
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}
	if info.Main.Version != "" && info.Main.Version != "(devel)" {
		version = info.Main.Version
	}
	for _, s := range info.Settings {
		if s.Key == "vcs.revision" && len(s.Value) >= 7 {
			commit = s.Value[:7]
		}
	}
}

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
	resolveVersion()
	rootCmd.PersistentFlags().StringVar(&stripeYAMLFile, "config", "stripe.yaml", "Path to Stripe configuration YAML file")
	rootCmd.AddCommand(versionCmd)
	rootCmd.Version = version
}
