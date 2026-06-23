/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/
import { Link } from '@tanstack/react-router'
import { cn } from '@/lib/utils'
import { ArrowRight } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Button } from '@/components/ui/button'
import { mkt } from '../../lib/marketing-theme'

interface HeroProps {
  className?: string
  isAuthenticated?: boolean
}

export function Hero(props: HeroProps) {
  const { t } = useTranslation()

  return (
    <section className='relative z-10 overflow-hidden px-6 pt-28 pb-10 md:pt-36 md:pb-14'>
      <div className='mx-auto flex max-w-3xl flex-col items-center text-center'>
        <p
          className={cn(
            'landing-animate-fade-up mb-5 text-[11px] font-semibold tracking-[0.22em] uppercase',
            mkt.eyebrow
          )}
          style={{ animationDelay: '0ms' }}
        >
          {t('Full-Spec AI · Within Reach')}
        </p>

        <h1
          className={cn(
            'landing-animate-fade-up text-[clamp(2rem,5vw,3.25rem)] leading-[1.12] font-bold tracking-tight',
            mkt.heading
          )}
          style={{ animationDelay: '60ms' }}
        >
          {t('Access frontier AI models')}
          <br />
          <span className='bg-gradient-to-r from-cyan-600 via-emerald-600 to-violet-600 bg-clip-text text-transparent dark:from-cyan-300 dark:via-emerald-300 dark:to-violet-300'>
            {t('at a fraction of the cost')}
          </span>
        </h1>

        <div
          className='landing-animate-fade-up mt-8'
          style={{ animationDelay: '180ms' }}
        >
          {props.isAuthenticated ? (
            <Button
              className='group h-12 rounded-lg px-8 text-sm font-medium'
              render={<Link to='/dashboard' />}
            >
              {t('Go to Dashboard')}
              <ArrowRight className='ml-1.5 size-4 transition-transform duration-200 group-hover:translate-x-0.5' />
            </Button>
          ) : (
            <Button
              className='group h-12 rounded-lg px-8 text-sm font-medium'
              render={<Link to='/sign-up' />}
            >
              {t('Get Started')}
              <ArrowRight className='ml-1.5 size-4 transition-transform duration-200 group-hover:translate-x-0.5' />
            </Button>
          )}
        </div>
      </div>
    </section>
  )
}
