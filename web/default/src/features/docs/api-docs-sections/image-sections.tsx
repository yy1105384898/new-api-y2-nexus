/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/
import { CodeBlock } from '../components/code-block'
import { DocsSection } from '../components/docs-section'
import { DocsTable } from '../components/docs-table'
import type { ApiDocsContext } from './context'
import { pricingNote } from './context'

export function GptImageSection(props: ApiDocsContext) {
  const { base } = props

  return (
    <DocsSection
      id='api-gpt-image'
      title='GPT-Image 系列'
      description='OpenAI GPT-Image 生图模型。参数与调用方式见通用文档，本节列可用档位与特性。'
    >
      <p className='text-muted-foreground text-sm'>
        完整参数说明与能力对照见{' '}
        <a href='#api-image-api' className='text-primary font-medium hover:underline'>
          图像生成 API（通用）
        </a>
        。
      </p>

      <DocsTable
        headers={['模型', '说明', '接口']}
        rows={[
          ['gpt-image-2', '标准文生图 / 图生图', '/v1/images/generations · /v1/images/edits'],
          ['gpt-image-2-4k', '4K 超清文生图（异步）', 'POST /v1/videos → 轮询'],
        ]}
      />
      <p className='text-muted-foreground text-sm'>{pricingNote()}</p>

      <h3 className='text-lg font-semibold'>厂商特性</h3>
      <ul className='list-disc space-y-2 pl-5'>
        <li>gpt-image-2 支持 generations（JSON）与 edits（multipart 图生图）</li>
        <li>gpt-image-2-4k 固定 2880×2880 输出，典型耗时 1–3 分钟，从轮询响应 video_url 取图</li>
        <li>size：1024x1024、1536x1024、1024x1536、auto</li>
      </ul>

      <CodeBlock
        title='4K 异步生图'
        code={`curl -X POST ${base}/videos \\
  -H "Authorization: Bearer sk-xxx" \\
  -H "Content-Type: application/json" \\
  -d '{"model":"gpt-image-2-4k","prompt":"一只橘猫坐在窗台上"}'`}
      />
    </DocsSection>
  )
}

export function GeminiImageSection(props: ApiDocsContext) {
  const { base } = props

  return (
    <DocsSection
      id='api-gemini-image'
      title='Gemini 图像系列'
      description='基于 Gemini 的图像生成。参数见通用文档，本节列可用模型与特性。'
    >
      <p className='text-muted-foreground text-sm'>
        完整参数说明与能力对照见{' '}
        <a href='#api-image-api' className='text-primary font-medium hover:underline'>
          图像生成 API（通用）
        </a>
        。
      </p>

      <DocsTable
        headers={['模型', '说明']}
        rows={[
          ['gemini-image', '标准生图'],
          ['gemini-image-pro', '高质量生图'],
        ]}
      />
      <p className='text-muted-foreground text-sm'>{pricingNote()}</p>

      <h3 className='text-lg font-semibold'>厂商特性</h3>
      <ul className='list-disc space-y-2 pl-5'>
        <li>参考图最多 5 张、每张 ≤5MB；image 传单张用字符串、多张用数组均可</li>
        <li>支持 mask 蒙版局部重绘</li>
      </ul>

      <CodeBlock
        title='多参考图融合'
        code={`curl -X POST ${base}/images/generations \\
  -H "Authorization: Bearer sk-xxx" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "gemini-image",
    "prompt": "把这两张图融合成一张海报",
    "image": ["https://cdn.example.com/a.jpg", "https://cdn.example.com/b.jpg"]
  }'`}
      />
    </DocsSection>
  )
}

export function GeminiMusicSection(props: ApiDocsContext) {
  const { base } = props

  return (
    <DocsSection
      id='api-gemini-music'
      title='Gemini 音乐生成'
      description='通过 Chat Completions 接口生成音乐。'
    >
      <DocsTable headers={['模型', '说明']} rows={[['gemini-music', '按次计费，见模型广场']]} />
      <p className='text-muted-foreground text-sm'>{pricingNote()}</p>
      <CodeBlock
        title='示例'
        code={`curl ${base}/chat/completions \\
  -H "Authorization: Bearer sk-xxx" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "gemini-music",
    "messages": [{"role":"user","content":"创作一首轻快的电子风格BGM，适合科技产品广告"}]
  }'`}
      />
    </DocsSection>
  )
}
