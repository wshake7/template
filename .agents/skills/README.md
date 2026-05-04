# Project Skills Index

本目录存放 template 仓库的项目级技能。后续对话开始做任务时，先按任务类型查看对应 `SKILL.md`，再动手修改代码。

## 分层索引

| 分类 | 技能 | 适用场景 |
|---|---|---|
| frontend | `frontend/admin-react` | 修改 React 管理后台页面、路由、API、状态、主题、Mock 和测试 |
| backend | `backend/admin-service` | 修改 GoFiber 管理后台服务、路由、业务逻辑、中间件、配置、Swagger |
| backend | `backend/orm-crud` | 修改 GORM CRUD、分页过滤、排序、proto/OpenAPI 辅助与 ORM 生成链路 |
| backend | `backend/go-common` | 修改通用 Go 工具库：日志、配置、ID、加密、集合、转换、IP/Geo 等 |

## 使用约定

- 优先选择最窄的技能，不要把所有技能一次性读入上下文。
- 涉及库、框架、SDK、API、CLI 或云服务用法时，先按仓库根 `AGENTS.md` 要求使用 `ctx7` CLI 获取当前文档。
- 只修改文档或 `.agents/skills/**` 时，不需要执行前端/后端构建测试；修改代码时按对应技能的验证命令执行。
- 遇到已有未提交改动时，只处理当前任务相关文件，不回滚用户改动。

## 自动优化

- `.github/workflows/ai-skill-optimizer.yml` 会在提交影响 `backend/**` 或 `front/**` 时触发。
- workflow 调用 `.github/scripts/optimize_skills.py`，读取对应提交 diff 和现有技能内容，让 AI 动态重写匹配的完整 `SKILL.md`。
- 每次重写前会把旧版技能归档到同目录的 `archive/`，当前 `SKILL.md` 只保留整理后的可执行指南，不保留流水账式“自动优化记录”。
- 默认 AI 服务商为 DeepSeek，使用 OpenAI-compatible chat completions 接口；仓库 secret 需要配置 `DEEPSEEK_API_KEY`，可通过 repository variables 覆盖 `DEEPSEEK_BASE_URL` 和 `DEEPSEEK_MODEL`。
- 后续切换服务商时，在脚本中新增 `ChatProvider` 实现并注册到 `PROVIDERS`，workflow 只需改 `AI_PROVIDER` 与对应密钥变量。
