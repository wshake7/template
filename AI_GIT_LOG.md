# AI Git 辅助记录

本文件用于记录与 AI 工作记录相关的 Git 证据，按提交粒度追踪“用户手动改动”与“AI 改动”的差异。

## 记录原则（精简模式）

- 仅在检测到新提交时新增记录（或用户明确要求追溯时）
- 未产生新提交时，不新增记录，`AI_WORK_LOG.md` 复用最近一条 Git 依据

## 记录模板

```md
## [YYYY-MM-DD HH:mm] <任务标题>
- 分支: <branch>
- 提交哈希: <latest-commit-hash>
- 提交说明: <git log -n 1 --oneline 的标题>
- 提交文件:
  - <git show --name-only 的输出>
- 对比基线: <上次记录的 commit-hash / 无>
- 备注: <如：用户手动修改了 AI 产物，后续记录以此为准>
```
