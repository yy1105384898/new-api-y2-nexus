/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/
import { Link } from '@tanstack/react-router'
import { ArrowRight } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { useMemo } from 'react'
import { cn } from '@/lib/utils'
import { Button } from '@/components/ui/button'
import { AnimateInView } from '@/components/animate-in-view'
import { usePricingData } from '@/features/pricing/hooks/use-pricing-data'
import { mkt } from '../../lib/marketing-theme'
import {
  buildHomePricingRows,
  formatUsdPerM,
} from '../../lib/pricing-preview'

const HOME_PRICING_LIMIT = 20
const HOME_PRICING_MAX_PER_PREFIX = 2

export function PricingPreview() {
  const { t } = useTranslation()
  const { models, isLoading } = usePricingData()

  const rows = useMemo(
    () =>
      buildHomePricingRows(models, {
        limit: HOME_PRICING_LIMIT,
        maxPerPrefix: HOME_PRICING_MAX_PER_PREFIX,
      }),
    [models]
  )

  return (
    <section
      id='pricing'
      className={cn('relative z-10 px-6 py-12 md:py-16', mkt.sectionBorder)}
    >
      <div className='mx-auto max-w-5xl'>
        <AnimateInView className='mb-8 text-center'>
          <h2 className={cn('text-lg font-semibold md:text-xl', mkt.heading)}>
            {t('Model Pricing')}
          </h2>
          <p className={cn('mt-1 text-sm', mkt.muted)}>
            {t('USD per 1M tokens')}
          </p>
          <p className={cn('mt-1 text-xs', mkt.muted)}>
            {t(
              'Showing up to {{count}} public models; at most {{max}} per model family. See the model marketplace for the full list.',
              {
                count: HOME_PRICING_LIMIT,
                max: HOME_PRICING_MAX_PER_PREFIX,
              }
            )}
          </p>
        </AnimateInView>

        <AnimateInView animation='fade-up'>
          <div className={cn('overflow-hidden rounded-2xl', mkt.card)}>
            <div className='overflow-x-auto'>
              <table className='w-full min-w-[640px] text-sm'>
                <thead>
                  <tr className={cn('border-b text-left text-xs', mkt.sectionBorder)}>
                    <th className={cn('px-5 py-3.5 font-medium', mkt.muted)}>
                      {t('Model')}
                    </th>
                    <th className={cn('px-4 py-3.5 font-medium', mkt.muted)}>
                      {t('Input')}
                    </th>
                    <th className={cn('px-4 py-3.5 font-medium', mkt.muted)}>
                      {t('Output')}
                    </th>
                    <th className={cn('hidden px-4 py-3.5 font-medium sm:table-cell', mkt.muted)}>
                      {t('Cache Read')}
                    </th>
                    <th className={cn('hidden px-4 py-3.5 font-medium md:table-cell', mkt.muted)}>
                      {t('Cache Write')}
                    </th>
                    <th className={cn('hidden px-4 py-3.5 font-medium lg:table-cell', mkt.muted)}>
                      {t('Official pricing')}
                    </th>
                    <th className={cn('px-4 py-3.5 text-right font-medium', mkt.muted)}>
                      {t('Save')}
                    </th>
                  </tr>
                </thead>
                <tbody>
                  {isLoading ? (
                    Array.from({ length: 8 }).map((_, i) => (
                      <tr key={i} className={cn('border-b last:border-0', mkt.sectionBorder)}>
                        {Array.from({ length: 7 }).map((__, j) => (
                          <td key={j} className='px-4 py-3.5'>
                            <div className='h-4 animate-pulse rounded bg-slate-200/80 dark:bg-white/10' />
                          </td>
                        ))}
                      </tr>
                    ))
                  ) : rows.length === 0 ? (
                    <tr>
                      <td
                        colSpan={7}
                        className={cn('px-5 py-8 text-center', mkt.muted)}
                      >
                        {t('Pricing data unavailable')}
                      </td>
                    </tr>
                  ) : (
                    rows.map((row) => (
                      <tr
                        key={row.modelName}
                        className={cn(
                          'border-b transition-colors last:border-0 hover:bg-slate-50/60 dark:hover:bg-white/[0.03]',
                          mkt.sectionBorder
                        )}
                      >
                        <td className={cn('px-5 py-3.5 font-medium', mkt.heading)}>
                          <span className='font-mono text-xs sm:text-sm'>
                            {row.display}
                          </span>
                        </td>
                        <td className={cn('px-4 py-3.5 tabular-nums', mkt.body)}>
                          {row.isRequestBased
                            ? row.requestPrice
                            : formatUsdPerM(row.input)}
                        </td>
                        <td className={cn('px-4 py-3.5 tabular-nums', mkt.body)}>
                          {row.isRequestBased ? '—' : formatUsdPerM(row.output)}
                        </td>
                        <td className={cn('hidden px-4 py-3.5 tabular-nums sm:table-cell', mkt.body)}>
                          {formatUsdPerM(row.cacheRead)}
                        </td>
                        <td className={cn('hidden px-4 py-3.5 tabular-nums md:table-cell', mkt.body)}>
                          {formatUsdPerM(row.cacheWrite)}
                        </td>
                        <td className={cn('hidden px-4 py-3.5 tabular-nums lg:table-cell', mkt.muted)}>
                          {row.officialInput != null && row.officialOutput != null
                            ? `${formatUsdPerM(row.officialInput)} / ${formatUsdPerM(row.officialOutput)}`
                            : '—'}
                        </td>
                        <td className='px-4 py-3.5 text-right'>
                          {row.savePercent != null ? (
                            <span className='inline-flex rounded-full bg-emerald-500/10 px-2 py-0.5 text-xs font-semibold text-emerald-600 dark:text-emerald-400'>
                              −{row.savePercent}%
                            </span>
                          ) : (
                            <span className={mkt.muted}>—</span>
                          )}
                        </td>
                      </tr>
                    ))
                  )}
                </tbody>
              </table>
            </div>
          </div>

          <p className={cn('mt-4 text-center text-xs leading-relaxed', mkt.muted)}>
            {t(
              'Fixed pricing — what you see is what you pay. No fluctuation with network load, no hidden coefficient.'
            )}
          </p>
          <p className={cn('mt-1 text-center text-xs', mkt.muted)}>
            {t('Save = discount vs. official Anthropic / OpenAI API pricing.')}
          </p>

          <div className='mt-6 flex justify-center'>
            <Button
              variant='outline'
              className={cn('group rounded-lg', mkt.btnGhost)}
              render={<Link to='/pricing' />}
            >
              {t('View all models')}
              <ArrowRight className='ml-1.5 size-3.5 transition-transform duration-200 group-hover:translate-x-0.5' />
            </Button>
          </div>
        </AnimateInView>
      </div>
    </section>
  )
}
