import { createFileRoute, Outlet } from '@tanstack/react-router'

export const Route = createFileRoute('/_app/account')({
  staticData: {
    menu: {
      name: '账号管理',
      menuType: 'catalog',
      order: 2,
    },
  },
  component: RouteComponent,
})

function RouteComponent() {
  return <Outlet />
}
