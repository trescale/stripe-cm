# stripe-cm

Declarative Stripe configuration manager CLI.

## Architecture

```
cmd/stripe-cm/          — CLI commands (init, sync, import, version)
internal/config/        — credential storage (file-based + macOS Keychain)
internal/stripe/        — Stripe API sync/import logic, one file per resource
  types.go              — all YAML config types
  business.go           — account profile, tax settings, portal configs
  webhooks.go           — webhook endpoints
  coupons.go            — coupons + promotion codes
  tax.go                — tax rates + tax registrations
  shipping.go           — shipping rates
  meters.go             — billing meters
  entitlements.go       — entitlement features
  payments.go           — Apple Pay domains
examples/full-config.yaml — comprehensive example with every resource type
```

## Key Patterns

- Every managed Stripe resource is tagged with `stripe_cm_managed: "true"` metadata
- Portal configs use `stripe_cm_profile: "<name>"` metadata for matching
- Sync functions: list existing managed resources, diff against config, create/update/delete
- Import functions: list from Stripe, filter to managed-only, sort deterministically
- Local-only sections (business, billing, checkout, tax.automatic) are preserved via `--merge`

## IMPORTANT

- **Always update `examples/full-config.yaml`** when adding new resource types, fields, or changing YAML schema. This is the primary reference for users and must demonstrate every configurable option with realistic values and comments.
- Use `stripe_cm_managed` / `stripe_cm_profile` as metadata keys (not retier)
- All types and functions in `internal/stripe` must be exported (capitalized)
- Run `golangci-lint run ./...` before pushing — CI enforces it
- The `version` and `commit` vars in `cmd/stripe-cm/root.go` are injected via ldflags at build time

## Build

```
go build -ldflags="-X main.version=dev -X main.commit=$(git rev-parse --short HEAD)" ./cmd/stripe-cm/
```
