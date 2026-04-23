# Ant Design 主题配置

本项目支持多种主题切换，使用 antd-style 和 CSS-in-JS 方案。

## 主题目录

位置: `front/apps/admin/src/config/themes/`

## 现有主题

### 1. useDefaultTheme (默认主题)
位置: `src/config/themes/useDefaultTheme.ts`

```typescript
import { theme } from 'antd'
// 使用默认算法 (light theme)
```

### 2. useCartoonTheme (卡通主题)
位置: `src/config/themes/useCartoonTheme.ts`

使用卡通风格的配色方案。

### 3. useShadcnTheme (Shadcn 风格主题)
位置: `src/config/themes/useShadcnTheme.ts`

基于 shadcn/ui 的设计风格。

## 主题使用方式

### 在 _app.tsx 中应用主题

```typescript
import useDefaultTheme from '~/config/themes/useDefaultTheme'
import useCartoonTheme from '~/config/themes/useCartoonTheme'
import useShadcnTheme from '~/config/themes/useShadcnTheme'

const themeConfig = useDefaultTheme() // 选择主题

export default function App() {
  return (
    <ConfigProvider {...themeConfig}>
      {/* 应用主题 */}
    </ConfigProvider>
  )
}
```

### 动态切换主题

通过 `src/stores/theme.ts` 中的 Zustand store 管理主题状态:

```typescript
import { useThemeStore } from '~/stores/theme'
import { useAntTheme } from '~/stores/theme'

// 获取当前主题
const currentTheme = useThemeStore.getState().themeType

// 切换主题
useThemeStore.getState().setThemeType('default') // 'default' | 'cartoon' | 'shadcn'

// 在组件中获取 antd 可用主题对象
const { theme, themeType, setThemeType } = useAntTheme()
```

## Ant Design 主题算法

```typescript
import { theme } from 'antd'

// 亮色主题
theme.defaultAlgorithm

// 暗色主题
theme.darkAlgorithm

// 紧凑主题
theme.compactAlgorithm
```

## 主题配置选项

```typescript
interface ConfigProviderProps {
  theme?: {
    algorithm?: typeof theme.defaultAlgorithm | typeof theme.darkAlgorithm
    token?: {
      colorPrimary?: string      // 主色
      borderRadius?: number       // 圆角
      fontSize?: number           // 字体大小
      // 更多 token 配置...
    }
    components?: {
      // 组件级别覆盖
      Button?: {
        colorPrimary?: string
      }
    }
  }
}
```

## CSS 变量 (用于全局样式)

Ant Design 的 CSS 变量定义在 `src/styles/index.css` 中，可以通过覆盖 CSS 变量来自定义主题。

## 样式工具

### class-variance-authority (CVA)
用于管理组件变体样式。

### clsx
条件类名组合工具。

### @emotion/css
CSS-in-JS 解决方案，用于编写组件样式。

## 注意事项

1. 主题配置应在应用顶层 (`_app.tsx`) 通过 `ConfigProvider` 提供
2. 动态切换主题需要结合状态管理 (Zustand store)
3. 组件库样式优先使用 token 配置，其次使用 CSS 变量覆盖
