import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_app/account/user')({
  staticData: {
    menu: {
      name: '用户管理',
      menuType: 'menu',
    },
  },
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/_app/account/user"!</div>
}
