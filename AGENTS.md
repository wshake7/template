# AGENTS.md

本文件是 AI 在本仓库工作的入口规则。开始任务时先读本文件，再按任务类型进入 `skills/**` 中最少必要的文档。

<!-- context7 -->
## Context7 文档规则

当用户询问库、框架、SDK、API、CLI 工具或云服务的用法时，必须使用 `ctx7` CLI 获取当前文档，即使是 React、Next.js、Prisma、Express、Tailwind、Django、Spring Boot 等常见技术也一样。适用范围包括 API 语法、配置项、版本迁移、库相关调试、安装配置、CLI 用法。

不适用于：重构、从零写脚本、业务逻辑调试、代码审查、通用编程概念。

执行步骤：

1. 解析库：`npx ctx7@latest library <name> "<user's question>"`
2. 从结果中按名称匹配度、描述相关性、代码片段数量、来源信誉、benchmark 分数选择最佳 `/org/project` ID
3. 拉取文档：`npx ctx7@latest docs <libraryId> "<user's question>"`
4. 若答案不足，再运行同一 docs 命令并追加 `--research`
5. 基于获取到的文档回答

约束：

- 用户未直接提供 `/org/project` ID 时，必须先执行 `library`
- 查询语句使用用户完整问题，避免包含 API Key、密码等敏感信息
- 每个问题最多执行 3 条 Context7 命令
- 版本限定优先使用 `library` 输出中的 `/org/project/version`
- 若遇到 quota 错误，告知用户可执行 `npx ctx7@latest login` 或设置 `CONTEXT7_API_KEY`
- Context7 CLI 请求需在 Codex 默认沙箱外执行；若出现 DNS、ENOTFOUND、host resolution、fetch failed 等网络错误，按沙箱外执行规则重试
<!-- context7 -->

## 项目入口

- Skills 总索引：`skills/README.md`
- 前端应用：`front/apps/admin-react`
- 前端工具库：`front/packages/utils`
- 后端工作区：`backend/go.work`
- 后端主服务：`backend/admin`

开始任务时遵循“先定位、后读取、最小充分验证”：

1. 优先用 `rg` / `rg --files` 缩小范围。
2. 再读取当前任务的入口文件、直接依赖和相近实现。
3. 不一次性展开全部 skills、生成文件或无关模块。

## Skills 任务入口

完整索引见 `skills/README.md`。常用入口如下：

### 全局与工具链

- Vite+ 工作流：`skills/viteplus.md`
- Skill 编写与重构：`skills/global/skill-authoring-style.md`
- AI 工作记录：`skills/global/work-log.md`
- Git 辅助记录：`skills/global/git-record.md`

### Front

- 工具库 `@vp/utils`：`skills/front/utils-library.md`
- Admin 框架：`skills/front/admin-react/admin-framework.md`
- CRUD 页面：`skills/front/admin-react/crud-page.md`
- API 请求层：`skills/front/admin-react/api.md`
- 状态管理：`skills/front/admin-react/stores.md`
- Ant Design 主题：`skills/front/admin-react/antd-theme.md`
- Playwright 测试：`skills/front/admin-react/playwright-e2e.md`
- 国际化：`skills/front/admin-react/i18n.md`

### Backend

- 工作区总览：`skills/backend/workspace-overview.md`
- Admin 服务开发：`skills/backend/admin-service.md`
- Router 风格：`skills/backend/router-style.md`
- FiberC 核心：`skills/backend/fiberc-core.md`
- Services 生命周期：`skills/backend/services-lifecycle.md`
- Swagger 风格：`skills/backend/swagger-style.md`
- ORM Models 与 Query：`skills/backend/orm-models.md`
- ORM 过滤语法：`skills/backend/orm-query.md`
- orm-crud 基础设施：`skills/backend/orm-crud.md`
- go-common 公共库：`skills/backend/go-common.md`

## 执行前检查

- 若仓库存在 `AI_REQUIREMENTS.md`，开始任务前先查看相关分区：总需求、前端需求、后端需求。
- 若需求执行后状态变化，同步更新对应条目的状态、备注或验收信息。
- 如果 `AI_REQUIREMENTS.md` 不存在，以用户当前需求和代码事实为准，不强行创建。

## 代码改动校验

### Front

- 只要修改 `front/**` 代码，提交结果前必须执行：`vp staged`
- 如果命令失败，先询问用户是否需要继续修复。

### Backend

- 只要修改 `backend/**` 代码，提交结果前必须在对应模块目录执行：`go fix ./...`、`go vet ./...`、`go test ./...`
- 若修改 `backend/admin` 中涉及 Swagger 注释的内容，必须在 `backend/admin` 目录执行 `make swagger`
- 涉及多个 backend 模块时，在每个受影响模块分别执行 `go fix ./...`、`go vet ./...`、`go test ./...`
- 如果命令失败，先询问用户是否需要继续修复。

### Docs / Skills

- 只修改 `skills/**`、`AGENTS.md` 或纯文档时，不需要执行前端或后端构建测试。
- 需要校验时优先检查链接、路径和 Markdown 可读性。

## 记录规则

- AI 工作记录采用里程碑记录，避免每次小改都追加新条目。
- 工作记录格式遵循 `skills/global/work-log.md`。
- Git 辅助记录格式遵循 `skills/global/git-record.md`。
- AI 写工作记录时，优先引用 `AI_GIT_LOG.md` 最近有效提交记录。
- 仅当出现新提交、用户要求追溯、需求状态变化或工作方式变化时，才更新记录文件。
