/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/
import { useTranslation } from 'react-i18next'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { CanvasKeySelectForm } from '@/features/canvas/components/canvas-key-select-form'

type CanvasKeySelectDialogProps = {
  open: boolean
  onOpenChange: (open: boolean) => void
  canvasBaseUrl: string
}

export function CanvasKeySelectDialog({
  open,
  onOpenChange,
  canvasBaseUrl,
}: CanvasKeySelectDialogProps) {
  const { t } = useTranslation()

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className='sm:max-w-md'>
        <DialogHeader>
          <DialogTitle>{t('Open Infinite Canvas')}</DialogTitle>
          <DialogDescription>选择 API Key 后将在新标签页打开无限画布，并自动填入网关地址。</DialogDescription>
        </DialogHeader>
        <CanvasKeySelectForm
          canvasBaseUrl={canvasBaseUrl}
          onSuccess={() => onOpenChange(false)}
        />
      </DialogContent>
    </Dialog>
  )
}
