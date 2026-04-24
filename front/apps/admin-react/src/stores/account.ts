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
        useAccountStore.persist.clearStorage()
        set((state) => { state.token = '' })
      },
    })),
    {
      name: 'account-store',
    },
  ),
)
