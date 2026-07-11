/*
Copyright (C) 2023-2026 QuantumNous
*/
import type { TFunction } from 'i18next'
import { VIDEO_API_MODES } from './model-params-types'

export type VideoApiMode = (typeof VIDEO_API_MODES)[number]

const VIDEO_API_MODE_META: Record<
  VideoApiMode,
  { labelKey: string; descKey: string }
> = {
  'videos-json-async': {
    labelKey: 'API mode: videos-json-async',
    descKey:
      'POST /v1/videos using the unified video-task contract, then poll GET /v1/videos/{id}.',
  },
}

export function getVideoApiModeLabel(t: TFunction, mode: string) {
  const meta = VIDEO_API_MODE_META[mode as VideoApiMode]
  return meta ? t(meta.labelKey) : mode
}

export function getVideoApiModeDescription(t: TFunction, mode: string) {
  const meta = VIDEO_API_MODE_META[mode as VideoApiMode]
  return meta ? t(meta.descKey) : ''
}

export function getVideoApiModeOptions(t: TFunction) {
  return VIDEO_API_MODES.map((mode) => ({
    value: mode,
    label: getVideoApiModeLabel(t, mode),
    description: getVideoApiModeDescription(t, mode),
  }))
}
