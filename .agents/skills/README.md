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
