# 状态管理 (Zustand)

## 场景与目标
- 适用场景：新增/修改全局状态时
- 目标：统一使用 Zustand，保持 store 定义风格一致

## 目录/文件位置
- `src/stores/account.ts` - 账户状态（token、登录信息）
- `src/stores/device.ts` - 设备信息
- `src/stores/theme.ts` - 主题状态
- `src/stores/menuTabs.ts` - 多标签页状态
- `src/stores/mock.ts` - Mock 开关状态

## 核心依赖
- `zustand`：轻量状态管理
- `immer`：不可变数据
- `persist`：持久化中间件

## Account Store

位置：`src/stores/account.ts`

```typescript
import { useAccountStore } from '~/stores/account'

// 获取状态
const { token, account } = useAccountStore()

// 登录
useAccountStore.getState().login(token)

// 登出
useAccountStore.getState().logout()
```

## Theme Store

位置：`src/stores/theme.ts`

```typescript
import { useThemeStore, useAntTheme } from '~/stores/theme'

// 获取当前主题类型
const currentTheme = useThemeStore.getState().themeType

// 切换主题
useThemeStore.getState().setThemeType('default') // 'default' | 'cartoon' | 'shadcn'

// 在组件中获取 antd 可用主题对象
const { theme, themeType, setThemeType } = useAntTheme()
```

## Device Store

位置：`src/stores/device.ts` - 设备相关信息（屏幕尺寸等）

## MenuTabs Store

位置：`src/stores/menuTabs.ts` - 多标签页状态管理

## Mock Store

位置：`src/stores/mock.ts` - Mock 服务开关

## 注意事项

1. 状态管理优先使用 Zustand，避免 Redux 过度复杂
2. 全局状态才放 store，组件局部状态用 `useState`/`useReducer`
3. 需要使用不可变更新时组合 `immer` 中间件
4. 需要持久化时组合 `persist` 中间件
