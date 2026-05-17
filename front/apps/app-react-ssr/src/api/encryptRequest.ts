import { createEncryptedRequestHelpers } from '@vp/request'

const encryptedRequestHelpers = createEncryptedRequestHelpers({
  getCachedPublicKey: () => useDeviceStore.getState().publicKey,
  setPublicKey: publicKey => useDeviceStore.getState().setPublicKey(publicKey),
  getPublicCryptoKey: () => useDeviceStore.getState().getPublicCryptoKey(),
})

export const ensurePublicKey = encryptedRequestHelpers.ensurePublicKey
export const encryptRequest = encryptedRequestHelpers.encryptRequest
export const createEncryptedRequestConfig = encryptedRequestHelpers.createEncryptedRequestConfig
export const decryptText = encryptedRequestHelpers.decryptText
