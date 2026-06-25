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

export function SeedanceVideoSection(props: ApiDocsContext) {
  const { base } = props

  return (
    <DocsSection
      id='api-seedance-video'
      title='Seedance 2.0 系列'
      description='多上游渠道接入的 Seedance 2.0 视频模型。参数与调用方式见通用文档，本节列可用档位与厂商特性。'
    >
      <p className='text-muted-foreground text-sm'>
        完整参数说明与能力对照见{' '}
        <a href='#api-video-api' className='text-primary font-medium hover:underline'>
          视频生成 API（通用）
        </a>
        。
      </p>

      <DocsTable
        headers={['模型', '分辨率', '定位', 'duration']}
        rows={[
          ['Seedance2.0-480p', '480p', '经济标准档', '4–15'],
          ['Seedance2.0-fast-480p', '480p', '经济快速档', '4–15'],
          ['Seedance2.0-720p', '720p', '高清标准', '4–15'],
          ['Seedance2.0-fast-720p', '720p', '高清快速', '4–15'],
        ]}
      />
      <p className='text-muted-foreground text-sm'>{pricingNote()} 按秒计费：费用 = 单价 × duration。</p>

      <h3 className='text-lg font-semibold'>厂商特性</h3>
      <ul className='list-disc space-y-2 pl-5'>
        <li>支持文生、图生、多参考图（≤9）、参考视频（≤3）、首尾帧过渡</li>
        <li>画幅：16:9、9:16、1:1、21:9、3:4、4:3</li>
        <li>输出 H.264 / 24fps，含 AAC 立体声，无水印</li>
        <li>首尾帧模式不接受额外 reference_image_urls</li>
        <li>模型名含 480p/720p 时分辨率档位已锁定，无需再传 resolution</li>
      </ul>

      <h3 className='mt-6 text-lg font-semibold'>常见错误码</h3>
      <DocsTable
        headers={['error_code', '含义 / 处理']}
        rows={[
          ['400017', '参数或参考图不合规——按提示修正后重试'],
          ['500341', '参考视频不符合要求——更换视频后重试'],
          ['GENERATION_FAILED', '生成失败或内容策略拦截——更换图片或调整提示词'],
          ['TIMEOUT', '生成超时——稍后重试'],
          ['PROMPT_BLOCKED', '提示词违禁——修改提示词（不消耗额度）'],
        ]}
      />

      <h3 className='mt-6 text-lg font-semibold'>示例</h3>
      <CodeBlock
        title='多参考图 + 参考视频'
        code={`curl -X POST ${base}/videos \\
  -H "Authorization: Bearer sk-xxx" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "Seedance2.0-fast-480p",
    "prompt": "把 @image1 的人物换进 @video1 的画面",
    "image_url": "https://cdn.example.com/person.jpg",
    "reference_videos": ["https://cdn.example.com/ref.mp4"],
    "duration": 5
  }'`}
      />
      <CodeBlock
        title='首尾帧过渡'
        code={`curl -X POST ${base}/videos \\
  -H "Authorization: Bearer sk-xxx" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "Seedance2.0-fast-480p",
    "prompt": "平滑电影感过渡",
    "first_image_url": "https://cdn.example.com/start.jpg",
    "last_image_url": "https://cdn.example.com/end.jpg",
    "duration": 5
  }'`}
      />
    </DocsSection>
  )
}
