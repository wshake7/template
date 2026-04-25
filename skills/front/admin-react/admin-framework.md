# Admin 开发框架

本项目是基于 React + TanStack Router + Ant Design 的管理后台应用。

## 技术栈

- **路由**: @tanstack/react-router (类型安全的路由)
- **UI 框架**: Ant Design v6 + antd-style
- **状态管理**: Zustand (支持 immer 和 persist)
- **数据请求**: Alova (类 Axios 的跨平台请求库)
- **加密**: AES + RSA 混合加密

## package.json 组件清单

### 根目录 package.json（工作区入口）

位置: `package.json`

- 主要脚本:
  - `vp run admin-react#dev`：启动 admin-react 开发
  - `vp run admin-react#build`：构建 admin-react
  - `vp run lint:fix`：修复 lint
  - `vp run build -r`：递归执行工作区构建
- 关键开发依赖:
  - `vite-plus`：统一工具链入口
  - `eslint` / `@antfu/eslint-config`：代码规范与校验

### Admin package.json（front/apps/admin-react）

位置: `front/apps/admin-react/package.json`

- UI 与样式:
  - `antd`、`@ant-design/pro-components`、`antd-style`
  - `@ant-design/icons`、`@emotion/css`
  - `class-variance-authority`、`clsx`
  - `lucide-react`、`motion`
- 路由与调试:
  - `@tanstack/react-router`
  - `@tanstack/react-router-devtools`
  - `@tanstack/react-devtools`（dev）
  - `@tanstack/router-plugin`（dev）
- 状态与数据请求:
  - `zustand`、`immer`
  - `alova`
  - `js-cookie`、`nprogress`
- 国际化与校验:
  - `i18next`、`react-i18next`
  - `i18next-browser-languagedetector`
  - `zod`
- Mock 与测试:
  - `msw`
  - `@playwright/test`
  - `@playwright/experimental-ct-react`
  - `@faker-js/faker`
- 构建与工程化:
  - `vite`、`@vitejs/plugin-react`
  - `vite-plugin-pwa`
  - `unplugin-auto-import`
  - `typescript`

### 依赖版本策略（catalog）

- `catalog:build`：运行时依赖（生产依赖）
- `catalog:dev`：开发依赖（构建、测试、类型、工具）
- 版本统一在工作区 catalog 管理，避免子包各自漂移

## 目录结构

```
front/apps/admin-react/src/
├── api/          # API 请求定义
├── components/   # React 组件
│   ├── business/ # 业务组件
│   └── lib/      # 工具组件
├── config/       # 配置
│   └── themes/   # 主题配置
├── domains/      # 领域模型 (类型定义、常量)
├── lib/          # 库代码 (SEO等)
├── mocks/        # MSW 模拟服务
├── routes/       # 路由页面组件
├── stores/       # Zustand 状态库
└── utils/        # 工具函数
```

## 路由开发

### 路由文件位置
- 路由配置: `src/routeTree.gen.ts` (自动生成)
- 路由注册: `src/router.ts`
- 页面组件: `src/routes/` 目录下

### 添加新页面
1. 在 `src/routes/` 下创建页面组件
2. 在对应父路由文件中定义子路由
3. 运行 `vp dev` 自动生成 `routeTree.gen.ts`

### 菜单配置
使用 `staticData.menu` 配置菜单属性:
```typescript
// 在路由的 staticData 中配置
staticData: {
  menu: {
    path: '/dashboard',
    name: '仪表盘',
    menuType: 'menu', // 'menu' | 'catalog' | 'title'
    roles: ['admin', 'user'],
  }
}
```

菜单工具函数位于 `src/utils/menu.ts`:
- `getMenu(route)` - 获取菜单配置
- `filterValidMenu(routes)` - 过滤有效菜单
- `groupByParentId(routes)` - 按父ID分组

## 状态管理 (Zustand)

### Account Store
位置: `src/stores/account.ts`

```typescript
import { useAccountStore } from '~/stores/account'

// 获取状态
const { token, account } = useAccountStore()

// 登录
useAccountStore.getState().login(token)

// 登出
useAccountStore.getState().logout()
```

### Device Store
位置: `src/stores/device.ts` - 设备相关信息

### Theme Store
位置: `src/stores/theme.ts` - 主题状态管理

### MenuTabs Store
位置: `src/stores/menuTabs.ts` - 多标签页状态

## API 请求

### API 实例配置
位置: `src/api/index.ts`

使用 Alova 创建，配置了:
- Token 认证
- AES + RSA 加密
- 响应拦截器
- NProgress 进度条

### 定义新的 API
```typescript
// src/api/xxx.ts
import type { ResXXX } from '~/domains/xxx'
import API from '~/api'

export const xxxAPI = {
  list: (params: ListParams) =>
    API.Post<Res<ResXXX[]>>('/api/xxx/list', params).send(),
}
```

### 加密工具
位置: `src/utils/encrypt.ts`

- `generateAesKey()` - 生成 AES 密钥
- `aesEncrypt(data, key, iv)` - AES 加密
- `aesDecrypt(encrypted, key, iv)` - AES 解密
- `rsaEncrypt(data, publicKey)` - RSA 加密

## 领域模型

位置: `src/domains/`

- `account.ts` - 账户相关类型
- `constant.ts` - 常量定义
- `encrypt.ts` - 加密相关类型

## Mock 服务

位置: `src/mocks/`

使用 MSW (Mock Service Worker) 进行接口模拟:
- `handlers/index.ts` - 合并所有 handler
- `browser.ts` - 浏览器环境 mock
- `node.ts` - Node 环境 mock

## 环境变量

位置: `front/apps/admin-react/.env.*`

- `.env` - 默认环境
- `.env.dev` - 开发环境
- `.env.test` - 测试环境
- `.env.prod` - 生产环境

## 常用命令

```bash
# 根目录启动 admin-react 开发
vp run admin-react#dev

# 根目录构建 admin-react
vp run admin-react#build

# E2E 测试
vp run admin-react#e2e:test
vp run admin-react#e2e:test-ui    # UI 模式
vp run admin-react#e2e:show       # 查看报告
vp run admin-react#e2e:codegen    # 生成测试代码
```

## 列表页前端编写风格（以 system/dict.tsx 为例）

### 1) 目录与职责分层

- 页面放在 `src/routes/**`，只负责页面编排与交互流程
- 接口调用放在 `src/api/**`（如 `DictApi`），页面内尽量不拼装底层请求细节
- 类型与分页请求结构放在 `src/domains/**`，避免页面中出现大量内联类型

### 2) 数据请求与分页规范

- 列表统一使用 `usePagination` 驱动 `ProTable` 的 `current/pageSize/total`
- **默认使用后端分页**，不要先拉全量再做前端切片
- 筛选条件通过 `PagingRequest.query` 传入，使用 `JSON.stringify(...)`
- 切换上下文（如切换左侧类型）时，使用 `watchingStates` 触发重拉
- 请求参数统一显式传递：`page`、`pageSize`、`orderBy`、`query`

示例：

```typescript
const { data, total, page, pageSize, loading, update, send } = usePagination(
  (nextPage, nextPageSize) =>
    API.Post<Res<PagingResult<DictEntry>>>('/api/dict/entry/list', {
      page: nextPage,
      pageSize: nextPageSize,
      orderBy: 'sort_order asc,id desc',
      query: selectedType ? JSON.stringify({ sysDictTypeId: selectedType.id }) : undefined,
    }, { cacheFor: 0 }),
  {
    initialPage: 1,
    initialPageSize: 10,
    immediate: false,
    watchingStates: [selectedType?.id],
    data: response => response.data?.items ?? [],
    total: response => response.data?.total ?? 0,
  },
)
```

### 3) 表单与校验规范

- 表单字段用 `zod` 定义 schema，提交前统一走 `useZodForm`
- `create` 与 `edit` 共用一个 `ModalForm`，通过 `editing` 状态切换模式
- 默认值集中维护（如 `dictTypeDefaults`、`dictEntryDefaults`），避免散落硬编码
- 提交时先做必要上下文校验（如“未选中字典类型不可新增字典项”）

### 4) 表格与交互规范

- 表格列定义放在 `useMemo`，编辑函数放在 `useCallback`，减少重复渲染
- 操作列统一包含：编辑、启停、删除，并在成功后 `await send()` 刷新列表
- 行点击用于切换上下文（如“选中类型 -> 刷新右侧字典项”）
- 工具栏按钮的禁用状态与上下文绑定（如无 `selectedType` 禁用新增）

### 5) 数据映射与稳健性

- 在 `data` 适配阶段处理后端空值（如 `numericValue ?? 0`）
- 页面内只处理“展示需要的最小映射”，复杂转换下沉到 API/domain 层
- 成功提示文案统一简短（如“保存成功”“删除成功”“操作成功”）

## 开发注意事项

1. 路由变更后由开发流程自动更新 `routeTree.gen.ts`，不要手改该文件
2. API 请求使用 Alova，不直接使用 fetch/axios
3. 状态管理优先使用 Zustand，避免 Redux 过于复杂
4. 组件放在 `components/business/` 或 `components/lib/` 下
5. 工具函数放在 `utils/` 下，领域逻辑放在 `domains/` 下

## 低 Token 执行清单（Admin）

用于减少前端任务中的上下文消耗，默认按以下顺序执行：

1. 先定位再读取
   - 先用 `Glob` 找页面/接口文件（如 `routes/**`、`api/**`）
   - 再用 `Grep` 找符号（组件名、接口路径、store 名）
   - 最后只 `Read` 命中的关键文件
2. 新增页面最小读取集
   - `src/routes/_app.tsx`（布局约束）
   - 目标父级路由文件（如 `src/routes/_app/account.tsx`）
   - 目标页面文件（如 `src/routes/_app/account/user.tsx`）
   - 对应 `src/api/*.ts` 与 `src/domains/*.ts`
3. 角色/列表页优先复用模式
   - 表格优先 `ProTable`
   - 弹窗表单优先 `ModalForm`
   - 请求优先复用 `API.Get/Post(...).send()` 封装
4. 避免高消耗文件反复读取
   - `src/routeTree.gen.ts` 仅在路由异常时查看
   - `auto-imports.d.ts` 仅在类型提示异常时查看
5. 命令执行节奏
   - 开发阶段优先局部验证（类型/页面行为）
   - 改动收敛后再执行一次：`vp run lint:fix`
   - 命令失败先汇报错误原因，不重复盲跑
