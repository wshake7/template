import type { ResLogin } from '~/domains/account'
import { createAlova } from 'alova'
import { createClientTokenAuthentication } from 'alova/client'
import adapterFetch from 'alova/fetch'
import reactHook from 'alova/react'
import Cookies from 'js-cookie'
import NProgress from 'nprogress'
import { decryptText, encryptRequest } from '~/api/encryptRequest'
import { HttpCode, XHeader } from '~/domains/http'
import { gEnv } from '~/env'
import { router } from '~/router'

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
      const decrypted = await decryptText(encryptedText, aesKey)
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
        router.navigate({ to: '/' })
      }
    }
  },
  logout() {
  },
})

interface RequestMethodMeta {
  meta?: {
    authRole?: string
  }
  url: string
}

function canRequestWithoutToken(method: RequestMethodMeta) {
  const authRole = (method.meta as { authRole?: string } | undefined)?.authRole
  return authRole === 'login'
    || authRole === 'logout'
    || authRole === 'visitor'
    || method.url === '/api/encrypt/public/key'
}

const API = createAlova({
  baseURL: gEnv.VITE_MOCK ? '' : '',
  statesHook: reactHook,
  cacheFor: null,
  requestAdapter: adapterFetch(),
  shareRequest: false,
  beforeRequest: onAuthRequired(async (method) => {
    const token = Cookies.get(XHeader.Token) || useAccountStore.getState().token
    if (!token && !canRequestWithoutToken(method)) {
      throw new Error('request canceled: unauthenticated')
    }
    NProgress.start()
    await encryptRequest(method)
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
        const decrypted = await decryptText(encryptedText, aesKey)
        response = new Response(decrypted, {
          status: response.status,
          statusText: response.statusText,
          headers: response.headers,
        })
      }
      if (contentType.includes('application/json')) {
        const json = await response.clone().json()
        const res = json as Res
        console.log(`${method.url}:response`, res)
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
