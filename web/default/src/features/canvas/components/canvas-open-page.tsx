/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/
import {
  DEFAULT_CANVAS_BASE_URL,
  getCanvasToolName,
} from '@/features/canvas/lib/canvas-config'
import { CanvasKeySelectForm } from '@/features/canvas/components/canvas-key-select-form'

type CanvasOpenPageProps = {
  redirect?: string
}

export function CanvasOpenPage({ redirect }: CanvasOpenPageProps) {
  const canvasBaseUrl = redirect?.trim() || DEFAULT_CANVAS_BASE_URL
  const toolName = getCanvasToolName(canvasBaseUrl)

  return (
    <div className='mx-auto flex min-h-[60vh] w-full max-w-lg flex-col justify-center px-4 py-10'>
      <div className='rounded-xl border bg-card p-6 shadow-sm'>
        <h1 className='text-xl font-semibold'>{`打开${toolName}`}</h1>
        <p className='mt-2 text-sm text-muted-foreground'>
          登录成功。请选择要带入{toolName}的 API Key，确认后将在新标签页打开。
        </p>
        <div className='mt-6'>
          <CanvasKeySelectForm
            canvasBaseUrl={canvasBaseUrl}
            submitLabel={`打开${toolName}`}
            toolName={toolName}
          />
        </div>
      </div>
    </div>
  )
}
