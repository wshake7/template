# Skill: Project Orientation

## 何时使用

当任务需要理解整个 template 仓库、选择工作目录、判断修改影响范围、恢复项目上下文，或不确定该读哪个模块时使用。

## 项目身份

- 仓库是一个全流程开发模板，目标是沉淀前端、后端、CI/CD、文档和 AI 技能库。
- AI 助手代号为 Claw；项目记忆来源优先级为 `AGENTS.md`、`MEMORY.md`、最近的 `memory/YYYY-MM-DD.md`、`memory/tasks.md`。
- 用户主要使用中文，偏好简洁文件、最小变更、可复用技能和持久化记录。

## 目录地图

```text
.
├── AGENTS.md                 # 助手身份、项目规则、ctx7 文档规则
├── MEMORY.md                 # 长期记忆
├── memory/                   # 每日日志与任务追踪
├── .agents/skills/           # 项目级技能库
├── front/                    # 前端 pnpm workspace
│   ├── apps/admin-react/     # 管理后台 React 应用
│   └── packages/utils/       # @vp/utils 工具包
├── backend/                  # Go workspace
│   ├── go.work
│   ├── admin/                # GoFiber 管理后台服务
│   ├── go-common/            # 通用 Go 基础库
│   ├── orm-crud/             # GORM CRUD、分页、proto 辅助
│   └── sa-token/rueidis/     # sa-token Redis 适配模块
└── .github/workflows/        # GitHub Actions
```

## 技术栈速查

- 前端：pnpm monorepo、Vite+、React 19、TypeScript、TanStack Router、Ant Design v6、Ant Design Pro Components、antd-style、Tailwind CSS v4、Zustand、Alova、Zod、i18next、MSW、Playwright、Vitest。
- 后端：Go workspace，主服务 `backend/admin`，GoFiber v3、GORM/GORM Gen、Casbin、sa-token-go、Redis、Asynq、Zap、Swagger。
- 通用库：`backend/go-common` 提供配置、日志、结果、集合、mapper、ID、加密、文件、字符串、转换、IP/Geo 等工具。
- ORM/CRUD：`backend/orm-crud/gormc` 负责 GORM 客户端、Repository、mixin、过滤、排序、分页；`pagination` 负责请求分页过滤表达式；`api` 负责 pagination proto 生成。

## 开始任务的固定流程

1. 读取 `MEMORY.md`，必要时读取最近的 `memory/*.md` 与 `memory/tasks.md`。
2. 查看 `.agents/skills/README.md`，选择最匹配的项目级技能。
3. 用 `rg --files`、`rg`、`sed -n` 定位相关代码，优先沿已有模式修改。
4. 对库/框架/SDK/API/CLI 用法不确定时，按 `AGENTS.md` 使用 `ctx7`：先 `library`，再 `docs`。
5. 修改前确认 `git status --short`，避免覆盖用户已有改动。
6. 完成后按影响范围验证，并更新日志/任务状态。

## 验证边界

- 只改 `.agents/skills/**`、`AGENTS.md` 或纯文档：通常无需构建测试。
- 修改 `front/**`：至少运行 `vp staged`；若涉及应用行为，优先运行对应 `vp run admin-react#build` 或 Playwright。
- 修改 `backend/**`：在受影响 Go module 执行 `go fix ./...`、`go vet ./...`、`go test ./...`。
- 修改 `backend/admin` 且影响 Swagger 注释：在 `backend/admin` 执行 `make swagger`。
