/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/
import { ENDPOINT_TYPES } from '../constants'
import type { PricingData, PricingModel } from '../types'

const DEFAULT_UI_PROFILE_PREFIX = 'default-'

function isBoundImageProfile(model: PricingModel): boolean {
  const profileId = (model.image_ui_params as { id?: string } | undefined)?.id
  return Boolean(
    profileId &&
      !profileId.startsWith(DEFAULT_UI_PROFILE_PREFIX) &&
      profileId.startsWith('image-tpl')
  )
}

/** 视频 API 文档模型：绑定 openai-video 端点，排除 Chat 默认 UI 参数误匹配。 */
export function isVideoDocModel(model: PricingModel): boolean {
  return (
    model.supported_endpoint_types?.includes(ENDPOINT_TYPES.OPENAI_VIDEO) ??
    false
  )
}

/** 图像 API 文档模型：绑定 image-generation 或显式 image profile。 */
export function isImageDocModel(model: PricingModel): boolean {
  if (
    model.supported_endpoint_types?.includes(ENDPOINT_TYPES.IMAGE_GENERATION)
  ) {
    return true
  }
  return isBoundImageProfile(model)
}

export function isModelDocCandidate(
  model: PricingModel,
  capability: 'video' | 'image' | 'all'
): boolean {
  if (capability === 'video') return isVideoDocModel(model)
  if (capability === 'image') return isImageDocModel(model)
  return isVideoDocModel(model) || isImageDocModel(model)
}

export function enrichPricingModels(data: PricingData): PricingModel[] {
  const vendorMap = new Map(data.vendors.map((vendor) => [vendor.id, vendor]))

  return data.data.map((model) => {
    const vendor = model.vendor_id ? vendorMap.get(model.vendor_id) : undefined
    return {
      ...model,
      key: model.model_name,
      display_name: model.model_name,
      vendor_name: vendor?.name,
      vendor_icon: vendor?.icon,
      vendor_description: vendor?.description,
      group_ratio: data.group_ratio,
    }
  })
}
