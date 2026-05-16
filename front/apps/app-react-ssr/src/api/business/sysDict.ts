import API from '../index'

export interface DictMatchedEntry {
  id: number
  labelComponent: string
  entryLabel: string
  entryValue: string
}

export type DictMatchedEntriesByCode = Record<string, DictMatchedEntry[]>

export interface ReqDictEntryMatch {
  codes: string[]
}

function entryMatch(req: ReqDictEntryMatch) {
  return API.Post<Res<DictMatchedEntriesByCode>>('/api/sys/dict/entry/match', req, {
    cacheFor: 0,
  })
}

export const DictApi = {
  entryMatch,
}
