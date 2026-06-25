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
  chatVideoParams,
  videoApiParams,
  videoModelCapabilities,
  videoMultipartFields,
} from './video-api-data'

export function UnifiedVideoApiSection(props: ApiDocsContext) {
  const { base } = props

  return (
    <DocsSection
      id='api-video-api'
      title='视频生成 API（通用）'
      description='所有视频模型共用同一套 OpenAI 兼容接口。先查下方模型能力对照表确认该模型支持哪些参数，再按需传参；不支持的字段会被忽略或返回 400。'
    >
      <p className='text-muted-foreground text-sm'>{pricingNote()}</p>

      <h3 className='text-lg font-semibold'>调用流程</h3>
      <DocsTable
        headers={['步骤', '方法', '说明']}
        rows={[
          ['1. 提交任务', 'POST /v1/videos', 'JSON 或 multipart/form-data'],
          ['2. 轮询进度', 'GET /v1/videos/{task_id}', 'status: queued / in_progress / completed / failed'],
          ['3. 下载成片', 'GET /v1/videos/{task_id}/content', '或取响应 data[0].url / video_url'],
        ]}
      />
      <p>
        鉴权：<code className='text-sm'>Authorization: Bearer sk-你的令牌</code>。模型名与{' '}
        <Link to='/pricing' className='text-primary font-medium hover:underline'>
          模型广场
        </Link>{' '}
        展示名一致。
      </p>

      <h3 className='mt-8 text-lg font-semibold'>POST /v1/videos 参数（全集）</h3>
      <p className='text-muted-foreground mb-4 text-sm'>
        以下为平台接受的全部字段；具体模型是否支持见下一节对照表。带 * 表示多数模型必填，特殊模型（如仅图生的
        1.5-cli）以对照表为准。
      </p>
      <DocsTable
        headers={['参数', '类型', '必填', '默认', '说明']}
        rows={videoApiParams.map((p) => [p.name, p.type, p.required, p.default, p.note])}
      />

      <h3 className='mt-8 text-lg font-semibold'>Multipart 字段</h3>
      <DocsTable
        headers={['字段', '说明']}
        rows={videoMultipartFields.map((f) => [f.name, f.note])}
      />

      <h3 className='mt-8 text-lg font-semibold'>POST /v1/chat/completions（Grok 视频）</h3>
      <p className='text-muted-foreground mb-4 text-sm'>
        <code>grok-imagine-video</code> 走 Chat 接口而非 /videos；也支持 /videos 异步，但 Chat 流式体验更好。
      </p>
      <DocsTable
        headers={['参数', '类型', '必填', '说明']}
        rows={chatVideoParams.map((p) => [p.name, p.type, p.required, p.note])}
      />

      <h3 className='mt-8 text-lg font-semibold'>模型能力对照表</h3>
      <p className='text-muted-foreground mb-4 text-sm'>
        多上游渠道接入同一套 API；切换 model 即可换供应商/档位。表中 ✓ 表示支持；— 表示不支持或无需传；数字为上限或范围。
        其他已在模型广场上架的视频模型（GZ Seedance、云雾 Veo/Sora、HappyHorse 等）同样走本接口，参数规则与相近档位一致。
      </p>
      <DocsTable
        headers={[
          '模型',
          '供应商',
          '接口',
          '计费',
          'prompt',
          '画幅',
          '时长',
          '参考图',
          '首尾帧',
          '参考视频',
          'V2V/编辑',
        ]}
        rows={videoModelCapabilities.map((m) => [
          m.model,
          m.vendor,
          m.api.replace('POST ', ''),
          m.billing,
          m.prompt,
          m.aspectRatio,
          m.duration,
          m.refImages,
          m.frameTransition,
          m.refVideo,
          m.v2vOrEdit,
        ])}
      />

      <h3 className='mt-8 text-lg font-semibold'>生成模式（自动判定）</h3>
      <p>无需传 mode 参数；服务端按所传素材字段自动选择文生 / 图生 / 全能参考 / 首尾帧 / V2V：</p>
      <DocsTable
        headers={['模式', '最少必传', '如何触发']}
        rows={[
          ['文生视频', 'prompt', '只传 prompt，不带素材'],
          ['图生视频', 'prompt + ≥1 张图', 'image_url / reference_image_urls / input_reference'],
          ['全能参考', 'prompt + ≥1 图 + ≥1 视频', 'reference_videos + 至少 1 张图（Seedance）'],
          ['首尾帧', 'prompt + 首帧 + 尾帧', 'first_image_url + last_image_url 成对'],
          ['视频转视频', 'prompt + 源视频', 'video_url / input_video / video（视模型）'],
        ]}
      />

      <h3 className='mt-8 text-lg font-semibold'>示例</h3>
      <CodeBlock
        title='文生视频（通用）'
        code={`curl -X POST ${base}/videos \\
  -H "Authorization: Bearer sk-xxx" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "Seedance2.0-720p",
    "prompt": "雨夜霓虹街道，镜头缓慢推进，电影感光影",
    "aspect_ratio": "16:9",
    "duration": 8
  }'`}
      />
      <CodeBlock
        title='图生视频 · multipart'
        code={`curl -X POST ${base}/videos \\
  -H "Authorization: Bearer sk-xxx" \\
  -F "model=omni-fast" \\
  -F "prompt=保持人物一致，缓慢走动" \\
  -F "aspect_ratio=16:9" \\
  -F "input_reference=@photo.jpg"`}
      />
      <CodeBlock
        title='轮询取片'
        code={`curl ${base}/videos/{task_id} \\
  -H "Authorization: Bearer sk-xxx"

# completed: {"status":"completed","data":[{"url":"/v1/videos/{task_id}/content"}]}`}
      />
      <CodeBlock
        title='Grok Chat 视频'
        code={`curl ${base}/chat/completions \\
  -H "Authorization: Bearer sk-xxx" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "grok-imagine-video",
    "stream": true,
    "messages": [{"role":"user","content":"灯塔在日落时分，海浪拍打礁石"}],
    "video_config": {"seconds": 10, "size": "1280x720", "public_url": true}
  }'`}
      />

      <ul className='list-disc space-y-2 pl-5'>
        <li>视频生成通常 30 秒–5 分钟，轮询间隔建议 5–10 秒，客户端超时 ≥300 秒</li>
        <li>仅成功出片才计费；失败不扣费</li>
        <li>
          Gemini 系列内容审查较严，见{' '}
          <a href='#api-video-guide' className='text-primary font-medium hover:underline'>
            视频避坑指南
          </a>
        </li>
      </ul>
    </DocsSection>
  )
}
