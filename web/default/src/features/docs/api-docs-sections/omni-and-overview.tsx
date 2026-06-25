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

export function OmniVideoSection(props: ApiDocsContext) {
  const { base } = props

  return (
    <DocsSection
      id='api-omni-video'
      title='Omni 视频生成'
      description='基于 Gemini Veo 的视频生成，支持文生视频、图生视频（最多 5 张参考图）、视频转视频。'
    >
      <DocsTable
        headers={['模型', '能力', '计费']}
        rows={[
          ['omni-fast', '文生视频 / 图生视频', '按次'],
          ['omni-fast-v2v', '视频转视频（V2V）', '按次'],
          ['omni-fast-no-water', '文生 / 图生（无水印）', '按次'],
          ['omni-fast-v2v-no-water', 'V2V（无水印）', '按次'],
        ]}
      />
      <p className='text-muted-foreground text-sm'>{pricingNote()}</p>
      <p>无水印模型输出经自动清洗，完成前可能多一个 processing 阶段，稍慢。</p>

      <DocsTable
        headers={['项', '说明']}
        rows={[
          ['提交任务', 'POST /v1/videos（JSON 或 multipart）'],
          ['轮询进度', 'GET /v1/videos/{task_id}'],
          ['下载成片', 'GET /v1/videos/{task_id}/content 或 data[0].url'],
          ['鉴权', 'Authorization: Bearer sk-你的令牌'],
        ]}
      />

      <DocsTable
        headers={['参数', '类型', '必填', '默认', '说明']}
        rows={[
          ['model', 'string', '是', '-', '见上表'],
          ['prompt', 'string', '是', '-', '视频描述'],
          ['aspect_ratio', 'string', '否', '16:9', '16:9（横）或 9:16（竖）'],
          ['seconds / duration', 'string/int', '否', '10', '时长（Gemini 固定输出约 10 秒）'],
          ['image_url', 'string', '否', '-', '单张参考图 URL 或 base64'],
          ['first_image_url', 'string', '否', '-', '首帧参考图'],
          ['last_image_url', 'string', '否', '-', '末帧参考图'],
          ['video_url', 'string', '否', '-', 'V2V 源视频 URL（≤5MB、1920×1080 内）'],
        ]}
      />

      <h3 className='text-lg font-semibold'>Multipart 文件上传</h3>
      <DocsTable
        headers={['字段', '说明']}
        rows={[
          ['input_reference', '参考图文件（最多 5 张，每张 ≤5MB）'],
          ['input_video', 'V2V 源视频文件（≤5MB）'],
        ]}
      />

      <CodeBlock
        title='文生视频'
        code={`curl -X POST ${base}/videos \\
  -H "Authorization: Bearer sk-xxx" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "omni-fast",
    "prompt": "雨夜霓虹街道，镜头缓慢推进，电影感光影",
    "aspect_ratio": "16:9"
  }'`}
      />
      <CodeBlock
        title='图生视频'
        code={`curl -X POST ${base}/videos \\
  -H "Authorization: Bearer sk-xxx" \\
  -H "Content-Type: application/json" \\
  -d '{
    "model": "omni-fast",
    "prompt": "保持人物一致，缓慢走动",
    "image_url": "https://your-cdn.com/photo.jpg",
    "aspect_ratio": "16:9"
  }'`}
      />
      <CodeBlock
        title='视频转视频（V2V）'
        code={`curl -X POST ${base}/videos \\
  -H "Authorization: Bearer sk-xxx" \\
  -F "model=omni-fast-v2v" \\
  -F "prompt=将画面风格转换为赛博朋克风" \\
  -F "input_video=@source.mp4"`}
      />
      <CodeBlock
        title='轮询取片'
        code={`curl ${base}/videos/{task_id} \\
  -H "Authorization: Bearer sk-xxx"

# 完成后: {"status":"completed","data":[{"url":"/v1/videos/{task_id}/content"}]}`}
      />
      <CodeBlock
        title='Python 完整示例'
        code={`import time, requests

BASE = "${base}"
H = {"Authorization": "Bearer sk-xxx", "Content-Type": "application/json"}

task = requests.post(f"{BASE}/videos", headers=H, json={
    "model": "omni-fast",
    "prompt": "雨夜霓虹街道，镜头缓慢推进",
    "aspect_ratio": "16:9"
}).json()
task_id = task["task_id"]

while True:
    time.sleep(8)
    s = requests.get(f"{BASE}/videos/{task_id}", headers=H).json()
    if s["status"] == "completed":
        print("下载:", s["data"][0]["url"])
        break
    if s["status"] == "failed":
        print("失败:", s.get("error"))
        break
    print(f"进度: {s.get('progress', 0)}%")`}
      />

      <ul className='list-disc space-y-2 pl-5'>
        <li>视频生成通常 1–5 分钟，轮询间隔建议 5–10 秒</li>
        <li>画幅仅支持 16:9 与 9:16，输出分辨率固定 720p</li>
        <li>含可识别真人面孔的参考图可能触发内容策略，见{' '}
          <a href='#api-video-guide' className='text-primary font-medium hover:underline'>
            视频避坑指南
          </a>
        </li>
      </ul>
    </DocsSection>
  )
}

export function VeoCleanSection(props: ApiDocsContext) {
  const { base } = props

  return (
    <DocsSection
      id='api-veo-clean'
      title='Veo-Clean 去水印'
      description='上传带水印的视频，系统自动去除水印后返回。异步任务流程与视频生成一致。'
    >
      <DocsTable
        headers={['模型', '计费']}
        rows={[['veo-clean', '按视频实际秒数计费']]}
      />
      <p className='text-muted-foreground text-sm'>{pricingNote()}</p>

      <DocsTable
        headers={['参数', '类型', '必填', '说明']}
        rows={[
          ['model', 'string', '是', '固定 veo-clean'],
          ['prompt', 'string', '否', '可省略，默认 remove watermark'],
          ['input_video', 'file', '是', '带水印视频（≤20MB，须 multipart 上传）'],
        ]}
      />

      <CodeBlock
        title='Multipart 上传'
        code={`curl -X POST ${base}/videos \\
  -H "Authorization: Bearer sk-xxx" \\
  -F "model=veo-clean" \\
  -F "prompt=remove watermark" \\
  -F "input_video=@watermarked.mp4"

curl ${base}/videos/{task_id} \\
  -H "Authorization: Bearer sk-xxx"`}
      />

      <ul className='list-disc space-y-2 pl-5'>
        <li>仅支持 multipart/form-data 提交</li>
        <li>处理时间通常 20–60 秒，取决于视频长度</li>
      </ul>
    </DocsSection>
  )
}

export function OverviewFaqSection(props: ApiDocsContext) {
  const { base, siteOrigin } = props

  return (
    <DocsSection
      id='api-overview-faq'
      title='通用说明 & FAQ'
      description='平台鉴权、错误码与常见问题。'
    >
      <DocsTable
        headers={['项', '说明']}
        rows={[
          ['Base URL', base],
          ['备用写法', siteOrigin || 'https://YOUR_BASE（部分客户端填根域名）'],
          ['鉴权', 'Authorization: Bearer sk-你的令牌'],
          ['模型名', '与模型广场展示名一致，如 omni-fast、video-pro-720p、grok-imagine-video'],
        ]}
      />
      <p>
        创建 API Key 见{' '}
        <Link to='/keys' className='text-primary font-medium hover:underline'>
          控制台 · API 密钥
        </Link>
        ；模型列表与单价见{' '}
        <Link to='/pricing' className='text-primary font-medium hover:underline'>
          模型广场
        </Link>
        。令牌分组须与模型匹配，错误分组会返回「无可用渠道」。
      </p>

      <h3 className='text-lg font-semibold'>错误码</h3>
      <DocsTable
        headers={['HTTP', '含义', '处理', '计费']}
        rows={[
          ['200', '成功', '正常取用', '成功才扣'],
          ['400', '参数/素材问题', '按 message 修正', '不计费'],
          ['401', '鉴权失败', '检查令牌', '不计费'],
          ['404', '路径错误', '检查 URL（勿重复 /v1）', '不计费'],
          ['429', '限速/额度不足', '降并发或充值', '不计费'],
          ['502/5xx', '上游临时故障', '稍后重试', '不计费'],
        ]}
      />

      <h3 className='mt-6 text-lg font-semibold'>FAQ</h3>
      <div className='space-y-4'>
        <div>
          <p className='font-medium'>Q: 视频生成需要多久？</p>
          <p className='text-muted-foreground'>
            Omni 约 1–5 分钟，Seedance 约 1–3 分钟，Grok 约 30s–3 分钟。建议客户端超时 ≥300 秒。
          </p>
        </div>
        <div>
          <p className='font-medium'>Q: 失败会扣费吗？</p>
          <p className='text-muted-foreground'>不会。仅成功出片/出图才计费。</p>
        </div>
        <div>
          <p className='font-medium'>Q: 异步视频怎么调用？</p>
          <p className='text-muted-foreground'>
            POST /v1/videos 提交 → GET /v1/videos/&#123;task_id&#125; 轮询 → 从 url / video_url 下载。
          </p>
        </div>
        <div>
          <p className='font-medium'>Q: 终端用户如何使用？</p>
          <p className='text-muted-foreground'>
            注册、充值、第三方工具配置见{' '}
            <Link to='/docs' className='text-primary font-medium hover:underline'>
              用户指南
            </Link>
            。
          </p>
        </div>
      </div>
    </DocsSection>
  )
}
