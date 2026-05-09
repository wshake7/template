import type { ResourceMenuNode } from '~/api/business/sysResourceMenu'
import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface ResourceMenuStore {
  dynamicMenuTree: ResourceMenuNode[]
  setDynamicMenuTree: (dynamicMenuTree: ResourceMenuNode[]) => void
}

export const useResourceMenuStore = create<ResourceMenuStore>()(
  persist(
    set => ({
      dynamicMenuTree: [],
      setDynamicMenuTree: dynamicMenuTree => set({ dynamicMenuTree }),
    }),
    {
      name: 'resource-menu-store',
    },
  ),
)
