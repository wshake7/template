# Skill: Frontend Admin React

## 何时使用

当任务涉及 `front/apps/admin-react` 或 `front/packages/utils`，包括页面、路由、菜单、API、登录鉴权、主题、Mock、表单、表格、测试、构建以及任务调度页面时使用。特别适合开发系统管理功能（如角色、用户、菜单、API 资源、字典、语言、API 日志、登录日志、SSE 事件、任务配置、执行记录等）。

## 核心路径

text
front/apps/admin-react/
├── src/main.tsx
├── src/router.ts
├── src/routeTree.gen.ts
├── src/routes/
│   ├── _app.tsx
│   ├── _app/account/
│   │   ├── role.tsx
│   │   └── user.tsx
│   ├── _app/system/
│   │   ├── language.tsx
│   │   ├── dict.tsx
│   │   ├── logger/
│   │   │   ├── api.log.tsx
│   │   │   └── login.log.tsx     # 登录日志页面
│   │   ├── resource.api.tsx
│   │   └── resource.menu.tsx
│   ├── _app/job/
│   │   ├── schedule.tsx          # 任务配置页面
│   │   └── execution.tsx         # 执行记录页面
│   └── _app/dashboard.tsx
├── src/api/
│   ├── index.ts              # Alova 主实例
│   ├── encryptRequest.ts     # 加密请求配置模块
│   └── business/
│       ├── account.ts
│       ├── sysResourceApi.ts
│       ├── sysResourceMenu.ts
│       ├── sysRole.ts
│       ├── sysUser.ts
│       ├── sysDict.ts
│       ├── sysLanguage.ts
│       ├── sysApiLog.ts
│       ├── sysLoginLog.ts
│       ├── jobSchedule.ts    # 任务调度 API，包含 options 接口
│       ├── jobExecution.ts   # 任务执行 API
│       └── eventStream.ts   # SSE 事件流
├── src/domains/               # 领域类型、HTTP 状态码
├── src/stores/
│   ├── account.ts             # 账号状态
│   ├── resourceMenu.ts        # 动态菜单树（持久化）
│   ├── menuTabs.ts            # 页签状态
│   └── ...
├── src/components/
│   ├── business/
│   │   └── system/
│   │       ├── dictPanels.tsx
│   │       └── languagePanels.tsx
│   └── common/
│       └── antIconPicker.tsx
├── src/hooks/
│   └── useDictMatch.tsx       # 字典批量匹配
├── src/utils/
│   ├── antIcons.tsx           # 可用图标列表
│   ├── date.ts                # 日期格式化工具
│   ├── zod.ts
│   └── ...
└── locales/, tests/, playwright*.ts


## 关键架构

- 路由使用 TanStack Router 文件路由，`src/routes/**` 通过插件生成 `src/routeTree.gen.ts`。
- 根路由 `src/routes/__root.tsx` 根据 token 做登录跳转。
- 应用壳 `src/routes/_app.tsx` 使用 `ProLayout`、`PageContainer`，侧边栏菜单来源于后端 `sys_resource_menu` 表，通过 `useResourceMenuStore` 拉取并持久化缓存。该 store 使用 zustand persist，刷新页面后快速恢复菜单树。
  - 菜单树请求的 effect 依赖 `accountToken`（store 中的 token）。当 token 为空时（如登出、未登录），会将动态菜单设为空数组，不再发起请求，避免未认证请求。
- 系统管理模块下的页面路由统一放在 `src/routes/_app/system/` 目录中；日志类页面（如 API 日志、登录日志）放在 `src/routes/_app/logger/` 目录中；账号管理相关页面（如用户、角色）放在 `src/routes/_app/account/` 目录中；任务调度相关页面放在 `src/routes/_app/job/` 目录中。
- HTTP 层在 `src/api/index.ts`，用 Alova + token auth + 请求加密/响应解密 + NProgress。
  - **认证拦截**：`beforeRequest` 中新增未认证检查。获取 token（优先 cookie，否则 store），若 token 为空且请求不是 login/logout/visitor 或公钥请求，则抛出错误，阻止请求发出。`canRequestWithoutToken` 函数通过 `authRole` 或 URL 判断是否允许无 token 请求。
  - 加密逻辑已抽离至 `src/api/encryptRequest.ts`，提供 `encryptRequest`、`createEncryptedRequestConfig`、`decryptText`、`ensurePublicKey` 等工具函数，可在 Alova `beforeRequest` 及其他场景（如 SSE）中复用。
  - `beforeRequest` 直接调用 `encryptRequest(method)` 完成加密配置。
  - 登录成功后的响应解密使用 `decryptText`。
  - 登录成功后默认跳转到根路径 `/`（而非 `/dashboard`）。登出时，先清空 store 中的 token 和公钥、移除 cookie，并调用 `router.update` 清空路由上下文中的 token，然后发送 logout 请求并导航到 `/login`。
- 根路径 `/_app/` 不再强制重定向到 `/dashboard`，仅渲染空组件，实际首页内容通过侧边栏菜单导航到具体页面。
- 字典数据使用 `useDictMatch` 钩子批量获取并缓存启用字典项，返回 `{ value, label }` 格式的映射，适用于下拉选项、表格列渲染。
  - 字典条目新增 `label_component` 字段，可存放 JSX 模板（如 `<Tag color="success">${EntryLabel}</Tag>`），前端渲染时通过安全的替换机制生成最终展示内容（例如读取 `label_component` 后替换 `${EntryLabel}` 为真实标签文本）。
- 图标选择使用 `AntIconPicker` 组件，可选图标列表由 `src/utils/antIcons.tsx` 提供。
- 表单校验优先使用 Zod 与 `src/utils/zod.ts` 的 `useZodForm`。
- 页签（tabs）行为由 `src/stores/menuTabs.ts` 的 `useMenuTabsStore` 与 `_app.tsx` 共同管理，支持动态打开、关闭、刷新、全部关闭、新窗口打开等操作，并通过右键菜单触发。
- SSE 事件流通过 `src/api/eventStream.ts` 连接后端 `/api/events` 端点。
  - SSE 实例有独立的 `beforeRequest` 钩子：检查 store 中的 token，若不存在则抛出错误；否则设置请求头 `X-Token`。
  - 连接前使用 `createEncryptedRequestConfig` 获取加密配置和 AES 密钥。采用 token 快照机制，连接开始时保存当前 token，后续校验中使用该快照，防止 token 变更导致竞态问题。若加密配置失败或 token 已变更，则放弃连接并重置状态。
  - token 为空时（如登出）会关闭 SSE 连接并清除所有引用（AES 密钥、headers）。组件卸载时同样清理引用。
  - 事件回调中使用 `decryptText` 解密 `payload` 字段后得到实际数据。

## 请求加密模块 `encryptRequest.ts`

该文件提供以下核心功能：

- **`encryptRequest(method)`**：处理普通 Alova 方法的加密。如果请求 URL 是 `/api/encrypt/public/key` 则仅添加时间戳和 Nonce；否则调用 `createEncryptedRequestConfig` 并将生成的 headers/data 写入 method 对象，同时将 `aesKey` 和 `nonce` 存入 `meta` 供后续解密使用。
- **`createEncryptedRequestConfig(options?)`**：独立于 Alova，可用于任何需要生成加密请求配置的场景。它执行完整的加密握手：
  1. 确保公钥存在（`ensurePublicKey()`）。
  2. 生成临时 AES 密钥。
  3. 用 RSA 加密 AES 密钥。
  4. 构建请求签名所需的附加数据（`uriSort` of timestamp, nonce, params）。
  5. 用 AES 加密请求体并生成签名。
  6. 返回 `{ headers, data, aesKey, nonce }`。
- **`decryptText(encryptedText, aesKey)`**：用于解密响应体，封装 `aesDecrypt`。
- **`ensurePublicKey()`**：检查 store 中的公钥，若不存在则从 `/api/encrypt/public/key` 获取并缓存。

这些函数已通过自动导入全局可用，无需显式 import。

## 标签页缓存渲染与导航优化

- 使用 `useReducer` 管理缓存状态，`cachedTabPaneReducer` 处理导航、删除、刷新等动作，维护每个标签页的版本号和上次隐藏时间。
- 当页面切换时，非活动标签页的内容通过样式隐藏而不是销毁，保留其 DOM 状态，提高切换性能。
- 渲染标签页内容统一使用 `CachedTabPaneContent` 组件（由 `React.memo` 包裹），该组件直接从路由表中获取对应的页面组件，若无对应组件则回退到 `<Outlet />`。
- 每个标签页容器的样式由 `getCachedTabPaneStyle(active)` 动态生成，利用 `opacity`、`transform` 和 `pointer-events` 实现淡入淡出的视觉过渡，同时设置 `willChange` 提示浏览器准备动画。
- 引入 `interactivePathname` 状态（`useReducer` 管理），导航时优先更新该值（例如菜单点击、标签点击、关闭全部时），用于立即驱动 UI 高亮（如 `ProLayout` 的 `location.pathname`、`PageContainer` 的 `activeKey`）；而实际路由跳转通过 `startTransition` 包裹，作为低优先级更新，避免阻塞用户交互。
- 关闭所有标签页时，先将 `interactivePathname` 设为 `'/'` 以即时反映 UI 状态，再执行导航。
- 路由状态变化时（`pathname` 改变），通过 `useEffect` 同步 `interactivePathname`，保持最终一致。

## 常见任务流程

### 新增系统管理页面

1. 在 `src/routes/_app/system/`（系统资源）、`src/routes/_app/account/`（账号管理）、`src/routes/_app/logger/`（日志）或 `src/routes/_app/job/`（任务调度）下创建文件路由，使用 `createFileRoute`。
2. 在文件路由的 `staticData.menu` 中配置菜单：`name` 为页面标题，`menuType: 'menu'`；父级目录已在对应 `system.tsx`、`account.tsx`、`logger.tsx` 或 `job.tsx` 中定义为 `catalog`。
3. 页面组件优先使用 Ant Design Pro Components，参照 `resource.menu.tsx`、`login.log.tsx`、`job/schedule.tsx` 等实现。
4. 在 `src/api/business/<resource>.ts` 中定义 API 方法和类型，URL 路径使用 `/api/sys/<resource>/*` 格式（与 Swagger 一致）。
5. 若列表需要状态下拉，使用 `useDictMatch('system:is_enabled')` 获取选项。
6. 若表单需要图标字段，使用 `AntIconPicker` 组件。
7. 实现完成后，启动开发服务验证路由、菜单展现和数据交互。

### 登录日志页面

- 页面组件 `src/routes/_app/logger/login.log.tsx` 展示所有登录尝试记录，支持按用户名、IP、状态等筛选。
- 使用 `alova/client` 的 `usePagination` 分页加载数据，表格列包含用户名、登录 IP、设备信息、状态码、是否成功、登录时间等。
- 状态列渲染使用 `Tag` 组件，成功显示“成功”并用绿色，失败显示“失败”并用红色。
- 详情功能通过 `POST /api/sys/login/log/detail` 获取单条日志完整信息。
- 实时事件通过 `src/api/eventStream.ts` 订阅，页面需要时可监听登录事件更新列表。

### 任务调度页面

- 任务配置页面 `src/routes/_app/job/schedule.tsx` 展示任务调度列表，支持创建、编辑、启用/停用、同步、触发等操作。
- 任务执行记录页面 `src/routes/_app/job/execution.tsx` 展示执行历史，支持查看详情、取消运行和重试。
- API 定义在 `src/api/business/jobSchedule.ts` 和 `src/api/business/jobExecution.ts` 中，遵循 Swagger 路径。现在已包含 `options` 接口 (`/api/sys/job/schedule/options`)，用于获取调度配置选项。
- 调度管理涉及 Temporal 状态，通过后端同步接口与 Temporal 交互。
- 任务状态使用字典或自定义枚举显示，如启用/停用、执行状态（运行中、已完成、失败等）。

### 使用字典匹配

tsx
import { useDictMatch } from '@/hooks/useDictMatch';

function MyPage() {
  const { dictMap } = useDictMatch(['system:is_enabled']);

  const statusOptions = dictMap['system:is_enabled'] || [];
  // statusOptions => [{ value: '1', label: '启用' }, { value: '0', label: '停用' }]
}


`useDictMatch` 内部会批量请求 `/api/sys/dict/entry/match`，并缓存结果。

### 使用图标选择器

tsx
import AntIconPicker from '@/components/common/antIconPicker';

<Form.Item name="icon" label="图标" valuePropName="value">
  <AntIconPicker />
</Form.Item>


需确保 `src/utils/antIcons.tsx` 已导出所有可选图标。

### 格式化日期和时间

表格列或详情中展示时间时，使用全局可用的 `formatDateYYYYMMDDHHmm(value)`。该函数按 `YYYY-MM-DD HH:mm` 格式化日期字符串，无效输入返回原值或 `-`。已通过自动导入全局注册，无需显式 `import`。

tsx
// 表格列中使用
{
  title: '创建时间',
  dataIndex: 'createdAt',
  render: (_, record) => formatDateYYYYMMDDHHmm(record.createdAt),
}


如果需要其他日期格式，可以使用已安装的 `dayjs` 库直接调用，`dayjs` 已在 `package.json` 中按 `catalog:build` 引入。

### 修改菜单资源后刷新侧边栏

- 在 `resource.menu.tsx` 页面增删改菜单后，调用 `useResourceMenuStore.getState().refresh()` 强制刷新侧边栏缓存。

### 修改 API 或鉴权

- 前端 API 路径变更后，同步更新 `src/api/business/` 对应文件的 URL。
- Header 名称保持与后端 `admin/internal/domains/headers.go` 对齐。
- 业务码处理集中在 `HttpCodeCheck`，新增业务码同步更新前后端常量。
- 加密相关功能统一使用 `encryptRequest.ts` 提供的函数，避免在业务代码中重复实现加密握手。

### 处理登出与未认证状态

- 登出时务必按顺序执行：
  1. 清空 account store 中的 token（调用 `useAccountStore.getState().logout()`）。
  2. 清空公钥（`setPublicKey('')`）。
  3. 移除 cookie 中的 token。
  4. 通过 `router.update` 清空路由上下文中的 token，防止路由守卫误判。
  5. 发送 `/api/account/logout` 请求（该请求会携带 logout 的 `authRole`，允许在无 token 时发送）。
  6. 导航到 `/login`。
- 菜单树和 SSE 连接依赖 store 中的 token。token 为空时，`_app.tsx` 会将动态菜单置空，`eventStream` 会关闭连接并清理资源，避免未认证请求。

### 标签页操作

管理后台支持通过 `PageContainer` 的 `tabList` 管理多个打开的页签，并提供右键菜单快捷操作。

- **打开页签**：点击侧边栏菜单时会自动添加对应页签（如果当前路径不是目录）。
- **关闭页签**：右键选择“关闭”或点击页签上的 × 按钮，会关闭该页签并自动切换到相邻页签或首页。
- **关闭所有页签**：右键选择“全部关闭”清空所有页签并导航到首页。该操作会立即将 `interactivePathname` 设为 `'/'` 以更新 UI，随后通过 `startTransition` 执行路由跳转。
- **刷新页签**：右键选择“刷新”会重新渲染对应页签的内容（通过递增版本号强制更新）。
- **新窗口打开**：右键选择“新窗口打开”在新标签页中打开该页面对应的 URL。

在 `_app.tsx` 中，这些操作由 `closeTab`、`closeAllTabs`、`refreshTab`、`openTabInNewWindow` 等回调实现，并通过 `tabList` 的 Dropdown 为每个标签项注入右键菜单。相关的状态管理使用 `useReducer`
