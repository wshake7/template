# AI 工作记录

本文件用于记录 AI 在本项目中的每次实际改动与验证结果。

## 记录规范

- 每次涉及代码或配置改动时，新增一条记录
- 记录时间使用本地时间，格式 `YYYY-MM-DD HH:mm`
- 内容必须包含：目标、改动文件、执行命令、结果、风险/待办
- 不记录与项目无关的对话

## 记录模板

```md
## [YYYY-MM-DD HH:mm] <任务标题>
- 目标: <本次任务目标>
- 改动文件:
  - <path1>
  - <path2>
- 执行命令:
  - <command1>
  - <command2>
- 结果: <通过/失败 + 简述>
- 风险/待办: <无/具体项>
- Git 依据: <对应 AI_GIT_LOG.md 的时间戳标题 / 无（未采集）>
```

## [2026-04-23 16:37] 新增 AI 工作记录机制
- 目标: 增加 AI 固定工作记录文件与对应 skill 规则，确保后续改动可追踪
- 改动文件:
  - AI_WORK_LOG.md
  - skills/global/work-log.md
  - AGENTS.md
- 执行命令:
  - Get-Date -Format "yyyy-MM-dd HH:mm"
- 结果: 通过，记录文件与 skill 已创建并接入 AGENTS 规则
- 风险/待办: 无
- Git 依据: 无（未采集）

## [2026-04-23 16:40] 新增 Git 辅助记录机制
- 目标: 增加 Git 辅助记录文件与 skill，支持用户手动改动后的记录追溯
- 改动文件:
  - AI_GIT_LOG.md
  - skills/global/git-record.md
  - skills/global/work-log.md
  - AGENTS.md
  - AI_WORK_LOG.md
- 执行命令:
  - Get-Date -Format "yyyy-MM-dd HH:mm"
  - git branch --show-current
  - git status --short
  - git diff --name-only
  - git log -n 5 --oneline
- 结果: 通过，机制已接入；当前环境不是 Git 仓库，未能采集分支/状态/日志
- 风险/待办: 若需 Git 证据，请在真实 Git 仓库目录执行同样命令
- Git 依据: [2026-04-23 16:40] 新增 Git 辅助记录机制
