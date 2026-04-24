import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { immer } from 'zustand/middleware/immer'
import { base64ToArrayBuffer } from '~/utils/encrypt'

interface DeviceState {
  publicKey: string
}

interface DeviceActions {
  setPublicKey: (key: string) => void
  getPublicCryptoKey: () => Promise<CryptoKey | undefined>
}

export const useDeviceStore = create<DeviceState & DeviceActions>()(
  persist(
    immer((set, get) => ({
      publicKey: '',
      setPublicKey(key) {
        set((state) => {
          state.publicKey = key
        })
      },
      getPublicCryptoKey: async () => {
        const state = get()
        if (state.publicKey !== '') {
          const keyData = base64ToArrayBuffer(state.publicKey)
          const publicKey = await window.crypto.subtle.importKey(
            'spki',
            keyData,
            { name: 'RSA-OAEP', hash: 'SHA-256' },
            false,
            ['encrypt'],
          )
          return publicKey
        }
        return undefined
      },
    })),

    {
      name: 'device-store',
    },
  ),
)
