/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/
import { useEffect, useState } from 'react'
import { ExternalLink, Loader2 } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { StatusBadge } from '@/components/status-badge'
import { API_KEY_STATUSES } from '@/features/keys/constants'
import type { ApiKey } from '@/features/keys/types'
import {
  buildCanvasRedirectUrl,
  listSelectableCanvasApiKeys,
  openCanvasInNewTab,
} from '@/features/canvas/api'

type CanvasKeySelectFormProps = {
  canvasBaseUrl: string
  onSuccess?: () => void
  submitLabel?: string
  toolName?: string
}

export function CanvasKeySelectForm({
  canvasBaseUrl,
  onSuccess,
  submitLabel,
  toolName = '无限画布',
}: CanvasKeySelectFormProps) {
  const { t } = useTranslation()
  const [loadingKeys, setLoadingKeys] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [keys, setKeys] = useState<ApiKey[]>([])
  const [selectedId, setSelectedId] = useState<string>('')

  useEffect(() => {
    let cancelled = false
    setLoadingKeys(true)
    void listSelectableCanvasApiKeys()
      .then((items) => {
        if (cancelled) return
        setKeys(items)
        setSelectedId(items[0] ? String(items[0].id) : '')
      })
      .catch((error) => {
        if (cancelled) return
        toast.error(error instanceof Error ? error.message : t('Failed to load API keys'))
        setKeys([])
        setSelectedId('')
      })
      .finally(() => {
        if (!cancelled) setLoadingKeys(false)
      })
    return () => {
      cancelled = true
    }
  }, [t])

  const handleSubmit = async () => {
    const tokenId = Number(selectedId)
    if (!tokenId) {
      toast.error('请先选择一个 API Key')
      return
    }
    setSubmitting(true)
    try {
      const url = await buildCanvasRedirectUrl(tokenId, canvasBaseUrl)
      openCanvasInNewTab(url)
      onSuccess?.()
    } catch (error) {
      toast.error(error instanceof Error ? error.message : '打开画布失败')
    } finally {
      setSubmitting(false)
    }
  }

  if (loadingKeys) {
    return (
      <div className='flex items-center justify-center gap-2 py-8 text-sm text-muted-foreground'>
        <Loader2 className='size-4 animate-spin' />
        正在加载 API Key…
      </div>
    )
  }

  if (!keys.length) {
    return (
      <div className='space-y-4 py-2'>
        <p className='text-sm text-muted-foreground'>
          暂无可用 API Key。请先在控制台创建并启用 Token，再打开{toolName}。
        </p>
        <Button asChild variant='outline'>
          <a href='/keys' target='_blank' rel='noopener noreferrer'>
            前往创建 API Key
          </a>
        </Button>
      </div>
    )
  }

  const selectedKey = keys.find((item) => String(item.id) === selectedId)

  return (
    <div className='space-y-4'>
      <p className='text-sm text-muted-foreground'>
        请选择用于连接{toolName}的 API Key。{toolName}将使用该 Key 直连网关并拉取可用模型。
      </p>
      <div className='space-y-2'>
        <Label htmlFor='canvas-api-key'>API Key</Label>
        <Select value={selectedId} onValueChange={setSelectedId}>
          <SelectTrigger id='canvas-api-key' className='w-full'>
            <SelectValue placeholder='选择 API Key' />
          </SelectTrigger>
          <SelectContent>
            {keys.map((item) => {
              const status = API_KEY_STATUSES[item.status]
              return (
                <SelectItem key={item.id} value={String(item.id)}>
                  <span className='flex items-center gap-2'>
                    <span>{item.name || `Token #${item.id}`}</span>
                    {status ? (
                      <StatusBadge variant={status.variant}>{t(status.label)}</StatusBadge>
                    ) : null}
                  </span>
                </SelectItem>
              )
            })}
          </SelectContent>
        </Select>
        {selectedKey ? (
          <p className='text-xs text-muted-foreground'>
            {selectedKey.group ? `分组：${selectedKey.group}` : '未设置分组'}
            {selectedKey.model_limits_enabled && selectedKey.model_limits
              ? ` · 模型限制：${selectedKey.model_limits}`
              : ''}
          </p>
        ) : null}
      </div>
      <Button
        type='button'
        className='w-full'
        disabled={submitting || !selectedId}
        onClick={() => void handleSubmit()}
      >
        {submitting ? (
          <>
            <Loader2 className='size-4 animate-spin' />
            正在跳转…
          </>
        ) : (
          <>
            <ExternalLink className='size-4' />
            {submitLabel || '打开无限画布'}
          </>
        )}
      </Button>
    </div>
  )
}
