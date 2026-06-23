/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/
import { QUOTA_TYPE_VALUES } from '@/features/pricing/constants'
import type { PricingModel, PriceType } from '@/features/pricing/types'

export interface FeaturedModelSpec {
  display: string
  pattern: RegExp
  official?: { input: number; output: number }
}

/** Frontier models shown on the home pricing table, in display order */
export const FEATURED_MODEL_ORDER: FeaturedModelSpec[] = [
  {
    display: 'Claude Opus 4',
    pattern: /claude-opus-4(?:\.\d|-)/i,
    official: { input: 15, output: 75 },
  },
  {
    display: 'Claude Sonnet 4',
    pattern: /claude-sonnet-4(?:\.\d|-)/i,
    official: { input: 3, output: 15 },
  },
  {
    display: 'Claude Haiku 4',
    pattern: /claude-haiku-4(?:\.\d|-)/i,
    official: { input: 0.8, output: 4 },
  },
  {
    display: 'GPT-4o',
    pattern: /^gpt-4o(?!-mini)(?:-2024|-audio|$)/i,
    official: { input: 2.5, output: 10 },
  },
  {
    display: 'GPT-4.1 Mini',
    pattern: /^gpt-4\.1-mini/i,
    official: { input: 0.4, output: 1.6 },
  },
  {
    display: 'GPT-4.1',
    pattern: /^gpt-4\.1(?!-mini|-nano)/i,
    official: { input: 2, output: 8 },
  },
  {
    display: 'DeepSeek V3',
    pattern: /deepseek-(?:chat|v3)/i,
    official: { input: 0.27, output: 1.1 },
  },
  {
    display: 'Gemini 2.5 Pro',
    pattern: /gemini-2\.5-pro/i,
    official: { input: 1.25, output: 10 },
  },
]

function getMinGroupRatio(
  enableGroups: string[],
  groupRatio: Record<string, number>
): number {
  if (enableGroups.length === 0) return 1
  let minRatio = Number.POSITIVE_INFINITY
  for (const group of enableGroups) {
    const ratio = groupRatio[group]
    if (ratio !== undefined && ratio < minRatio) minRatio = ratio
  }
  return minRatio === Number.POSITIVE_INFINITY ? 1 : minRatio
}

function hasRatio(value: number | null | undefined): boolean {
  return value !== undefined && value !== null && Number.isFinite(Number(value))
}

function calculateTokenPriceUSD(
  model: PricingModel,
  type: PriceType,
  ratio: number
): number | null {
  const base = model.model_ratio * 2 * ratio
  switch (type) {
    case 'input':
      return base
    case 'output':
      return base * model.completion_ratio
    case 'cache':
      return hasRatio(model.cache_ratio)
        ? base * Number(model.cache_ratio)
        : null
    case 'create_cache':
      return hasRatio(model.create_cache_ratio)
        ? base * Number(model.create_cache_ratio)
        : null
    default:
      return null
  }
}

export interface HomePricingRow {
  display: string
  modelName: string
  input: number
  output: number
  cacheRead: number | null
  cacheWrite: number | null
  officialInput: number | null
  officialOutput: number | null
  savePercent: number | null
}

export function formatUsdPerM(value: number | null): string {
  if (value === null || !Number.isFinite(value)) return '—'
  if (value >= 1) return `$${value.toFixed(2)}`
  if (value >= 0.01) return `$${value.toFixed(3)}`
  return `$${value.toFixed(4)}`
}

export function pickFeaturedPricingRows(
  models: PricingModel[]
): HomePricingRow[] {
  const tokenModels = models.filter(
    (m) => m.quota_type === QUOTA_TYPE_VALUES.TOKEN
  )
  const rows: HomePricingRow[] = []

  for (const spec of FEATURED_MODEL_ORDER) {
    const model = tokenModels.find((m) => spec.pattern.test(m.model_name))
    if (!model) continue

    const enableGroups = Array.isArray(model.enable_groups)
      ? model.enable_groups
      : []
    const groupRatio = model.group_ratio || {}
    const minRatio = getMinGroupRatio(enableGroups, groupRatio)

    const input = calculateTokenPriceUSD(model, 'input', minRatio)!
    const output = calculateTokenPriceUSD(model, 'output', minRatio)!
    const cacheRead = calculateTokenPriceUSD(model, 'cache', minRatio)
    const cacheWrite = calculateTokenPriceUSD(model, 'create_cache', minRatio)

    let savePercent: number | null = null
    if (spec.official && input > 0) {
      const avgOur = (input + output) / 2
      const avgOfficial = (spec.official.input + spec.official.output) / 2
      savePercent = Math.round((1 - avgOur / avgOfficial) * 100)
      if (savePercent < 0) savePercent = null
    }

    rows.push({
      display: spec.display,
      modelName: model.model_name,
      input,
      output,
      cacheRead,
      cacheWrite,
      officialInput: spec.official?.input ?? null,
      officialOutput: spec.official?.output ?? null,
      savePercent,
    })
  }

  return rows
}
