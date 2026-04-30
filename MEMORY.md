# MEMORY.md - 长期记忆

本文件保存跨对话需要记住的持久信息。每次新对话读取此文件，重要信息在对话结束时更新到这里。

---

## 用户与项目背景

- **仓库名称**: template
- **仓库用途**: 开发一套全流程的开发模板，包含代码、文档、CI/CD、技能库等，支持快速启动新项目并持续演进。
- **初始化日期**: 2026-04-30
- **AI 助手代号**: Claw

---

## 用户偏好

- 使用中文作为主要交流语言。
- 文件保持简洁，不做过度设计。
- 优先使用已有工具和库，避免引入不必要的依赖。
- 使用技能（skill）时，统一保存到 `.agents/skills/<skill-name>/`，主入口为 `SKILL.md`。
- 用户重视技能的可复用性、持久化和可发现性。

---

## 技术栈决策

### 仓库与工具链

- 采用前后端同仓模板结构：前端位于 `front/**`，后端位于 `backend/**`。
- 前端工作区使用 `pnpm` monorepo，根包名为 `vp`，包管理器版本为 `pnpm@10.33.0`。
- Node.js 要求 `>=22.12.0`。
- 前端任务统一通过 `vite-plus` / `vp` 调度，例如 `vp run admin-react#dev`、`vp run admin-react#build`、`vp staged`。
- 代码规范使用 ESLint 9，配置基于 `@antfu/eslint-config` 与 `@tanstack/eslint-config`。

### Front

- 主应用为 `front/apps/admin-react`，工具库为 `front/packages/utils`（包名 `@vp/utils`）。
- 主应用技术栈：React 19、TypeScript、Vite+、TanStack Router、Ant Design v6、Ant Design Pro Components。
- 样式与主题：Ant Design token、`antd-style`、Tailwind CSS v4、`@tailwindcss/vite`、`@fontsource-variable/geist`。
- 状态与数据：Zustand、Alova、React hooks、Immer、Zod。
- 国际化：i18next、react-i18next、i18next-browser-languagedetector。
- 交互与工具：ahooks、lucide-react、motion、nprogress、js-cookie、lodash、class-variance-authority、clsx。
- Mock 与测试：MSW、Playwright、`@playwright/experimental-ct-react`、Vitest（通过 Vite+ 测试包）。
- 构建能力：`vite-plugin-pwa`、`unplugin-auto-import`、TanStack Router Vite 插件、TanStack/Vite DevTools。

### Backend

- 后端使用 Go workspace，入口文件为 `backend/go.work`，当前 workspace Go 版本为 `1.25.7`。
- 主要模块包括：`backend/admin`、`backend/go-common`、`backend/orm-crud/api`、`backend/orm-crud/gormc`、`backend/orm-crud/pagination`、`backend/sa-token/rueidis`。
- 主服务为 `backend/admin`，HTTP 框架使用 GoFiber v3。
- 后端核心库：GORM、GORM Gen、GORM datatypes、dbresolver、soft_delete、go-playground/validator、Resty v3、Zap、Bytedance Sonic。
- 权限与认证：Casbin v3、casbin gorm-adapter、sa-token-go。
- 数据与任务：PostgreSQL 为当前默认数据库配置，同时依赖 MySQL、SQLite、SQL Server 等 GORM driver；Redis 使用 `go-redis/v9` 与 `rueidis`；异步任务使用 Asynq。
- API 文档与观测：Swagger 使用 swaggo / gofiber contrib swaggo；监控使用 Prometheus client 与 gofiber monitor；日志使用 Zap / zapgorm2。
- `go-common` 提供通用能力，包括配置（Viper）、日志（Zap）、ID 生成、复制、地理/IP 库、MongoDB driver、Goja 脚本运行等。
- `orm-crud` 相关模块承担 ORM、分页、OpenAPI / protobuf 辅助能力。

### CI/CD

- GitHub Actions 工作流位于 `.github/workflows/deploy.yml`。
- 当前 CI 主要面向前端部署：push 到 `main` / `master` 且 `front/**` 变更时触发，使用 `voidzero-dev/setup-vp@v1` 安装构建，并将 `front/apps/admin-react/dist` 部署到 GitHub Pages。
- 现有工作流中的测试 job 目前为占位：`echo "no tests, skipping"`。
- AI 技能优化工作流位于 `.github/workflows/ai-skill-optimizer.yml`：提交影响 `backend/**` 或 `front/**` 时触发，默认通过 DeepSeek 分析 diff，并将可复用经验追加到 `.agents/skills/**/SKILL.md` 的“自动优化记录”。
- AI 技能优化器脚本位于 `.github/scripts/optimize_skills.py`，采用 `ChatProvider` 抽象；默认 secret 为 `DEEPSEEK_API_KEY`，可用 repository variables 覆盖 `DEEPSEEK_BASE_URL` 和 `DEEPSEEK_MODEL`。

### 项目约束

- 修改 `front/**` 后，提交结果前需执行 `vp staged`。
- 修改 `backend/**` 后，需在受影响模块执行 `go fix ./...`、`go vet ./...`、`go test ./...`。
- 修改 `backend/admin` 且涉及 Swagger 注释时，需在 `backend/admin` 执行 `make swagger`。
- 只修改 `skills/**`、`AGENTS.md` 或纯文档时，不需要执行前端或后端构建测试。
- 询问库、框架、SDK、API、CLI 工具或云服务用法时，必须优先用 `ctx7` CLI 获取当前文档。

---

## 重要决策记录

| 日期 | 决策 | 原因 |
|---|---|---|

---

*最后更新: 2026-04-30*
