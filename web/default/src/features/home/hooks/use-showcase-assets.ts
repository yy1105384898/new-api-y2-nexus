/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/
import { useQuery } from '@tanstack/react-query'
import { DEFAULT_CANVAS_BASE_URL } from '@/features/canvas/lib/canvas-config'
import { INSPIRATION_SLIDES } from '../lib/site-assets'
import type { ShowcaseAsset } from '../types'

interface PromptListItem {
  id: string
  title: string
  coverUrl: string
  tags?: string[]
}

interface PromptListResponse {
  items: PromptListItem[]
}

const CDN_HOST =
  (import.meta.env.VITE_STATIC_CDN || 'https://assets.cangyuansuanli.cn')
    .replace(/^https?:\/\//, '')
    .replace(/\/$/, '')

function isCdnAssetUrl(url: string): boolean {
  try {
    const host = new URL(url).hostname
    return host === CDN_HOST || host.endsWith('.cangyuansuanli.cn')
  } catch {
    return false
  }
}

function fallbackAssets(): ShowcaseAsset[] {
  return INSPIRATION_SLIDES.map((slide) => ({
    id: slide.id,
    image: slide.image,
    title: slide.title,
    tags: slide.tags,
    aspectRatio: slide.width / slide.height,
  }))
}

function dedupeByImage(assets: ShowcaseAsset[]): ShowcaseAsset[] {
  const seen = new Set<string>()
  return assets.filter((asset) => {
    if (seen.has(asset.image)) return false
    seen.add(asset.image)
    return true
  })
}

async function fetchShowcaseAssets(): Promise<ShowcaseAsset[]> {
  const url = new URL('/api/prompts', DEFAULT_CANVAS_BASE_URL)
  url.searchParams.set('modality', 'image')
  url.searchParams.set('previewType', 'image')
  url.searchParams.set('pageSize', '48')

  const res = await fetch(url.toString())
  if (!res.ok) return fallbackAssets()

  const data = (await res.json()) as PromptListResponse
  const fromPrompts = (data.items ?? [])
    .filter((item) => item.coverUrl && isCdnAssetUrl(item.coverUrl))
    .map((item) => ({
      id: item.id,
      image: item.coverUrl,
      title: item.title,
      tags: (item.tags ?? []).slice(0, 3).join(' / '),
    }))

  const merged = dedupeByImage([...fromPrompts, ...fallbackAssets()])
  return merged.length >= 6 ? merged : fallbackAssets()
}

export function useShowcaseAssets() {
  return useQuery({
    queryKey: ['home-showcase-assets'],
    queryFn: fetchShowcaseAssets,
    staleTime: 30 * 60 * 1000,
    placeholderData: fallbackAssets(),
  })
}

/** Split assets into two rows for opposite marquee directions */
export function splitShowcaseRows(assets: ShowcaseAsset[]) {
  const rowA: ShowcaseAsset[] = []
  const rowB: ShowcaseAsset[] = []
  assets.forEach((asset, index) => {
    if (index % 2 === 0) rowA.push(asset)
    else rowB.push(asset)
  })
  if (rowB.length === 0 && rowA.length > 1) {
    const half = Math.ceil(rowA.length / 2)
    return {
      rowA: rowA.slice(0, half),
      rowB: rowA.slice(half),
    }
  }
  return { rowA, rowB }
}

export function loopShowcaseRow(assets: ShowcaseAsset[]) {
  if (assets.length === 0) return []
  return [...assets, ...assets]
}
