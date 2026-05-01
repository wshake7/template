import type { ResLogin } from '~/domains/account'
import { createAlova } from 'alova'
import { createClientTokenAuthentication } from 'alova/client'
import adapterFetch from 'alova/fetch'
import reactHook from 'alova/react'
import Cookies from 'js-cookie'
import NProgress from 'nprogress'
import { HttpCode, XHeader } from '~/domains/http'

import { gEnv } from '~/env'
import { router } from '~/router'
import { gMessage } from '~/utils/antd'
import { aesDecrypt, aesEncrypt, generateAesKey, rsaEncrypt, uriSort } from '~/utils/encrypt'

const { onAuthRequired, onResponseRefreshToken } = createClientTokenAuthentication<typeof reactHook>({
  visitorMeta: {
    visitor: true,
  },
  assignToken(method) {
    let token = Cookies.get(XHeader.Token)
    if (token === undefined || token === '') {
      token = useAccountStore.getState().token
    }
    method.config.headers[XHeader.Token] = token
  },
  async login(response, method) {
    if (response.ok) {
      const encryptedText = await response.clone().text()
      const aesKey = method.meta.aesKey
      const decrypted = await aesDecrypt(encryptedText, aesKey, '')
      const json = JSON.parse(decrypted)
      const res = json as Res<ResLogin>
      const data = res.data
      if (res.code === HttpCode.SUCCESS && data) {
        useAccountStore.getState().login(data.token)
        useDeviceStore.getState().setPublicKey(data.publicKey)
        Cookies.set(XHeader.Token, data.token, {
          path: '/',
          sameSite: 'Lax',
        })
        router.update({
          context: {
            account: {
              token: data.token,
            },
          },
        })
        router.navigate({ to: '/dashboard' })
      }
    }
  },
  logout() {
  },
})

function normalizeParams(params: Record<string, any> | string): Record<string, any> {
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

const API = createAlova({
  baseURL: gEnv.VITE_MOCK ? '' : '',
  statesHook: reactHook,
  cacheFor: null,
  requestAdapter: adapterFetch(),
  shareRequest: false,
  beforeRequest: onAuthRequired(async (method) => {
    NProgress.start()
    let publicKey = useDeviceStore.getState().publicKey
    if (publicKey === '' && method.url !== '/api/encrypt/public/key') {
      publicKey = await EncryptApi.publicKey() || ''
      if (publicKey === '') {
        gMessage.error('系统异常')
        return
      }
      useDeviceStore.getState().setPublicKey(publicKey)
    }
    const timestamp = Date.now()
    const nonce = Math.random().toString(36).substring(2, 18)
    method.config.headers[XHeader.XRequestTimestamp] = timestamp
    method.config.headers[XHeader.XRequestID] = nonce
    method.config.headers[XHeader.XRequestEncryptedKey] = publicKey
    if (method.url !== '/api/encrypt/public/key') {
      const publicCryptoKey = await useDeviceStore.getState().getPublicCryptoKey()
      if (!publicCryptoKey) {
        gMessage.error('系统异常')
        return
      }
      const { key, keyBase64 } = await generateAesKey()
      method.meta = {
        ...method.meta,
        aesKey: key,
        nonce,
      }
      const encryptedKey = await rsaEncrypt(keyBase64, publicCryptoKey)
      method.config.headers[XHeader.XRequestEncryptedKey] = encryptedKey
      const queryParams = method.config.params || {}

      const sort = uriSort({
        [XHeader.XRequestTimestamp]: timestamp,
        [XHeader.XRequestID]: nonce,
        ...normalizeParams(queryParams),
      })
      const aesData = await aesEncrypt(key, sort, method.data)
      if (aesData.Ciphertext !== '') {
        method.data = aesData.Ciphertext
      }
      method.config.headers[XHeader.XRequestSignature] = aesData.TagIv
    }
  }),
  responded: onResponseRefreshToken({
    onSuccess: async (response, method) => {
      if (!response.ok) {
        gMessage.error('请求错误')
        throw new Error(`[${response.status}]${response.statusText}`)
      }
      const contentType = response.headers.get('Content-Type') || ''
      if (response.headers.get(XHeader.XResponseIsEncrypt) === 'true') {
        const encryptedText = await response.clone().text()
        const aesKey = method.meta.aesKey
        const decrypted = await aesDecrypt(encryptedText, aesKey, '')
        response = new Response(decrypted, {
          status: response.status,
          statusText: response.statusText,
          headers: response.headers,
        })
      }
      if (contentType.includes('application/json')) {
        const json = await response.clone().json()
        const res = json as Res
        console.log('response', res)
        if (method.url !== '/api/account/logout') {
          await HttpCodeCheck(res)
        }
        return json
      }
      return response
    },
    onError: async (error) => {
      console.error('[API Error]', error)
      throw error
    },
    onComplete: async () => {
      NProgress.done()
    },
  }),
})

export default API
