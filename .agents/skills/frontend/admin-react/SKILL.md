# Skill: Frontend App Admin React

## 何时使用

当任务涉及 `front/apps/admin-react` 的后台管理功能时使用，包括登录鉴权、动态菜单、系统管理、账号管理、日志、任务调度、Ant Design / Pro Components 页面、请求加密、SSE、Mock、Playwright 和构建配置。

## 模块职责

`admin-react` 是后台管理应用，不是通用模板。它承载业务最重的管理端实现：

- 路由：`src/routes/**`，使用 TanStack Router 文件路由。
- 管理壳：`src/routes/_app.tsx`，使用 `ProLayout`、`PageContainer`、动态菜单和标签页。
- 业务页面：`src/routes/_app/account/**`、`system/**`、`logger/**`、`job/**`。
- 业务 API：`src/api/business/**`。
- API 客户端：`src/api/index.ts`，基于 `@vp/request` 的 `createVpApiClient` 构建，注入 token 管理、加密、HttpCode 检查和通知器。
- 加密请求：`src/api/encryptRequest.ts`，通过 `@vp/request` 的 `createEncryptedRequestHelpers` 创建，绑定设备存储和通知器。
- 状态：`src/stores/**`，包含账号、动态菜单、标签页、设备、Mock 等状态。
- 通用组件：`src/components/**`，包含后台专用业务组件和 Ant Design 封装。
- 工具：消息提示（`src/utils/message.ts`）、通知器（`src/utils/notifier.ts`）、加密（`src/utils/encrypt.ts`）等。
- 共享类型：优先使用 `@vp/core` 提供的通用类型（如 `PagingRequest`, `PagingResult`, `DEFAULT_PAGE_SIZE`, 加密相关类型），减少应用内重复定义。

## 调用关系

text
admin-react
├── 调用 @vp/core
│   └── 账号、字典、加密、HTTP、通知器、页面等类型和常量（如 DEFAULT_PAGE_SIZE）
├── 调用 @vp/react-core
│   └── env、i18n、字典匹配 hooks、store factory、Mock helper
├── 调用 @vp/request
│   └── createVpApiClient：构建 API 客户端
│        createEncryptedRequestHelpers：生成加密请求助手
│        Alova 基础能力、NProgress、通知注入
├── 调用 @vp/utils
│   └── cn、Web Crypto、日期格式化等浏览器安全工具
└── 调用 @vp/build-config
    └── Vite 分包、Playwright 配置工厂（createPlaywrightConfig / createPlaywrightCtConfig）、MSW 清理


`admin-react` 可以调用共享包，但共享包不能反向调用 `admin-react`。

## 常见任务

### 新增后台页面

1. 在 `src/routes/_app/**` 下按业务域创建文件路由。
2. 在 route 的 `staticData.menu` 配置菜单元数据。
3. 在 `src/api/business/<resource>.ts` 定义业务 API，分页参数和返回类型统一使用 `@vp/core` 的 `PagingRequest`、`PagingResult`。
4. 表格、表单、字典、图标、日期、通知、加密请求优先复用现有封装。
5. 修改菜单资源后，必要时调用 `useResourceMenuStore.getState().refresh()` 刷新侧边栏缓存。
6. 使用 `@vp/core` 中的 `DEFAULT_PAGE_SIZE` 作为分页默认值，不要硬编码。

### 修改登录/鉴权/请求

1. API 客户端核心逻辑在 `src/api/index.ts` 中，基于 `createVpApiClient` 配置，不要绕过该函数自行组装 Alova 实例。
2. 应用特有的 token 管理、通知、加密回调通过配置对象传入（如 `getToken`, `setToken`, `encryptRequest`, `notifier` 等）。
3. 加密请求相关的公共密钥获取、加密逻辑集中在 `src/api/encryptRequest.ts`，通过 `createEncryptedRequestHelpers` 生成，并绑定 `useDeviceStore` 和 `appNotifier`。
4. 共享包 `@vp/request` 提供了通用的请求加密、响应解密、Http 状态码检查等能力；只在必要时在应用层做适配。
5. 消息提示统一使用 `src/utils/message.ts`（`gMessage` / `message`），通知器使用 `src/utils/notifier.ts` 的 `appNotifier`，不要直接使用 Ant Design 的静态方法。

### 修改标签页或动态菜单

1. 标签页逻辑优先看 `src/stores/menuTabs.ts` 和 `src/routes/_app.tsx`。
2. 动态菜单逻辑优先看 `src/stores/resourceMenu.ts`。
3. 保持路由状态、交互高亮状态和持久化缓存的一致性。

### 调整 Playwright 配置

1. 不再手写冗长的 `playwright.config.ts` / `playwright-ct.config.ts`。
2. 直接调用 `@vp/build-config` 暴露的 `createPlaywrightConfig()` 或 `createPlaywrightCtConfig()` 并导出。
3. 共享包内已包含 CI 环境变量、浏览器、端口（ctPort: 3100）等通用配置，无需在应用内重复声明。

### 使用共享类型和工具

- 页面相关的常量（如 `DEFAULT_PAGE_SIZE`）、分页请求/返回类型、加密相关类型（`ResPublicKey`）、字典匹配类型等均从 `@vp/core` 导入。
- 避免在 `src/domains/` 中重复定义已由共享包提供的类型，逐步迁移到共享包。

## 验证

bash
vp staged
vp run admin-react#build


只改文档或 skill 不需要跑前端构建；改 `front/**` 代码后至少执行 `vp staged`。
