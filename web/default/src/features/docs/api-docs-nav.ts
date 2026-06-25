/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/
export type ApiDocsNavGroup = {
  titleKey: string
  items: { id: string; titleKey: string }[]
}

/** 按供应商/产品分类，与下游 API 模型名一致（非内部渠道路由名）。 */
export const apiDocsNavGroups: ApiDocsNavGroup[] = [
  {
    titleKey: 'apiDocs.nav.groupGeminiVideo',
    items: [
      { id: 'api-omni-video', titleKey: 'apiDocs.nav.omniVideo' },
      { id: 'api-veo-clean', titleKey: 'apiDocs.nav.veoClean' },
    ],
  },
  {
    titleKey: 'apiDocs.nav.groupSeedance',
    items: [{ id: 'api-seedance-video', titleKey: 'apiDocs.nav.seedanceVideo' }],
  },
  {
    titleKey: 'apiDocs.nav.groupGrok',
    items: [
      { id: 'api-grok-image-video', titleKey: 'apiDocs.nav.grokImageVideo' },
      { id: 'api-grok-cli-video', titleKey: 'apiDocs.nav.grokCliVideo' },
    ],
  },
  {
    titleKey: 'apiDocs.nav.groupImage',
    items: [
      { id: 'api-gpt-image', titleKey: 'apiDocs.nav.gptImage' },
      { id: 'api-gemini-image', titleKey: 'apiDocs.nav.geminiImage' },
    ],
  },
  {
    titleKey: 'apiDocs.nav.groupOther',
    items: [
      { id: 'api-overview-faq', titleKey: 'apiDocs.nav.overviewFaq' },
      { id: 'api-gemini-music', titleKey: 'apiDocs.nav.geminiMusic' },
      { id: 'api-video-guide', titleKey: 'apiDocs.nav.videoGuide' },
    ],
  },
]

export const apiDocsNavItems = apiDocsNavGroups.flatMap((group) => group.items)
