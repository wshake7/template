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
  - `vp run admin#dev`：启动 admin 开发
  - `vp run admin#build`：构建 admin
  - `vp run lint:fix`：修复 lint
  - `vp run test -r`：递归执行工作区测试
  - `vp run build -r`：递归执行工作区构建
- 关键开发依赖:
  - `vite-plus`：统一工具链入口
  - `eslint` / `@antfu/eslint-config`：代码规范与校验

### Admin package.json（front/apps/admin）

位置: `front/apps/admin/package.json`

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
front/apps/admin/src/
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

位置: `front/apps/admin/.env.*`

- `.env` - 默认环境
- `.env.dev` - 开发环境
- `.env.test` - 测试环境
- `.env.prod` - 生产环境

## 常用命令

```bash
# 根目录启动 admin 开发
vp run admin#dev

# 根目录构建 admin
vp run admin#build

# E2E 测试
vp run admin#e2e:test
vp run admin#e2e:test-ui    # UI 模式
vp run admin#e2e:show       # 查看报告
vp run admin#e2e:codegen    # 生成测试代码
```

## 开发注意事项

1. 路由变更后由开发流程自动更新 `routeTree.gen.ts`，不要手改该文件
2. API 请求使用 Alova，不直接使用 fetch/axios
3. 状态管理优先使用 Zustand，避免 Redux 过于复杂
4. 组件放在 `components/business/` 或 `components/lib/` 下
5. 工具函数放在 `utils/` 下，领域逻辑放在 `domains/` 下
