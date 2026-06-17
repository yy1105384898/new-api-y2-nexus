/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/
import { z } from 'zod'
import { createFileRoute } from '@tanstack/react-router'
import { CanvasOpenPage } from '@/features/canvas/components/canvas-open-page'

const canvasOpenSearchSchema = z.object({
  redirect: z.string().optional(),
})

export const Route = createFileRoute('/_authenticated/canvas/open')({
  component: RouteComponent,
  validateSearch: canvasOpenSearchSchema,
})

function RouteComponent() {
  const { redirect } = Route.useSearch()
  return <CanvasOpenPage redirect={redirect} />
}
