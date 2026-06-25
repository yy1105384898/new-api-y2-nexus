/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
import { useMemo, useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { AlertTriangle, Pencil, Plus, Trash2 } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { ConfirmDialog } from '@/components/confirm-dialog'
import { Dialog } from '@/components/dialog'
import { StatusBadge } from '@/components/status-badge'
import {
  createModelChannelPrefix,
  createModelPublicAlias,
  deleteModelChannelPrefix,
  deleteModelPublicAlias,
  getModelPublicNameRegistryStatus,
  listModelChannelPrefixes,
  listModelPublicAliases,
  updateModelChannelPrefix,
  updateModelPublicAlias,
} from './model-naming-api'
import {
  modelNamingQueryKeys,
  type ModelChannelPrefix,
  type ModelPublicAlias,
} from './model-naming-types'

type PrefixFormState = {
  prefix: string
  note: string
  enabled: boolean
  sort_order: string
}

type AliasFormState = {
  internal_name: string
  public_name: string
}

const emptyPrefixForm = (): PrefixFormState => ({
  prefix: '',
  note: '',
  enabled: true,
  sort_order: '0',
})

const emptyAliasForm = (): AliasFormState => ({
  internal_name: '',
  public_name: '',
})

export function ModelNamingSection() {
  const { t } = useTranslation()
  const queryClient = useQueryClient()

  const [prefixDialogOpen, setPrefixDialogOpen] = useState(false)
  const [editingPrefix, setEditingPrefix] = useState<ModelChannelPrefix | null>(
    null
  )
  const [prefixForm, setPrefixForm] = useState<PrefixFormState>(emptyPrefixForm)
  const [deletePrefixId, setDeletePrefixId] = useState<number | null>(null)

  const [aliasDialogOpen, setAliasDialogOpen] = useState(false)
  const [editingAlias, setEditingAlias] = useState<ModelPublicAlias | null>(null)
  const [aliasForm, setAliasForm] = useState<AliasFormState>(emptyAliasForm)
  const [deleteAliasId, setDeleteAliasId] = useState<number | null>(null)

  const invalidateAll = async () => {
    await Promise.all([
      queryClient.invalidateQueries({ queryKey: modelNamingQueryKeys.prefixes() }),
      queryClient.invalidateQueries({ queryKey: modelNamingQueryKeys.aliases() }),
      queryClient.invalidateQueries({ queryKey: modelNamingQueryKeys.status() }),
    ])
  }

  const { data: prefixes = [], isLoading: prefixesLoading } = useQuery({
    queryKey: modelNamingQueryKeys.prefixes(),
    queryFn: listModelChannelPrefixes,
  })

  const { data: aliases = [], isLoading: aliasesLoading } = useQuery({
    queryKey: modelNamingQueryKeys.aliases(),
    queryFn: listModelPublicAliases,
  })

  const { data: registryStatus } = useQuery({
    queryKey: modelNamingQueryKeys.status(),
    queryFn: getModelPublicNameRegistryStatus,
  })

  const collisionEntries = useMemo(
    () => Object.entries(registryStatus?.collisions ?? {}),
    [registryStatus?.collisions]
  )

  const prefixMutation = useMutation({
    mutationFn: async () => {
      const payload = {
        prefix: prefixForm.prefix.trim(),
        note: prefixForm.note.trim(),
        enabled: prefixForm.enabled,
        sort_order: Number(prefixForm.sort_order) || 0,
      }
      if (editingPrefix) {
        return updateModelChannelPrefix({ id: editingPrefix.id, ...payload })
      }
      return createModelChannelPrefix(payload)
    },
    onSuccess: async (response) => {
      if (!response.success) {
        toast.error(response.message || t('Operation failed'))
        return
      }
      toast.success(
        editingPrefix
          ? t('Channel prefix updated')
          : t('Channel prefix created')
      )
      setPrefixDialogOpen(false)
      setEditingPrefix(null)
      setPrefixForm(emptyPrefixForm())
      await invalidateAll()
    },
    onError: (error: Error) => {
      toast.error(error.message || t('Operation failed'))
    },
  })

  const deletePrefixMutation = useMutation({
    mutationFn: (id: number) => deleteModelChannelPrefix(id),
    onSuccess: async (response) => {
      if (!response.success) {
        toast.error(response.message || t('Operation failed'))
        return
      }
      toast.success(t('Channel prefix deleted'))
      setDeletePrefixId(null)
      await invalidateAll()
    },
    onError: (error: Error) => {
      toast.error(error.message || t('Operation failed'))
    },
  })

  const aliasMutation = useMutation({
    mutationFn: async () => {
      const payload = {
        internal_name: aliasForm.internal_name.trim(),
        public_name: aliasForm.public_name.trim(),
      }
      if (editingAlias) {
        return updateModelPublicAlias({ id: editingAlias.id, ...payload })
      }
      return createModelPublicAlias(payload)
    },
    onSuccess: async (response) => {
      if (!response.success) {
        toast.error(response.message || t('Operation failed'))
        return
      }
      toast.success(
        editingAlias ? t('Public alias updated') : t('Public alias created')
      )
      setAliasDialogOpen(false)
      setEditingAlias(null)
      setAliasForm(emptyAliasForm())
      await invalidateAll()
    },
    onError: (error: Error) => {
      toast.error(error.message || t('Operation failed'))
    },
  })

  const deleteAliasMutation = useMutation({
    mutationFn: (id: number) => deleteModelPublicAlias(id),
    onSuccess: async (response) => {
      if (!response.success) {
        toast.error(response.message || t('Operation failed'))
        return
      }
      toast.success(t('Public alias deleted'))
      setDeleteAliasId(null)
      await invalidateAll()
    },
    onError: (error: Error) => {
      toast.error(error.message || t('Operation failed'))
    },
  })

  const openCreatePrefix = () => {
    setEditingPrefix(null)
    setPrefixForm(emptyPrefixForm())
    setPrefixDialogOpen(true)
  }

  const openEditPrefix = (item: ModelChannelPrefix) => {
    setEditingPrefix(item)
    setPrefixForm({
      prefix: item.prefix,
      note: item.note ?? '',
      enabled: item.enabled,
      sort_order: String(item.sort_order ?? 0),
    })
    setPrefixDialogOpen(true)
  }

  const openCreateAlias = () => {
    setEditingAlias(null)
    setAliasForm(emptyAliasForm())
    setAliasDialogOpen(true)
  }

  const openEditAlias = (item: ModelPublicAlias) => {
    setEditingAlias(item)
    setAliasForm({
      internal_name: item.internal_name,
      public_name: item.public_name,
    })
    setAliasDialogOpen(true)
  }

  return (
    <div className='flex h-full min-h-0 flex-col gap-4 overflow-y-auto pb-4'>
      {collisionEntries.length > 0 ? (
        <Alert variant='destructive'>
          <AlertTriangle className='h-4 w-4' />
          <AlertTitle>{t('Public name collisions detected')}</AlertTitle>
          <AlertDescription>
            <p className='mb-2'>
              {t(
                'Multiple internal models map to the same public name. Add aliases to resolve ambiguity.'
              )}
            </p>
            <ul className='list-disc space-y-1 pl-5 text-sm'>
              {collisionEntries.map(([publicName, internals]) => (
                <li key={publicName}>
                  <span className='font-medium'>{publicName}</span>:{' '}
                  {internals.join(', ')}
                </li>
              ))}
            </ul>
          </AlertDescription>
        </Alert>
      ) : (
        <Alert>
          <AlertTitle>{t('No public name collisions')}</AlertTitle>
          <AlertDescription>
            {t(
              'Registry is ready. Prefix stripping and aliases are applied to inbound and outbound model names.'
            )}
          </AlertDescription>
        </Alert>
      )}

      <Card>
        <CardHeader className='flex flex-row items-center justify-between space-y-0'>
          <div>
            <CardTitle>{t('Channel prefixes')}</CardTitle>
            <p className='text-muted-foreground mt-1 text-sm'>
              {t(
                'Strip these prefixes from internal ability names when exposing public model names.'
              )}
            </p>
          </div>
          <Button size='sm' onClick={openCreatePrefix}>
            <Plus className='h-4 w-4' />
            {t('Add prefix')}
          </Button>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>{t('Prefix')}</TableHead>
                <TableHead>{t('Note')}</TableHead>
                <TableHead>{t('Sort')}</TableHead>
                <TableHead>{t('Status')}</TableHead>
                <TableHead className='w-[100px]'>{t('Actions')}</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {prefixesLoading ? (
                <TableRow>
                  <TableCell colSpan={5} className='text-muted-foreground'>
                    {t('Loading...')}
                  </TableCell>
                </TableRow>
              ) : prefixes.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={5} className='text-muted-foreground'>
                    {t('No channel prefixes configured')}
                  </TableCell>
                </TableRow>
              ) : (
                prefixes.map((item) => (
                  <TableRow key={item.id}>
                    <TableCell className='font-mono text-sm'>
                      {item.prefix}
                    </TableCell>
                    <TableCell>{item.note || '—'}</TableCell>
                    <TableCell>{item.sort_order}</TableCell>
                    <TableCell>
                      <StatusBadge
                        label={item.enabled ? t('Enabled') : t('Disabled')}
                        variant={item.enabled ? 'success' : 'neutral'}
                        showDot
                        copyable={false}
                      />
                    </TableCell>
                    <TableCell>
                      <div className='flex gap-1'>
                        <Button
                          variant='ghost'
                          size='icon'
                          onClick={() => openEditPrefix(item)}
                        >
                          <Pencil className='h-4 w-4' />
                        </Button>
                        <Button
                          variant='ghost'
                          size='icon'
                          onClick={() => setDeletePrefixId(item.id)}
                        >
                          <Trash2 className='h-4 w-4 text-destructive' />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <Card>
        <CardHeader className='flex flex-row items-center justify-between space-y-0'>
          <div>
            <CardTitle>{t('Public aliases')}</CardTitle>
            <p className='text-muted-foreground mt-1 text-sm'>
              {t(
                'Override automatic prefix stripping when multiple internals share a public name.'
              )}
            </p>
          </div>
          <Button size='sm' onClick={openCreateAlias}>
            <Plus className='h-4 w-4' />
            {t('Add alias')}
          </Button>
        </CardHeader>
        <CardContent>
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>{t('Internal name')}</TableHead>
                <TableHead>{t('Public name')}</TableHead>
                <TableHead className='w-[100px]'>{t('Actions')}</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {aliasesLoading ? (
                <TableRow>
                  <TableCell colSpan={3} className='text-muted-foreground'>
                    {t('Loading...')}
                  </TableCell>
                </TableRow>
              ) : aliases.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={3} className='text-muted-foreground'>
                    {t('No public aliases configured')}
                  </TableCell>
                </TableRow>
              ) : (
                aliases.map((item) => (
                  <TableRow key={item.id}>
                    <TableCell className='font-mono text-sm'>
                      {item.internal_name}
                    </TableCell>
                    <TableCell className='font-mono text-sm'>
                      {item.public_name}
                    </TableCell>
                    <TableCell>
                      <div className='flex gap-1'>
                        <Button
                          variant='ghost'
                          size='icon'
                          onClick={() => openEditAlias(item)}
                        >
                          <Pencil className='h-4 w-4' />
                        </Button>
                        <Button
                          variant='ghost'
                          size='icon'
                          onClick={() => setDeleteAliasId(item.id)}
                        >
                          <Trash2 className='h-4 w-4 text-destructive' />
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <Dialog
        open={prefixDialogOpen}
        onOpenChange={setPrefixDialogOpen}
        title={
          editingPrefix ? t('Edit channel prefix') : t('Add channel prefix')
        }
        footer={
          <>
            <Button variant='outline' onClick={() => setPrefixDialogOpen(false)}>
              {t('Cancel')}
            </Button>
            <Button
              onClick={() => prefixMutation.mutate()}
              disabled={prefixMutation.isPending || !prefixForm.prefix.trim()}
            >
              {t('Save')}
            </Button>
          </>
        }
      >
        <div className='space-y-4'>
          <div className='space-y-2'>
            <Label htmlFor='prefix'>{t('Prefix')}</Label>
            <Input
              id='prefix'
              placeholder='go2api-'
              value={prefixForm.prefix}
              onChange={(e) =>
                setPrefixForm((prev) => ({ ...prev, prefix: e.target.value }))
              }
            />
          </div>
          <div className='space-y-2'>
            <Label htmlFor='prefix-note'>{t('Note')}</Label>
            <Input
              id='prefix-note'
              value={prefixForm.note}
              onChange={(e) =>
                setPrefixForm((prev) => ({ ...prev, note: e.target.value }))
              }
            />
          </div>
          <div className='space-y-2'>
            <Label htmlFor='prefix-sort'>{t('Sort order')}</Label>
            <Input
              id='prefix-sort'
              type='number'
              value={prefixForm.sort_order}
              onChange={(e) =>
                setPrefixForm((prev) => ({
                  ...prev,
                  sort_order: e.target.value,
                }))
              }
            />
          </div>
          <div className='flex items-center gap-2'>
            <Checkbox
              id='prefix-enabled'
              checked={prefixForm.enabled}
              onCheckedChange={(checked) =>
                setPrefixForm((prev) => ({
                  ...prev,
                  enabled: checked === true,
                }))
              }
            />
            <Label htmlFor='prefix-enabled'>{t('Enabled')}</Label>
          </div>
        </div>
      </Dialog>

      <Dialog
        open={aliasDialogOpen}
        onOpenChange={setAliasDialogOpen}
        title={editingAlias ? t('Edit public alias') : t('Add public alias')}
        footer={
          <>
            <Button variant='outline' onClick={() => setAliasDialogOpen(false)}>
              {t('Cancel')}
            </Button>
            <Button
              onClick={() => aliasMutation.mutate()}
              disabled={
                aliasMutation.isPending ||
                !aliasForm.internal_name.trim() ||
                !aliasForm.public_name.trim()
              }
            >
              {t('Save')}
            </Button>
          </>
        }
      >
        <div className='space-y-4'>
          <div className='space-y-2'>
            <Label htmlFor='internal-name'>{t('Internal name')}</Label>
            <Input
              id='internal-name'
              placeholder='go2api-gpt-image-2-1k'
              value={aliasForm.internal_name}
              onChange={(e) =>
                setAliasForm((prev) => ({
                  ...prev,
                  internal_name: e.target.value,
                }))
              }
            />
          </div>
          <div className='space-y-2'>
            <Label htmlFor='public-name'>{t('Public name')}</Label>
            <Input
              id='public-name'
              placeholder='gpt-image-2-1k'
              value={aliasForm.public_name}
              onChange={(e) =>
                setAliasForm((prev) => ({
                  ...prev,
                  public_name: e.target.value,
                }))
              }
            />
          </div>
        </div>
      </Dialog>

      <ConfirmDialog
        open={deletePrefixId !== null}
        onOpenChange={(open) => !open && setDeletePrefixId(null)}
        title={t('Delete channel prefix')}
        desc={t(
          'Removing a prefix stops automatic stripping for matching internal names.'
        )}
        confirmText={t('Delete')}
        destructive
        handleConfirm={() => {
          if (deletePrefixId !== null) {
            deletePrefixMutation.mutate(deletePrefixId)
          }
        }}
      />

      <ConfirmDialog
        open={deleteAliasId !== null}
        onOpenChange={(open) => !open && setDeleteAliasId(null)}
        title={t('Delete public alias')}
        desc={t('This internal model will fall back to prefix stripping rules.')}
        confirmText={t('Delete')}
        destructive
        handleConfirm={() => {
          if (deleteAliasId !== null) {
            deleteAliasMutation.mutate(deleteAliasId)
          }
        }}
      />
    </div>
  )
}
