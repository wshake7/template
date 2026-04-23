# Git 辅助记录规范

用于在 AI 写工作记录时按“提交粒度”记录 Git 证据，帮助识别用户手动改动并校正后续记录。

## 记录文件

- 固定路径: `AI_GIT_LOG.md`

## 何时记录

- 检测到**新的 commit**（`git log -n 1 --oneline` 与上次记录不同）时
- 用户明确表示“我手动改过 AI 改动内容”后
- 发生记录冲突、需追溯改动来源时

## 记录步骤

1. 执行 `git log -n 1 --oneline` 获取最新提交哈希
2. 与 `AI_GIT_LOG.md` 最近一条记录的哈希比对
3. 仅当哈希变化时继续记录；否则跳过写入
4. 执行 `git branch --show-current`、`git show --name-only --pretty=oneline <hash>` 补充证据
5. 将结果追加到 `AI_GIT_LOG.md`

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

## 与 AI_WORK_LOG 联动

- `AI_WORK_LOG.md` 中每条记录增加 `Git 依据` 字段
- `Git 依据` 内容填写对应 `AI_GIT_LOG.md` 的时间戳标题
- 若本次没有新提交，`AI_WORK_LOG.md` 中写 `Git 依据: 复用上次提交记录`
