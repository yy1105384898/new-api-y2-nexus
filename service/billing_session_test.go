package service

import (
	"testing"

	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/stretchr/testify/assert"
)

func TestBillingSessionNeedsRefund_UnlimitedTokenWalletPreConsume(t *testing.T) {
	session := &BillingSession{
		relayInfo: &relaycommon.RelayInfo{
			TokenUnlimited: true,
		},
		funding:          &WalletFunding{userId: 1, consumed: 1_950_000},
		preConsumedQuota: 1_950_000,
		tokenConsumed:    0,
	}

	assert.True(t, session.NeedsRefund(), "wallet pre-consume with unlimited token should be refundable")
}

func TestBillingSessionNeedsRefund_NoPreConsume(t *testing.T) {
	session := &BillingSession{
		relayInfo: &relaycommon.RelayInfo{},
		funding:   &WalletFunding{userId: 1},
	}

	assert.False(t, session.NeedsRefund())
}

func TestBillingSessionNeedsRefund_AlreadySettled(t *testing.T) {
	session := &BillingSession{
		relayInfo:        &relaycommon.RelayInfo{},
		funding:          &WalletFunding{userId: 1, consumed: 100},
		preConsumedQuota: 100,
		settled:          true,
	}

	assert.False(t, session.NeedsRefund())
}
