import API from './index'

export interface DictType {
  id: number
  typeCode: string
  typeName: string
  isEnabled: boolean
  sortOrder: number
  description: string
}

export interface DictEntry {
  id: number
  entryLabel: string
  entryValue: string
  numericValue: number
  languageCode: string
  sysDictTypeId: number
  sortOrder: number
  isEnabled: boolean
  remark: string
}

export interface ReqDictTypeCreate {
  typeCode: string
  typeName: string
  isEnabled: boolean
  sortOrder: number
  description: string
}

export interface ReqDictTypeUpdate extends Partial<ReqDictTypeCreate> {
  id: number
}

export interface ReqDictTypeSwitchStatus {
  id: number
  isEnabled: boolean
}

export interface ReqDictTypeBatchDelete {
  ids: number[]
}

export interface ReqDictEntryCreate {
  entryLabel: string
  entryValue: string
  numericValue: number
  languageCode: string
  sysDictTypeId: number
  sortOrder: number
  isEnabled: boolean
  remark: string
}

export interface ReqDictEntryUpdate extends Partial<ReqDictEntryCreate> {
  id: number
}

export interface ReqDictEntrySwitchStatus {
  id: number
  isEnabled: boolean
}

export interface ReqDictEntryBatchDelete {
  ids: number[]
}

export interface ReqDictEntryBatchCopy {
  entryIds: number[]
  targetTypeId: number
}

async function typeList(req: PagingRequest) {
  return await API.Post<Res<PagingResult<DictType>>>('/api/sys/dict/type/list', req, {
    cacheFor: 0,
  }).send()
}

async function typeCreate(req: ReqDictTypeCreate) {
  await API.Post<Res>('/api/sys/dict/type/create', req, {
    cacheFor: 0,
  }).send()
}

async function typeUpdate(req: ReqDictTypeUpdate) {
  await API.Post<Res>('/api/sys/dict/type/update', req, {
    cacheFor: 0,
  }).send()
}

async function typeSwitch(req: ReqDictTypeSwitchStatus) {
  await API.Post<Res>('/api/sys/dict/type/switch', req, {
    cacheFor: 0,
  }).send()
}

async function typeDel(req: ReqDictTypeBatchDelete) {
  await API.Post<Res>('/api/sys/dict/type/del', req, {
    cacheFor: 0,
  }).send()
}

async function entryList(req: PagingRequest) {
  return await API.Post<Res<PagingResult<DictEntry>>>('/api/sys/dict/entry/list', req, {
    cacheFor: 0,
  }).send()
}

async function entryCreate(req: ReqDictEntryCreate) {
  await API.Post<Res>('/api/sys/dict/entry/create', req, {
    cacheFor: 0,
  }).send()
}

async function entryUpdate(req: ReqDictEntryUpdate) {
  await API.Post<Res>('/api/sys/dict/entry/update', req, {
    cacheFor: 0,
  }).send()
}

async function entrySwitch(req: ReqDictEntrySwitchStatus) {
  await API.Post<Res>('/api/sys/dict/entry/switch', req, {
    cacheFor: 0,
  }).send()
}

async function entryDel(req: ReqDictEntryBatchDelete) {
  await API.Post<Res>('/api/sys/dict/entry/del', req, {
    cacheFor: 0,
  }).send()
}

async function entryBatchCopy(req: ReqDictEntryBatchCopy) {
  await API.Post<Res>('/api/sys/dict/entry/batch/copy', req, {
    cacheFor: 0,
  }).send()
}

export const DictApi = {
  typeList,
  typeCreate,
  typeUpdate,
  typeSwitch,
  typeDel,
  entryList,
  entryCreate,
  entryUpdate,
  entrySwitch,
  entryDel,
  entryBatchCopy,
}
