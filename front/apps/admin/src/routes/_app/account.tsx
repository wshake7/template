import { createFileRoute, Outlet } from '@tanstack/react-router'

export const Route = createFileRoute('/_app/account')({
  staticData: {
    menu: {
      name: '账号管理',
      menuType: 'catalog',
    },
  },
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <>
      <div>Account</div>
      <Outlet />
    </>
  )
}
