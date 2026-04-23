import type { MenuDataItem } from '@ant-design/pro-components'
import type { AnyRoute } from '@tanstack/react-router'

/** ---------------- 工具函数 ---------------- */

/**
 * 获取菜单配置
 */
export const getMenu = (route: AnyRoute) =>
  route.options.staticData?.menu

/**
 * 获取排序权重
 */
export const getOrder = (route: AnyRoute) =>
  getMenu(route)?.order ?? 99

/**
 * 排序函数
 */
export const sortByOrder = (a: AnyRoute, b: AnyRoute) =>
  getOrder(a) - getOrder(b)

/**
 * 过滤有效菜单
 */
export const filterValidMenu = (routes: AnyRoute[]) =>
  routes.filter(
    route =>
      getMenu(route)
      && !getMenu(route)!.hidden
      && (getMenu(route)!.menuType === 'menu' || getMenu(route)!.menuType === 'catalog'),
  )

/**
 * 将菜单按 parentId 分组
 */
export const groupByParentId = (routes: AnyRoute[]) => {
  const groupMap = new Map<string, AnyRoute[]>()
  for (const route of routes) {
    const parentId = route.parentRoute?.id
    if (!parentId) { continue }
    if (!groupMap.has(parentId)) { groupMap.set(parentId, []) }
    groupMap.get(parentId)!.push(route)
  }
  return groupMap
}

/**
 * 构建两级菜单（仅 _app + catalog 的一级二级）
 */
export function buildMenuItems(router: any): MenuDataItem[] {
  const routes = Object.values(router.routesByPath) as AnyRoute[]
  const validRoutes = filterValidMenu(routes)
  const groupMap = groupByParentId(validRoutes)

  const rootRoutes = groupMap.get('/_app') || []

  return rootRoutes
    .sort(sortByOrder)
    .map((route) => {
      const menu = getMenu(route)!

      const children = (groupMap.get(route.id) || [])
        .filter(child => getMenu(child)?.menuType === 'menu')
        .sort(sortByOrder)
        .map((child) => {
          const childMenu = getMenu(child)!
          return {
            path: child.fullPath,
            icon: childMenu.icon,
            name: childMenu.name,
          } as MenuDataItem
        })

      // 条件解构，只有 children.length > 0 才生成 routes
      const item: MenuDataItem = {
        path: route.fullPath,
        icon: menu.icon,
        name: menu.name,
        routes: children.length > 0 ? children as any : [],
      }

      return item
    })
}

function buildTree(parentId: string, groupMap: Map<string, AnyRoute[]>): MenuDataItem[] {
  return (groupMap.get(parentId) || [])
    .sort(sortByOrder)
    .map((route) => {
      const menu = getMenu(route)!
      const children = buildTree(route.id, groupMap)
      return {
        path: route.fullPath,
        icon: menu.icon,
        name: menu.name,
        ...(children.length ? { routes: children } : {}),
      } as MenuDataItem
    })
}

/**
 * 构建无限级菜单（递归）
 */
export function buildMenuTree(router: any): MenuDataItem[] {
  const routes = Object.values(router.routesByPath) as AnyRoute[]
  const validRoutes = filterValidMenu(routes)
  const groupMap = groupByParentId(validRoutes)
  return buildTree('/_app', groupMap)
}
