/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/

/** POST /v1/images/generations（JSON）参数全集。 */
export const imageGenerationsParams = [
  { name: 'model', type: 'string', required: '是', default: '—', note: '模型广场展示名，见下方对照表' },
  { name: 'prompt', type: 'string', required: '是', default: '—', note: '图像描述' },
  { name: 'n', type: 'integer', required: '否', default: '1', note: '生成数量，见对照表上限' },
  { name: 'size', type: 'string', required: '否', default: '因模型而异', note: '1024x1024、1536x1024、1024x1536、auto 等' },
  { name: 'quality', type: 'string', required: '否', default: 'auto', note: 'auto / low / medium / high（部分模型支持）' },
  { name: 'response_format', type: 'string', required: '否', default: 'b64_json', note: 'b64_json 或 url' },
  { name: 'image', type: 'string | string[]', required: '否', default: '—', note: '参考图 URL / base64；单张传字符串，多张传数组' },
  { name: 'images', type: 'string[]', required: '否', default: '—', note: '多参考图数组（与 image 等效）' },
  { name: 'mask', type: 'string', required: '否', default: '—', note: '蒙版图，局部重绘（Gemini 等）' },
  { name: 'background', type: 'string', required: '否', default: 'auto', note: 'auto / opaque / transparent（GPT 扩展）' },
  { name: 'output_format', type: 'string', required: '否', default: '—', note: 'png / jpeg / webp（GPT 扩展）' },
  { name: 'output_compression', type: 'integer', required: '否', default: '100', note: '0–100，JPEG/WebP 压缩率' },
] as const

/** POST /v1/images/edits（multipart）字段。 */
export const imageEditsFields = [
  { name: 'image', note: '参考图文件（必填，可多张）' },
  { name: 'prompt', note: '编辑描述（必填）' },
  { name: 'model', note: '模型名' },
  { name: 'mask', note: '蒙版文件，局部重绘可选' },
  { name: 'n', note: '生成数量' },
  { name: 'size', note: '输出尺寸' },
  { name: 'quality', note: '质量档位' },
] as const

/** POST /v1/chat/completions 生图参数。 */
export const chatImageParams = [
  { name: 'model', type: 'string', required: '是', note: 'grok-imagine-image* 等' },
  { name: 'messages', type: 'array', required: '是', note: 'user 消息；图生时在 content 中附 image_url' },
  { name: 'stream', type: 'boolean', required: '否', note: 'Grok 等可选流式' },
] as const

/** gpt-image-2-4k 异步生图（走 /v1/videos）。 */
export const asyncImageParams = [
  { name: 'model', type: 'string', required: '是', note: 'gpt-image-2-4k' },
  { name: 'prompt', type: 'string', required: '是', note: '图像描述' },
] as const

export type ImageModelCapability = {
  model: string
  vendor: string
  api: string
  billing: string
  prompt: string
  size: string
  quality: string
  count: string
  refImages: string
  mask: string
  asyncTask: string
}

/** Public 模型名 → 参数支持。完整列表见模型广场。 */
export const imageModelCapabilities: ImageModelCapability[] = [
  {
    model: 'gpt-image-2',
    vendor: 'OpenAI GPT-Image',
    api: 'POST /v1/images/*',
    billing: '按次',
    prompt: '✓',
    size: '1024² / 1536×1024 / 1024×1536 / auto',
    quality: '—',
    count: '1–4',
    refImages: 'edits 端点',
    mask: 'edits',
    asyncTask: '—',
  },
  {
    model: 'gpt-image-2-4k',
    vendor: 'OpenAI GPT-Image',
    api: 'POST /v1/videos',
    billing: '按次',
    prompt: '✓',
    size: '固定 2880×2880',
    quality: '—',
    count: '1',
    refImages: '—',
    mask: '—',
    asyncTask: '✓（轮询取 video_url）',
  },
  {
    model: 'gemini-image',
    vendor: 'Gemini',
    api: 'POST /v1/images/generations',
    billing: '按次',
    prompt: '✓',
    size: '1024² / 1536×1024 / 1024×1536 / auto',
    quality: '—',
    count: '1–4',
    refImages: '≤5（image/images）',
    mask: '✓',
    asyncTask: '—',
  },
  {
    model: 'gemini-image-pro',
    vendor: 'Gemini',
    api: 'POST /v1/images/generations',
    billing: '按次',
    prompt: '✓',
    size: '1024² / 1536×1024 / 1024×1536 / auto',
    quality: '—',
    count: '1–4',
    refImages: '≤5',
    mask: '✓',
    asyncTask: '—',
  },
  {
    model: 'grok-imagine-image',
    vendor: 'Grok Imagine',
    api: 'POST /v1/chat/completions',
    billing: '按次',
    prompt: '✓',
    size: '1024×1024',
    quality: '—',
    count: '1–4',
    refImages: 'Chat image_url',
    mask: '—',
    asyncTask: '—',
  },
  {
    model: 'grok-imagine-image-lite',
    vendor: 'Grok Imagine',
    api: 'POST /v1/chat/completions',
    billing: '按次',
    prompt: '✓',
    size: '1024×1024',
    quality: '—',
    count: '1–4',
    refImages: 'Chat image_url',
    mask: '—',
    asyncTask: '—',
  },
  {
    model: 'grok-imagine-image-pro',
    vendor: 'Grok Imagine',
    api: 'POST /v1/chat/completions',
    billing: '按次',
    prompt: '✓',
    size: '1024×1024',
    quality: '—',
    count: '1–4',
    refImages: 'Chat image_url',
    mask: '—',
    asyncTask: '—',
  },
  {
    model: 'grok-imagine-image-edit',
    vendor: 'Grok Imagine',
    api: 'POST /v1/chat/completions',
    billing: '按次',
    prompt: '✓',
    size: '1024×1024',
    quality: '—',
    count: '1–4',
    refImages: '必填',
    mask: '—',
    asyncTask: '—',
  },
]
