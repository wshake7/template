import type { MenuDataItem } from '@ant-design/pro-components/es/layout/typing'
import { createRouter } from '@tanstack/react-router'
import { routeTree } from './routeTree.gen'

// Create a new router instance
export const router = createRouter({
  routeTree,
  defaultPreload: 'intent',
  scrollRestoration: true,
  defaultStructuralSharing: true,
  context: {
    account: useAccountStore.getState().account(),
  },
})

type MenuType = 'menu' | 'catalog' | 'title'
// Register the router instance for type safety
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router
  }
  interface StaticDataRouteOption {
    menu?: {
      menuType?: MenuType
      order?: number
    } & MenuDataItem
  }
}
