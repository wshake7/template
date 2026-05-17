# Skill: Frontend App React SSR

## 何时使用

当任务涉及 `front/apps/app-react-ssr` 的 SSR React 应用模板时使用，包括 TanStack Start、Nitro、服务端入口、客户端入口、SSR 安全初始化、路由、Mock、i18n、env、请求层和构建。

## 模块职责

`app-react-ssr` 是通用 SSR React 应用模板。

- 客户端入口：`src/client.tsx`。
- 服务端入口：`src/server.ts`。
- 路由：`src/router.ts`、`src/routes/**`、`src/routeTree.gen.ts`。
- SSR 插件：`@tanstack/react-start/plugin/vite`。
- 运行时：`nitro`。
- 基础能力：`src/env.ts`、`src/i18n.ts`、`src/api/**`、`src/stores/**`、`src/mocks/**`。

## 调用关系

```text
app-react-ssr
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

`app-react-ssr` 与 `app-react` 应尽量保持同构目录；SSR 差异只放在入口、Vite 插件、服务端运行时和必要的环境判断中。

## SSR 注意事项

- 共享包和服务端入口不要在模块顶层直接访问 `window`、`document`、`localStorage`、`navigator`。
- 浏览器专属初始化放到 client entry 或运行时判断之后。
- 与 `app-react` 共享的逻辑优先放入 `@vp/react-core`、`@vp/request`、`@vp/core` 或 `@vp/utils`。
- 修改路由、登录页、Mock、env、i18n 时，同步对照 `app-react`，避免 CSR / SSR 模板漂移。

## 验证

```bash
vp staged
vp run app-react-ssr#build
```

