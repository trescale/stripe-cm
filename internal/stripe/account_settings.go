package stripe

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/account"
	"github.com/stripe/stripe-go/v81/file"
)

func uploadFile(path string, purpose stripe.FilePurpose) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()
	params := &stripe.FileParams{
		Purpose:    stripe.String(string(purpose)),
		Filename:   stripe.String(filepath.Base(path)),
		FileReader: f,
	}
	result, err := file.New(params)
	if err != nil {
		return "", err
	}
	return result.ID, nil
}

func SyncBranding(cfg BrandingConfig) {
	if cfg.IconFile != "" {
		id, err := uploadFile(cfg.IconFile, stripe.FilePurposeBusinessIcon)
		if err != nil {
			log.Printf("  branding: icon upload warning — %v", err)
		} else {
			cfg.Icon = id
			log.Printf("  branding: uploaded icon %s → %s", cfg.IconFile, id)
		}
	}
	if cfg.LogoFile != "" {
		id, err := uploadFile(cfg.LogoFile, stripe.FilePurposeBusinessLogo)
		if err != nil {
			log.Printf("  branding: logo upload warning — %v", err)
		} else {
			cfg.Logo = id
			log.Printf("  branding: uploaded logo %s → %s", cfg.LogoFile, id)
		}
	}

	if cfg.Icon == "" && cfg.Logo == "" && cfg.PrimaryColor == "" && cfg.SecondaryColor == "" {
		log.Println("  branding: no changes")
		return
	}

	params := &stripe.AccountParams{
		Settings: &stripe.AccountSettingsParams{
			Branding: &stripe.AccountSettingsBrandingParams{},
		},
	}

	if cfg.Icon != "" {
		params.Settings.Branding.Icon = stripe.String(cfg.Icon)
		log.Printf("  branding: setting icon=%s", cfg.Icon)
	}
	if cfg.Logo != "" {
		params.Settings.Branding.Logo = stripe.String(cfg.Logo)
		log.Printf("  branding: setting logo=%s", cfg.Logo)
	}
	if cfg.PrimaryColor != "" {
		params.Settings.Branding.PrimaryColor = stripe.String(cfg.PrimaryColor)
		log.Printf("  branding: setting primary_color=%s", cfg.PrimaryColor)
	}
	if cfg.SecondaryColor != "" {
		params.Settings.Branding.SecondaryColor = stripe.String(cfg.SecondaryColor)
		log.Printf("  branding: setting secondary_color=%s", cfg.SecondaryColor)
	}

	acct, err := account.Get()
	if err != nil {
		log.Printf("  branding: warning — %v", err)
		return
	}
	_, err = account.Update(acct.ID, params)
	if err != nil {
		if strings.Contains(err.Error(), "live") || strings.Contains(err.Error(), "test") {
			log.Printf("  branding: warning — account settings require live keys (%v)", err)
			return
		}
		log.Printf("  branding: warning — %v", err)
		return
	}
	log.Println("  branding: updated")
}

func ImportBranding(merged BrandingConfig) BrandingConfig {
	acct, err := account.Get()
	if err != nil {
		log.Printf("  branding: warning — %v", err)
		return merged
	}
	bc := merged
	if acct.Settings != nil && acct.Settings.Branding != nil {
		b := acct.Settings.Branding
		if b.Icon != nil && b.Icon.ID != "" {
			bc.Icon = b.Icon.ID
		}
		if b.Logo != nil && b.Logo.ID != "" {
			bc.Logo = b.Logo.ID
		}
		if b.PrimaryColor != "" {
			bc.PrimaryColor = b.PrimaryColor
		}
		if b.SecondaryColor != "" {
			bc.SecondaryColor = b.SecondaryColor
		}
	}
	log.Printf("  branding: imported (icon=%q logo=%q primary=%q secondary=%q)",
		bc.Icon, bc.Logo, bc.PrimaryColor, bc.SecondaryColor)
	return bc
}

func SyncPayoutSchedule(cfg PayoutScheduleConfig) {
	if cfg.Interval == "" && cfg.MonthlyAnchor == 0 && cfg.WeeklyAnchor == "" && cfg.DelayDays == 0 {
		log.Println("  payout_schedule: no changes")
		return
	}

	params := &stripe.AccountParams{
		Settings: &stripe.AccountSettingsParams{
			Payouts: &stripe.AccountSettingsPayoutsParams{
				Schedule: &stripe.AccountSettingsPayoutsScheduleParams{},
			},
		},
	}

	if cfg.Interval != "" {
		params.Settings.Payouts.Schedule.Interval = stripe.String(cfg.Interval)
		log.Printf("  payout_schedule: setting interval=%s", cfg.Interval)
	}
	if cfg.MonthlyAnchor != 0 {
		params.Settings.Payouts.Schedule.MonthlyAnchor = stripe.Int64(cfg.MonthlyAnchor)
		log.Printf("  payout_schedule: setting monthly_anchor=%d", cfg.MonthlyAnchor)
	}
	if cfg.WeeklyAnchor != "" {
		params.Settings.Payouts.Schedule.WeeklyAnchor = stripe.String(cfg.WeeklyAnchor)
		log.Printf("  payout_schedule: setting weekly_anchor=%s", cfg.WeeklyAnchor)
	}
	if cfg.DelayDays != 0 {
		params.Settings.Payouts.Schedule.DelayDays = stripe.Int64(cfg.DelayDays)
		log.Printf("  payout_schedule: setting delay_days=%d", cfg.DelayDays)
	}

	acct, err := account.Get()
	if err != nil {
		log.Printf("  payout_schedule: warning — %v", err)
		return
	}
	_, err = account.Update(acct.ID, params)
	if err != nil {
		if strings.Contains(err.Error(), "live") || strings.Contains(err.Error(), "test") {
			log.Printf("  payout_schedule: warning — account settings require live keys (%v)", err)
			return
		}
		log.Printf("  payout_schedule: warning — %v", err)
		return
	}
	log.Println("  payout_schedule: updated")
}

func ImportPayoutSchedule(merged PayoutScheduleConfig) PayoutScheduleConfig {
	acct, err := account.Get()
	if err != nil {
		log.Printf("  payout_schedule: warning — %v", err)
		return merged
	}
	ps := merged
	if acct.Settings != nil && acct.Settings.Payouts != nil && acct.Settings.Payouts.Schedule != nil {
		s := acct.Settings.Payouts.Schedule
		if string(s.Interval) != "" {
			ps.Interval = string(s.Interval)
		}
		if s.MonthlyAnchor != 0 {
			ps.MonthlyAnchor = s.MonthlyAnchor
		}
		if string(s.WeeklyAnchor) != "" {
			ps.WeeklyAnchor = string(s.WeeklyAnchor)
		}
		if s.DelayDays != 0 {
			ps.DelayDays = s.DelayDays
		}
	}
	log.Printf("  payout_schedule: imported (interval=%q monthly_anchor=%d weekly_anchor=%q delay_days=%d)",
		ps.Interval, ps.MonthlyAnchor, ps.WeeklyAnchor, ps.DelayDays)
	return ps
}
