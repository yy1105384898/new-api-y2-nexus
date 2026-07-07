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
import { useTranslation } from 'react-i18next'
import { CopyButton } from '@/components/copy-button'
import { useStatus } from '@/hooks/use-status'
import { isEuDeployment } from '@/i18n/region'

const DEFAULT_API_ENDPOINT = 'https://yynewapi.yangyangnj.top/v1'

function normalizeApiEndpoint(value?: string) {
  const endpoint = value?.trim().replace(/\/$/, '')
  if (!endpoint) return ''
  return endpoint.endsWith('/v1') ? endpoint : `${endpoint}/v1`
}

function EndpointBlock(props: {
  url: string
  audience: string
  description: string
  whenToUse: string
  copyLabel: string
  recommended?: boolean
}) {
  const { t } = useTranslation()

  return (
    <li className='border-border bg-muted/20 space-y-1.5 rounded-lg border px-3 py-2.5'>
      <div className='flex flex-wrap items-center gap-2'>
        <p className='text-foreground/90 text-xs font-medium sm:text-sm'>
          {props.audience}
        </p>
        {props.recommended ? (
          <span className='bg-primary/10 text-primary rounded px-1.5 py-0.5 text-[10px] font-medium sm:text-[11px]'>
            {t('Recommended')}
          </span>
        ) : null}
      </div>
      <p className='text-muted-foreground text-[11px] leading-relaxed sm:text-xs'>
        {props.description}
      </p>
      <p className='text-muted-foreground/80 text-[11px] leading-relaxed sm:text-xs'>
        <span className='text-foreground/70 font-medium'>
          {t('Choose this when')}
        </span>
        {props.whenToUse}
      </p>
      <div className='flex flex-wrap items-center gap-x-1.5 gap-y-1 pt-0.5'>
        <code className='bg-muted rounded px-1 py-0.5 text-[11px] break-all sm:text-xs'>
          {props.url}
        </code>
        <CopyButton
          value={props.url}
          size='icon'
          variant='ghost'
          className='size-6 shrink-0'
          iconClassName='size-3'
          tooltip={props.copyLabel}
          aria-label={props.copyLabel}
        />
      </div>
    </li>
  )
}

export function ApiEndpointHints() {
  const { t } = useTranslation()
  const { status } = useStatus()
  const euDeployment = isEuDeployment()

  const apiEndpoint = (() => {
    const fromStatus = normalizeApiEndpoint(
      typeof status?.server_address === 'string'
        ? status.server_address
        : undefined
    )
    return fromStatus || DEFAULT_API_ENDPOINT
  })()

  return (
    <div className='text-muted-foreground w-full space-y-3 text-xs sm:text-sm'>
      <p className='text-muted-foreground text-[11px] leading-relaxed sm:text-xs'>
        {t(
          'All endpoints below accept the same API key. Pick the one that matches where you deploy and how much traffic you send — details for each option follow.'
        )}
      </p>

      <div className='space-y-2'>
        <p className='text-foreground/80 font-medium'>
          {euDeployment ? t('European API access') : t('API base URL (OpenAI-compatible)')}
        </p>
        <ul className='space-y-2'>
          <EndpointBlock
            url={apiEndpoint}
            audience={t('Primary API endpoint')}
            description={t(
              'OpenAI-compatible HTTPS endpoint for this site. Use it as the base URL in API clients and integrations.'
            )}
            whenToUse={t(
              'You are using an API key from this console — this is the recommended base URL for this deployment.'
            )}
            copyLabel={t('Copy API base URL')}
            recommended
          />
        </ul>
      </div>
    </div>
  )
}
