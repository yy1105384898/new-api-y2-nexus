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

interface ShowcaseMarqueeRowProps {
  assets: ShowcaseAsset[]
  direction?: 'left' | 'right'
  durationSec?: number
  onSelect?: () => void
}

export function ShowcaseMarqueeRow(props: ShowcaseMarqueeRowProps) {
  const { assets, direction = 'left', durationSec = 45, onSelect } = props
  if (assets.length === 0) return null

  return (
    <div className='showcase-marquee-container overflow-hidden'>
      <div
        className={cn(
          'showcase-marquee-track flex w-max gap-3 sm:gap-4',
          direction === 'left' ? 'animate-marquee-left' : 'animate-marquee-right'
        )}
        style={{ animationDuration: `${durationSec}s` }}
      >
        {assets.map((asset, index) => (
          <button
            key={`${asset.id}-${index}`}
            type='button'
            onClick={onSelect}
            className={cn(
              'group relative shrink-0 overflow-hidden rounded-xl border transition-all duration-300',
              'hover:border-cyan-500/40 hover:shadow-[0_0_32px_-8px_rgba(6,182,212,0.35)]',
              'focus-visible:ring-2 focus-visible:ring-cyan-500/50 focus-visible:outline-none',
              mkt.mediaCard
            )}
            style={{
              width: 'clamp(140px, 18vw, 220px)',
              aspectRatio: asset.aspectRatio ? `${asset.aspectRatio}` : '3/4',
            }}
          >
            <img
              src={asset.image}
              alt={asset.title}
              loading='lazy'
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
        ))}
      </div>
    </div>
  )
}
