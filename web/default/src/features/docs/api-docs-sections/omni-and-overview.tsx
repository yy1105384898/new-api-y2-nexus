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

function videoApiLink() {
  return (
    <a href='#api-video-api' className='text-primary font-medium hover:underline'>
      视频生成 API（通用）
    </a>
  )
}

export function OmniVideoSection(props: ApiDocsContext) {
  const { base } = props

  return (
    <DocsSection
      id='api-omni-video'
      title='Gemini Veo · Omni 系列'
      description='基于 Gemini Veo 的视频生成。参数与调用方式见通用文档，本节仅列可用模型与厂商特性。'
    >
      <p className='text-muted-foreground text-sm'>
        完整参数说明与能力对照见 {videoApiLink()}。
      </p>

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

      <h3 className='text-lg font-semibold'>厂商特性</h3>
      <ul className='list-disc space-y-2 pl-5'>
        <li>画幅仅 16:9 与 9:16，输出固定 720p、约 10 秒</li>
        <li>图生最多 5 张参考图；V2V 源视频 ≤5MB、1920×1080 内</li>
        <li>无水印版输出经自动清洗，完成前可能多一个 processing 阶段，稍慢</li>
        <li>含可识别真人面孔的参考图可能触发内容策略，见{' '}
          <a href='#api-video-guide' className='text-primary font-medium hover:underline'>
            视频避坑指南
          </a>
        </li>
      </ul>

      <h3 className='mt-6 text-lg font-semibold'>示例</h3>
      <CodeBlock
        title='V2V · multipart'
        code={`curl -X POST ${base}/videos \\
  -H "Authorization: Bearer sk-xxx" \\
  -F "model=omni-fast-v2v" \\
  -F "prompt=将画面风格转换为赛博朋克风" \\
  -F "input_video=@source.mp4"`}
      />
    </DocsSection>
  )
}

export function VeoCleanSection(props: ApiDocsContext) {
  const { base } = props

  return (
    <DocsSection
      id='api-veo-clean'
      title='Veo-Clean 去水印'
      description='上传带水印的视频，系统自动去除水印。异步流程与视频生成一致，参数见通用文档。'
    >
      <DocsTable
        headers={['模型', '计费', '必传素材']}
        rows={[['veo-clean', '按视频实际秒数', 'input_video（multipart，≤20MB）']]}
      />
      <p className='text-muted-foreground text-sm'>{pricingNote()}</p>

      <CodeBlock
        title='Multipart 上传'
        code={`curl -X POST ${base}/videos \\
  -H "Authorization: Bearer sk-xxx" \\
  -F "model=veo-clean" \\
  -F "prompt=remove watermark" \\
  -F "input_video=@watermarked.mp4"`}
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
          ['模型名', '与模型广场展示名一致，如 omni-fast、Seedance2.0-720p、grok-imagine-video'],
          ['视频 API', '见 #api-video-api 统一参数与模型能力对照表'],
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
          <p className='font-medium'>Q: 不同模型参数怎么传？</p>
          <p className='text-muted-foreground'>
            先查{' '}
            <a href='#api-video-api' className='text-primary font-medium hover:underline'>
              视频生成 API（通用）
            </a>{' '}
            或{' '}
            <a href='#api-image-api' className='text-primary font-medium hover:underline'>
              图像生成 API（通用）
            </a>{' '}
            中的参数全集与模型能力对照表，只传该模型支持的字段。
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
