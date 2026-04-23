# Git 辅助记录规范

用于在 AI 写工作记录时同步记录 Git 证据，帮助识别用户手动改动并校正后续记录。

## 记录文件

- 固定路径: `AI_GIT_LOG.md`

## 何时记录

- AI 完成任何实际文件改动后
- 用户明确表示“我手动改过 AI 改动内容”后
- 发生记录冲突、需追溯改动来源时

## 记录步骤

1. 执行 `git branch --show-current` 获取当前分支
2. 执行 `git status --short` 获取工作区状态
3. 执行 `git diff --name-only` 获取本次未提交变更文件
4. 执行 `git log -n 5 --oneline` 获取最近提交参考
5. 将结果追加到 `AI_GIT_LOG.md`

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

## 与 AI_WORK_LOG 联动

- `AI_WORK_LOG.md` 中每条记录增加 `Git 依据` 字段
- `Git 依据` 内容填写对应 `AI_GIT_LOG.md` 的时间戳标题
- 若未生成 Git 辅助记录，`AI_WORK_LOG.md` 中必须写 `Git 依据: 无（未采集）`
