package model

import (
	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/setting/operation_setting"

	"github.com/shopspring/decimal"
)

// EpayDisplayExchangeRate returns the CNY display rate used for 1:1 epay top-ups.
func EpayDisplayExchangeRate() float64 {
	rate := operation_setting.USDExchangeRate
	if rate <= 0 {
		rate = 1
	}
	return rate
}

// EpayQuotaFromCny credits quota so that CNY display increases by cnyAmount (1:1).
func EpayQuotaFromCny(cnyAmount int64) int {
	if cnyAmount <= 0 {
		return 0
	}
	rate := EpayDisplayExchangeRate()
	return int(decimal.NewFromInt(cnyAmount).
		Div(decimal.NewFromFloat(rate)).
		Mul(decimal.NewFromFloat(common.QuotaPerUnit)).
		Round(0).
		IntPart())
}

// IsLegacyEpayUsdTopUpAmount detects older epay orders that stored USD in Amount.
func IsLegacyEpayUsdTopUpAmount(topUp *TopUp) bool {
	if topUp == nil || topUp.PaymentProvider != PaymentProviderEpay {
		return false
	}
	if topUp.Amount <= 0 || topUp.Money <= 0 {
		return false
	}
	if topUp.Amount >= 10 && float64(topUp.Amount) >= topUp.Money*0.5 {
		return false
	}
	rate := EpayDisplayExchangeRate()
	return float64(topUp.Amount)*rate < topUp.Money*0.75
}

// EpayTopUpQuota resolves quota to add for an epay order (new CNY face or legacy USD).
func EpayTopUpQuota(topUp *TopUp) int {
	if topUp == nil {
		return 0
	}
	if topUp.PaymentProvider == PaymentProviderEpay && !IsLegacyEpayUsdTopUpAmount(topUp) {
		return EpayQuotaFromCny(topUp.Amount)
	}
	return int(decimal.NewFromInt(topUp.Amount).
		Mul(decimal.NewFromFloat(common.QuotaPerUnit)).
		Round(0).
		IntPart())
}
