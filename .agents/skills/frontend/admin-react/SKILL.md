# Skill: Frontend App Admin React

## 何时使用

当任务涉及 `front/apps/admin-react` 的后台管理功能时使用，包括登录鉴权、动态菜单、系统管理、账号管理、日志、任务调度、Ant Design / Pro Components 页面、请求加密、SSE、Mock、Playwright 和构建配置。

## 模块职责

`admin-react` 是后台管理应用，不是通用模板。它承载业务最重的管理端实现：

- 路由：`src/routes/**`，使用 TanStack Router 文件路由。
- 管理壳：`src/routes/_app.tsx`，使用 `ProLayout`、`PageContainer`、动态菜单和标签页。
- 业务页面：`src/routes/_app/account/**`、`system/**`、`logger/**`、`job/**`。
- 业务 API：`src/api/business/**`。
- 状态：`src/stores/**`，包含账号、动态菜单、标签页、设备、Mock 等状态。
- 通用组件：`src/components/**`，包含后台专用业务组件和 Ant Design 封装。

## 调用关系

```text
admin-react
├── 调用 @vp/core
│   └── 账号、字典、加密、HTTP、通知器、页面等类型和常量
├── 调用 @vp/react-core
│   └── env、i18n、字典匹配 hooks、store factory、Mock helper
├── 调用 @vp/request
│   └── Alova client、token/cookie、请求加密、响应解密、NProgress、通知注入
├── 调用 @vp/utils
│   └── cn、Web Crypto、日期格式化等浏览器安全工具
└── 调用 @vp/build-config
    └── Vite 分包、MSW 清理、Playwright 配置
```

`admin-react` 可以调用共享包，但共享包不能反向调用 `admin-react`。

## 常见任务

### 新增后台页面

1. 在 `src/routes/_app/**` 下按业务域创建文件路由。
2. 在 route 的 `staticData.menu` 配置菜单元数据。
3. 在 `src/api/business/<resource>.ts` 定义业务 API。
4. 表格、表单、字典、图标、日期、通知、加密请求优先复用现有封装。
5. 修改菜单资源后，必要时调用 `useResourceMenuStore.getState().refresh()` 刷新侧边栏缓存。

### 修改登录/鉴权/请求

1. 先检查 `@vp/request` 是否已有可复用能力。
2. 应用独有适配放在 `src/api/**`、`src/stores/account.ts`、`src/utils/notifier.ts`。
3. 不要把 Ant Design 的 `message`、`notification` 或后台业务状态写入共享包。

### 修改标签页或动态菜单

1. 标签页逻辑优先看 `src/stores/menuTabs.ts` 和 `src/routes/_app.tsx`。
2. 动态菜单逻辑优先看 `src/stores/resourceMenu.ts`。
3. 保持路由状态、交互高亮状态和持久化缓存的一致性。

## 验证

```bash
vp staged
vp run admin-react#build
```

只改文档或 skill 不需要跑前端构建；改 `front/**` 代码后至少执行 `vp staged`。

