# Admin 状态管理

## 场景与目标

- 适用场景：新增或修改 `front/apps/admin-react/src/stores/**`。
- 目标：全局状态统一使用 Zustand，组件局部状态留在组件内。

## 目录/文件位置

- 账户：`front/apps/admin-react/src/stores/account.ts`
- 设备：`front/apps/admin-react/src/stores/device.ts`
- 主题：`front/apps/admin-react/src/stores/theme.ts`
- 多标签页：`front/apps/admin-react/src/stores/menuTabs.ts`
- Mock 开关：`front/apps/admin-react/src/stores/mock.ts`

## 当前约定

- 全局共享、跨页面复用、需要持久化的状态才进 store。
- 局部表单、弹窗、表格选择等状态优先 `useState` / `useReducer`。
- 复杂不可变更新可组合 `immer`。
- 需要持久化时使用 Zustand `persist`。

## 常用模式

```typescript
import { useAccountStore } from '~/stores/account'

const { token } = useAccountStore()
useAccountStore.getState().login(token)
useAccountStore.getState().logout()
```

```typescript
import { useAntTheme, useThemeStore } from '~/stores/theme'

const current = useThemeStore.getState().themeType
useThemeStore.getState().setThemeType('default')

const { theme, themeType, setThemeType } = useAntTheme()
```

## 新增 store 步骤

1. 先确认状态是否跨页面共享或需要持久化。
2. 在 `src/stores/<name>.ts` 定义 state、actions 和 hook。
3. action 命名使用动词，如 `setThemeType`、`login`、`logout`。
4. 需要在路由上下文或顶层组件使用时，同步检查 `src/router.ts` 和 `_app.tsx`。
5. 补充使用处并执行 `vp staged`。

## 注意事项

1. 不为单页面弹窗、搜索词、分页状态建立 store。
2. store 中不直接写 UI 组件逻辑。
3. 持久化字段要谨慎，避免保存临时或敏感信息。
4. 账户、设备和主题已有约定，优先复用已有 store。
