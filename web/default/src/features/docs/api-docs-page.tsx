/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/
import { useEffect, useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { useSystemConfig } from '@/hooks/use-system-config'
import { apiDocsNavGroups } from './api-docs-nav'
import { ApiDocsSections } from './api-docs-sections'
import { DocsNavLink, DocsShell } from './docs-shell'

export function ApiDocsPage() {
  const { t } = useTranslation()
  const { systemName } = useSystemConfig()

  const siteOrigin = useMemo(() => {
    if (typeof window === 'undefined') return ''
    return window.location.origin
  }, [])

  const displayName = systemName?.trim() || '沧元算力'

  useEffect(() => {
    document.title = t('apiDocs.pageTitle', { siteName: displayName })
  }, [displayName, t])

  return (
    <DocsShell
      mode='api'
      eyebrow={t('apiDocs.eyebrow')}
      title={t('apiDocs.title', { siteName: displayName })}
      subtitle={t('apiDocs.subtitle')}
      sidebarLabel={t('apiDocs.sidebarLabel')}
      nav={
        <>
          {apiDocsNavGroups.map((group) => (
            <div key={group.titleKey} className='mb-4 last:mb-0'>
              <p className='text-muted-foreground/80 mb-1 hidden px-3 text-[11px] font-semibold tracking-wide uppercase lg:block'>
                {t(group.titleKey)}
              </p>
              {group.items.map((item) => (
                <DocsNavLink key={item.id} href={`#${item.id}`}>
                  {t(item.titleKey)}
                </DocsNavLink>
              ))}
            </div>
          ))}
        </>
      }
    >
      <ApiDocsSections siteOrigin={siteOrigin} />
    </DocsShell>
  )
}
