package stripe

import (
	"log"
	"sort"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/coupon"
	"github.com/stripe/stripe-go/v81/promotioncode"
)

func SyncCoupons(coupons []CouponConfig) {
	existing := listManagedCoupons()

	desired := make(map[string]bool, len(coupons))
	for _, c := range coupons {
		desired[c.ID] = true

		if _, ok := existing[c.ID]; ok {
			params := &stripe.CouponParams{
				Name: stripe.String(c.Name),
				Metadata: map[string]string{
					"stripe_cm_managed": "true",
				},
			}
			_, err := coupon.Update(c.ID, params)
			if err != nil {
				log.Printf("  coupons [%s]: update warning — %v", c.ID, err)
				continue
			}
			log.Printf("  coupons [%s]: updated", c.ID)
		} else {
			params := &stripe.CouponParams{
				ID:       stripe.String(c.ID),
				Duration: stripe.String(c.Duration),
				Metadata: map[string]string{
					"stripe_cm_managed": "true",
				},
			}
			if c.Name != "" {
				params.Name = stripe.String(c.Name)
			}
			if c.PercentOff > 0 {
				params.PercentOff = stripe.Float64(c.PercentOff)
			}
			if c.AmountOff > 0 {
				params.AmountOff = stripe.Int64(c.AmountOff)
				if c.Currency != "" {
					params.Currency = stripe.String(c.Currency)
				}
			}
			if c.DurationInMonths > 0 {
				params.DurationInMonths = stripe.Int64(c.DurationInMonths)
			}
			if c.MaxRedemptions > 0 {
				params.MaxRedemptions = stripe.Int64(c.MaxRedemptions)
			}
			_, err := coupon.New(params)
			if err != nil {
				log.Printf("  coupons [%s]: create warning — %v", c.ID, err)
				continue
			}
			log.Printf("  coupons [%s]: created", c.ID)
		}
	}

	for id := range existing {
		if desired[id] {
			continue
		}
		c, err := coupon.Get(id, nil)
		if err != nil {
			log.Printf("  coupons [%s]: fetch warning — %v", id, err)
			continue
		}
		if c.TimesRedeemed > 0 {
			log.Printf("  coupons [%s]: skipped deletion (redeemed %d times, may be referenced by active subscriptions)", id, c.TimesRedeemed)
			continue
		}
		_, err = coupon.Del(id, nil)
		if err != nil {
			log.Printf("  coupons [%s]: delete warning — %v", id, err)
			continue
		}
		log.Printf("  coupons [%s]: deleted (never redeemed)", id)
	}
}

func SyncPromotionCodes(promos []PromotionConfig) {
	existing := listManagedPromotionCodes()

	desired := make(map[string]bool, len(promos))
	for _, p := range promos {
		desired[p.Code] = true

		if ex, ok := existing[p.Code]; ok {
			params := &stripe.PromotionCodeParams{
				Active: stripe.Bool(p.Active),
				Metadata: map[string]string{
					"stripe_cm_managed": "true",
				},
			}
			_, err := promotioncode.Update(ex.ID, params)
			if err != nil {
				log.Printf("  promotion_codes [%s]: update warning — %v", p.Code, err)
				continue
			}
			log.Printf("  promotion_codes [%s]: updated", p.Code)
		} else {
			params := &stripe.PromotionCodeParams{
				Code:   stripe.String(p.Code),
				Coupon: stripe.String(p.Coupon),
				Active: stripe.Bool(p.Active),
				Metadata: map[string]string{
					"stripe_cm_managed": "true",
				},
			}
			if p.MaxRedemptions > 0 {
				params.MaxRedemptions = stripe.Int64(p.MaxRedemptions)
			}
			_, err := promotioncode.New(params)
			if err != nil {
				log.Printf("  promotion_codes [%s]: create warning — %v", p.Code, err)
				continue
			}
			log.Printf("  promotion_codes [%s]: created", p.Code)
		}
	}

	for code, ex := range existing {
		if desired[code] {
			continue
		}
		params := &stripe.PromotionCodeParams{
			Active: stripe.Bool(false),
		}
		_, err := promotioncode.Update(ex.ID, params)
		if err != nil {
			log.Printf("  promotion_codes [%s]: deactivate warning — %v", code, err)
			continue
		}
		log.Printf("  promotion_codes [%s]: deactivated", code)
	}
}

func ImportCoupons() []CouponConfig {
	params := &stripe.CouponListParams{}
	params.Filters.AddFilter("limit", "", "100")

	var configs []CouponConfig
	iter := coupon.List(params)
	for iter.Next() {
		c := iter.Coupon()
		if c.Metadata["stripe_cm_managed"] != "true" {
			continue
		}
		cfg := CouponConfig{
			ID:               c.ID,
			Name:             c.Name,
			PercentOff:       c.PercentOff,
			AmountOff:        c.AmountOff,
			Currency:         string(c.Currency),
			Duration:         string(c.Duration),
			DurationInMonths: c.DurationInMonths,
			MaxRedemptions:   c.MaxRedemptions,
		}
		configs = append(configs, cfg)
	}
	if err := iter.Err(); err != nil {
		log.Printf("  coupons: warning — %v", err)
	}
	sort.Slice(configs, func(i, j int) bool {
		return configs[i].ID < configs[j].ID
	})
	log.Printf("  coupons: imported %d managed coupon(s)", len(configs))
	return configs
}

func ImportPromotionCodes() []PromotionConfig {
	params := &stripe.PromotionCodeListParams{}
	params.Filters.AddFilter("limit", "", "100")

	var configs []PromotionConfig
	iter := promotioncode.List(params)
	for iter.Next() {
		pc := iter.PromotionCode()
		if pc.Metadata["stripe_cm_managed"] != "true" {
			continue
		}
		couponID := ""
		if pc.Coupon != nil {
			couponID = pc.Coupon.ID
		}
		cfg := PromotionConfig{
			Code:           pc.Code,
			Coupon:         couponID,
			MaxRedemptions: pc.MaxRedemptions,
			Active:         pc.Active,
		}
		configs = append(configs, cfg)
	}
	if err := iter.Err(); err != nil {
		log.Printf("  promotion_codes: warning — %v", err)
	}
	sort.Slice(configs, func(i, j int) bool {
		return configs[i].Code < configs[j].Code
	})
	log.Printf("  promotion_codes: imported %d managed promotion code(s)", len(configs))
	return configs
}

func listManagedCoupons() map[string]bool {
	result := make(map[string]bool)
	params := &stripe.CouponListParams{}
	params.Filters.AddFilter("limit", "", "100")
	iter := coupon.List(params)
	for iter.Next() {
		c := iter.Coupon()
		if c.Metadata["stripe_cm_managed"] == "true" {
			result[c.ID] = true
		}
	}
	if err := iter.Err(); err != nil {
		log.Printf("  coupons: warning listing existing — %v", err)
	}
	return result
}

func listManagedPromotionCodes() map[string]*stripe.PromotionCode {
	result := make(map[string]*stripe.PromotionCode)
	params := &stripe.PromotionCodeListParams{}
	params.Filters.AddFilter("limit", "", "100")
	iter := promotioncode.List(params)
	for iter.Next() {
		pc := iter.PromotionCode()
		if pc.Metadata["stripe_cm_managed"] == "true" {
			result[pc.Code] = pc
		}
	}
	if err := iter.Err(); err != nil {
		log.Printf("  promotion_codes: warning listing existing — %v", err)
	}
	return result
}
