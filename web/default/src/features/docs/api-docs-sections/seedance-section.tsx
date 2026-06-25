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
      title='Seedance 2.0 视频生成'
      description='基于 Seedance 2.0 的视频生成。各模型调用方式一致，切换 model 即可。支持文生、图生、多参考图（≤9）、参考视频（≤3）、首尾帧过渡。'
    >
      <DocsTable
        headers={['模型', '版本', '定位', 'duration 范围']}
        rows={[
          ['video-pro-720p', '满血', '最佳质量，正式成片首选', '4–15（任意整数）'],
          ['video-fast-720p', '满血', '快速出片，质量接近 Pro', '4–15（任意整数）'],
          ['video-lite-720p', '非满血', '经济档，复杂真人一致性略弱', '仅 4 / 8 / 12'],
          ['Seedance2.0-480p', '480p', '经济标准档', '4–15（任意整数）'],
          ['Seedance2.0-fast-480p', '480p', '经济快速档', '4–15（任意整数）'],
          ['Seedance2.0-720p', '720p', '高清标准', '4–15（任意整数）'],
          ['Seedance2.0-fast-720p', '720p', '高清快速', '4–15（任意整数）'],
        ]}
      />
      <p className='text-muted-foreground text-sm'>{pricingNote()} 按秒计费：费用 = 单价 × duration。</p>
      <p>
        <strong>满血 vs 非满血：</strong>Pro/Fast 支持 4–15 秒任意整数及 @Image/@Video/@Audio 高级引用；Lite
        时长仅 4/8/12 秒三档，高级参考字段名不同。
      </p>

      <h3 className='text-lg font-semibold'>四种生成模式</h3>
      <p>无需传「模式」参数，服务端按素材字段自动判定：</p>
      <DocsTable
        headers={['模式', '最少必传', '如何触发']}
        rows={[
          ['文生视频', 'prompt', '只传 prompt，不带素材'],
          ['图生视频', 'prompt + ≥1 张图', 'image_url 和/或 reference_image_urls，不带视频'],
          ['全能参考', 'prompt + ≥1 图 + ≥1 视频', '传 reference_videos，并至少配 1 张图'],
          ['首尾帧', 'prompt + first + last（成对）', '同时传 first_image_url 与 last_image_url'],
        ]}
      />

      <DocsTable
        headers={['参数', '类型', '说明']}
        rows={[
          ['model', 'string', '上表模型名之一'],
          ['prompt', 'string', '视频描述；多素材用 @image1/@video1 引用'],
          ['aspect_ratio', 'string', '16:9、9:16、1:1、21:9、3:4、4:3'],
          ['duration', 'integer', 'Pro/Fast/Seedance2.0: 4–15；Lite: 仅 4/8/12'],
          ['image_url', 'string', '主参考图 URL / base64 / multipart'],
          ['reference_image_urls', 'array', '多参考图，与 image_url 合计 ≤9'],
          ['reference_videos', 'array', '参考视频 ≤3（mp4/mov，2–15s，≤50MB）'],
          ['first_image_url / last_image_url', 'string', '首尾帧过渡（须成对）'],
          ['extra_images / extra_videos / extra_audios', 'array', 'Pro/Fast 高级引用，@Image1…@Video1…'],
        ]}
      />

      <h3 className='text-lg font-semibold'>常见错误码</h3>
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

      <CodeBlock
        title='文生视频'
        code={`curl -X POST ${base}/videos \\
  -H "Authorization: Bearer sk-xxx" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "video-pro-720p",
    "prompt": "雨夜霓虹街道，镜头缓慢推进，电影感光影",
    "aspect_ratio": "16:9",
    "duration": 8
  }'`}
      />
      <CodeBlock
        title='480p 图生 · base64 直传'
        code={`curl -X POST ${base}/videos \\
  -H "Authorization: Bearer sk-xxx" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "Seedance2.0-fast-480p",
    "prompt": "让画面动起来",
    "duration": 5,
    "image_url": "data:image/png;base64,iVBORw0KGgo..."
  }'`}
      />
      <CodeBlock
        title='480p 图生 · multipart'
        code={`curl -X POST ${base}/videos \\
  -H "Authorization: Bearer sk-xxx" \\
  -F "model=Seedance2.0-fast-480p" \\
  -F "prompt=让画面动起来" \\
  -F "duration=5" \\
  -F "image=@/path/to/photo.jpg"`}
      />
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

      <ul className='list-disc space-y-2 pl-5'>
        <li>输出 H.264 / 24fps，含 AAC 立体声，无水印</li>
        <li>下载地址在轮询响应的 video_url 字段</li>
        <li>首尾帧模式不接受额外 reference_image_urls</li>
      </ul>
    </DocsSection>
  )
}
