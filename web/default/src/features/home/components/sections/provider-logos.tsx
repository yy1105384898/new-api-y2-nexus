/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/
import { useTranslation } from 'react-i18next'
import { cn } from '@/lib/utils'
import { AnimateInView } from '@/components/animate-in-view'
import { getLobeIcon } from '@/lib/lobe-icon'
import { mkt } from '../../lib/marketing-theme'
import { PROVIDER_ICONS } from '../../constants'

const COMPACT_PROVIDERS = PROVIDER_ICONS.slice(0, 12)

export function ProviderLogos() {
  const { t } = useTranslation()

  return (
    <section className={cn('relative z-10 px-6 py-10 md:py-12', mkt.sectionBorder)}>
      <div className='mx-auto max-w-4xl'>
        <AnimateInView
          animation='fade-up'
          className='flex flex-wrap items-center justify-center gap-5 sm:gap-6'
        >
          {COMPACT_PROVIDERS.map((iconName) => (
            <div
              key={iconName}
              className='flex size-9 items-center justify-center opacity-70 transition-opacity hover:opacity-100 sm:size-10'
            >
              {getLobeIcon(iconName, 32)}
            </div>
          ))}
          <div className={cn('text-sm font-medium', mkt.muted)}>30+</div>
        </AnimateInView>
        <p className={cn('mt-4 text-center text-xs', mkt.muted)}>
          {t('Supports a wide range of model providers')}
        </p>
      </div>
    </section>
  )
}
