import API from './index'

export interface SysOperationLog {
  id: number
  requestID: string
  method: string
  module: string
  path: string
  referer: string
  beforeChange: string
  afterChange: string
  formatChange: string
  requestURI: string
  requestBody: string
  requestHeader: string
  response: string
  costTime: number
  userID: number
  username: string
  clientIP: string
  statusCode: number
  reason: string
  success: boolean
  location: string
  userAgent: string
  browserName: string
  browserVersion: string
  clientID: string
  clientName: string
  osName: string
  osVersion: string
  createdAt: string
}

export interface ReqLogDetail {
  id: number
}

async function list(req: PagingRequest) {
  return await API.Post<Res<PagingResult<SysOperationLog>>>('/api/sys/operation/log/list', req, {
    cacheFor: 0,
  }).send()
}

async function detail(req: ReqLogDetail) {
  return await API.Post<Res<SysOperationLog>>('/api/sys/operation/log/detail', req, {
    cacheFor: 0,
  }).send()
}

export const OperationLogApi = {
  list,
  detail,
}
