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
import { docsNavItems } from './docs-nav'
import { DocsNavLink, DocsShell } from './docs-shell'
import { UserDocsSections } from './user-docs-sections'

export function UserDocsPage() {
  const { t } = useTranslation()
  const { systemName } = useSystemConfig()

  const siteOrigin = useMemo(() => {
    if (typeof window === 'undefined') return ''
    return window.location.origin
  }, [])

  const displayName = systemName?.trim() || '沧元算力'

  useEffect(() => {
    document.title = t('userDocs.pageTitle', { siteName: displayName })
  }, [displayName, t])

  return (
    <DocsShell
      mode='user'
      eyebrow={t('userDocs.eyebrow')}
      title={t('userDocs.title', { siteName: displayName })}
      subtitle={t('userDocs.subtitle')}
      sidebarLabel={t('userDocs.sidebarLabel')}
      nav={docsNavItems.map((item) => (
        <DocsNavLink key={item.id} href={`#${item.id}`}>
          {t(item.titleKey)}
        </DocsNavLink>
      ))}
    >
      <UserDocsSections siteOrigin={siteOrigin} siteName={displayName} />
    </DocsShell>
  )
}
