import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_app/dashboard2')({
  staticData: {
    menu: {
      name: '控制台2',
      menuType: 'menu',
    },
  },
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/_app/dashboard2"!</div>
}
