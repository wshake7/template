import { HttpCode, HttpCodeSet, XHeader } from '@vp/core'

export { HttpCode, XHeader }
export type { CodeType } from '@vp/core'

const errorHandlers: Partial<Record<number, (res: Res) => Promise<void> | void>> = {
  [HttpCode.FailLogin]: (res) => {
    gMessage.error(res.msg)
    AccountApi.logout()
    throw new Error(res.msg)
  },
  [HttpCode.FailRequestKey]: async (res) => {
    gMessage.error(res.msg)
    useDeviceStore.getState().setPublicKey('')
    AccountApi.logout()
    const publicKey = await EncryptApi.publicKey() || ''
    if (publicKey === '') {
      gMessage.error('系统异常')
      return
    }
    useDeviceStore.getState().setPublicKey(publicKey)
    throw new Error(res.msg)
  },
  [HttpCode.UN_KNOW]: (res) => {
    gMessage.error('请求错误')
    throw new Error(JSON.stringify(res))
  },
}

export async function HttpCodeCheck(res: Res) {
  const { code, msg } = res

  if (code === HttpCode.SUCCESS) { return }

  const handler = errorHandlers[code]

  if (handler) {
    await handler(res)
  }
  else if (HttpCodeSet.has(code)) {
    gMessage.error(msg)
    throw new Error(msg)
  }
  else {
    console.error('未识别状态码', code)
    throw new Error(msg)
  }
}
