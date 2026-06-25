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

export interface ModelChannelPrefix {
  id: number
  prefix: string
  note?: string
  enabled: boolean
  sort_order: number
  created_time: number
  updated_time: number
}

export interface ModelPublicAlias {
  id: number
  internal_name: string
  public_name: string
  created_time: number
  updated_time: number
}

export interface ModelPublicNameRegistryStatus {
  ready: boolean
  collisions: Record<string, string[]>
}

export const modelNamingQueryKeys = {
  all: ['model-naming'] as const,
  prefixes: () => [...modelNamingQueryKeys.all, 'prefixes'] as const,
  aliases: () => [...modelNamingQueryKeys.all, 'aliases'] as const,
  status: () => [...modelNamingQueryKeys.all, 'status'] as const,
}
