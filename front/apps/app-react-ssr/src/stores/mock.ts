import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { immer } from 'zustand/middleware/immer'

interface MockStore {
  account: {
    username: string
    pwd: string
    token: string
  }
  updatePwd: (pwd: string) => void
  _hasHydrated: boolean
  setHasHydrated: (state: boolean) => void
}

export const useMockStore = create<MockStore>()(
  persist(
    immer(set => ({
      account: {
        username: 'admin',
        pwd: '123456',
        token: 'mock_token',
      },
      updatePwd: (pwd: string) =>
        set((state) => {
          state.account.pwd = pwd
        }),
      _hasHydrated: false,
      setHasHydrated: (state: boolean) =>
        set((s) => {
          s._hasHydrated = state
        }),
    })),
    {
      name: 'mock-store',
      onRehydrateStorage: () => (state) => {
        state?.setHasHydrated(true)
      },
    },
  ),
)

// 导出一个等待水合完成的 Promise
export const waitForHydration = (): Promise<void> => {
  return new Promise((resolve) => {
    // 已经水合完成，直接 resolve
    if (useMockStore.getState()._hasHydrated) {
      resolve()
      return
    }
    // 订阅状态变化，等待水合完成
    const unsub = useMockStore.subscribe((state) => {
      if (state._hasHydrated) {
        unsub()
        resolve()
      }
    })
  })
}
