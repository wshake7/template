# Git 辅助记录规范

## 场景与目标

- 适用场景：用户要求追溯、出现新提交、需要识别用户手动改动来源。
- 目标：用提交粒度记录证据，辅助 `AI_WORK_LOG.md` 保持可信。

## 记录文件

- 固定路径：`AI_GIT_LOG.md`

## 何时记录

1. `git log -n 1 --oneline` 与 `AI_GIT_LOG.md` 最近记录不同。
2. 用户明确表示手动改过 AI 产物，需要以后以手动版本为准。
3. 工作记录出现冲突，需要追溯改动来源。

## 记录步骤

1. 执行 `git log -n 1 --oneline` 获取最新提交。
2. 执行 `git branch --show-current` 获取分支。
3. 执行 `git show --name-only --pretty=oneline <hash>` 获取文件证据。
4. 与 `AI_GIT_LOG.md` 最近一条记录比对，只有 hash 变化或用户要求时追加。

## 记录模板

```md
## [YYYY-MM-DD HH:mm] <任务标题>
- 分支: <branch>
- 提交哈希: <hash>
- 提交说明: <git log -n 1 --oneline>
- 提交文件:
  - <path>
- 对比基线: <上次记录 hash / 无>
- 备注: <必要说明>
```

## 注意事项

1. Git 记录是证据，不是交付总结。
2. 不因普通未提交改动自动写 `AI_GIT_LOG.md`。
3. `AI_WORK_LOG.md` 的 `Git 依据` 引用这里的标题；无新提交时写 `复用上次提交记录`。
