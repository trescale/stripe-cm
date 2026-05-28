package stripe

import (
	"log"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/paymentmethodconfiguration"
)

func toPreference(available string) string {
	switch available {
	case "active":
		return "on"
	case "inactive":
		return "off"
	default:
		return ""
	}
}

func fromPreference(pref string) string {
	switch pref {
	case "on":
		return "active"
	case "off":
		return "inactive"
	default:
		return ""
	}
}

func SyncPaymentMethodsConfig(cfg PaymentMethodsConfig) {
	existing := findDefaultPaymentMethodConfig()

	params := &stripe.PaymentMethodConfigurationParams{}
	count := 0

	if cfg.Card != nil {
		if p := toPreference(cfg.Card.Available); p != "" {
			params.Card = &stripe.PaymentMethodConfigurationCardParams{
				DisplayPreference: &stripe.PaymentMethodConfigurationCardDisplayPreferenceParams{
					Preference: stripe.String(p),
				},
			}
			count++
		}
	}
	if cfg.ApplePay != nil {
		if p := toPreference(cfg.ApplePay.Available); p != "" {
			params.ApplePay = &stripe.PaymentMethodConfigurationApplePayParams{
				DisplayPreference: &stripe.PaymentMethodConfigurationApplePayDisplayPreferenceParams{
					Preference: stripe.String(p),
				},
			}
			count++
		}
	}
	if cfg.GooglePay != nil {
		if p := toPreference(cfg.GooglePay.Available); p != "" {
			params.GooglePay = &stripe.PaymentMethodConfigurationGooglePayParams{
				DisplayPreference: &stripe.PaymentMethodConfigurationGooglePayDisplayPreferenceParams{
					Preference: stripe.String(p),
				},
			}
			count++
		}
	}
	if cfg.Link != nil {
		if p := toPreference(cfg.Link.Available); p != "" {
			params.Link = &stripe.PaymentMethodConfigurationLinkParams{
				DisplayPreference: &stripe.PaymentMethodConfigurationLinkDisplayPreferenceParams{
					Preference: stripe.String(p),
				},
			}
			count++
		}
	}
	if cfg.CashApp != nil {
		if p := toPreference(cfg.CashApp.Available); p != "" {
			params.CashApp = &stripe.PaymentMethodConfigurationCashAppParams{
				DisplayPreference: &stripe.PaymentMethodConfigurationCashAppDisplayPreferenceParams{
					Preference: stripe.String(p),
				},
			}
			count++
		}
	}
	if cfg.USBankAccount != nil {
		if p := toPreference(cfg.USBankAccount.Available); p != "" {
			params.USBankAccount = &stripe.PaymentMethodConfigurationUSBankAccountParams{
				DisplayPreference: &stripe.PaymentMethodConfigurationUSBankAccountDisplayPreferenceParams{
					Preference: stripe.String(p),
				},
			}
			count++
		}
	}
	if cfg.SEPADebit != nil {
		if p := toPreference(cfg.SEPADebit.Available); p != "" {
			params.SEPADebit = &stripe.PaymentMethodConfigurationSEPADebitParams{
				DisplayPreference: &stripe.PaymentMethodConfigurationSEPADebitDisplayPreferenceParams{
					Preference: stripe.String(p),
				},
			}
			count++
		}
	}
	if cfg.Klarna != nil {
		if p := toPreference(cfg.Klarna.Available); p != "" {
			params.Klarna = &stripe.PaymentMethodConfigurationKlarnaParams{
				DisplayPreference: &stripe.PaymentMethodConfigurationKlarnaDisplayPreferenceParams{
					Preference: stripe.String(p),
				},
			}
			count++
		}
	}
	if cfg.Affirm != nil {
		if p := toPreference(cfg.Affirm.Available); p != "" {
			params.Affirm = &stripe.PaymentMethodConfigurationAffirmParams{
				DisplayPreference: &stripe.PaymentMethodConfigurationAffirmDisplayPreferenceParams{
					Preference: stripe.String(p),
				},
			}
			count++
		}
	}
	if cfg.AfterpayClearpay != nil {
		if p := toPreference(cfg.AfterpayClearpay.Available); p != "" {
			params.AfterpayClearpay = &stripe.PaymentMethodConfigurationAfterpayClearpayParams{
				DisplayPreference: &stripe.PaymentMethodConfigurationAfterpayClearpayDisplayPreferenceParams{
					Preference: stripe.String(p),
				},
			}
			count++
		}
	}
	if cfg.AmazonPay != nil {
		if p := toPreference(cfg.AmazonPay.Available); p != "" {
			params.AmazonPay = &stripe.PaymentMethodConfigurationAmazonPayParams{
				DisplayPreference: &stripe.PaymentMethodConfigurationAmazonPayDisplayPreferenceParams{
					Preference: stripe.String(p),
				},
			}
			count++
		}
	}
	if cfg.Paypal != nil {
		if p := toPreference(cfg.Paypal.Available); p != "" {
			params.Paypal = &stripe.PaymentMethodConfigurationPaypalParams{
				DisplayPreference: &stripe.PaymentMethodConfigurationPaypalDisplayPreferenceParams{
					Preference: stripe.String(p),
				},
			}
			count++
		}
	}
	if cfg.IDEAL != nil {
		if p := toPreference(cfg.IDEAL.Available); p != "" {
			params.IDEAL = &stripe.PaymentMethodConfigurationIDEALParams{
				DisplayPreference: &stripe.PaymentMethodConfigurationIDEALDisplayPreferenceParams{
					Preference: stripe.String(p),
				},
			}
			count++
		}
	}
	if cfg.Bancontact != nil {
		if p := toPreference(cfg.Bancontact.Available); p != "" {
			params.Bancontact = &stripe.PaymentMethodConfigurationBancontactParams{
				DisplayPreference: &stripe.PaymentMethodConfigurationBancontactDisplayPreferenceParams{
					Preference: stripe.String(p),
				},
			}
			count++
		}
	}
	if cfg.EPS != nil {
		if p := toPreference(cfg.EPS.Available); p != "" {
			params.EPS = &stripe.PaymentMethodConfigurationEPSParams{
				DisplayPreference: &stripe.PaymentMethodConfigurationEPSDisplayPreferenceParams{
					Preference: stripe.String(p),
				},
			}
			count++
		}
	}
	if cfg.Giropay != nil {
		if p := toPreference(cfg.Giropay.Available); p != "" {
			params.Giropay = &stripe.PaymentMethodConfigurationGiropayParams{
				DisplayPreference: &stripe.PaymentMethodConfigurationGiropayDisplayPreferenceParams{
					Preference: stripe.String(p),
				},
			}
			count++
		}
	}
	if cfg.P24 != nil {
		if p := toPreference(cfg.P24.Available); p != "" {
			params.P24 = &stripe.PaymentMethodConfigurationP24Params{
				DisplayPreference: &stripe.PaymentMethodConfigurationP24DisplayPreferenceParams{
					Preference: stripe.String(p),
				},
			}
			count++
		}
	}
	if cfg.Sofort != nil {
		if p := toPreference(cfg.Sofort.Available); p != "" {
			params.Sofort = &stripe.PaymentMethodConfigurationSofortParams{
				DisplayPreference: &stripe.PaymentMethodConfigurationSofortDisplayPreferenceParams{
					Preference: stripe.String(p),
				},
			}
			count++
		}
	}
	if cfg.Boleto != nil {
		if p := toPreference(cfg.Boleto.Available); p != "" {
			params.Boleto = &stripe.PaymentMethodConfigurationBoletoParams{
				DisplayPreference: &stripe.PaymentMethodConfigurationBoletoDisplayPreferenceParams{
					Preference: stripe.String(p),
				},
			}
			count++
		}
	}
	if cfg.OXXO != nil {
		if p := toPreference(cfg.OXXO.Available); p != "" {
			params.OXXO = &stripe.PaymentMethodConfigurationOXXOParams{
				DisplayPreference: &stripe.PaymentMethodConfigurationOXXODisplayPreferenceParams{
					Preference: stripe.String(p),
				},
			}
			count++
		}
	}
	if cfg.Alipay != nil {
		if p := toPreference(cfg.Alipay.Available); p != "" {
			params.Alipay = &stripe.PaymentMethodConfigurationAlipayParams{
				DisplayPreference: &stripe.PaymentMethodConfigurationAlipayDisplayPreferenceParams{
					Preference: stripe.String(p),
				},
			}
			count++
		}
	}
	if cfg.GrabPay != nil {
		if p := toPreference(cfg.GrabPay.Available); p != "" {
			params.Grabpay = &stripe.PaymentMethodConfigurationGrabpayParams{
				DisplayPreference: &stripe.PaymentMethodConfigurationGrabpayDisplayPreferenceParams{
					Preference: stripe.String(p),
				},
			}
			count++
		}
	}
	if cfg.PayNow != nil {
		if p := toPreference(cfg.PayNow.Available); p != "" {
			params.PayNow = &stripe.PaymentMethodConfigurationPayNowParams{
				DisplayPreference: &stripe.PaymentMethodConfigurationPayNowDisplayPreferenceParams{
					Preference: stripe.String(p),
				},
			}
			count++
		}
	}
	if cfg.PromptPay != nil {
		if p := toPreference(cfg.PromptPay.Available); p != "" {
			params.PromptPay = &stripe.PaymentMethodConfigurationPromptPayParams{
				DisplayPreference: &stripe.PaymentMethodConfigurationPromptPayDisplayPreferenceParams{
					Preference: stripe.String(p),
				},
			}
			count++
		}
	}
	if cfg.WeChatPay != nil {
		if p := toPreference(cfg.WeChatPay.Available); p != "" {
			params.WeChatPay = &stripe.PaymentMethodConfigurationWeChatPayParams{
				DisplayPreference: &stripe.PaymentMethodConfigurationWeChatPayDisplayPreferenceParams{
					Preference: stripe.String(p),
				},
			}
			count++
		}
	}
	if cfg.BLIK != nil {
		if p := toPreference(cfg.BLIK.Available); p != "" {
			params.BLIK = &stripe.PaymentMethodConfigurationBLIKParams{
				DisplayPreference: &stripe.PaymentMethodConfigurationBLIKDisplayPreferenceParams{
					Preference: stripe.String(p),
				},
			}
			count++
		}
	}
	if cfg.ACSSDebit != nil {
		log.Printf("  payment_methods: acss_debit not supported for sync, skipping")
	}

	if count == 0 {
		log.Printf("  payment_methods: no methods configured, skipping")
		return
	}

	if existing != nil {
		_, err := paymentmethodconfiguration.Update(existing.ID, params)
		if err != nil {
			log.Printf("  payment_methods: update warning — %v", err)
			return
		}
		log.Printf("  payment_methods: updated config %s (%d method(s))", existing.ID, count)
	} else {
		created, err := paymentmethodconfiguration.New(params)
		if err != nil {
			log.Printf("  payment_methods: create warning — %v", err)
			return
		}
		log.Printf("  payment_methods: created config %s (%d method(s))", created.ID, count)
	}
}

func ImportPaymentMethodsConfig() PaymentMethodsConfig {
	result := PaymentMethodsConfig{}

	pmc := findDefaultPaymentMethodConfig()
	if pmc == nil {
		log.Printf("  payment_methods: no configuration found")
		return result
	}

	count := 0

	if pmc.Card != nil && pmc.Card.DisplayPreference != nil {
		if v := fromPreference(string(pmc.Card.DisplayPreference.Preference)); v != "" {
			result.Card = &PMToggle{Available: v}
			count++
		}
	}
	if pmc.ApplePay != nil && pmc.ApplePay.DisplayPreference != nil {
		if v := fromPreference(string(pmc.ApplePay.DisplayPreference.Preference)); v != "" {
			result.ApplePay = &PMToggle{Available: v}
			count++
		}
	}
	if pmc.GooglePay != nil && pmc.GooglePay.DisplayPreference != nil {
		if v := fromPreference(string(pmc.GooglePay.DisplayPreference.Preference)); v != "" {
			result.GooglePay = &PMToggle{Available: v}
			count++
		}
	}
	if pmc.Link != nil && pmc.Link.DisplayPreference != nil {
		if v := fromPreference(string(pmc.Link.DisplayPreference.Preference)); v != "" {
			result.Link = &PMToggle{Available: v}
			count++
		}
	}
	if pmc.CashApp != nil && pmc.CashApp.DisplayPreference != nil {
		if v := fromPreference(string(pmc.CashApp.DisplayPreference.Preference)); v != "" {
			result.CashApp = &PMToggle{Available: v}
			count++
		}
	}
	if pmc.USBankAccount != nil && pmc.USBankAccount.DisplayPreference != nil {
		if v := fromPreference(string(pmc.USBankAccount.DisplayPreference.Preference)); v != "" {
			result.USBankAccount = &PMToggle{Available: v}
			count++
		}
	}
	if pmc.SEPADebit != nil && pmc.SEPADebit.DisplayPreference != nil {
		if v := fromPreference(string(pmc.SEPADebit.DisplayPreference.Preference)); v != "" {
			result.SEPADebit = &PMToggle{Available: v}
			count++
		}
	}
	if pmc.Klarna != nil && pmc.Klarna.DisplayPreference != nil {
		if v := fromPreference(string(pmc.Klarna.DisplayPreference.Preference)); v != "" {
			result.Klarna = &PMToggle{Available: v}
			count++
		}
	}
	if pmc.Affirm != nil && pmc.Affirm.DisplayPreference != nil {
		if v := fromPreference(string(pmc.Affirm.DisplayPreference.Preference)); v != "" {
			result.Affirm = &PMToggle{Available: v}
			count++
		}
	}
	if pmc.AfterpayClearpay != nil && pmc.AfterpayClearpay.DisplayPreference != nil {
		if v := fromPreference(string(pmc.AfterpayClearpay.DisplayPreference.Preference)); v != "" {
			result.AfterpayClearpay = &PMToggle{Available: v}
			count++
		}
	}
	if pmc.AmazonPay != nil && pmc.AmazonPay.DisplayPreference != nil {
		if v := fromPreference(string(pmc.AmazonPay.DisplayPreference.Preference)); v != "" {
			result.AmazonPay = &PMToggle{Available: v}
			count++
		}
	}
	if pmc.Paypal != nil && pmc.Paypal.DisplayPreference != nil {
		if v := fromPreference(string(pmc.Paypal.DisplayPreference.Preference)); v != "" {
			result.Paypal = &PMToggle{Available: v}
			count++
		}
	}
	if pmc.IDEAL != nil && pmc.IDEAL.DisplayPreference != nil {
		if v := fromPreference(string(pmc.IDEAL.DisplayPreference.Preference)); v != "" {
			result.IDEAL = &PMToggle{Available: v}
			count++
		}
	}
	if pmc.Bancontact != nil && pmc.Bancontact.DisplayPreference != nil {
		if v := fromPreference(string(pmc.Bancontact.DisplayPreference.Preference)); v != "" {
			result.Bancontact = &PMToggle{Available: v}
			count++
		}
	}
	if pmc.EPS != nil && pmc.EPS.DisplayPreference != nil {
		if v := fromPreference(string(pmc.EPS.DisplayPreference.Preference)); v != "" {
			result.EPS = &PMToggle{Available: v}
			count++
		}
	}
	if pmc.Giropay != nil && pmc.Giropay.DisplayPreference != nil {
		if v := fromPreference(string(pmc.Giropay.DisplayPreference.Preference)); v != "" {
			result.Giropay = &PMToggle{Available: v}
			count++
		}
	}
	if pmc.P24 != nil && pmc.P24.DisplayPreference != nil {
		if v := fromPreference(string(pmc.P24.DisplayPreference.Preference)); v != "" {
			result.P24 = &PMToggle{Available: v}
			count++
		}
	}
	if pmc.Sofort != nil && pmc.Sofort.DisplayPreference != nil {
		if v := fromPreference(string(pmc.Sofort.DisplayPreference.Preference)); v != "" {
			result.Sofort = &PMToggle{Available: v}
			count++
		}
	}
	if pmc.Boleto != nil && pmc.Boleto.DisplayPreference != nil {
		if v := fromPreference(string(pmc.Boleto.DisplayPreference.Preference)); v != "" {
			result.Boleto = &PMToggle{Available: v}
			count++
		}
	}
	if pmc.OXXO != nil && pmc.OXXO.DisplayPreference != nil {
		if v := fromPreference(string(pmc.OXXO.DisplayPreference.Preference)); v != "" {
			result.OXXO = &PMToggle{Available: v}
			count++
		}
	}
	if pmc.Alipay != nil && pmc.Alipay.DisplayPreference != nil {
		if v := fromPreference(string(pmc.Alipay.DisplayPreference.Preference)); v != "" {
			result.Alipay = &PMToggle{Available: v}
			count++
		}
	}
	if pmc.Grabpay != nil && pmc.Grabpay.DisplayPreference != nil {
		if v := fromPreference(string(pmc.Grabpay.DisplayPreference.Preference)); v != "" {
			result.GrabPay = &PMToggle{Available: v}
			count++
		}
	}
	if pmc.PayNow != nil && pmc.PayNow.DisplayPreference != nil {
		if v := fromPreference(string(pmc.PayNow.DisplayPreference.Preference)); v != "" {
			result.PayNow = &PMToggle{Available: v}
			count++
		}
	}
	if pmc.PromptPay != nil && pmc.PromptPay.DisplayPreference != nil {
		if v := fromPreference(string(pmc.PromptPay.DisplayPreference.Preference)); v != "" {
			result.PromptPay = &PMToggle{Available: v}
			count++
		}
	}
	if pmc.WeChatPay != nil && pmc.WeChatPay.DisplayPreference != nil {
		if v := fromPreference(string(pmc.WeChatPay.DisplayPreference.Preference)); v != "" {
			result.WeChatPay = &PMToggle{Available: v}
			count++
		}
	}
	if pmc.BLIK != nil && pmc.BLIK.DisplayPreference != nil {
		if v := fromPreference(string(pmc.BLIK.DisplayPreference.Preference)); v != "" {
			result.BLIK = &PMToggle{Available: v}
			count++
		}
	}
	if pmc.ACSSDebit != nil && pmc.ACSSDebit.DisplayPreference != nil {
		if v := fromPreference(string(pmc.ACSSDebit.DisplayPreference.Preference)); v != "" {
			result.ACSSDebit = &PMToggle{Available: v}
			count++
		}
	}

	log.Printf("  payment_methods: imported %d method(s) from config %s", count, pmc.ID)
	return result
}

func findDefaultPaymentMethodConfig() *stripe.PaymentMethodConfiguration {
	iter := paymentmethodconfiguration.List(&stripe.PaymentMethodConfigurationListParams{})
	var first *stripe.PaymentMethodConfiguration
	for iter.Next() {
		pmc := iter.PaymentMethodConfiguration()
		if pmc.IsDefault {
			return pmc
		}
		if first == nil {
			first = pmc
		}
	}
	if err := iter.Err(); err != nil {
		log.Printf("  payment_methods: warning listing configs — %v", err)
		return nil
	}
	return first
}
