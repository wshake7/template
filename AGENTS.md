# Skills

## 全局
- Vite+ 构建工具: [skills/viteplus.md](skills/viteplus.md)
- AI 工作记录规范: [skills/global/work-log.md](skills/global/work-log.md)
- Git 辅助记录规范: [skills/global/git-record.md](skills/global/git-record.md)
- Skills 编写风格规范: [skills/global/skill-authoring-style.md](skills/global/skill-authoring-style.md)

## Front 通用
- 工具库开发: [skills/front/utils-library.md](skills/front/utils-library.md)

## Front Admin (admin-react)
- Admin 开发框架: [skills/front/admin-react/admin-framework.md](skills/front/admin-react/admin-framework.md)
- CRUD 页面开发: [skills/front/admin-react/crud-page.md](skills/front/admin-react/crud-page.md)
- API 请求层: [skills/front/admin-react/api.md](skills/front/admin-react/api.md)
- 状态管理: [skills/front/admin-react/stores.md](skills/front/admin-react/stores.md)
- Ant Design 主题配置: [skills/front/admin-react/antd-theme.md](skills/front/admin-react/antd-theme.md)
- E2E 测试 (Playwright): [skills/front/admin-react/playwright-e2e.md](skills/front/admin-react/playwright-e2e.md)
- 国际化 (i18next): [skills/front/admin-react/i18n.md](skills/front/admin-react/i18n.md)

## Backend
- 工作区总览: [skills/backend/workspace-overview.md](skills/backend/workspace-overview.md)
- Admin 服务开发: [skills/backend/admin-service.md](skills/backend/admin-service.md)
- FiberC 核心结构: [skills/backend/fiberc-core.md](skills/backend/fiberc-core.md)
- Router 编写风格: [skills/backend/router-style.md](skills/backend/router-style.md)
- Services 生命周期与依赖编排: [skills/backend/services-lifecycle.md](skills/backend/services-lifecycle.md)
- Swagger 编写风格: [skills/backend/swagger-style.md](skills/backend/swagger-style.md)
- ORM 数据库操作与 Models 维护: [skills/backend/orm-models.md](skills/backend/orm-models.md)
- ORM 数据库 List 分页查询语法: [skills/backend/orm-query.md](skills/backend/orm-query.md)
- ORM CRUD 能力: [skills/backend/orm-crud.md](skills/backend/orm-crud.md)
- go-common 公共库: [skills/backend/go-common.md](skills/backend/go-common.md)

## Front 代码改动检查规则
- AI 只要修改 `front/**` 代码，提交结果前必须先执行：`vp run lint:fix`
- 如果命令失败，先询问是否需要修复

## Backend 代码改动检查规则
- AI 只要修改 `backend/**` 代码，提交结果前必须在对应模块目录执行：`go fix ./...`,`go test ./...`
- **Swagger 更新**：如果 AI 修改了 `backend/admin` 中涉及 Swagger 注释的内容，必须在 `backend/admin` 目录执行 `make swagger` 重新生成文档。
- 如果涉及多个模块改动，需要在每个受影响模块分别执行：`go fix ./...`,`go test ./...`
- 如果命令失败，先询问是否需要修复

## AI 工作记录规则
- AI 工作记录采用"里程碑记录"，避免每次小改都追加新条目
- 记录格式与字段必须遵循：`skills/global/work-log.md`
- AI 写工作记录时，Git 依据优先引用 `AI_GIT_LOG.md` 最近有效提交记录
- 仅当出现新提交或用户要求追溯时，才更新 `AI_GIT_LOG.md`

## 需求清单规则
- 需求总入口：`AI_REQUIREMENTS.md`
- AI 开始任务前，应先查看 `AI_REQUIREMENTS.md` 中相关分区（总需求/前端需求/后端需求）
- 若需求执行后有状态变化，AI 应同步更新 `AI_REQUIREMENTS.md` 对应条目状态

## 低 Token 执行规则
- AI 默认遵循"先定位、后读取、最小充分验证"原则，减少无关文件读取与重复检索
- 优先使用 `Glob/Grep` 缩小范围，再按需 `Read` 关键文件，避免一次性读取大文件
- 前端任务默认按 [skills/front/admin-react/admin-framework.md](skills/front/admin-react/admin-framework.md) 的"低 Token 执行清单（Admin）"执行
