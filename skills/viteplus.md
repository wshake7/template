# Vite+ 工作流

## 场景与目标

- 适用场景：运行前端任务、管理依赖、构建、测试、发布工具库。
- 目标：统一使用 Vite+ 的 `vp` 命令，避免绕开项目工具链。

## 当前项目事实

- 根包：`package.json`
- 工作区：`pnpm-workspace.yaml`
- 包管理器：`pnpm@10.33.0`，由 Vite+ 识别和封装。
- Node 要求：`>=22.12.0`
- 前端工作区：
  - `front/apps/admin-react`
  - `front/packages/utils`

## 常用命令

```bash
# 根目录
vp run dev
vp run build
vp staged
vp run ready

# 指定包任务
vp run admin-react#dev
vp run admin-react#build
vp run admin-react#e2e:test
```

## 命令选择

1. 运行根 `package.json` 脚本时使用 `vp run <script>`。
2. 运行工作区包脚本时使用 `vp run <package>#<script>`。
3. 依赖管理优先使用 `vp add`、`vp remove`、`vp install`，不要直接用 `pnpm`。
4. 一次性执行包二进制优先使用 `vp dlx`；项目规则特别要求 `npx ctx7@latest` 时按 `AGENTS.md` 执行。

## 注意事项

1. `vp dev` 是 Vite+ 内置命令，不等同于根脚本 `dev`；根脚本要写 `vp run dev`。
2. 不安装或直接升级 Vitest、Oxlint、Oxfmt、tsdown 等被 Vite+ 封装的工具。
3. 测试工具从 `vite-plus/test` 导入，工具库测试不要从 `vitest` 导入。
4. 修改 `front/**` 后，交付前必须按 `AGENTS.md` 执行 `vp staged`。
