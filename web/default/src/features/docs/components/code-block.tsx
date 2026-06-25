/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.
*/
import { Copy } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'

type CodeBlockProps = {
  title?: string
  code: string
  className?: string
}

export function CodeBlock(props: CodeBlockProps) {
  const { t } = useTranslation()

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(props.code)
      toast.success(t('userDocs.copySuccess'))
    } catch {
      toast.error(t('userDocs.copyFailed'))
    }
  }

  return (
    <div className={cn('border-border/60 bg-muted/20 overflow-hidden rounded-xl border', props.className)}>
      <div className='border-border/50 flex items-center justify-between gap-3 border-b px-4 py-2'>
        <p className='text-muted-foreground text-xs font-medium tracking-wide uppercase'>
          {props.title || 'Example'}
        </p>
        <Button type='button' variant='ghost' size='sm' className='h-8 px-2' onClick={handleCopy}>
          <Copy className='size-4' />
          {t('userDocs.copy')}
        </Button>
      </div>
      <pre className='text-foreground overflow-x-auto p-4 font-mono text-[13px] leading-6 whitespace-pre'>
        {props.code}
      </pre>
    </div>
  )
}
