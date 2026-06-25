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
      description='基于 xAI Grok Imagine 的图像与视频生成，封装为 OpenAI 兼容接口。'
    >
      <DocsTable
        headers={['模型', '能力']}
        rows={[
          ['grok-imagine-image', '标准文生图'],
          ['grok-imagine-image-lite', '快速文生图'],
          ['grok-imagine-image-pro', '高质量文生图'],
          ['grok-imagine-image-edit', '图生图 / 编辑'],
          ['grok-imagine-video', '文生 / 图生视频'],
        ]}
      />
      <p className='text-muted-foreground text-sm'>{pricingNote()}</p>

      <DocsTable
        headers={['项', '说明']}
        rows={[
          ['统一接口', 'POST /v1/chat/completions（按模型自动分发图像/视频）'],
          ['视频异步', '也支持 POST /v1/videos（multipart + 轮询）'],
          ['鉴权', 'Authorization: Bearer sk-xxx'],
        ]}
      />

      <h3 className='text-lg font-semibold'>video_config 参数</h3>
      <DocsTable
        headers={['字段', '取值', '默认', '说明']}
        rows={[
          ['seconds', '6 / 10 / 12 / 16 / 20', '6', '时长，最长 20 秒'],
          ['size', '720x1280 / 1280x720 / 960x960', '720x1280', '画幅（仅 3 种）'],
          ['public_url', 'true / false', 'false', '建议 true，返回完整下载链接'],
        ]}
      />

      <CodeBlock
        title='文生视频'
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
      <CodeBlock
        title='图生视频'
        code={`curl ${base}/chat/completions \\
  -H "Authorization: Bearer sk-xxx" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "grok-imagine-video",
    "stream": true,
    "messages": [{"role":"user","content":[
      {"type":"text","text":"让画面动起来，光线柔和"},
      {"type":"image_url","image_url":{"url":"https://your-image.jpg"}}
    ]}],
    "video_config": {"seconds": 6, "size": "1280x720", "public_url": true}
  }'`}
      />
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

      <ul className='list-disc space-y-2 pl-5'>
        <li>视频提示词不超过 1500 字符</li>
        <li>参考图最多 7 张，支持 @IMAGE1…@IMAGE7 占位符</li>
        <li>stream: true 时推送「视频正在生成 NN%」进度</li>
        <li>客户端超时建议 ≥300 秒</li>
      </ul>
    </DocsSection>
  )
}

export function GrokCliVideoSection(props: ApiDocsContext) {
  const { base } = props

  return (
    <DocsSection
      id='api-grok-cli-video'
      title='Grok CLI 视频专线'
      description='独立视频专线，官方通道稳定性更好。仅走 POST /v1/videos 异步接口。'
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
      <p>
        <strong>grok-imagine-video-1.5-cli</strong> 仅支持单张首帧图生视频，不支持文生、多参考图与视频编辑。
      </p>

      <DocsTable
        headers={['参数', '取值', '说明']}
        rows={[
          ['prompt', '文本', '必填'],
          ['seconds', '整数', '文生/单图最长 15 秒；多图最长 10 秒'],
          ['size', '720x1280 / 1280x720', '画幅'],
          ['aspect_ratio', '1:1/16:9/9:16/4:3/3:4/3:2/2:3', '可替代 size'],
          ['resolution', '480p / 720p', '默认 720p'],
          ['input_reference', 'URL / data URL', '单张首帧参考图'],
          ['reference_images', '数组 ≤10', '多参考图（与 input_reference 二选一）'],
        ]}
      />

      <h3 className='text-lg font-semibold'>编辑参数（grok-imagine-video-cli-edit）</h3>
      <DocsTable
        headers={['参数', '说明']}
        rows={[
          ['prompt', '只写要改的内容，如 add a gold necklace'],
          ['video', '源视频：公网 URL、base64 或 multipart 上传（≤8.7s、≤25MB、H.264 MP4）'],
        ]}
      />

      <CodeBlock
        title='文生视频'
        code={`curl ${base}/videos \\
  -H "Authorization: Bearer sk-xxx" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "grok-imagine-video-cli",
    "prompt": "hot air balloon rising over green hills at sunrise",
    "seconds": "10",
    "size": "1280x720"
  }'`}
      />
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
