/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/
import { cn } from '@/lib/utils'
import { mkt } from '../lib/marketing-theme'
import type { ShowcaseAsset } from '../types'

const ROW_HEIGHT = 'clamp(160px, 20vh, 220px)'
const MIN_CARD_WIDTH = '80px'
const DEFAULT_ASPECT_RATIO = 3 / 4

interface ShowcaseMarqueeRowProps {
  assets: ShowcaseAsset[]
  direction?: 'left' | 'right'
  durationSec?: number
  ready?: boolean
  onSelect?: () => void
}

function cardWidth(aspectRatio: number | undefined): string {
  const ratio =
    aspectRatio && aspectRatio > 0 ? aspectRatio : DEFAULT_ASPECT_RATIO
  return `max(${MIN_CARD_WIDTH}, calc(${ROW_HEIGHT} * ${ratio}))`
}

function ShowcaseCard(props: {
  asset: ShowcaseAsset
  onSelect?: () => void
}) {
  const { asset, onSelect } = props

  return (
    <button
      type='button'
      onClick={onSelect}
      className={cn(
        'group relative shrink-0 overflow-hidden rounded-xl border transition-all duration-300',
        'hover:border-cyan-500/40 hover:shadow-[0_0_32px_-8px_rgba(6,182,212,0.35)]',
        'focus-visible:ring-2 focus-visible:ring-cyan-500/50 focus-visible:outline-none',
        mkt.mediaCard
      )}
      style={{
        height: ROW_HEIGHT,
        width: cardWidth(asset.aspectRatio),
      }}
    >
      <img
        src={asset.image}
        alt={asset.title}
        loading='eager'
        decoding='async'
        className='size-full object-cover transition-transform duration-500 group-hover:scale-105'
      />
      <div
        className={cn(
          'absolute inset-x-0 bottom-0 translate-y-full p-3 transition-transform duration-300 group-hover:translate-y-0',
          mkt.mediaOverlay
        )}
      >
        <p className='line-clamp-1 text-xs font-medium text-white'>
          {asset.title}
        </p>
        {asset.tags ? (
          <p className='mt-0.5 line-clamp-1 text-[10px] text-white/70'>
            {asset.tags}
          </p>
        ) : null}
      </div>
    </button>
  )
}

function MarqueeGroup(props: {
  assets: ShowcaseAsset[]
  groupKey: string
  onSelect?: () => void
  hidden?: boolean
}) {
  const { assets, groupKey, onSelect, hidden } = props

  return (
    <div
      className='flex shrink-0 items-center gap-3 sm:gap-4'
      aria-hidden={hidden || undefined}
    >
      {assets.map((asset) => (
        <ShowcaseCard
          key={`${groupKey}-${asset.id}`}
          asset={asset}
          onSelect={onSelect}
        />
      ))}
    </div>
  )
}

export function ShowcaseMarqueeRow(props: ShowcaseMarqueeRowProps) {
  const {
    assets,
    direction = 'left',
    durationSec = 45,
    ready = true,
    onSelect,
  } = props

  if (assets.length === 0) return null

  return (
    <div className='showcase-marquee-container overflow-hidden'>
      <div
        className={cn(
          'showcase-marquee-track flex w-max gap-0',
          ready &&
            (direction === 'left'
              ? 'animate-marquee-left'
              : 'animate-marquee-right')
        )}
        style={{ animationDuration: `${durationSec}s` }}
      >
        <MarqueeGroup assets={assets} groupKey='a' onSelect={onSelect} />
        <MarqueeGroup
          assets={assets}
          groupKey='b'
          onSelect={onSelect}
          hidden
        />
      </div>
    </div>
  )
}
