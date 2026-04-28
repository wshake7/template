# Admin React 开发框架

## 场景与目标

- 适用场景：开始任何 `front/apps/admin-react` 页面、路由、布局、菜单、API、store、主题或测试任务。
- 目标：快速定位结构、技术栈和最小执行路径。

## 技术栈

- 路由：`@tanstack/react-router`
- UI：`antd` v6、`@ant-design/pro-components`、`antd-style`
- 请求：`alova`
- 状态：`zustand`、`immer`、`persist`
- 表单校验：`zod` + `useZodForm`
- 国际化：`i18next`、`react-i18next`
- Mock：`msw`
- 测试：Playwright E2E / React CT
- 工具链：Vite+、TypeScript、ESLint

## 目录/文件位置

```text
front/apps/admin-react/
├── locales/                 # i18n 资源
├── public/                  # 静态资源、MSW worker
├── tests/                   # E2E 测试
└── src/
    ├── api/                 # Alova API
    ├── components/          # 通用和业务组件
    ├── config/themes/       # 主题配置
    ├── domains/             # 领域类型、HTTP 常量
    ├── mocks/               # MSW handlers
    ├── routes/              # TanStack 文件路由
    ├── stores/              # Zustand stores
    ├── styles/              # 全局样式
    └── utils/               # 应用工具
```

## 路由与菜单

- 路由文件位于 `src/routes/`。
- `src/routeTree.gen.ts` 是生成文件，不手动编辑。
- `src/router.ts` 扩展了 `staticData.menu` 类型。
- 菜单信息写在路由的 `staticData.menu` 中。

```typescript
export const Route = createFileRoute('/_app/system/resource')({
  staticData: {
    menu: {
      name: '资源管理',
      menuType: 'menu',
    },
  },
  staleTime: 1000 * 60 * 2,
  component: RouteComponent,
})
```

## 开发步骤

1. 用 `rg --files front/apps/admin-react/src/routes` 找目标路由或相近页面。
2. 页面任务优先读父级路由、相近页面、相关 API 文件和 `src/domains/page.ts`。
3. 新增路由时只写 `src/routes/**` 文件，不手动改 `routeTree.gen.ts`。
4. 页面内局部状态使用 React state；跨页面共享状态才进入 `src/stores/`。
5. 请求统一走 `src/api/index.ts` 导出的 Alova 实例。
6. 完成 `front/**` 改动后执行 `vp staged`。

## 常用命令

```bash
vp run admin-react#dev
vp run admin-react#build
vp run admin-react#e2e:test
vp staged
```

## 注意事项

1. 文案当前以中文为主；新增多语言需求时同步 `locales/zh.json`。
2. 生成文件只作为校验结果查看，避免手改。
3. 新增库或 API 用法问题必须按 `AGENTS.md` 使用 Context7 查询当前文档。
4. CRUD 细节见 `crud-page.md`，请求层见 `api.md`。
