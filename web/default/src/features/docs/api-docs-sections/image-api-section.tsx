/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/
import { Link } from '@tanstack/react-router'
import { CodeBlock } from '../components/code-block'
import { DocsSection } from '../components/docs-section'
import { DocsTable } from '../components/docs-table'
import type { ApiDocsContext } from './context'
import { pricingNote } from './context'
import {
  asyncImageParams,
  chatImageParams,
  imageEditsFields,
  imageGenerationsParams,
  imageModelCapabilities,
} from './image-api-data'

export function UnifiedImageApiSection(props: ApiDocsContext) {
  const { base } = props

  return (
    <DocsSection
      id='api-image-api'
      title='图像生成 API（通用）'
      description='多上游渠道共用 OpenAI 兼容图像接口。先查下方模型能力对照表，再按需传参；不支持的字段会被忽略或返回 400。'
    >
      <p className='text-muted-foreground text-sm'>{pricingNote()}</p>

      <h3 className='text-lg font-semibold'>可用端点</h3>
      <DocsTable
        headers={['端点', 'Content-Type', '用途']}
        rows={[
          ['POST /v1/images/generations', 'application/json', '文生图 / 带参考图的 JSON 生图'],
          ['POST /v1/images/edits', 'multipart/form-data', '图生图 / 编辑（上传参考图文件）'],
          ['POST /v1/chat/completions', 'application/json', 'Grok Imagine、Gemini Banana 等 Chat 生图'],
          ['POST /v1/videos → 轮询', 'application/json', 'gpt-image-2-4k 等异步超清生图'],
        ]}
      />
      <p>
        鉴权：<code className='text-sm'>Authorization: Bearer sk-你的令牌</code>。模型名与{' '}
        <Link to='/pricing' className='text-primary font-medium hover:underline'>
          模型广场
        </Link>{' '}
        展示名一致。
      </p>

      <h3 className='mt-8 text-lg font-semibold'>POST /v1/images/generations 参数（全集）</h3>
      <DocsTable
        headers={['参数', '类型', '必填', '默认', '说明']}
        rows={imageGenerationsParams.map((p) => [p.name, p.type, p.required, p.default, p.note])}
      />

      <h3 className='mt-8 text-lg font-semibold'>POST /v1/images/edits（multipart）</h3>
      <DocsTable
        headers={['字段', '说明']}
        rows={imageEditsFields.map((f) => [f.name, f.note])}
      />

      <h3 className='mt-8 text-lg font-semibold'>POST /v1/chat/completions（Grok 生图）</h3>
      <DocsTable
        headers={['参数', '类型', '必填', '说明']}
        rows={chatImageParams.map((p) => [p.name, p.type, p.required, p.note])}
      />

      <h3 className='mt-8 text-lg font-semibold'>异步 4K 生图（POST /v1/videos）</h3>
      <p className='text-muted-foreground mb-4 text-sm'>
        <code>gpt-image-2-4k</code> 走视频异步接口提交，轮询完成后从 <code>video_url</code> 取 4K 图片直链。
      </p>
      <DocsTable
        headers={['参数', '类型', '必填', '说明']}
        rows={asyncImageParams.map((p) => [p.name, p.type, p.required, p.note])}
      />

      <h3 className='mt-8 text-lg font-semibold'>模型能力对照表</h3>
      <p className='text-muted-foreground mb-4 text-sm'>
        多上游接入同一套 API；切换 model 即可换供应商/档位。其他已在模型广场上架的生图模型（Banana、Flux 等）同样走本接口，参数规则与相近档位一致。
      </p>
      <DocsTable
        headers={[
          '模型',
          '供应商',
          '接口',
          '计费',
          'prompt',
          '尺寸',
          'quality',
          '数量',
          '参考图',
          '蒙版',
          '异步',
        ]}
        rows={imageModelCapabilities.map((m) => [
          m.model,
          m.vendor,
          m.api.replace('POST ', ''),
          m.billing,
          m.prompt,
          m.size,
          m.quality,
          m.count,
          m.refImages,
          m.mask,
          m.asyncTask,
        ])}
      />

      <h3 className='mt-8 text-lg font-semibold'>示例</h3>
      <CodeBlock
        title='文生图'
        code={`curl -X POST ${base}/images/generations \\
  -H "Authorization: Bearer sk-xxx" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "gpt-image-2",
    "prompt": "一只橘猫趴在窗台上晒太阳，水彩画风格",
    "size": "1024x1024"
  }'`}
      />
      <CodeBlock
        title='图生图 · multipart'
        code={`curl -X POST ${base}/images/edits \\
  -H "Authorization: Bearer sk-xxx" \\
  -F "image=@reference.png" \\
  -F "prompt=把背景改成海边日落" \\
  -F "model=gpt-image-2"`}
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
        <li>标准档响应 15–60 秒，4K 异步 1–3 分钟；客户端超时建议 ≥120 秒（4K ≥300 秒）</li>
        <li>response_format: url 返回完整图片地址，可直接使用</li>
        <li>仅成功出图才计费</li>
      </ul>
    </DocsSection>
  )
}
