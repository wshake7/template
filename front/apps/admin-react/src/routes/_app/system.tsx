import { createFileRoute, Outlet } from '@tanstack/react-router'

export const Route = createFileRoute('/_app/system')({
  staticData: {
    menu: {
      name: '系统管理',
      menuType: 'catalog',
      order: 1,
    },
  },
  component: RouteComponent,
})

function RouteComponent() {
  return <Outlet />
}
