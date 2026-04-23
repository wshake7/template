# AI Git 辅助记录

本文件用于记录与 AI 工作记录相关的 Git 证据，辅助判断“用户手动改动”与“AI 改动”的差异。

## 记录模板

```md
## [YYYY-MM-DD HH:mm] <任务标题>
- 分支: <branch>
- 工作区状态:
  - <git status --short 的关键输出>
- 本次变更文件:
  - <git diff --name-only 的输出>
- 最近提交:
  - <git log -n 5 --oneline 的关键输出>
- 备注: <如：用户手动修改了 AI 产物，后续记录以此为准>
```

## [2026-04-23 16:40] 新增 Git 辅助记录机制
- 分支: 无法采集（当前目录非 Git 仓库）
- 工作区状态:
  - `git status --short` 执行失败：not a git repository
- 本次变更文件:
  - `git diff --name-only` 执行失败：not a git repository
- 最近提交:
  - `git log -n 5 --oneline` 执行失败：not a git repository
- 备注: 机制已建立；进入真实 Git 仓库后可直接按模板采集

## [2026-04-23 16:44] 新增需求清单文件机制
- 分支: master
- 工作区状态:
  - `git status --short` 显示当前文件均为未跟踪状态（`??`）
- 本次变更文件:
  - 无（`git diff --name-only` 输出为空，未加入索引的新增文件不在该输出中）
- 最近提交:
  - `git log -n 5 --oneline` 执行失败：当前分支尚无提交记录
- 备注: 仓库可用，但尚未有首个 commit；后续提交后可获得完整 Git 追溯链路
