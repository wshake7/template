import API from './index'

export interface Resource {
  id: number
  type: string
  code: string
  name: string
  isEnabled: boolean
  remark: string
}

export interface ReqResourceCreate {
  type: string
  code: string
  name: string
  isEnabled: boolean
  remark: string
}

export interface ReqResourceUpdate extends Partial<ReqResourceCreate> {
  id: number
}

export interface ReqResourceBatchDelete {
  ids: number[]
}

function resourceList(req: PagingRequest) {
  return API.Post<Res<PagingResult<Resource>>>('/api/sys/resource/list', req, {
    cacheFor: 0,
  })
}

async function resourceCreate(req: ReqResourceCreate) {
  await API.Post<Res>('/api/sys/resource/create', req, {
    cacheFor: 0,
  }).send()
}

async function resourceUpdate(req: ReqResourceUpdate) {
  await API.Post<Res>('/api/sys/resource/update', req, {
    cacheFor: 0,
  }).send()
}

async function resourceDel(req: ReqResourceBatchDelete) {
  await API.Post<Res>('/api/sys/resource/del', req, {
    cacheFor: 0,
  }).send()
}

export const ResourceApi = {
  resourceList,
  resourceCreate,
  resourceUpdate,
  resourceDel,
}
