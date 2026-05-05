import type { ReactNode } from 'react'
import type { DictMatchedEntriesByCode, DictMatchedEntry } from '~/api/business/sysDict'
import { useCallback, useEffect, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { DictApi } from '~/api/business/sysDict'
import { renderDictEntryLabel } from '~/components/dictEntryLabel'

function normalizeDictCodes(codes: string[]) {
  return Array.from(new Set(codes.map(code => code.trim()).filter(Boolean)))
}

export function useDictMatches(codes: string[]) {
  const codesKey = useMemo(() => normalizeDictCodes(codes).join('\u0000'), [codes])
  const [entriesByCode, setEntriesByCode] = useState<DictMatchedEntriesByCode>({})

  useEffect(() => {
    let ignore = false
    const requestCodes = codesKey ? codesKey.split('\u0000') : []

    if (requestCodes.length === 0) {
      setEntriesByCode({})
      return () => {
        ignore = true
      }
    }

    DictApi.entryMatch({ codes: requestCodes })
      .send()
      .then((res) => {
        if (!ignore) {
          setEntriesByCode(res.data ?? {})
        }
      })
      .catch(() => {
        if (!ignore) {
          setEntriesByCode({})
        }
      })

    return () => {
      ignore = true
    }
  }, [codesKey])

  return entriesByCode
}

export function useDictMatch(code: string) {
  const { i18n, t } = useTranslation()
  const codes = useMemo(() => [code], [code])
  const entriesByCode = useDictMatches(codes)
  const entries = useMemo<DictMatchedEntry[]>(() => entriesByCode[code] ?? [], [code, entriesByCode])

  const entryByValue = useMemo(() => {
    return new Map(entries.map(item => [item.entryValue, item]))
  }, [entries])

  const getEntry = useCallback((value: string | number | boolean) => {
    return entryByValue.get(String(value))
  }, [entryByValue])

  const resolveEntryLabel = useCallback((entryLabel: string) => {
    return i18n.exists(entryLabel) ? t(entryLabel) : entryLabel
  }, [i18n, t])

  const getLabel = useCallback((value: string | number | boolean, fallback = '') => {
    const entry = getEntry(value)
    return entry ? resolveEntryLabel(entry.entryLabel) : fallback
  }, [getEntry, resolveEntryLabel])

  const renderLabel = useCallback((value: string | number | boolean, fallback: ReactNode = '未知状态') => {
    const entry = getEntry(value)
    if (!entry) {
      return fallback
    }
    return renderDictEntryLabel(entry.labelComponent, resolveEntryLabel(entry.entryLabel))
  }, [getEntry, resolveEntryLabel])

  return useMemo(() => ({
    entries,
    getEntry,
    getLabel,
    renderLabel,
  }), [entries, getEntry, getLabel, renderLabel])
}
