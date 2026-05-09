import type { MenuDataItem, ProSettings } from '@ant-design/pro-components'
import type { FormListFieldData } from 'antd'
import type { ComponentType } from 'react'
import type { ResourceMenuNode } from '~/api/business/sysResourceMenu'

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

import { useEffect, useMemo, useReducer, useRef, useState } from 'react'
import { ResourceMenuApi } from '~/api/business/sysResourceMenu'
import ChangePwdModal from '~/components/business/account/changePwdModal'
import { TAB_REFRESH_INTERVAL } from '~/config/tabs'
import { useResourceMenuStore } from '~/stores/resourceMenu'
import { renderAntIcon } from '~/utils/antIcons'
import { buildMenuTree } from '~/utils/menu'

export const Route = createFileRoute('/_app')({
  component: AppLayout,
})

interface CachedTabPane {
  lastHiddenAt?: number
  version: number
}

interface CachedTabPaneNavigateAction {
  type: 'NAVIGATE'
  payload: {
    previousPath: string
    pathname: string
    currentMenuTab: { key: string } | null
  }
}

interface CachedTabPaneRemoveAction {
  type: 'REMOVE'
  payload: {
    key: string
  }
}

type CachedTabPaneAction = CachedTabPaneNavigateAction | CachedTabPaneRemoveAction

function toDynamicMenuItems(nodes: ResourceMenuNode[], router: ReturnType<typeof useRouter>): MenuDataItem[] {
  const routesByPath = router.routesByPath as unknown as Record<string, unknown>

  return nodes
    .filter(node => !node.hidden)
    .map((node) => {
      const children = toDynamicMenuItems(node.children ?? [], router)
      const isUrl = node.isUrl || /^https?:\/\//.test(node.path)
      const routeExists = isUrl || Boolean(routesByPath[node.path])
      if (!routeExists && children.length === 0) {
        return null
      }

      return {
        path: node.path,
        name: node.name,
        icon: renderAntIcon(node.icon),
        isUrl,
        routes: children.length > 0 ? children : undefined,
      } as MenuDataItem
    })
    .filter(Boolean) as MenuDataItem[]
}

function flattenMenuItems(items: MenuDataItem[]) {
  const map = new Map<string, MenuDataItem>()
  const walk = (nodes: MenuDataItem[]) => {
    for (const node of nodes) {
      if (node.path) {
        map.set(node.path, node)
      }
      const routes = (node as any).routes as MenuDataItem[] | undefined
      if (routes?.length) {
        walk(routes)
      }
    }
  }
  walk(items)
  return map
}

function cachedTabPaneReducer(
  state: Record<string, CachedTabPane>,
  action: CachedTabPaneAction,
): Record<string, CachedTabPane> {
  if (action.type === 'NAVIGATE') {
    const { previousPath, pathname, currentMenuTab } = action.payload
    const now = Date.now()
    const next = { ...state }

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
  }

  if (action.type === 'REMOVE' && action.payload.key) {
    const next = { ...state }
    delete next[action.payload.key]
    return next
  }

  return state
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
  const dynamicMenuTree = useResourceMenuStore(state => state.dynamicMenuTree)
  const setDynamicMenuTree = useResourceMenuStore(state => state.setDynamicMenuTree)
  const staticMenuItems = useMemo(() => buildMenuTree(router), [router])
  const cachedDynamicMenuItems = useMemo(
    () => dynamicMenuTree.length > 0 ? toDynamicMenuItems(dynamicMenuTree, router) : undefined,
    [dynamicMenuTree, router],
  )
  const [dynamicMenuItems, setDynamicMenuItems] = useState<MenuDataItem[]>()
  const menuItems = dynamicMenuItems?.length
    ? dynamicMenuItems
    : cachedDynamicMenuItems?.length
      ? cachedDynamicMenuItems
      : staticMenuItems
  const menuItemMap = useMemo(() => flattenMenuItems(menuItems), [menuItems])
  const previousPathRef = useRef(pathname)
  const [cachedTabPanes, dispatch] = useReducer(cachedTabPaneReducer, {})

  useEffect(() => {
    let disposed = false
    ResourceMenuApi.menuTree()
      .send()
      .then((res) => {
        if (disposed) {
          return
        }
        const nodes = res.data ?? []
        setDynamicMenuTree(nodes)
        const nextItems = toDynamicMenuItems(nodes, router)
        setDynamicMenuItems(nextItems.length > 0 ? nextItems : undefined)
      })
      .catch(() => {
        if (!disposed) {
          setDynamicMenuItems(undefined)
        }
      })
    return () => {
      disposed = true
    }
  }, [router, setDynamicMenuTree])

  /** ---------------- tab 初始化（路由驱动） ---------------- */

  const currentMenuTab = useMemo(() => {
    const allRoutes = Object.values(router.routesByPath)
    const currentRoute = allRoutes.find(r => r.fullPath === pathname)
    const menu = currentRoute?.options.staticData?.menu
    const dynamicMenu = menuItemMap.get(pathname)

    if (dynamicMenu) {
      return {
        key: pathname,
        label: dynamicMenu.name,
        icon: dynamicMenu.icon,
      }
    }

    if (!menu || menu.menuType === 'catalog') { return null }

    return {
      key: currentRoute.fullPath,
      label: menu.name,
      icon: menu.icon,
    }
  }, [menuItemMap, pathname, router.routesByPath])

  useEffect(() => {
    if (!currentMenuTab) { return }
    if (!items.some(item => item.key === currentMenuTab.key)) {
      add(currentMenuTab)
    }
  }, [currentMenuTab, items, add])

  useEffect(() => {
    dispatch({
      type: 'NAVIGATE',
      payload: {
        previousPath: previousPathRef.current,
        pathname,
        currentMenuTab,
      },
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
    if (!item.path) {
      return
    }
    if (item.isUrl || /^https?:\/\//.test(item.path)) {
      window.open(item.path, '_blank', 'noopener,noreferrer')
      return
    }
    if (!(router.routesByPath as unknown as Record<string, unknown>)[item.path]) {
      return
    }
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
                      dispatch({
                        type: 'REMOVE',
                        payload: { key: e },
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
