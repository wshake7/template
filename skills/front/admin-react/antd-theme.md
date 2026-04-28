# Admin Ant Design 主题

## 场景与目标

- 适用场景：调整 Ant Design token、主题切换、全局样式或主题配置。
- 目标：让主题配置集中在 `config/themes` 和 `stores/theme.ts`，避免散落覆盖。

## 目录/文件位置

- 主题配置：`front/apps/admin-react/src/config/themes/`
- 主题状态：`front/apps/admin-react/src/stores/theme.ts`
- 全局样式：`front/apps/admin-react/src/styles/index.css`
- 应用入口：`front/apps/admin-react/src/routes/_app.tsx`

## 现有主题

- `useDefaultTheme.ts`：默认主题。
- `useCartoonTheme.ts`：卡通风格主题。
- `useShadcnTheme.ts`：shadcn 风格主题。

## 当前约定

1. Ant Design 配置通过顶层 `ConfigProvider` 注入。
2. 主题类型由 `useThemeStore` 管理。
3. 组件库视觉优先通过 token 和 component token 调整。
4. 全局 CSS 只放通用基础样式和必要变量，不承载页面级样式。

## 示例

```typescript
import { theme } from 'antd'

export default function useDefaultTheme() {
  return {
    theme: {
      algorithm: theme.defaultAlgorithm,
      token: {
        colorPrimary: '#1677ff',
        borderRadius: 6,
      },
    },
  }
}
```

## 修改步骤

1. 判断是全局 token、组件 token，还是单页面局部样式。
2. 全局视觉改动优先进入 `src/config/themes/*.ts`。
3. 主题切换逻辑进入 `src/stores/theme.ts`。
4. 页面专属布局或间距留在页面组件或局部样式中。
5. 改动后执行 `vp staged`；影响明显 UI 时建议启动 Admin 预览。

## 注意事项

1. 不在多个页面复制同一套 token 覆盖。
2. 主题名称、store 联动和可选值要保持同步。
3. 新增主题时同时检查 `_app.tsx` 的选择逻辑。
4. Ant Design API 用法问题必须按 `AGENTS.md` 使用 Context7。
