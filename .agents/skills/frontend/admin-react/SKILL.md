# Skill: Frontend Admin React

## 何时使用

当任务涉及 `front/apps/admin-react` 或 `front/packages/utils`，包括页面、路由、菜单、API、登录鉴权、主题、Mock、表单、表格、测试和构建时使用。

## 核心路径

```text
front/apps/admin-react/
├── src/main.tsx
├── src/router.ts
├── src/routeTree.gen.ts
├── src/routes/              # TanStack file routes
├── src/api/                 # Alova API 封装
├── src/domains/             # 领域类型、HTTP 状态码、Header
├── src/stores/              # Zustand store
├── src/components/          # 通用/业务组件
├── src/config/themes/       # Ant Design 主题
├── src/mocks/               # MSW mock
├── src/utils/               # 表单、加密、菜单、antd helpers
├── locales/                 # i18n 资源
└── tests/, playwright*.ts    # E2E / component test
```

## 关键架构

- 路由使用 TanStack Router 文件路由，`src/routes/**` 通过插件生成 `src/routeTree.gen.ts`。
- 根路由 `src/routes/__root.tsx` 根据 `useAccountStore` 的 token 做登录跳转。
- 应用壳在 `src/routes/_app.tsx`，使用 `ProLayout`、`PageContainer`、菜单树和页签。
- 菜单数据来自路由 `staticData.menu`，由 `src/utils/menu.ts` 转为 ProLayout 菜单。
- HTTP 层在 `src/api/index.ts`，用 Alova + token auth + 请求加密/响应解密 + NProgress。
- 业务 API 按资源拆在 `src/api/*.ts`，领域常量与响应码在 `src/domains/*.ts`。
- 状态使用 Zustand，账号状态持久化在 `src/stores/account.ts`。
- 表单校验优先使用 Zod 与 `src/utils/zod.ts` 的 `useZodForm`。

## 常见任务流程

### 新增管理页面

1. 在 `src/routes/_app/**` 新增文件路由，使用 `createFileRoute`。
2. 写入 `staticData.menu`，至少包含 `name` 和 `menuType: 'menu'`；目录型节点使用 `menuType: 'catalog'`。
3. 页面组件优先沿用现有 Ant Design Pro Components 模式，例如 `ProTable`、`ModalForm`、`ProFormText`。
4. 如需后端数据，在 `src/api/<resource>.ts` 增加 API 方法和类型，在页面里通过 Alova hooks 调用。
5. 如需 Mock，在 `src/mocks/handlers/**` 增加处理器，并在聚合文件中导出。

### 修改 API 或鉴权

1. 先读 `src/api/index.ts`、`src/domains/http.ts`、`src/utils/encrypt.ts`。
2. 保持 Header 名称与后端 `admin/internal/domains/headers.go` 对齐。
3. 登录成功路径会写入 token、publicKey、Cookie，并更新 router context。
4. 状态码处理集中在 `HttpCodeCheck`；新增业务码时同步更新前后端常量。

### 修改主题和布局

1. 主题入口在 `src/config/themes/*`，布局入口在 `_app.tsx`。
2. 优先用 Ant Design token、`antd-style` 和现有 CSS 变量，不新增孤立样式体系。
3. 管理后台界面应保持高信息密度、克制、可扫描；避免营销页式的大 hero 和装饰性布局。

## 命令

在仓库根目录执行：

```bash
vp run admin-react#dev
vp run admin-react#build
vp staged
```

在应用目录也可执行：

```bash
pnpm --filter admin-react dev
pnpm --filter admin-react build
pnpm --filter admin-react e2e:test
```

## 验证

- 修改前端代码后，提交前执行 `vp staged`。
- 页面/路由/API 行为变化优先执行 `vp run admin-react#build`。
- 交互复杂或布局敏感时，启动开发服务并用浏览器检查桌面和移动视口。
- 只改纯文档或技能文件时不需要运行前端验证。

## 自动优化记录

<!-- ai-skill-optimizer:6eab4249bb74:2 -->
### 字典类型页面权限按钮控制

- 新建管理页面时，若列表项需按权限控制操作按钮，API 响应应包含 canWrite/canDelete 等 boolean 字段
- 在 ProTable 的 columns 中，通过 render 或 valueType 根据行数据的权限标记决定按钮 disabled 或隐藏
- 前端 API 封装在 src/api/dict.ts 中，保持与后端 Swagger 定义同步
- 页面路由文件放置于 src/routes/_app/system/dict.tsx，配合 staticData.menu 注册菜单

