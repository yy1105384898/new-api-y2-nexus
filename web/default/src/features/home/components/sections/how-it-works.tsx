/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/
import { KeyRound, Rocket, UserPlus } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { cn } from '@/lib/utils'
import { AnimateInView } from '@/components/animate-in-view'
import { mkt } from '../../lib/marketing-theme'

export function HowItWorks() {
  const { t } = useTranslation()

  const steps = [
    {
      num: '1',
      title: t('Sign up'),
      desc: t('Email + verification code. No credit card, no KYC.'),
      icon: UserPlus,
    },
    {
      num: '2',
      title: t('Add balance'),
      desc: t('Top up your balance. Pay as little as you need.'),
      icon: KeyRound,
    },
    {
      num: '3',
      title: t('Start using'),
      desc: t('Use our API or compatible tools. Same as Claude & OpenAI.'),
      icon: Rocket,
    },
  ]

  return (
    <section
      id='get-started'
      className={cn('relative z-10 px-6 py-20 md:py-28', mkt.sectionBorder)}
    >
      <div className='mx-auto max-w-4xl'>
        <AnimateInView className='mb-12 text-center md:mb-16'>
          <h2 className={cn('text-2xl font-bold tracking-tight md:text-3xl', mkt.heading)}>
            {t('Get started in 30 seconds')}
          </h2>
          <p className={cn('mx-auto mt-4 max-w-lg text-sm leading-relaxed md:text-base', mkt.muted)}>
            {t(
              'Three simple steps. No credit card, no phone number, no hassle.'
            )}
          </p>
        </AnimateInView>

        <div className='grid gap-8 md:grid-cols-3 md:gap-6'>
          {steps.map((step, i) => (
            <AnimateInView
              key={step.num}
              delay={i * 100}
              animation='fade-up'
              className='flex flex-col items-center text-center'
            >
              <div className='relative mb-5'>
                <div className={cn(mkt.cardIcon, 'size-14', mkt.iconAccent)}>
                  <step.icon className='size-6' strokeWidth={1.5} />
                </div>
                <span
                  className={cn(
                    'absolute -top-2 -right-2 flex size-6 items-center justify-center rounded-full bg-cyan-600 text-xs font-bold text-white dark:bg-cyan-500',
                  )}
                >
                  {step.num}
                </span>
              </div>
              <h3 className={cn('mb-2 text-base font-semibold', mkt.heading)}>
                {step.title}
              </h3>
              <p className={cn('max-w-[220px] text-sm leading-relaxed', mkt.muted)}>
                {step.desc}
              </p>
            </AnimateInView>
          ))}
        </div>
      </div>
    </section>
  )
}
