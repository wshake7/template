import type { ReactNode } from 'react'
import type { DictMatchedEntriesByCode, DictMatchedEntry } from '~/api/business/sysDict'
import { useCallback, useEffect, useMemo, useSyncExternalStore } from 'react'
import { useTranslation } from 'react-i18next'
import { DictApi } from '~/api/business/sysDict'
import { renderDictEntryLabel } from '~/components/dictEntryLabel'

const dictMatchCache = new Map<string, DictMatchedEntry[]>()
const loadingDictCodes = new Set<string>()
const pendingDictCodes = new Set<string>()
const dictMatchListeners = new Set<() => void>()
let flushScheduled = false
let dictMatchSnapshot: DictMatchedEntriesByCode = {}

function normalizeDictCodes(codes: string[]) {
  return Array.from(new Set(codes.map(code => code.trim()).filter(Boolean))).sort()
}

function subscribeDictMatches(listener: () => void) {
  dictMatchListeners.add(listener)
  return () => {
    dictMatchListeners.delete(listener)
  }
}

function emitDictMatchesChange() {
  dictMatchSnapshot = Object.fromEntries(dictMatchCache)
  dictMatchListeners.forEach(listener => listener())
}

function getDictMatchSnapshot() {
  return dictMatchSnapshot
}

function flushPendingDictMatches() {
  flushScheduled = false
  const requestCodes = Array.from(pendingDictCodes)
    .filter(code => !dictMatchCache.has(code) && !loadingDictCodes.has(code))
  pendingDictCodes.clear()

  if (requestCodes.length === 0) {
    return
  }

  requestCodes.forEach(code => loadingDictCodes.add(code))
  DictApi.entryMatch({ codes: requestCodes })
    .send()
    .then((res) => {
      const entriesByCode = res.data ?? {}
      requestCodes.forEach((code) => {
        dictMatchCache.set(code, entriesByCode[code] ?? [])
      })
    })
    .catch(() => {
      requestCodes.forEach((code) => {
        dictMatchCache.set(code, [])
      })
    })
    .finally(() => {
      requestCodes.forEach(code => loadingDictCodes.delete(code))
      emitDictMatchesChange()
    })
}

function scheduleDictMatches(codes: string[]) {
  const missingCodes = normalizeDictCodes(codes)
    .filter(code => !dictMatchCache.has(code) && !loadingDictCodes.has(code))

  if (missingCodes.length === 0) {
    return
  }

  missingCodes.forEach(code => pendingDictCodes.add(code))
  if (!flushScheduled) {
    flushScheduled = true
    queueMicrotask(flushPendingDictMatches)
  }
}

export function useDictMatches(codes: string[]) {
  const codesKey = useMemo(() => normalizeDictCodes(codes).join('\u0000'), [codes])
  const cachedEntriesByCode = useSyncExternalStore(
    subscribeDictMatches,
    getDictMatchSnapshot,
    getDictMatchSnapshot,
  )

  useEffect(() => {
    const requestCodes = codesKey ? codesKey.split('\u0000') : []
    scheduleDictMatches(requestCodes)
  }, [codesKey])

  return useMemo(() => {
    const requestCodes = codesKey ? codesKey.split('\u0000') : []
    return Object.fromEntries(requestCodes.map(code => [code, cachedEntriesByCode[code] ?? []]))
  }, [cachedEntriesByCode, codesKey])
}

export function useDictMatch(code: string) {
  const { t } = useTranslation()
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
    const translated = t(entryLabel, { defaultValue: entryLabel })
    return typeof translated === 'string' ? translated : entryLabel
  }, [t])

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
