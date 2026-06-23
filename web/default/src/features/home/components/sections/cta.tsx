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
import { Button } from '@/components/ui/button'
import { AnimateInView } from '@/components/animate-in-view'

interface CTAProps {
  className?: string
  isAuthenticated?: boolean
}

export function CTA(props: CTAProps) {
  const { t } = useTranslation()

  if (props.isAuthenticated) {
    return null
  }

  return (
    <section className='relative z-10 overflow-hidden px-6 py-20 md:py-28'>
      <AnimateInView
        className='mx-auto max-w-2xl text-center'
        animation='scale-in'
      >
        <Button
          className='group h-12 rounded-lg px-8 text-sm font-medium'
          render={<Link to='/sign-up' />}
        >
          {t("Start now — it's free to sign up")}
          <ArrowRight className='ml-1.5 size-4 transition-transform duration-200 group-hover:translate-x-0.5' />
        </Button>
      </AnimateInView>
    </section>
  )
}
