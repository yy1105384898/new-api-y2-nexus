/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/
function normalizeOrigin(url: string) {
  try {
    return new URL(url).origin
  } catch {
    return ''
  }
}

export function isCanvasRedirectUrl(
  redirectTo: string,
  canvasBaseUrl: string
) {
  const targetOrigin = normalizeOrigin(redirectTo)
  const canvasOrigin = normalizeOrigin(canvasBaseUrl)
  return Boolean(targetOrigin && canvasOrigin && targetOrigin === canvasOrigin)
}

export function isExternalRedirect(redirectTo: string) {
  if (!redirectTo.startsWith('http://') && !redirectTo.startsWith('https://')) {
    return false
  }
  try {
    const target = new URL(redirectTo)
    return target.origin !== window.location.origin
  } catch {
    return false
  }
}

/** 登录后跳转画布时，先进入 Key 选择页，不再自动使用第一个 Token。 */
export function buildCanvasOpenPagePath(redirectTo: string) {
  return `/canvas/open?redirect=${encodeURIComponent(redirectTo)}`
}
