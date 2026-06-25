/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/
import { createApiDocsContext } from './context'
import { GrokCliVideoSection, GrokImageVideoSection } from './grok-sections'
import { GeminiImageSection, GeminiMusicSection, GptImageSection } from './image-sections'
import { OmniVideoSection, OverviewFaqSection, VeoCleanSection } from './omni-and-overview'
import { SeedanceVideoSection } from './seedance-section'
import { VideoGuideSection } from './video-guide-section'

type ApiDocsSectionsProps = {
  siteOrigin: string
}

export function ApiDocsSections(props: ApiDocsSectionsProps) {
  const ctx = createApiDocsContext(props.siteOrigin)

  return (
    <>
      <OmniVideoSection {...ctx} />
      <VeoCleanSection {...ctx} />
      <SeedanceVideoSection {...ctx} />
      <GrokImageVideoSection {...ctx} />
      <GrokCliVideoSection {...ctx} />
      <GptImageSection {...ctx} />
      <GeminiImageSection {...ctx} />
      <GeminiMusicSection {...ctx} />
      <OverviewFaqSection {...ctx} />
      <VideoGuideSection />
    </>
  )
}
