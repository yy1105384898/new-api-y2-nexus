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
      title='GPT-Image-2'
      description='文生图 / 图生图 / Chat 生图，支持多种调用方式。'
    >
      <DocsTable
        headers={['模型', '说明']}
        rows={[
          ['gpt-image-2', '标准文生图 / 图生图'],
          ['gpt-image-2-4k', '4K 超清文生图（异步任务）'],
        ]}
      />
      <p className='text-muted-foreground text-sm'>{pricingNote()}</p>

      <DocsTable
        headers={['端点', '方式', '说明']}
        rows={[
          ['/v1/images/generations', 'JSON', '文生图'],
          ['/v1/images/edits', 'multipart', '图生图（参考图 + 描述）'],
          ['/v1/chat/completions', 'JSON', 'Chat 对话生图'],
        ]}
      />

      <DocsTable
        headers={['参数', '类型', '说明']}
        rows={[
          ['prompt', 'string', '图片描述（必填）'],
          ['model', 'string', '默认 gpt-image-2'],
          ['n', 'integer', '生成数量 1–4'],
          ['size', 'string', '1024x1024、1536x1024、1024x1536、auto'],
          ['quality', 'string', 'auto / low / medium / high'],
          ['response_format', 'string', 'b64_json（默认）或 url'],
        ]}
      />

      <CodeBlock
        title='文生图'
        code={`curl -X POST ${base}/images/generations \\
  -H "Authorization: Bearer sk-xxx" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "gpt-image-2",
    "prompt": "一只橘猫趴在窗台上晒太阳，水彩画风格",
    "size": "1024x1024",
    "quality": "high"
  }'`}
      />
      <CodeBlock
        title='图生图'
        code={`curl -X POST ${base}/images/edits \\
  -H "Authorization: Bearer sk-xxx" \\
  -F "image=@reference.png" \\
  -F "prompt=把背景改成海边日落" \\
  -F "model=gpt-image-2"`}
      />

      <h3 className='mt-8 text-lg font-semibold'>gpt-image-2-4k（异步 4K）</h3>
      <p>4K 超清文生图，异步任务制，出图分辨率 2880×2880，典型耗时 1–3 分钟。</p>
      <DocsTable
        headers={['步骤', '端点', '说明']}
        rows={[
          ['① 提交', 'POST /v1/videos', '返回 {"id":"task_xxx","status":"queued"}'],
          ['② 轮询', 'GET /v1/videos/{id}', 'queued → in_progress → completed'],
          ['③ 取图', '—', '读取 video_url 字段（4K 图片直链）'],
        ]}
      />
      <CodeBlock
        title='4K 异步生图'
        code={`curl -X POST ${base}/videos \\
  -H "Authorization: Bearer sk-xxx" \\
  -H "Content-Type: application/json" \\
  -d '{"model":"gpt-image-2-4k","prompt":"一只橘猫坐在窗台上"}'

curl ${base}/videos/task_xxxx \\
  -H "Authorization: Bearer sk-xxx"
# → {"status":"completed","video_url":"https://.../xxxx.png"}`}
      />

      <ul className='list-disc space-y-2 pl-5'>
        <li>响应时间 15–60 秒（标准档），超时建议 ≥120 秒</li>
        <li>response_format: url 返回完整图片地址，可直接使用</li>
      </ul>
    </DocsSection>
  )
}

export function GeminiImageSection(props: ApiDocsContext) {
  const { base } = props

  return (
    <DocsSection
      id='api-gemini-image'
      title='Gemini 图像生成'
      description='基于 Gemini 的图像生成服务。'
    >
      <DocsTable
        headers={['模型', '说明']}
        rows={[
          ['gemini-image', '标准生图'],
          ['gemini-image-pro', '高质量生图'],
        ]}
      />
      <p className='text-muted-foreground text-sm'>{pricingNote()}</p>

      <CodeBlock
        title='文生图'
        code={`curl -X POST ${base}/images/generations \\
  -H "Authorization: Bearer sk-xxx" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "gemini-image",
    "prompt": "赛博朋克风格的东京夜景",
    "size": "1024x1024"
  }'`}
      />

      <h3 className='text-lg font-semibold'>参考图（图生图 / 编辑）</h3>
      <DocsTable
        headers={['字段', '类型', '说明']}
        rows={[
          ['image', 'string 或 string[]', '单张或多张参考图'],
          ['images', 'string[]', '多张参考图数组（与 image 等效）'],
          ['mask', 'string', '蒙版图，局部重绘可选'],
        ]}
      />
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
      <p className='text-muted-foreground text-sm'>参考图最多 5 张、每张 ≤5MB。image 传单张用字符串、多张用数组均可。</p>
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
