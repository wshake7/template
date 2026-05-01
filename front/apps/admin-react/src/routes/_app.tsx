import type { MenuDataItem, ProSettings } from '@ant-design/pro-components'
import type { FormListFieldData } from 'antd'
import type { ComponentType } from 'react'
import type { ChangePwdFormValues } from '~/components/business/account/changePwdModal'

import { LockOutlined, LogoutOutlined } from '@ant-design/icons'
import {
  PageContainer,
  ProCard,
  ProConfigProvider,
  ProLayout,
  SettingDrawer,
} from '@ant-design/pro-components'

import {
  createFileRoute,
  Outlet,
  useNavigate,
  useRouter,
  useRouterState,
} from '@tanstack/react-router'

import { ConfigProvider, Dropdown } from 'antd'

import useApp from 'antd/es/app/useApp'

import { useRef, useState } from 'react'

import ChangePwdModal from '~/components/business/account/changePwdModal'
import { TAB_REFRESH_INTERVAL } from '~/config/tabs'
import { buildMenuTree } from '~/utils/menu'

export const Route = createFileRoute('/_app')({
  component: AppLayout,
})

interface CachedTabPane {
  lastHiddenAt?: number
  version: number
}

function AppLayout() {
  const [config, setConfig] = useState<Partial<ProSettings>>({
    fixSiderbar: true,
    layout: 'mix',
    fixedHeader: true,
    menu: {
      defaultOpenAll: false,
    },
  })

  const app = useApp()
  const navigate = useNavigate()
  const router = useRouter()

  // ⭐ 唯一数据源
  const pathname = useRouterState({
    select: s => s.location.pathname,
  })

  const items = useMenuTabsStore(state => state.items)
  const add = useMenuTabsStore(state => state.add)
  const remove = useMenuTabsStore(state => state.remove)
  const menuItems = buildMenuTree(router)
  const previousPathRef = useRef(pathname)
  const [cachedTabPanes, setCachedTabPanes] = useState<Record<string, CachedTabPane>>({})

  /** ---------------- tab 初始化（路由驱动） ---------------- */

  const currentMenuTab = useMemo(() => {
    const allRoutes = Object.values(router.routesByPath)
    const currentRoute = allRoutes.find(r => r.fullPath === pathname)
    const menu = currentRoute?.options.staticData?.menu

    if (!menu || menu.menuType === 'catalog') { return null }

    return {
      key: currentRoute.fullPath,
      label: menu.name,
      icon: menu.icon,
    }
  }, [pathname, router.routesByPath])

  useEffect(() => {
    if (!currentMenuTab) { return }
    if (!items.some(item => item.key === currentMenuTab.key)) {
      add(currentMenuTab)
    }
  }, [currentMenuTab, items, add])

  useEffect(() => {
    const previousPath = previousPathRef.current
    const now = Date.now()

    setCachedTabPanes((panes) => {
      const next = { ...panes }

      if (previousPath !== pathname && next[previousPath]) {
        next[previousPath] = {
          ...next[previousPath],
          lastHiddenAt: now,
        }
      }

      if (currentMenuTab) {
        const currentPane = next[currentMenuTab.key]
        if (!currentPane) {
          next[currentMenuTab.key] = { version: 0 }
        }
        else if (currentPane.lastHiddenAt && now - currentPane.lastHiddenAt > TAB_REFRESH_INTERVAL) {
          next[currentMenuTab.key] = {
            version: currentPane.version + 1,
          }
        }
        else {
          next[currentMenuTab.key] = {
            ...currentPane,
            lastHiddenAt: undefined,
          }
        }
      }

      return next
    })

    previousPathRef.current = pathname
  }, [currentMenuTab, pathname])

  const renderedTabs = useMemo(() => {
    const tabMap = new Map(items.map(item => [item.key, item]))
    if (currentMenuTab) {
      tabMap.set(currentMenuTab.key, currentMenuTab)
    }
    return [...tabMap.values()]
  }, [currentMenuTab, items])
  /** ---------------- handlers ---------------- */

  function onMenuClick(item: MenuDataItem & { isUrl: boolean, onClick: () => void }) {
    navigate({ to: item.path })
  }

  function onTabClick(key: string) {
    navigate({ to: key })
  }

  function renderCachedTabPane(path: string) {
    const routesByPath = router.routesByPath as unknown as Record<string, { options: { component?: ComponentType } }>
    const route = routesByPath[path]
    const Component = route?.options.component as ComponentType | undefined

    if (!Component) {
      return path === pathname ? <Outlet /> : null
    }

    return <Component />
  }

  async function submitChangePwd(values?: ChangePwdFormValues, error?: FormListFieldData) {
    if (error || !values) {
      app.message.error('提交失败,请检查输入')
      return false
    }

    try {
      await AccountApi.changePwd({
        oldPwd: values.oldPwd,
        newPwd: values.newPwd,
      })
    }
    catch {
      return false
    }
    app.message.success('修改密码成功')
    return true
  }

  if (typeof document === 'undefined') {
    return <div />
  }

  return (
    <div
      id="test-pro-layout"
      style={{
        height: '100vh',
        overflow: 'auto',
      }}
    >
      <ProConfigProvider hashed={false}>
        <ConfigProvider
          getTargetContainer={() => document.getElementById('test-pro-layout') || document.body}
        >
          <ProLayout
            route={{ routes: menuItems }}
            location={{ pathname }}
            token={{
              header: {
                // colorBgMenuItemSelected: 'rgba(0,0,0,0.04)',
              },
            }}
            menu={{
              collapsedShowGroupTitle: true,
            }}
            title="Wshake"
            menuItemRender={(item, dom) => (
              <div
                className="w-full"
                onClick={() => {
                  onMenuClick(item)
                }}
              >
                {dom}
              </div>
            )}
            avatarProps={{
              src: 'https://gw.alipayobjects.com/zos/antfincdn/efFD%24IOql2/weixintupian_20170331104822.jpg',
              size: 'small',
              title: '七妮妮',
              render: (_, dom) => (
                <Dropdown
                  menu={{
                    items: [
                      {
                        key: 'changePwd',
                        label: (
                          <ChangePwdModal username="七妮妮" onSubmit={submitChangePwd}>
                            <div>
                              <LockOutlined className="ant-dropdown-menu-item-icon" />
                              <span>修改密码</span>
                            </div>
                          </ChangePwdModal>
                        ),
                      },
                      {
                        key: 'logout',
                        icon: <LogoutOutlined />,
                        label: '退出登录',
                        onClick: async () => {
                          AccountApi.logout()
                        },
                      },
                    ],
                  }}
                >
                  {dom}
                </Dropdown>
              ),
            }}
            {...config}
          >
            <PageContainer
              fixedHeader
              title={false}
              tabList={items}
              tabProps={{
                activeKey: pathname,
                hideAdd: true,
                onChange: onTabClick,
                tabBarStyle: {
                  margin: 0,
                },
                onEdit(e, action) {
                  if (action === 'remove') {
                    if (typeof e === 'string') {
                      setCachedTabPanes((panes) => {
                        const next = { ...panes }
                        delete next[e]
                        return next
                      })
                    }
                    remove(e)
                  }
                },
                type: 'editable-card',
              }}
            >
              <ProCard style={{ minHeight: 1000 }}>
                {renderedTabs.length > 0
                  ? renderedTabs.map((item) => {
                      const pane = cachedTabPanes[item.key] ?? { version: 0 }
                      const active = item.key === pathname

                      return (
                        <div
                          key={`${item.key}:${pane.version}`}
                          style={{ display: active ? 'block' : 'none' }}
                        >
                          {renderCachedTabPane(item.key)}
                        </div>
                      )
                    })
                  : <Outlet />}
              </ProCard>
            </PageContainer>

            <SettingDrawer
              pathname={pathname}
              enableDarkTheme
              settings={config}
              onSettingChange={setConfig}
            />
          </ProLayout>
        </ConfigProvider>
      </ProConfigProvider>
    </div>
  )
}
