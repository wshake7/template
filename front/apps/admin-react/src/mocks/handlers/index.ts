import { HttpCode } from '~/domains/http'
import { gEnv } from '~/env'

const baseUrl = `http://localhost:${gEnv.VITE_PORT}/mock`

export function url(path: string) {
  return `${baseUrl}${path}`
}

export function fail(msg: string): Res {
  return {
    code: HttpCode.ERROR,
    msg,
  }
}

export function success<T>(data: T): Res<T> {
  return {
    code: HttpCode.SUCCESS,
    msg: 'success',
    data,
  }
}
