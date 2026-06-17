/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/
import { useTranslation } from 'react-i18next'
import { cn } from '@/lib/utils'
import { DEFAULT_CANVAS_BASE_URL } from '@/features/canvas/lib/canvas-config'
import { useCanvasKeyPicker } from '@/features/canvas/hooks/use-canvas-key-picker'

type CanvasTopNavLinkProps = {
  className?: string
  style?: React.CSSProperties
  onClick?: () => void
}

export function CanvasTopNavLink({ className, style, onClick }: CanvasTopNavLinkProps) {
  const { t } = useTranslation()
  const { requestOpen, dialog } = useCanvasKeyPicker(DEFAULT_CANVAS_BASE_URL)

  return (
    <>
      <button
        type='button'
        style={style}
        onClick={(event) => {
          onClick?.()
          event.preventDefault()
          requestOpen()
        }}
        className={cn(
          'hover:text-primary text-sm font-medium transition-colors text-muted-foreground',
          className
        )}
      >
        {t('Cangyuan Image to Video')}
      </button>
      {dialog}
    </>
  )
}
