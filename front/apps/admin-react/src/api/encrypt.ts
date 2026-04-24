import type { ResPublicKey } from '~/domains/encrypt'
import API from './index'

async function publicKey() {
  const res = await API.Get<Res<ResPublicKey>>('/api/encrypt/public/key', {
    cacheFor: 0,
  }).send()
  return res.data?.publicKey
}

export const EncryptApi = {
  publicKey,
}
