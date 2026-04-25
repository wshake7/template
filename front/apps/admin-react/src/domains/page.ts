export interface PagingRequest {
  page?: number
  pageSize?: number
  noPaging?: boolean
  orderBy?: string
}

export interface PagingResult<T> {
  items: T[]
  total: number
}
