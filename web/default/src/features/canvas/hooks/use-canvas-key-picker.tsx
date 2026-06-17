/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/
import { useCallback, useState } from 'react'
import { useAuthStore } from '@/stores/auth-store'
import { DEFAULT_CANVAS_BASE_URL } from '@/features/canvas/lib/canvas-config'
import { openCanvasInNewTab } from '@/features/canvas/api'
import { CanvasKeySelectDialog } from '@/features/canvas/components/canvas-key-select-dialog'

export function useCanvasKeyPicker(defaultCanvasUrl = DEFAULT_CANVAS_BASE_URL) {
  const isAuthenticated = !!useAuthStore((state) => state.auth.user)
  const [open, setOpen] = useState(false)
  const [targetUrl, setTargetUrl] = useState(defaultCanvasUrl)

  const requestOpen = useCallback(
    (canvasBaseUrl = defaultCanvasUrl) => {
      if (!isAuthenticated) {
        openCanvasInNewTab(canvasBaseUrl)
        return
      }
      setTargetUrl(canvasBaseUrl)
      setOpen(true)
    },
    [defaultCanvasUrl, isAuthenticated]
  )

  const dialog = (
    <CanvasKeySelectDialog open={open} onOpenChange={setOpen} canvasBaseUrl={targetUrl} />
  )

  return { requestOpen, dialog }
}
