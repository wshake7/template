import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { immer } from 'zustand/middleware/immer'

interface AccountActions {
  account: () => AccountState
  login: (token: string) => void
  logout: () => void
}

export const useAccountStore = create<AccountState & AccountActions>()(
  persist(
    immer((set, get) => ({
      token: '',
      account() {
        const { token } = get()
        return {
          token,
        }
      },
      login(token: string) {
        set((state) => { state.token = token })
      },
      logout() {
        set((state) => { state.token = '' })
        useAccountStore.persist.clearStorage()
        useMenuTabsStore.getState().removeAll()
        useDeviceStore.getState().clear()
        useResourceMenuStore.getState().clear()
      },
    })),
    {
      name: 'account-store',
    },
  ),
)
