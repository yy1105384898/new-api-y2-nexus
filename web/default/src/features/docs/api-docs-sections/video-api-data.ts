/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/

/** Unified POST /v1/videos (JSON or multipart) parameter reference. */
export const videoApiParams = [
  {
    name: 'model',
    type: 'string',
    required: '是',
    default: '—',
    note: '模型广场展示名，见下方能力对照表',
  },
  {
    name: 'prompt',
    type: 'string',
    required: '是*',
    default: '—',
    note: '视频描述；多素材可用 @image1/@video1 占位（Seedance 等）',
  },
  {
    name: 'aspect_ratio',
    type: 'string',
    required: '否',
    default: '因模型而异',
    note: '16:9、9:16、1:1、4:3、3:4、21:9、3:2、2:3 等，见对照表',
  },
  {
    name: 'duration',
    type: 'integer',
    required: '否',
    default: '—',
    note: '时长（秒）；Seedance 等按秒计费模型必填或可调',
  },
  {
    name: 'seconds',
    type: 'string / integer',
    required: '否',
    default: '—',
    note: 'duration 别名；Grok CLI 等使用此字段',
  },
  {
    name: 'size',
    type: 'string',
    required: '否',
    default: '—',
    note: '画幅像素，如 1280x720、720x1280；可与 aspect_ratio 二选一',
  },
  {
    name: 'resolution',
    type: 'string',
    required: '否',
    default: '720p',
    note: '480p / 720p 等清晰度档位',
  },
  {
    name: 'image_url',
    type: 'string',
    required: '否',
    default: '—',
    note: '主参考图：HTTPS URL、data URL 或 multipart 字段 image',
  },
  {
    name: 'reference_image_urls',
    type: 'string[]',
    required: '否',
    default: '—',
    note: '多张参考图；与 image_url 合计上限见对照表',
  },
  {
    name: 'input_reference',
    type: 'string / file',
    required: '否',
    default: '—',
    note: '单张首帧参考图（URL / base64 / multipart）',
  },
  {
    name: 'reference_images',
    type: 'string[] / file[]',
    required: '否',
    default: '—',
    note: '多参考图；与 input_reference 通常二选一',
  },
  {
    name: 'first_image_url',
    type: 'string',
    required: '否',
    default: '—',
    note: '首帧图；须与 last_image_url 成对（首尾帧过渡）',
  },
  {
    name: 'last_image_url',
    type: 'string',
    required: '否',
    default: '—',
    note: '尾帧图；须与 first_image_url 成对',
  },
  {
    name: 'reference_videos',
    type: 'string[]',
    required: '否',
    default: '—',
    note: '参考视频 URL 数组（mp4/mov，2–15s，≤50MB）',
  },
  {
    name: 'video_url',
    type: 'string',
    required: '否',
    default: '—',
    note: 'V2V 源视频 URL（Omni V2V 等）',
  },
  {
    name: 'video',
    type: 'string / file',
    required: '否',
    default: '—',
    note: '源视频：公网 URL、base64 或 multipart（编辑 / V2V）',
  },
  {
    name: 'input_video',
    type: 'file',
    required: '否',
    default: '—',
    note: 'multipart 源视频字段（V2V / 去水印 / 编辑）',
  },
] as const

/** Multipart-only fields (POST /v1/videos, Content-Type: multipart/form-data). */
export const videoMultipartFields = [
  { name: 'image', note: '参考图文件，可多次出现' },
  { name: 'input_reference', note: '单张首帧参考图' },
  { name: 'input_video', note: 'V2V / 去水印源视频' },
  { name: 'video', note: '视频编辑源文件' },
] as const

/** POST /v1/chat/completions — Grok Imagine 视频（stream 推荐）。 */
export const chatVideoParams = [
  { name: 'model', type: 'string', required: '是', note: 'grok-imagine-video' },
  { name: 'stream', type: 'boolean', required: '推荐 true', note: '推送生成进度' },
  { name: 'messages', type: 'array', required: '是', note: 'user 消息；图生时在 content 中附 image_url' },
  {
    name: 'video_config.seconds',
    type: 'integer',
    required: '否',
    note: '6 / 10 / 12 / 16 / 20，默认 6',
  },
  {
    name: 'video_config.size',
    type: 'string',
    required: '否',
    note: '720x1280 / 1280x720 / 960x960',
  },
  {
    name: 'video_config.public_url',
    type: 'boolean',
    required: '否',
    note: '建议 true，返回完整下载链接',
  },
] as const

type SupportCell = string

export type VideoModelCapability = {
  model: string
  vendor: string
  api: string
  billing: string
  prompt: SupportCell
  aspectRatio: SupportCell
  duration: SupportCell
  refImages: SupportCell
  frameTransition: SupportCell
  refVideo: SupportCell
  v2vOrEdit: SupportCell
}

/** Public model names → parameter support. 完整列表见模型广场。 */
export const videoModelCapabilities: VideoModelCapability[] = [
  {
    model: 'omni-fast',
    vendor: 'Gemini Veo',
    api: 'POST /v1/videos',
    billing: '按次',
    prompt: '✓',
    aspectRatio: '16:9 / 9:16',
    duration: '固定 ~10s',
    refImages: '≤5',
    frameTransition: '首/尾帧',
    refVideo: '—',
    v2vOrEdit: '—',
  },
  {
    model: 'omni-fast-no-water',
    vendor: 'Gemini Veo',
    api: 'POST /v1/videos',
    billing: '按次',
    prompt: '✓',
    aspectRatio: '16:9 / 9:16',
    duration: '固定 ~10s',
    refImages: '≤5',
    frameTransition: '首/尾帧',
    refVideo: '—',
    v2vOrEdit: '—',
  },
  {
    model: 'omni-fast-v2v',
    vendor: 'Gemini Veo',
    api: 'POST /v1/videos',
    billing: '按次',
    prompt: '✓',
    aspectRatio: '16:9 / 9:16',
    duration: '固定 ~10s',
    refImages: '—',
    frameTransition: '—',
    refVideo: '—',
    v2vOrEdit: 'V2V（input_video，≤5MB）',
  },
  {
    model: 'omni-fast-v2v-no-water',
    vendor: 'Gemini Veo',
    api: 'POST /v1/videos',
    billing: '按次',
    prompt: '✓',
    aspectRatio: '16:9 / 9:16',
    duration: '固定 ~10s',
    refImages: '—',
    frameTransition: '—',
    refVideo: '—',
    v2vOrEdit: 'V2V（input_video，≤5MB）',
  },
  {
    model: 'veo-clean',
    vendor: 'Gemini Veo',
    api: 'POST /v1/videos',
    billing: '按秒',
    prompt: '可选',
    aspectRatio: '—',
    duration: '—',
    refImages: '—',
    frameTransition: '—',
    refVideo: '—',
    v2vOrEdit: '去水印（input_video，≤20MB）',
  },
  {
    model: 'Seedance2.0-480p',
    vendor: 'Seedance 2.0',
    api: 'POST /v1/videos',
    billing: '按秒',
    prompt: '✓',
    aspectRatio: '16:9~4:3 共 6 种',
    duration: '4–15',
    refImages: '≤9',
    frameTransition: '首/尾帧',
    refVideo: '≤3',
    v2vOrEdit: '—',
  },
  {
    model: 'Seedance2.0-fast-480p',
    vendor: 'Seedance 2.0',
    api: 'POST /v1/videos',
    billing: '按秒',
    prompt: '✓',
    aspectRatio: '16:9~4:3 共 6 种',
    duration: '4–15',
    refImages: '≤9',
    frameTransition: '首/尾帧',
    refVideo: '≤3',
    v2vOrEdit: '—',
  },
  {
    model: 'Seedance2.0-720p',
    vendor: 'Seedance 2.0',
    api: 'POST /v1/videos',
    billing: '按秒',
    prompt: '✓',
    aspectRatio: '16:9~4:3 共 6 种',
    duration: '4–15',
    refImages: '≤9',
    frameTransition: '首/尾帧',
    refVideo: '≤3',
    v2vOrEdit: '—',
  },
  {
    model: 'Seedance2.0-fast-720p',
    vendor: 'Seedance 2.0',
    api: 'POST /v1/videos',
    billing: '按秒',
    prompt: '✓',
    aspectRatio: '16:9~4:3 共 6 种',
    duration: '4–15',
    refImages: '≤9',
    frameTransition: '首/尾帧',
    refVideo: '≤3',
    v2vOrEdit: '—',
  },
  {
    model: 'grok-imagine-video',
    vendor: 'Grok Imagine',
    api: 'POST /v1/chat/completions',
    billing: '按次',
    prompt: '✓',
    aspectRatio: 'size 三档',
    duration: '6–20（video_config）',
    refImages: '≤7',
    frameTransition: '—',
    refVideo: '—',
    v2vOrEdit: '—',
  },
  {
    model: 'grok-imagine-video-cli',
    vendor: 'Grok CLI',
    api: 'POST /v1/videos',
    billing: '按次',
    prompt: '✓',
    aspectRatio: '7 种比例 / size',
    duration: '4–15',
    refImages: '≤10',
    frameTransition: '—',
    refVideo: '—',
    v2vOrEdit: '—',
  },
  {
    model: 'grok-imagine-video-cli-edit',
    vendor: 'Grok CLI',
    api: 'POST /v1/videos',
    billing: '按次',
    prompt: '✓',
    aspectRatio: '—',
    duration: '—',
    refImages: '—',
    frameTransition: '—',
    refVideo: '—',
    v2vOrEdit: '编辑（video，≤8.7s）',
  },
  {
    model: 'grok-imagine-video-1.5-cli',
    vendor: 'Grok CLI',
    api: 'POST /v1/videos',
    billing: '按次',
    prompt: '✓',
    aspectRatio: '16:9 / 9:16',
    duration: '4–15',
    refImages: '1（必填）',
    frameTransition: '—',
    refVideo: '—',
    v2vOrEdit: '—',
  },
]
