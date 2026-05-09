import { createFileRoute, Outlet } from '@tanstack/react-router'
import { renderAntIcon } from '~/utils/antIcons'

export const Route = createFileRoute('/_app/system')({
  staticData: {
    menu: {
      name: '系统管理',
      menuType: 'catalog',
      icon: renderAntIcon('SettingOutlined'),
      order: 1,
    },
  },
  component: RouteComponent,
})

function RouteComponent() {
  return <Outlet />
}
