# stripe-cm

Declarative Stripe configuration manager. Define your Stripe settings in a YAML file and sync them across environments.

## Features

- **Declarative config**: Define webhooks, coupons, portal profiles, tax rates, shipping rates, billing meters, entitlements, and more in YAML
- **Bidirectional sync**: Push config to Stripe (`sync`) or pull current state (`import`)
- **Round-trip safe**: `sync` -> `import --merge` produces identical values
- **Selective management**: Only manages resources tagged with `stripe-cm` metadata -- won't touch manually created Stripe resources
- **Credential storage**: File-based or macOS Keychain

## Install

### Homebrew (macOS/Linux)

```bash
# Coming soon
```

### Go install

```bash
go install github.com/trescale/stripe-cm/cmd/stripe-cm@latest
```

### Docker

```bash
docker run --rm -v $(pwd):/config ghcr.io/trescale/stripe-cm sync --config /config/stripe.yaml
```

### Binary releases

Download from [GitHub Releases](https://github.com/trescale/stripe-cm/releases).

## Quick Start

```bash
# 1. Initialize credentials
stripe-cm init

# 2. Create your config file (or import existing)
stripe-cm import -o stripe.yaml

# 3. Edit stripe.yaml to your needs

# 4. Sync to Stripe
stripe-cm sync --config stripe.yaml
```

## Commands

### `stripe-cm init`

Interactive setup -- stores your Stripe API key in `~/.stripe-cm/config.yaml` or macOS Keychain.

### `stripe-cm sync`

Applies the YAML config to Stripe. Creates, updates, or removes resources to match the declared state.

```bash
stripe-cm sync --config stripe.yaml
```

### `stripe-cm import`

Reads current Stripe configuration and writes it to a YAML file.

```bash
# Fresh import
stripe-cm import -o stripe.yaml

# Merge with existing (preserves local-only sections)
stripe-cm import -o stripe.yaml --merge stripe.yaml
```

## Configuration Reference

```yaml
business:
  name: "My Company"
  support_email: "support@example.com"
  support_url: "https://example.com/support"

tax:
  automatic: true
  default_tax_code: "txcd_10000000"

billing:
  subscriptions:
    default_payment_method: "card"
    proration_behavior: "create_prorations"
  invoices:
    days_until_due: 30
    footer: "Thank you for your business"

customer_portal:
  profiles:
    - name: "default"
      default_return_url: "https://app.example.com"
      features:
        subscription_cancel:
          enabled: true
          mode: "at_period_end"
        subscription_update:
          enabled: true
        payment_method_update:
          enabled: true
        invoice_history:
          enabled: true

checkout:
  payment_methods: [card]
  tax_id_collection: true

webhooks:
  - url: "https://api.example.com/webhooks/stripe"
    enabled_events:
      - checkout.session.completed
      - customer.subscription.updated
      - invoice.paid

coupons:
  - id: "LAUNCH50"
    name: "Launch discount"
    percent_off: 50
    duration: "once"

promotion_codes:
  - code: "EARLY_ACCESS"
    coupon: "LAUNCH50"
    max_redemptions: 100
    active: true

tax_rates:
  - display_name: "CA Sales Tax"
    percentage: 8.625
    country: "US"
    state: "CA"
    tax_type: "sales_tax"
    inclusive: false
    active: true

shipping_rates:
  - display_name: "Free Shipping"
    amount: 0
    currency: "usd"
    active: true

billing_meters:
  - event_name: "api_calls"
    display_name: "API Calls"
    aggregation: "count"

entitlements:
  - name: "advanced_analytics"
    lookup_key: "advanced_analytics"

apple_pay_domains:
  - "app.example.com"
```

## Resource Management

stripe-cm uses metadata tags (`retier_managed: "true"`) to track which Stripe resources it manages. Resources created outside of stripe-cm are never modified or deleted.

| Resource | Unique Key | Sync | Import |
|----------|-----------|------|--------|
| Business profile | (singleton) | Update | Read |
| Tax settings | (singleton) | Update | Read |
| Portal configs | `name` metadata | Create/Update/Deactivate | List managed |
| Webhooks | `url` | Create/Update/Delete | List managed |
| Coupons | `id` | Create/Update/Delete | List managed |
| Promotion codes | `code` | Create/Deactivate | List managed |
| Tax rates | `display_name+country+state` | Create/Deactivate | List managed |
| Tax registrations | `country+state` | Create/Expire | List active |
| Shipping rates | `display_name` | Create/Update/Deactivate | List managed |
| Billing meters | `event_name` | Create/Update/Deactivate | List active |
| Entitlements | `lookup_key` | Create/Update/Archive | List managed |
| Apple Pay domains | `domain_name` | Create/Delete | List all |

## Development

```bash
# Clone
git clone https://github.com/trescale/stripe-cm.git
cd stripe-cm

# Build
go build ./cmd/stripe-cm/

# Test
go test ./... -v

# Lint
golangci-lint run
```

## TODO

- [ ] Add `stripe-cm diff` command to show what would change without applying
- [ ] Add `stripe-cm validate` command to check YAML syntax and Stripe API key
- [ ] Support multiple environments (staging/production) via profiles
- [ ] Add Homebrew tap for easy installation
- [ ] Add `--dry-run` flag to sync command
- [ ] Support Stripe Connect (manage connected account configs)
- [ ] Add payment method configurations sync
- [ ] Add payment method domains sync
- [ ] Add invoice rendering templates sync
- [ ] Add Radar rules management
- [ ] Add product + price catalog sync (optional, for standalone use)
- [ ] Add `stripe-cm plan` command (Terraform-style plan before apply)
- [ ] Support YAML anchors and includes for DRY configs
- [ ] Add JSON output format option
- [ ] Add Terraform provider wrapper
- [ ] Publish to Homebrew
- [ ] Add shell completions (bash/zsh/fish)
- [ ] State file for tracking managed resources (alternative to metadata tags)

## License

MIT
