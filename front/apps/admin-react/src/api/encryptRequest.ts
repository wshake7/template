import { XHeader } from '~/domains/http'
import { gMessage } from '~/utils/antd'
import { aesDecrypt, aesEncrypt, generateAesKey, rsaEncrypt, uriSort } from '~/utils/encrypt'

interface EncryptableMethod {
  url: string
  data?: any
  config: {
    headers?: Record<string, any>
    params?: Record<string, any> | string | URLSearchParams
  }
  meta?: Record<string, any>
}

interface EncryptedRequestConfig {
  headers: Record<string, any>
  data?: any
  aesKey?: CryptoKey
  nonce: string
}

function normalizeParams(params: Record<string, any> | string | URLSearchParams | undefined): Record<string, any> {
  if (!params) { return {} }

  if (typeof params === 'string') {
    return Object.fromEntries(new URLSearchParams(params))
  }

  if (params instanceof URLSearchParams) {
    return Object.fromEntries(params.entries())
  }

  if (typeof params === 'object') {
    return params as Record<string, any>
  }

  return {}
}

export async function ensurePublicKey() {
  let publicKey = useDeviceStore.getState().publicKey
  if (publicKey !== '') {
    return publicKey
  }

  const response = await fetch('/api/encrypt/public/key')
  if (!response.ok) {
    return ''
  }
  const res = await response.json() as Res<{ publicKey: string }>
  publicKey = res.data?.publicKey || ''
  if (publicKey !== '') {
    useDeviceStore.getState().setPublicKey(publicKey)
  }
  return publicKey
}

export async function encryptRequest(method: EncryptableMethod) {
  if (method.url === '/api/encrypt/public/key') {
    const timestamp = Date.now()
    const nonce = Math.random().toString(36).substring(2, 18)
    method.config.headers = method.config.headers || {}
    method.config.headers[XHeader.XRequestTimestamp] = timestamp
    method.config.headers[XHeader.XRequestID] = nonce
    return
  }

  const encryptedConfig = await createEncryptedRequestConfig({
    data: method.data,
    headers: method.config.headers,
    params: method.config.params,
  })
  method.config.headers = encryptedConfig.headers
  method.data = encryptedConfig.data
  method.meta = {
    ...method.meta,
    aesKey: encryptedConfig.aesKey,
    nonce: encryptedConfig.nonce,
  }
}

export async function createEncryptedRequestConfig(options: {
  data?: any
  headers?: Record<string, any>
  params?: Record<string, any> | string | URLSearchParams
} = {}): Promise<EncryptedRequestConfig> {
  const headers = {
    ...(options.headers || {}),
  }
  const timestamp = Date.now()
  const nonce = Math.random().toString(36).substring(2, 18)
  headers[XHeader.XRequestTimestamp] = timestamp
  headers[XHeader.XRequestID] = nonce

  const publicKey = await ensurePublicKey()
  if (publicKey === '') {
    gMessage.error('系统异常')
    return { headers, nonce }
  }

  const publicCryptoKey = await useDeviceStore.getState().getPublicCryptoKey()
  if (!publicCryptoKey) {
    gMessage.error('系统异常')
    return { headers, nonce }
  }

  const { key, keyBase64 } = await generateAesKey()
  headers[XHeader.XRequestEncryptedKey] = await rsaEncrypt(keyBase64, publicCryptoKey)

  const sort = uriSort({
    [XHeader.XRequestTimestamp]: timestamp,
    [XHeader.XRequestID]: nonce,
    ...normalizeParams(options.params),
  })
  const aesData = await aesEncrypt(key, sort, options.data)
  headers[XHeader.XRequestSignature] = aesData.TagIv

  return {
    headers,
    data: aesData.Ciphertext !== '' ? aesData.Ciphertext : options.data,
    aesKey: key,
    nonce,
  }
}

export async function decryptText(encryptedText: string, aesKey: CryptoKey) {
  return aesDecrypt(encryptedText, aesKey, '')
}
