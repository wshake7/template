# Skills

## 全局
- Vite+ 构建工具: [skills/viteplus.md](skills/viteplus.md)
- AI 工作记录规范: [skills/global/work-log.md](skills/global/work-log.md)
- Git 辅助记录规范: [skills/global/git-record.md](skills/global/git-record.md)
- Skills 编写风格规范: [skills/global/skill-authoring-style.md](skills/global/skill-authoring-style.md)

## Front 通用
- 工具库开发: [skills/front/utils-library.md](skills/front/utils-library.md)

## Front Admin
- Admin 开发框架: [skills/front/admin/admin-framework.md](skills/front/admin/admin-framework.md)
- Ant Design 主题配置: [skills/front/admin/antd-theme.md](skills/front/admin/antd-theme.md)
- E2E 测试 (Playwright): [skills/front/admin/playwright-e2e.md](skills/front/admin/playwright-e2e.md)
- 国际化 (i18next): [skills/front/admin/i18n.md](skills/front/admin/i18n.md)

## Backend
- Backend 工作区总览: [skills/backend/workspace/workspace-overview.md](skills/backend/workspace/workspace-overview.md)
- Backend Admin 服务开发: [skills/backend/admin/admin-service.md](skills/backend/admin/admin-service.md)
- Backend Services 生命周期与依赖编排: [skills/backend/framework/services-lifecycle.md](skills/backend/framework/services-lifecycle.md)
- Backend FiberC 核心结构: [skills/backend/framework/fiberc-core.md](skills/backend/framework/fiberc-core.md)
- Backend Router 编写风格: [skills/backend/framework/router-style.md](skills/backend/framework/router-style.md)
- Backend ORM Models 编写与维护: [skills/backend/libraries/orm-models.md](skills/backend/libraries/orm-models.md)
- Backend go-common 公共库: [skills/backend/libraries/go-common.md](skills/backend/libraries/go-common.md)
- Backend orm-crud 能力: [skills/backend/libraries/orm-crud.md](skills/backend/libraries/orm-crud.md)

## Front 代码改动检查规则
- AI 只要修改 `front/**` 代码，提交结果前必须先执行：`vp run lint:fix`
- 然后执行：`vp run test -r`
- 如果命令失败，先询问是否需要修复

## Backend 代码改动检查规则
- AI 只要修改 `backend/**` 代码，提交结果前必须在对应模块目录执行：`go fix ./...`,`go test ./...`
- 如果涉及多个模块改动，需要在每个受影响模块分别执行：`go fix ./...`,`go test ./...`
- 如果命令失败，先询问是否需要修复

## AI 工作记录规则
- AI 工作记录采用“里程碑记录”，避免每次小改都追加新条目
- 记录格式与字段必须遵循：`skills/global/work-log.md`
- AI 写工作记录时，Git 依据优先引用 `AI_GIT_LOG.md` 最近有效提交记录
- 仅当出现新提交或用户要求追溯时，才更新 `AI_GIT_LOG.md`

## 需求清单规则
- 需求总入口：`AI_REQUIREMENTS.md`
- AI 开始任务前，应先查看 `AI_REQUIREMENTS.md` 中相关分区（总需求/前端需求/后端需求）
- 若需求执行后有状态变化，AI 应同步更新 `AI_REQUIREMENTS.md` 对应条目状态
