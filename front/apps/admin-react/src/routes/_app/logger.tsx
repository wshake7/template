import { createFileRoute, Outlet } from '@tanstack/react-router'

export const Route = createFileRoute('/_app/logger')({
  staticData: {
    menu: {
      name: '日志管理',
      menuType: 'catalog',
      order: 1,
    },
  },
  component: RouteComponent,
})

function RouteComponent() {
  return <Outlet />
}
