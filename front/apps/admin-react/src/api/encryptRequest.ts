import { createEncryptedRequestHelpers } from '@vp/request'
import { gMessage } from '~/utils/message'

const encryptedRequestHelpers = createEncryptedRequestHelpers({
  getCachedPublicKey: () => useDeviceStore.getState().publicKey,
  setPublicKey: publicKey => useDeviceStore.getState().setPublicKey(publicKey),
  getPublicCryptoKey: () => useDeviceStore.getState().getPublicCryptoKey(),
  onSystemError: () => gMessage.error('系统异常'),
})

export const ensurePublicKey = encryptedRequestHelpers.ensurePublicKey
export const encryptRequest = encryptedRequestHelpers.encryptRequest
export const createEncryptedRequestConfig = encryptedRequestHelpers.createEncryptedRequestConfig
export const decryptText = encryptedRequestHelpers.decryptText
