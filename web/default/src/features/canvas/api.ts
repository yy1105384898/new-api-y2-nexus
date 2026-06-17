/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/
import { API_KEY_STATUS } from '@/features/keys/constants'
import { fetchTokenKey, getApiKeys } from '@/features/keys/api'
import type { ApiKey } from '@/features/keys/types'
import { api } from '@/lib/api'

export function appendCanvasConfigToUrl(
  targetUrl: string,
  apiKey: string,
  baseUrl: string
) {
  if (!apiKey && !baseUrl) return targetUrl
  try {
    const url = new URL(targetUrl)
    if (apiKey) url.searchParams.set('apiKey', apiKey)
    if (baseUrl) url.searchParams.set('baseUrl', baseUrl)
    return url.toString()
  } catch {
    const params = new URLSearchParams()
    if (apiKey) params.set('apiKey', apiKey)
    if (baseUrl) params.set('baseUrl', baseUrl)
    const separator = targetUrl.includes('?') ? '&' : '?'
    return `${targetUrl}${separator}${params.toString()}`
  }
}

export async function resolveCanvasGatewayBaseUrl() {
  try {
    const res = await api.get('/api/status')
    const address = res?.data?.data?.server_address
    if (typeof address === 'string' && address.trim()) {
      return address.trim().replace(/\/+$/, '')
    }
  } catch {
    // ignore
  }
  if (typeof window !== 'undefined') {
    return window.location.origin.replace(/\/+$/, '')
  }
  return ''
}

export async function listSelectableCanvasApiKeys(): Promise<ApiKey[]> {
  const list = await getApiKeys({ p: 1, size: 100 })
  const items = list?.data?.items || []
  return items.filter((item) => item.status === API_KEY_STATUS.ENABLED)
}

export async function buildCanvasRedirectUrl(
  tokenId: number,
  canvasBaseUrl: string
) {
  const [revealed, baseUrl] = await Promise.all([
    fetchTokenKey(tokenId),
    resolveCanvasGatewayBaseUrl(),
  ])
  const apiKey = revealed?.data?.key || ''
  if (!apiKey) {
    throw new Error('Failed to reveal API key')
  }
  if (!baseUrl) {
    throw new Error('Failed to resolve gateway base URL')
  }
  return appendCanvasConfigToUrl(canvasBaseUrl, apiKey, baseUrl)
}

export function openCanvasInNewTab(url: string) {
  window.open(url, '_blank', 'noopener,noreferrer')
}
