import { createFileRoute } from '@tanstack/react-router'
import { renderAntIcon } from '~/utils/antIcons'

export const Route = createFileRoute('/_app/account/role')({
  staticData: {
    menu: {
      name: '角色管理',
      menuType: 'menu',
      icon: renderAntIcon('TeamOutlined'),
    },
  },
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/_app/account/role"!</div>
}
