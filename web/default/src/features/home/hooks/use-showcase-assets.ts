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

const DEFAULT_ASPECT_RATIO = 3 / 4
const SHOWCASE_PAGE_SIZE = 100

function fallbackAssets(): ShowcaseAsset[] {
  return INSPIRATION_SLIDES.map((slide) => ({
    id: slide.id,
    image: slide.image,
    title: slide.title,
    tags: slide.tags,
    aspectRatio: slide.width / slide.height,
  }))
}

function isValidCoverUrl(url: string): boolean {
  return url.startsWith('http://') || url.startsWith('https://')
}

function probeImageAspectRatio(url: string): Promise<number> {
  return new Promise((resolve) => {
    const img = new Image()
    img.onload = () => {
      const { naturalWidth, naturalHeight } = img
      if (naturalWidth > 0 && naturalHeight > 0) {
        resolve(naturalWidth / naturalHeight)
      } else {
        resolve(DEFAULT_ASPECT_RATIO)
      }
    }
    img.onerror = () => resolve(DEFAULT_ASPECT_RATIO)
    img.src = url
  })
}

async function enrichAssetsWithAspectRatio(
  assets: Omit<ShowcaseAsset, 'aspectRatio'>[]
): Promise<ShowcaseAsset[]> {
  const ratios = await Promise.all(
    assets.map((asset) => probeImageAspectRatio(asset.image))
  )
  return assets.map((asset, index) => ({
    ...asset,
    aspectRatio: ratios[index],
  }))
}

async function fetchShowcaseAssets(): Promise<ShowcaseAsset[]> {
  try {
    const url = new URL('/api/prompts', DEFAULT_CANVAS_BASE_URL)
    url.searchParams.set('modality', 'image')
    url.searchParams.set('previewType', 'image')
    url.searchParams.set('pageSize', String(SHOWCASE_PAGE_SIZE))

    const res = await fetch(url.toString())
    if (!res.ok) return fallbackAssets()

    const data = (await res.json()) as PromptListResponse
    const baseAssets = (data.items ?? [])
      .filter((item) => item.coverUrl && isValidCoverUrl(item.coverUrl))
      .map((item) => ({
        id: item.id,
        image: item.coverUrl,
        title: item.title,
        tags: (item.tags ?? []).slice(0, 3).join(' / '),
      }))

    if (baseAssets.length === 0) return fallbackAssets()

    return enrichAssetsWithAspectRatio(baseAssets)
  } catch {
    return fallbackAssets()
  }
}

export function useShowcaseAssets() {
  return useQuery({
    queryKey: ['home-showcase-assets'],
    queryFn: fetchShowcaseAssets,
    staleTime: 30 * 60 * 1000,
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

/** Scale animation duration by item count so scroll speed stays comfortable */
export function showcaseMarqueeDuration(itemCount: number, baseSec = 45): number {
  return Math.max(baseSec, Math.round(itemCount * 2.5))
}
