package stripe

type Config struct {
	Business        BusinessConfig       `yaml:"business"`
	Tax             TaxConfig            `yaml:"tax"`
	Billing         BillingConfig        `yaml:"billing"`
	CustomerPortal  CustomerPortalConfig `yaml:"customer_portal"`
	Checkout        CheckoutConfig       `yaml:"checkout"`
	Webhooks        []WebhookConfig      `yaml:"webhooks,omitempty"`
	Coupons         []CouponConfig       `yaml:"coupons,omitempty"`
	PromotionCodes  []PromotionConfig    `yaml:"promotion_codes,omitempty"`
	TaxRates        []TaxRateConfig      `yaml:"tax_rates,omitempty"`
	TaxRegs         []TaxRegConfig       `yaml:"tax_registrations,omitempty"`
	ShippingRates   []ShippingRateConfig `yaml:"shipping_rates,omitempty"`
	BillingMeters   []BillingMeterConfig `yaml:"billing_meters,omitempty"`
	Entitlements    []EntitlementConfig  `yaml:"entitlements,omitempty"`
	ApplePayDomains []string             `yaml:"apple_pay_domains,omitempty"`
}

type BusinessConfig struct {
	Name         string `yaml:"name"`
	SupportEmail string `yaml:"support_email"`
	SupportURL   string `yaml:"support_url"`
	SupportPhone string `yaml:"support_phone"`
}

type TaxConfig struct {
	Automatic      bool   `yaml:"automatic"`
	DefaultTaxCode string `yaml:"default_tax_code"`
}

type BillingConfig struct {
	Subscriptions SubscriptionConfig `yaml:"subscriptions"`
	Invoices      InvoiceConfig      `yaml:"invoices"`
}

type SubscriptionConfig struct {
	DefaultPaymentMethod string `yaml:"default_payment_method"`
	ProrationBehavior    string `yaml:"proration_behavior"`
	MissingPaymentMethod string `yaml:"missing_payment_method"`
}

type InvoiceConfig struct {
	DaysUntilDue int    `yaml:"days_until_due"`
	Footer       string `yaml:"footer"`
}

type CustomerPortalConfig struct {
	Profiles []PortalProfile `yaml:"profiles"`
}

type PortalProfile struct {
	Name             string         `yaml:"name"`
	DefaultReturnURL string         `yaml:"default_return_url"`
	Features         PortalFeatures `yaml:"features"`
}

type PortalFeatures struct {
	SubscriptionCancel  PortalSubscriptionCancel `yaml:"subscription_cancel"`
	SubscriptionUpdate  PortalSubscriptionUpdate `yaml:"subscription_update"`
	PaymentMethodUpdate PortalToggle             `yaml:"payment_method_update"`
	InvoiceHistory      PortalToggle             `yaml:"invoice_history"`
}

type PortalToggle struct {
	Enabled bool `yaml:"enabled"`
}

type PortalSubscriptionCancel struct {
	Enabled           bool   `yaml:"enabled"`
	Mode              string `yaml:"mode,omitempty"`
	ProrationBehavior string `yaml:"proration_behavior,omitempty"`
}

type PortalSubscriptionUpdate struct {
	Enabled           bool   `yaml:"enabled"`
	ProrationBehavior string `yaml:"proration_behavior,omitempty"`
}

type CheckoutConfig struct {
	PaymentMethods      []string `yaml:"payment_methods"`
	AllowPromotionCodes bool     `yaml:"allow_promotion_codes"`
	CollectPhoneNumber  bool     `yaml:"collect_phone_number"`
	TaxIDCollection     bool     `yaml:"tax_id_collection"`
}

type WebhookConfig struct {
	URL           string   `yaml:"url"`
	Description   string   `yaml:"description,omitempty"`
	EnabledEvents []string `yaml:"enabled_events"`
	APIVersion    string   `yaml:"api_version,omitempty"`
	Disabled      bool     `yaml:"disabled,omitempty"`
}

type CouponConfig struct {
	ID               string  `yaml:"id"`
	Name             string  `yaml:"name,omitempty"`
	PercentOff       float64 `yaml:"percent_off,omitempty"`
	AmountOff        int64   `yaml:"amount_off,omitempty"`
	Currency         string  `yaml:"currency,omitempty"`
	Duration         string  `yaml:"duration"`
	DurationInMonths int64   `yaml:"duration_in_months,omitempty"`
	MaxRedemptions   int64   `yaml:"max_redemptions,omitempty"`
}

type PromotionConfig struct {
	Code           string `yaml:"code"`
	Coupon         string `yaml:"coupon"`
	MaxRedemptions int64  `yaml:"max_redemptions,omitempty"`
	Active         bool   `yaml:"active"`
}

type TaxRateConfig struct {
	DisplayName string  `yaml:"display_name"`
	Description string  `yaml:"description,omitempty"`
	Percentage  float64 `yaml:"percentage"`
	Country     string  `yaml:"country,omitempty"`
	State       string  `yaml:"state,omitempty"`
	TaxType     string  `yaml:"tax_type,omitempty"`
	Inclusive   bool    `yaml:"inclusive"`
	Active      bool    `yaml:"active"`
}

type TaxRegConfig struct {
	Country string `yaml:"country"`
	State   string `yaml:"state,omitempty"`
	Type    string `yaml:"type"`
}

type ShippingRateConfig struct {
	DisplayName string `yaml:"display_name"`
	Amount      int64  `yaml:"amount"`
	Currency    string `yaml:"currency"`
	TaxBehavior string `yaml:"tax_behavior,omitempty"`
	Active      bool   `yaml:"active"`
}

type BillingMeterConfig struct {
	EventName       string `yaml:"event_name"`
	DisplayName     string `yaml:"display_name"`
	Aggregation     string `yaml:"aggregation"`
	EventTimeWindow string `yaml:"event_time_window,omitempty"`
	CustomerMapping string `yaml:"customer_mapping,omitempty"`
	ValueKey        string `yaml:"value_key,omitempty"`
}

type EntitlementConfig struct {
	Name   string `yaml:"name"`
	Lookup string `yaml:"lookup_key,omitempty"`
}
