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
      description='Gemini 原生图像模型。Imagen 类走 OpenAI 图像接口；Banana 2 / Banana Pro 走 Chat 或 Gemini 原生 generateContent。'
    >
      <p className='text-muted-foreground text-sm'>
        完整参数说明与能力对照见{' '}
        <a href='#api-image-api' className='text-primary font-medium hover:underline'>
          图像生成 API（通用）
        </a>
        。模型名与模型广场展示名一致；计费与可用分组以模型广场为准。
      </p>

      <DocsTable
        headers={['模型', '说明', '接口']}
        rows={[
          ['gemini-image', '标准 Gemini 生图（Imagen）', 'POST /v1/images/generations'],
          ['gemini-image-pro', '高质量 Gemini 生图（Imagen）', 'POST /v1/images/generations'],
          ['gemini-banana-2.0', 'Banana 2 标准档', 'POST /v1/chat/completions · POST /v1beta/models/{model}:generateContent'],
          ['gemini-banana-2.0-pro', 'Banana Pro 高质量档', 'POST /v1/chat/completions · POST /v1beta/models/{model}:generateContent'],
        ]}
      />
      <p className='text-muted-foreground text-sm'>{pricingNote()}</p>

      <h3 className='text-lg font-semibold'>Banana 系列（gemini-banana-*）</h3>
      <p className='text-muted-foreground mb-4 text-sm'>
        Banana 上游为 Gemini generateContent 协议，<strong>不支持</strong> POST /v1/images/generations。请使用下列两种方式之一，响应格式与上游一致（Chat 返回
        Markdown 图片，Gemini 原生返回 candidates.parts 中的 inlineData）。
      </p>
      <DocsTable
        headers={['端点', 'Content-Type', '说明']}
        rows={[
          [
            'POST /v1beta/models/{model}:generateContent',
            'application/json',
            '推荐。Gemini 原生透传；model 用模型广场展示名，如 gemini-banana-2.0',
          ],
          [
            'POST /v1/chat/completions',
            'application/json',
            'OpenAI Chat 兼容；图片在 choices[0].message.content 的 Markdown 中返回',
          ],
        ]}
      />

      <h3 className='mt-8 text-lg font-semibold'>Imagen 系列（gemini-image*）</h3>
      <ul className='list-disc space-y-2 pl-5'>
        <li>走 POST /v1/images/generations，返回 OpenAI 图像 JSON（b64_json / url）</li>
        <li>参考图最多 5 张、每张 ≤5MB；image 传单张字符串、多张用数组</li>
        <li>支持 mask 蒙版局部重绘</li>
      </ul>

      <CodeBlock
        title='Banana — Gemini 原生 generateContent（推荐）'
        code={`curl -X POST ${base}/v1beta/models/gemini-banana-2.0:generateContent \\
  -H "Authorization: Bearer sk-xxx" \\
  -H "Content-Type: application/json" \\
  -d '{
    "contents": [{
      "role": "user",
      "parts": [{"text": "一只红苹果，白底产品图"}]
    }],
    "generationConfig": {
      "responseModalities": ["TEXT", "IMAGE"],
      "imageConfig": {
        "aspectRatio": "1:1"
      }
    }
  }'`}
      />

      <CodeBlock
        title='Banana — Chat Completions'
        code={`curl -X POST ${base}/chat/completions \\
  -H "Authorization: Bearer sk-xxx" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "gemini-banana-2.0",
    "messages": [{"role": "user", "content": "一只红苹果，白底产品图"}],
    "stream": false
  }'`}
      />

      <CodeBlock
        title='Banana — 图生图（Chat，参考图放 content）'
        code={`curl -X POST ${base}/chat/completions \\
  -H "Authorization: Bearer sk-xxx" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "gemini-banana-2.0-pro",
    "messages": [{
      "role": "user",
      "content": [
        {"type": "text", "text": "把这两张图融合成一张海报"},
        {"type": "image_url", "image_url": {"url": "https://cdn.example.com/a.jpg"}},
        {"type": "image_url", "image_url": {"url": "https://cdn.example.com/b.jpg"}}
      ]
    }],
    "stream": false
  }'`}
      />

      <CodeBlock
        title='Imagen — OpenAI 图像接口'
        code={`curl -X POST ${base}/images/generations \\
  -H "Authorization: Bearer sk-xxx" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "gemini-image",
    "prompt": "一只红苹果，白底产品图",
    "size": "1024x1024"
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
