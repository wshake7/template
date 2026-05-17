# Skill: Frontend App React

## 何时使用

当任务涉及 `front/apps/app-react` 的 CSR React 应用模板时使用，包括 SPA/PWA 初始化、TanStack Router、登录页、Mock、i18n、env、请求层、基础 store、Playwright 和 Vite+ 构建。

## 模块职责

`app-react` 是通用 CSR 应用模板，适合作为新前端 SPA / PWA 项目的起点。

- 入口：`src/main.tsx`。
- 路由：`src/router.ts`、`src/routes/**`、`src/routeTree.gen.ts`。
- 基础页面：`src/routes/index.tsx`、`src/routes/__root.tsx`、`src/routes/(login)/login.tsx`。
- 基础能力：`src/env.ts`、`src/i18n.ts`、`src/api/**`、`src/stores/**`、`src/mocks/**`。
- 样式：`src/styles/index.css`。

## 调用关系

```text
app-react
├── 调用 @vp/core
│   └── 共享类型、HTTP 常量、通知器契约、页面契约
├── 调用 @vp/react-core
│   └── env 工厂、i18n 工厂、字典 hooks 工厂、store factory、Mock helper
├── 调用 @vp/request
│   └── Alova client 和加密请求辅助
├── 调用 @vp/utils
│   └── cn、日期、Web Crypto 等工具
└── 调用 @vp/build-config
    └── Vite 分包、MSW 清理、Playwright 配置
```

`app-react` 与 `app-react-ssr` 应尽量保持同构目录和同名基础文件；差异主要集中在入口和 SSR 插件。

## 常见任务

### 增加模板基础能力

1. 如果能力对 CSR 和 SSR 都通用，优先同步检查 `app-react-ssr`。
2. 如果能力可被多个 app 复用，先沉淀到对应共享包，再在 app 内做薄适配。
3. 浏览器专属逻辑可以放在 CSR app 初始化路径，但如果未来 SSR 也要复用，应避免直接污染共享包。

### 修改路由

1. 改 `src/routes/**` 后确认 `src/routeTree.gen.ts` 会由 TanStack Router 插件生成。
2. 根路由、登录路由、错误页和 Not Found 属于模板基础结构，改动时同步考虑 SSR 模板。

## 验证

```bash
vp staged
vp run app-react#build
```

