import { createFileRoute } from '@tanstack/react-router'
import { renderAntIcon } from '~/utils/antIcons'

export const Route = createFileRoute('/_app/dashboard')({
  staticData: {
    menu: {
      name: '控制台',
      menuType: 'menu',
      icon: renderAntIcon('DashboardOutlined'),
    },
  },
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <>
      <p>---------------</p>
      <div>Test Action12 Hello "/_app/dashboard"!</div>
      <p>---------------</p>
    </>
  )
}
