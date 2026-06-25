/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/
import type { ReactNode } from 'react'
import { Link } from '@tanstack/react-router'
import { ArrowRight } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { PublicLayout } from '@/components/layout'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'

export type DocsMode = 'user' | 'api'

type DocsShellProps = {
  mode: DocsMode
  eyebrow: string
  title: string
  subtitle: string
  sidebarLabel: string
  nav: ReactNode
  children: ReactNode
}

export function DocsShell(props: DocsShellProps) {
  const { t } = useTranslation()

  return (
    <PublicLayout>
      <div className='mx-auto flex max-w-6xl flex-col gap-10 pb-16 lg:flex-row lg:gap-12'>
        <aside className='lg:w-56 lg:shrink-0'>
          <div className='lg:sticky lg:top-24 space-y-6'>
            <div className='bg-muted/40 border-border/50 flex rounded-xl border p-1'>
              <DocsTabLink to='/docs' active={props.mode === 'user'}>
                {t('userDocs.tabUserGuide')}
              </DocsTabLink>
              <DocsTabLink to='/docs/api' active={props.mode === 'api'}>
                {t('apiDocs.tabApiReference')}
              </DocsTabLink>
            </div>
            <div>
              <p className='text-muted-foreground mb-3 text-xs font-semibold tracking-[0.16em] uppercase'>
                {props.sidebarLabel}
              </p>
              <nav className='flex flex-wrap gap-2 lg:flex-col lg:gap-0.5'>{props.nav}</nav>
            </div>
          </div>
        </aside>

        <div className='min-w-0 flex-1'>
          <header className='mb-10 space-y-4'>
            <p className='text-primary text-xs font-semibold tracking-[0.2em] uppercase'>{props.eyebrow}</p>
            <h1 className='text-3xl font-bold tracking-tight md:text-4xl'>{props.title}</h1>
            <p className='text-muted-foreground max-w-2xl text-base leading-relaxed'>{props.subtitle}</p>
            <div className='flex flex-wrap gap-3 pt-1'>
              <Button render={<Link to='/sign-up' search={{ redirect: undefined }} />}>
                {t('Get API Key')}
                <ArrowRight className='size-4' />
              </Button>
              <Button variant='outline' render={<Link to='/keys' />}>
                {t('userDocs.nav.apiKey')}
              </Button>
              <Button variant='outline' render={<Link to='/pricing' />}>
                {t('Model Square')}
              </Button>
            </div>
          </header>
          {props.children}
        </div>
      </div>
    </PublicLayout>
  )
}

function DocsTabLink(props: { to: string; active: boolean; children: ReactNode }) {
  return (
    <Link
      to={props.to}
      className={cn(
        'flex-1 rounded-lg px-3 py-2 text-center text-sm font-medium transition-colors',
        props.active ? 'bg-background text-foreground shadow-sm' : 'text-muted-foreground hover:text-foreground'
      )}
    >
      {props.children}
    </Link>
  )
}

export function DocsNavLink(props: { href: string; children: ReactNode }) {
  return (
    <a
      href={props.href}
      className={cn(
        'text-muted-foreground hover:text-foreground rounded-lg px-3 py-2 text-sm transition-colors lg:block',
        'hover:bg-muted/60'
      )}
    >
      {props.children}
    </a>
  )
}
