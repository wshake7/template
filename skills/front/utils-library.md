# 前端工具库 `@vp/utils`

## 场景与目标

- 适用场景：修改 `front/packages/utils`，新增通用 TypeScript 工具函数。
- 目标：保持工具库轻量、可测试、可被工作区应用直接复用。

## 目录/文件位置

- 入口：`front/packages/utils/src/index.ts`
- 测试：`front/packages/utils/tests/index.test.ts`
- 构建配置：`front/packages/utils/vite.config.ts`
- 包配置：`front/packages/utils/package.json`

## 当前约定

- 包名：`@vp/utils`
- 模块类型：ESM
- 导出入口：`dist/index.mjs`
- 使用 Vite+ 的 `vp pack` 构建。
- `front/apps/admin-react` 通过 `workspace:*` 依赖该包。

## 操作步骤

1. 先确认目标能力是否真正通用；业务专属逻辑留在应用内。
2. 在 `src/index.ts` 导出函数或类型。
3. 为非平凡逻辑补充 `tests/index.test.ts` 或相近测试文件。
4. 在 `front/packages/utils` 下运行构建或测试。
5. 如影响 Admin 应用，再执行根目录 `vp staged`。

## 常用命令

```bash
cd front/packages/utils
vp build
vp test
vp check
```

## 注意事项

1. 工具库优先保持零运行时依赖。
2. 测试从 `vite-plus/test` 导入 `describe`、`test`、`expect`。
3. 导出类型时使用 `export type`。
4. 不把 Admin 页面、请求、store 等业务语义放入 `@vp/utils`。
