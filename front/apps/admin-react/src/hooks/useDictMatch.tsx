import type { ReactNode } from 'react'
import type { DictMatchedEntry } from '~/api/sysDict'
import { useCallback, useEffect, useMemo, useState } from 'react'
import { DictApi } from '~/api/sysDict'
import { renderDictEntryLabel } from '~/components/dictEntryLabel'

export function useDictMatch(code: string) {
  const [entries, setEntries] = useState<DictMatchedEntry[]>([])

  useEffect(() => {
    if (!code) {
      setEntries([])
      return
    }

    let ignore = false

    DictApi.entryMatch({ code })
      .send()
      .then((res) => {
        if (!ignore) {
          setEntries(res.data ?? [])
        }
      })
      .catch(() => {
        if (!ignore) {
          setEntries([])
        }
      })

    return () => {
      ignore = true
    }
  }, [code])

  const entryByValue = useMemo(() => {
    return new Map(entries.map(item => [item.entryValue, item]))
  }, [entries])

  const getEntry = useCallback((value: string | number | boolean) => {
    return entryByValue.get(String(value))
  }, [entryByValue])

  const getLabel = useCallback((value: string | number | boolean, fallback = '') => {
    return getEntry(value)?.entryLabel ?? fallback
  }, [getEntry])

  const renderLabel = useCallback((value: string | number | boolean, fallback: ReactNode = '未知状态') => {
    const entry = getEntry(value)
    if (!entry) {
      return fallback
    }
    return renderDictEntryLabel(entry.labelComponent, entry.entryLabel)
  }, [getEntry])

  return useMemo(() => ({
    entries,
    getEntry,
    getLabel,
    renderLabel,
  }), [entries, getEntry, getLabel, renderLabel])
}
