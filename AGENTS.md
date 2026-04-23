# Skills

## 全局
- Vite+ 构建工具: [skills/viteplus.md](skills/viteplus.md)
- AI 工作记录规范: [skills/global/work-log.md](skills/global/work-log.md)
- Git 辅助记录规范: [skills/global/git-record.md](skills/global/git-record.md)

## Front 通用
- 工具库开发: [skills/front/utils-library.md](skills/front/utils-library.md)

## Front Admin
- Admin 开发框架: [skills/front/admin/admin-framework.md](skills/front/admin/admin-framework.md)
- Ant Design 主题配置: [skills/front/admin/antd-theme.md](skills/front/admin/antd-theme.md)
- E2E 测试 (Playwright): [skills/front/admin/playwright-e2e.md](skills/front/admin/playwright-e2e.md)
- 国际化 (i18next): [skills/front/admin/i18n.md](skills/front/admin/i18n.md)

## Front 代码改动检查规则
- AI 只要修改 `front/**` 代码，提交结果前必须先执行：`vp run lint:fix`
- 然后执行：`vp run test -r`
- 如果命令失败，先询问是否需要修复

## AI 工作记录规则
- AI 发生实际改动后，必须把本次工作记录追加到 `AI_WORK_LOG.md`
- 记录格式与字段必须遵循：`skills/global/work-log.md`
- AI 写工作记录时，应同步采集 Git 证据并记录到 `AI_GIT_LOG.md`

## 需求清单规则
- 需求总入口：`AI_REQUIREMENTS.md`
- AI 开始任务前，应先查看 `AI_REQUIREMENTS.md` 中相关分区（总需求/前端需求/后端需求）
- 若需求执行后有状态变化，AI 应同步更新 `AI_REQUIREMENTS.md` 对应条目状态
