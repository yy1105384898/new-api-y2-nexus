/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/
import {
  Globe,
  Layers,
  Shield,
  Zap,
} from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { cn } from '@/lib/utils'
import { AnimateInView } from '@/components/animate-in-view'
import { mkt } from '../../lib/marketing-theme'
import { SITE_BRAND } from '../../lib/site-assets'

export function WhyUs() {
  const { t } = useTranslation()

  const items = [
    {
      icon: Globe,
      title: t('No barriers, no restrictions'),
      desc: t(
        'Access Claude, GPT and more from anywhere. Email signup only, flexible payment, privacy protected.'
      ),
      tags: [
        t('Email signup only'),
        t('Flexible payment'),
        t('Privacy protected'),
      ],
    },
    {
      icon: Zap,
      title: t('Full-spec models, uncut'),
      desc: t(
        'Official model capabilities — same context window, same output quality. No tweaks, no trimming.'
      ),
      tags: [],
    },
    {
      icon: Layers,
      title: t('All models, one account'),
      desc: t(
        'Claude Opus, Sonnet, GPT, Gemini — switch freely between models. No separate subscriptions.'
      ),
      tags: [],
    },
    {
      icon: Shield,
      title: t('Standard API compatible'),
      desc: t(
        'Drop-in replacement for Anthropic & OpenAI APIs. Change one URL in your code and start saving.'
      ),
      tags: [],
    },
  ]

  return (
    <section className='relative z-10 px-6 py-20 md:py-28'>
      <div className='mx-auto max-w-6xl'>
        <AnimateInView className='mb-12 text-center md:mb-16'>
          <h2 className={cn('text-2xl font-bold tracking-tight md:text-3xl', mkt.heading)}>
            {t('Why {{brand}}', { brand: SITE_BRAND.name })}
          </h2>
          <p className={cn('mx-auto mt-4 max-w-2xl text-sm leading-relaxed md:text-base', mkt.muted)}>
            {t(
              'The simplest way to use frontier AI models — no barriers, no subscriptions, no restrictions.'
            )}
          </p>
        </AnimateInView>

        <div className='grid gap-6 md:grid-cols-2 md:gap-8'>
          {items.map((item, i) => (
            <AnimateInView
              key={item.title}
              delay={i * 80}
              animation='fade-up'
              className={cn(mkt.card, 'p-6 md:p-8')}
            >
              <div className={cn(mkt.cardIcon, 'mb-4 size-11', mkt.iconAccent)}>
                <item.icon className='size-5' strokeWidth={1.5} />
              </div>
              <h3 className={cn('mb-2 text-base font-semibold tracking-tight md:text-lg', mkt.heading)}>
                {item.title}
              </h3>
              <p className={cn('text-sm leading-relaxed', mkt.muted)}>
                {item.desc}
              </p>
              {item.tags.length > 0 && (
                <div className='mt-4 flex flex-wrap gap-2'>
                  {item.tags.map((tag) => (
                    <span
                      key={tag}
                      className='rounded-full border border-slate-200/70 bg-slate-50/80 px-2.5 py-1 text-xs font-medium text-slate-600 dark:border-white/10 dark:bg-white/5 dark:text-[#98b4c1]'
                    >
                      {tag}
                    </span>
                  ))}
                </div>
              )}
            </AnimateInView>
          ))}
        </div>
      </div>
    </section>
  )
}
