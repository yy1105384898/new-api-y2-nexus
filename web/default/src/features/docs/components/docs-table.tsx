/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/
import type { ReactNode } from 'react'
import { cn } from '@/lib/utils'

type DocsTableProps = {
  headers: string[]
  rows: ReactNode[][]
  className?: string
}

export function DocsTable(props: DocsTableProps) {
  return (
    <div className={cn('border-border/50 overflow-x-auto rounded-xl border', props.className)}>
      <table className='w-full min-w-[32rem] text-left text-sm'>
        <thead className='bg-muted/40 border-border/50 border-b'>
          <tr>
            {props.headers.map((header) => (
              <th key={header} className='text-muted-foreground px-4 py-3 font-semibold'>
                {header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {props.rows.map((row, index) => (
            <tr key={index} className='border-border/40 border-b last:border-b-0'>
              {row.map((cell, cellIndex) => (
                <td key={cellIndex} className='text-foreground/90 px-4 py-3 align-top leading-6'>
                  {cell}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
