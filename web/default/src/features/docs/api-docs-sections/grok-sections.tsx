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

export function GrokImageVideoSection(props: ApiDocsContext) {
  const { base } = props

  return (
    <DocsSection
      id='api-grok-image-video'
      title='Grok 图像 & 视频'
      description='基于 xAI Grok Imagine 的图像与视频生成。视频参数见通用文档，图像走 Chat 接口。'
    >
      <DocsTable
        headers={['模型', '能力', '接口']}
        rows={[
          ['grok-imagine-image', '标准文生图', 'POST /v1/chat/completions'],
          ['grok-imagine-image-lite', '快速文生图', 'POST /v1/chat/completions'],
          ['grok-imagine-image-pro', '高质量文生图', 'POST /v1/chat/completions'],
          ['grok-imagine-image-edit', '图生图 / 编辑', 'POST /v1/chat/completions'],
          ['grok-imagine-video', '文生 / 图生视频', 'POST /v1/chat/completions（推荐）'],
        ]}
      />
      <p className='text-muted-foreground text-sm'>{pricingNote()}</p>
      <p className='text-muted-foreground text-sm'>
        图像参数见{' '}
        <a href='#api-image-api' className='text-primary font-medium hover:underline'>
          图像生成 API（通用）
        </a>
        ；视频参数见{' '}
        <a href='#api-video-api' className='text-primary font-medium hover:underline'>
          视频生成 API（通用）
        </a>
        。
      </p>

      <h3 className='text-lg font-semibold'>厂商特性</h3>
      <ul className='list-disc space-y-2 pl-5'>
        <li>视频提示词不超过 1500 字符</li>
        <li>参考图最多 7 张，支持 @IMAGE1…@IMAGE7 占位符</li>
        <li>stream: true 时推送「视频正在生成 NN%」进度</li>
        <li>客户端超时建议 ≥300 秒</li>
      </ul>

      <CodeBlock
        title='文生图'
        code={`curl ${base}/chat/completions \\
  -H "Authorization: Bearer sk-xxx" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "grok-imagine-image",
    "messages": [{"role":"user","content":"一只穿宇航服的橘猫，电影质感"}]
  }'`}
      />
    </DocsSection>
  )
}

export function GrokCliVideoSection(props: ApiDocsContext) {
  const { base } = props

  return (
    <DocsSection
      id='api-grok-cli-video'
      title='Grok CLI 视频专线'
      description='独立视频专线，稳定性更好。仅走 POST /v1/videos 异步接口，参数见通用文档。'
    >
      <DocsTable
        headers={['模型', '能力', '计费']}
        rows={[
          ['grok-imagine-video-cli', '文生 / 图生视频', '按次（一口价）'],
          ['grok-imagine-video-cli-edit', '视频局部编辑', '按次（一口价）'],
          ['grok-imagine-video-1.5-cli', '1.5 单图生视频', '按次（一口价）'],
        ]}
      />
      <p className='text-muted-foreground text-sm'>{pricingNote()} 一口价与时长、画幅无关。</p>

      <h3 className='text-lg font-semibold'>厂商特性</h3>
      <ul className='list-disc space-y-2 pl-5'>
        <li>
          <strong>grok-imagine-video-1.5-cli</strong> 仅支持单张首帧图生视频，须上传 1 张参考图，不支持纯文生
        </li>
        <li>
          <strong>grok-imagine-video-cli-edit</strong> 编辑模式：prompt 只写要改的内容；源视频 ≤8.7s、≤25MB、H.264 MP4
        </li>
        <li>文生/单图最长 15 秒；多参考图最长 10 秒</li>
      </ul>

      <CodeBlock
        title='视频编辑 · multipart'
        code={`curl ${base}/videos \\
  -H "Authorization: Bearer sk-xxx" \\
  -F "model=grok-imagine-video-cli-edit" \\
  -F "prompt=add a gold necklace" \\
  -F "video=@源视频.mp4"`}
      />
      <CodeBlock
        title='1.5 单图生视频'
        code={`curl ${base}/videos \\
  -H "Authorization: Bearer sk-xxx" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "grok-imagine-video-1.5-cli",
    "prompt": "gentle camera push-in, water flowing",
    "input_reference": "https://your-image.jpg",
    "seconds": "4",
    "size": "1280x720"
  }'`}
      />
    </DocsSection>
  )
}
