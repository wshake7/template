export const HttpCode = {
  UN_KNOW: 0,
  SUCCESS: 1,
  ERROR: 2,
  FailRequestExpired: 3,
  FailRequestNonce: 4,
  FailRequestKey: 5,
  FailLogin: 100,
  FailAuth: 200,
} as const

export type CodeType = (typeof HttpCode)[keyof typeof HttpCode]

const HttpCodeSet = new Set(Object.values(HttpCode))

export const Header = {
  XRequestTimestamp: 'X-Request-Timestamp',
  XRequestID: 'X-Request-ID',
  XRequestEncryptedKey: 'X-Request-Encrypted-Key',
  XRequestSignature: 'X-Request-Signature',
  XResponseIsEncrypt: 'X-Response-Is-Encrypt',
  Token: 'Token',
}

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
