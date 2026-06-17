/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/
import { useCallback, useEffect, useRef, useState } from 'react'
import { ArrowRight } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { cn } from '@/lib/utils'
import {
  Carousel,
  CarouselContent,
  CarouselItem,
  CarouselNext,
  CarouselPrevious,
  type CarouselApi,
} from '@/components/ui/carousel'
import { DEFAULT_CANVAS_BASE_URL } from '@/features/canvas/lib/canvas-config'
import { useCanvasKeyPicker } from '@/features/canvas/hooks/use-canvas-key-picker'
import { mkt } from '../lib/marketing-theme'
import { INSPIRATION_SLIDES } from '../lib/site-assets'

const MAX_SLIDE_HEIGHT_PX = 580
const MAX_SLIDE_HEIGHT_VH = 0.65

function slideFrameSize(
  containerWidth: number,
  slide: (typeof INSPIRATION_SLIDES)[number]
) {
  const maxHeight = Math.min(
    window.innerHeight * MAX_SLIDE_HEIGHT_VH,
    MAX_SLIDE_HEIGHT_PX
  )
  const naturalHeight = containerWidth * (slide.height / slide.width)

  if (naturalHeight <= maxHeight) {
    return { width: containerWidth, height: naturalHeight }
  }

  return {
    width: maxHeight * (slide.width / slide.height),
    height: maxHeight,
  }
}

export function FreeImageCarousel() {
  const { t } = useTranslation()
  const { requestOpen, dialog } = useCanvasKeyPicker(DEFAULT_CANVAS_BASE_URL)
  const containerRef = useRef<HTMLDivElement>(null)
  const [api, setApi] = useState<CarouselApi>()
  const [active, setActive] = useState(0)
  const [frame, setFrame] = useState({ width: 0, height: 400 })

  const onSelect = useCallback((carouselApi: CarouselApi) => {
    if (!carouselApi) return
    setActive(carouselApi.selectedScrollSnap())
  }, [])

  const updateFrame = useCallback(() => {
    const container = containerRef.current
    if (!container) return
    const containerWidth = container.offsetWidth
    if (!containerWidth) return
    setFrame(slideFrameSize(containerWidth, INSPIRATION_SLIDES[active]))
  }, [active])

  useEffect(() => {
    if (!api) return
    onSelect(api)
    api.on('select', onSelect)
    return () => {
      api.off('select', onSelect)
    }
  }, [api, onSelect])

  useEffect(() => {
    if (!api) return
    const timer = window.setInterval(() => {
      if (api.canScrollNext()) api.scrollNext()
      else api.scrollTo(0)
    }, 5000)
    return () => window.clearInterval(timer)
  }, [api])

  useEffect(() => {
    updateFrame()
    const container = containerRef.current
    if (!container) return

    const observer = new ResizeObserver(updateFrame)
    observer.observe(container)
    window.addEventListener('resize', updateFrame)

    return () => {
      observer.disconnect()
      window.removeEventListener('resize', updateFrame)
    }
  }, [updateFrame])

  useEffect(() => {
    if (!api) return
    api.reInit()
  }, [api, frame])

  return (
    <div ref={containerRef} className='relative w-full'>
      <div
        className='mx-auto overflow-hidden transition-[width,height] duration-500 ease-out'
        style={{
          width: frame.width || '100%',
          height: frame.height,
        }}
      >
        <Carousel
          setApi={setApi}
          opts={{ align: 'center', loop: true, duration: 28 }}
          className='size-full'
        >
          <CarouselContent className='ml-0 h-full'>
            {INSPIRATION_SLIDES.map((slide) => (
              <CarouselItem
                key={slide.id}
                className='h-full basis-full pl-0'
              >
                <div
                  role='button'
                  tabIndex={0}
                  onClick={() => requestOpen()}
                  onKeyDown={(event) => {
                    if (event.key === 'Enter' || event.key === ' ') {
                      event.preventDefault()
                      requestOpen()
                    }
                  }}
                  className={cn(
                    'group relative block size-full cursor-pointer overflow-hidden rounded-2xl',
                    mkt.mediaCard
                  )}
                >
                  <img
                    src={slide.image}
                    alt={slide.title}
                    loading='lazy'
                    className='size-full transition-transform duration-700 ease-out group-hover:scale-[1.02]'
                  />
                  <div
                    className={cn(
                      'pointer-events-none absolute inset-x-0 bottom-0 h-[42%] min-h-[7rem]',
                      mkt.mediaOverlay
                    )}
                  />
                  <div className='absolute inset-x-0 bottom-0 flex flex-col gap-2 p-4 sm:flex-row sm:items-end sm:justify-between sm:gap-3 sm:p-6'>
                    <div>
                      <span className='mb-1.5 inline-flex rounded-full border border-cyan-400/40 bg-cyan-500/15 px-2.5 py-0.5 text-[10px] font-semibold tracking-wide text-cyan-100 uppercase'>
                        {t('Free trial')}
                      </span>
                      <p className='text-[11px] font-medium tracking-[0.2em] text-cyan-200/90 uppercase'>
                        {t('Cangyuan Image to Video')}
                      </p>
                      <h3 className='mt-0.5 text-lg font-semibold text-white sm:text-xl'>
                        {slide.title}
                      </h3>
                      <p className='mt-0.5 text-sm text-white/75'>{slide.tags}</p>
                    </div>
                    <span className='inline-flex w-fit items-center gap-1.5 rounded-lg bg-white px-3.5 py-2 text-sm font-semibold text-slate-900 transition-colors group-hover:bg-cyan-50 sm:px-4 sm:py-2.5'>
                      {t('Try for free')}
                      <ArrowRight className='size-4' />
                    </span>
                  </div>
                </div>
              </CarouselItem>
            ))}
          </CarouselContent>
          <CarouselPrevious className={cn('left-2 sm:left-3', mkt.carouselBtn)} />
          <CarouselNext className={cn('right-2 sm:right-3', mkt.carouselBtn)} />
        </Carousel>
      </div>

      <div className='mt-3 flex items-center justify-center gap-1.5'>
        {INSPIRATION_SLIDES.map((slide, index) => (
          <button
            key={slide.id}
            type='button'
            aria-label={`${slide.title} ${index + 1}`}
            onClick={() => api?.scrollTo(index)}
            className={cn(
              'h-1.5 rounded-full transition-all',
              active === index
                ? cn('w-7', mkt.carouselDotActive)
                : cn('w-1.5', mkt.carouselDot)
            )}
          />
        ))}
      </div>
      {dialog}
    </div>
  )
}
