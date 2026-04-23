import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_app/dashboard')({
  staticData: {
    menu: {
      name: '控制台',
      menuType: 'menu',
    },
  },
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <>
      <p>---------------</p>
      <div>Hello "/_app/dashboard"!</div>
      <p>---------------</p>
    </>
  )
}
