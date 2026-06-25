/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/
import { DocsSection } from '../components/docs-section'
import { DocsTable } from '../components/docs-table'

export function VideoGuideSection() {
  return (
    <DocsSection
      id='api-video-guide'
      title='视频生成 · 内容审查避坑指南'
      description='适用于 omni-fast / veo 系列等 Gemini 视频模型，提高一次出片成功率。'
    >
      <p>
        若收到内容审查拒绝（identifiable real person / unsafe content / protected IP），这是上游的确定性拒绝——
        <strong>重试或换账号无效</strong>，须修改提示词或更换参考图后重新提交。
      </p>

      <h3 className='text-lg font-semibold'>六大高危雷区</h3>
      <DocsTable
        headers={['雷区', '典型触发']}
        rows={[
          ['① 可识别真人 / 名人', '写实正脸、名人名字、真人照片参考图；photorealistic + 真人 = 高危'],
          ['② 未成年人', '儿童/婴儿/青少年出现在任何危险或暧昧情境'],
          ['③ 性 / 裸露 / 暧昧', '裸露、性感、情色、内衣、亲密场景'],
          ['④ 暴力 / 危险 / 血腥', '武器、打斗、战争、自残、伤口'],
          ['⑤ 版权 / IP', '迪士尼、漫威、游戏角色、品牌 Logo、知名影视场景'],
          ['⑥ 政治 / 敏感', '政治人物、抗议、国旗敏感组合、宗教冲突'],
        ]}
      />

      <h3 className='mt-6 text-lg font-semibold'>安全改写对照</h3>
      <DocsTable
        headers={['❌ 高风险写法', '✅ 建议改写']}
        rows={[
          ['Taylor Swift 在演唱会唱歌', '一位虚构的流行歌手在舞台表演，非写实动画风格'],
          ['上传真人正脸照片做参考', '使用侧面/背影/远景，或非写实插画风格'],
          ['蜘蛛侠在城市荡秋千', '一个穿红蓝紧身衣的虚构超级英雄（非漫威官方形象）'],
          ['photorealistic 8K 写实人像', 'cinematic stylized portrait, semi-realistic, fictional character'],
          ['儿童在泳池玩耍', '避免未成年人；改为成年人在安全休闲场景'],
        ]}
      />

      <h3 className='mt-6 text-lg font-semibold'>参考图建议</h3>
      <ul className='list-disc space-y-2 pl-5'>
        <li>优先使用插画、3D 渲染、非写实风格</li>
        <li>人物用侧面、背影、远景或面部不可辨识</li>
        <li>避免上传含 Logo、水印、真人正脸的图片</li>
        <li>被拒绝后不要重复提交同一请求，立刻换图或改 prompt</li>
      </ul>
    </DocsSection>
  )
}
