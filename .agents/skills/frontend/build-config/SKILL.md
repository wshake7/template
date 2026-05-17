# Skill: Frontend Package Build Config

## 何时使用

当任务涉及 `front/packages/build-config`、前端 Vite 配置复用、Rollup 分包、MSW 生产清理、Playwright / CT 配置时使用。

## 模块职责

`@vp/build-config` 是前端构建和测试配置共享包。

- `src/vite/chunks.ts`：导出 `createManualChunks(options)`。
- `src/vite/plugins.ts`：导出 `createRemoveMswPlugin(mode)`。
- `src/playwright/config.ts`：导出 `createPlaywrightConfig()` 和 `createPlaywrightCtConfig()`。
- `src/index.ts`：统一导出公共 API。

## 调用关系

```text
@vp/build-config
├── 被 admin-react 调用
├── 被 app-react 调用
└── 被 app-react-ssr 调用
```

`@vp/build-config` 不应调用任何 app，也不应依赖 `@vp/core`、`@vp/react-core`、`@vp/request`、`@vp/utils` 的业务能力。它只负责构建/测试配置。

## 变更规则

- 通用 Vite 分包策略放这里，不要在三个 app 中复制维护。
- Ant Design、Monaco、Shiki 等可选分包通过 `CreateManualChunksOptions` 控制。
- Playwright 配置变更要考虑 E2E 和组件测试两类场景。
- 公共 API 变更后同步检查三个 app 的 `vite.config.ts` 和 `playwright*.ts`。

## 验证

```bash
vp staged
vp run @vp/build-config#build
vp run admin-react#build
vp run app-react#build
vp run app-react-ssr#build
```

