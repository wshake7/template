# Skill 编写与重构规范

## 场景与目标

- 适用场景：新增、合并、拆分、重写 `skills/**` 文档。
- 目标：让 skill 成为“任务执行手册”，而不是项目流水账或资料堆。

## 目录边界

- `skills/README.md`：总索引和任务入口。
- `skills/global/`：跨任务规则，如记录、Git、文档写法。
- `skills/front/`：前端通用能力。
- `skills/front/admin-react/`：Admin React 专属能力。
- `skills/backend/`：Go 后端、Admin 服务、ORM 和公共库能力。
- `skills/viteplus.md`：Vite+ 工具链。

## 单个 skill 结构

优先使用以下结构，按需要删减：

```md
# <Skill 名称>

## 场景与目标
## 目录/文件位置
## 当前约定
## 操作步骤
## 常用命令
## 注意事项
```

## 写作规则

1. 文件名使用小写中划线，如 `admin-service.md`。
2. 一个 skill 只承载一个主能力；索引文档负责路由，不承载细节。
3. 内容写稳定规则和当前项目事实，不写一次性执行过程。
4. 路径使用项目相对路径，并与仓库实际文件一致。
5. 命令使用本项目约定，如 `vp run ...`、在具体 Go module 目录执行 `go test ./...`。
6. 已废弃或历史兼容入口要明确标注，不继续添加新规则。
7. 长协议、长语法参考可独立成文，例如 `backend/orm-query.md`。

## 低上下文执行规则

1. 先用 `rg` / `rg --files` 定位，再读关键文件。
2. 一次任务通常只读入口文件、直接依赖、相近实现和待修改文件。
3. 不反复读取生成文件，如 `routeTree.gen.ts`、`auto-imports.d.ts`、`query/*.gen.go`，除非任务正涉及生成结果。
4. 同类文档只保留一个权威入口；历史入口只做跳转说明。
5. 完成后只汇报改动结果、校验情况和必要风险。

## 维护联动

1. 新增或移动 skill 后同步更新 `skills/README.md` 和 `AGENTS.md`。
2. 如果 skill 改变了执行方式，按需更新 `AI_WORK_LOG.md`。
3. 如果用户要求追溯提交证据，按 `global/git-record.md` 维护 `AI_GIT_LOG.md`。
