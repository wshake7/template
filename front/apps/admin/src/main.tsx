import { RouterProvider } from '@tanstack/react-router'
import { App, ConfigProvider } from 'antd'
import { StrictMode } from 'react'
import ReactDOM from 'react-dom/client'
import { router } from './router'
import GlobalMessage from './utils/antd'
import '~/styles/index.css'

async function prepare() {
  if (import.meta.env.MODE === 'dev') {
    const { worker } = await import('./mocks/browser')
    await waitForHydration()
    return worker.start({ onUnhandledRequest: 'bypass' })
  }
}

const Root = () => {
  const { theme } = useAntTheme()
  return (
    <ConfigProvider {...theme}>
      <App>
        <GlobalMessage />
        <RouterProvider router={router} />
      </App>
    </ConfigProvider>
  )
}

prepare().then(() => {
  // Render the app
  const rootElement = document.getElementById('root')!
  if (!rootElement.innerHTML) {
    const root = ReactDOM.createRoot(rootElement)
    root.render(
      <StrictMode>
        <Root />
      </StrictMode>,
    )
  }
})
