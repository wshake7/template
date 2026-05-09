import { createFileRoute, Outlet } from '@tanstack/react-router'
import { renderAntIcon } from '~/utils/antIcons'

export const Route = createFileRoute('/_app/account')({
  staticData: {
    menu: {
      name: '账号管理',
      menuType: 'catalog',
      icon: renderAntIcon('UserOutlined'),
      order: 2,
    },
  },
  component: RouteComponent,
})

function RouteComponent() {
  return <Outlet />
}
