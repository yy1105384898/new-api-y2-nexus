/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
import { formatTimestampToDate } from '@/lib/format'
import { formatCurrencyFromUSD, formatLocalCurrencyAmount } from '@/lib/currency'
import type { StatusBadgeProps } from '@/components/status-badge'
import type { TopupRecord, TopupStatus } from '../types'

// ============================================================================
// Billing Utility Functions
// ============================================================================

interface StatusConfig {
  variant: StatusBadgeProps['variant']
  label: string
}

/**
 * Status badge configuration
 */
export const STATUS_CONFIG: Record<TopupStatus, StatusConfig> = {
  success: {
    variant: 'success',
    label: 'Success',
  },
  pending: {
    variant: 'warning',
    label: 'Pending',
  },
  expired: {
    variant: 'danger',
    label: 'Expired',
  },
}

/**
 * Get status badge configuration
 */
export function getStatusConfig(status: TopupStatus): StatusConfig {
  return STATUS_CONFIG[status] || STATUS_CONFIG.pending
}

/**
 * Payment method display names
 */
export const PAYMENT_METHOD_NAMES: Record<string, string> = {
  stripe: 'Stripe',
  alipay: 'Alipay',
  wxpay: 'WeChat Pay',
  waffo: 'Waffo',
}

/**
 * Get payment method display name
 */
export function getPaymentMethodName(
  method: string,
  t?: (key: string) => string
): string {
  const name = PAYMENT_METHOD_NAMES[method] || method
  return t ? t(name) : name
}

/**
 * Format timestamp to readable date string
 */
export function formatTimestamp(timestamp: number): string {
  return formatTimestampToDate(timestamp)
}

/** 新易支付订单 Amount 存人民币面值；旧订单存 USD。 */
export function isEpayCnyCredit(
  record: Pick<TopupRecord, 'amount' | 'money' | 'payment_provider'>
): boolean {
  if (record.payment_provider !== 'epay') return false
  if (record.amount >= record.money * 0.5 && record.amount >= 1) return true
  return false
}

export function formatTopupCreditDisplay(record: TopupRecord): string {
  if (isEpayCnyCredit(record)) {
    return formatLocalCurrencyAmount(record.amount, {
      digitsLarge: 2,
      digitsSmall: 2,
      abbreviate: false,
    })
  }
  return formatCurrencyFromUSD(record.amount, {
    digitsLarge: 2,
    digitsSmall: 2,
    abbreviate: false,
  })
}

export function formatTopupPaidDisplay(record: TopupRecord): string {
  if (record.payment_provider === 'epay' || isEpayCnyCredit(record)) {
    return formatLocalCurrencyAmount(record.money, {
      digitsLarge: 2,
      digitsSmall: 2,
      abbreviate: false,
    })
  }
  return formatCurrencyFromUSD(record.money, {
    digitsLarge: 2,
    digitsSmall: 2,
    abbreviate: false,
  })
}
