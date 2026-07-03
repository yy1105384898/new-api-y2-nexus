/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/

export const DEFAULT_API_BASE_URL = 'https://yynewapi.yangyangnj.top'
export const DEFAULT_CANVAS_BASE_URL = 'https://canvas.yangyangnj.top'
export const DEFAULT_HUIYING_BASE_URL = 'https://huiying.yangyangnj.top'
export const DEFAULT_ECOMPIC_BASE_URL = 'https://ecompic.yangyangnj.top'
export const DEFAULT_CANVAS_DOCS_URL = `${DEFAULT_CANVAS_BASE_URL}/docs`

function normalizeOrigin(url: string) {
  try {
    return new URL(url).origin
  } catch {
    return ''
  }
}

export function getCanvasToolName(targetUrl: string) {
  const targetOrigin = normalizeOrigin(targetUrl)
  if (targetOrigin === normalizeOrigin(DEFAULT_HUIYING_BASE_URL)) return '绘影'
  if (targetOrigin === normalizeOrigin(DEFAULT_ECOMPIC_BASE_URL)) return '竞品图工具'
  return '无限画布'
}
