# Skills 总索引

本目录是 AI 在本仓库工作的任务手册。使用原则是：先读入口规则，再按任务选择最少的 skill，最后按改动范围执行校验。

## 使用流程

1. 先读根目录 `AGENTS.md`，确认 Context7、需求、校验和记录规则。
2. 如果存在 `AI_REQUIREMENTS.md`，先定位相关分区；文件不存在时继续按用户当前需求执行。
3. 按“任务入口”选择 1 到 3 个 skill 阅读，不一次性展开全部。
4. 完成代码改动后按 `AGENTS.md` 执行前端或后端校验。
5. 仅当用户要求、需求状态变化、产生新提交或工作方式变化时，更新工作记录。

## 任务入口

| 任务 | 必读 | 视情况追加 |
| --- | --- | --- |
| 前端 Admin 页面、路由、布局 | `front/admin-react/admin-framework.md` | `front/admin-react/crud-page.md`、`front/admin-react/i18n.md` |
| 前端 CRUD 列表、表单、批量操作 | `front/admin-react/crud-page.md`、`front/admin-react/api.md` | `backend/orm-query.md` |
| 前端请求层或接口类型 | `front/admin-react/api.md` | `backend/router-style.md`、`backend/admin-service.md` |
| Zustand 全局状态 | `front/admin-react/stores.md` | `front/admin-react/admin-framework.md` |
| 主题、全局样式、Ant Design token | `front/admin-react/antd-theme.md` | `front/admin-react/stores.md` |
| Playwright E2E 或组件测试 | `front/admin-react/playwright-e2e.md` | `front/admin-react/admin-framework.md` |
| `@vp/utils` 工具库 | `front/utils-library.md` | `viteplus.md` |
| 后端接口、业务逻辑、Swagger | `backend/workspace-overview.md`、`backend/router-style.md`、`backend/admin-service.md` | `backend/swagger-style.md` |
| 后端模型、分页、过滤、query 生成 | `backend/orm-models.md` | `backend/orm-query.md`、`backend/orm-crud.md` |
| Fiber 中间件、错误处理、优雅关闭 | `backend/fiberc-core.md` | `backend/router-style.md` |
| 服务初始化、依赖顺序、健康检查 | `backend/services-lifecycle.md` | `backend/fiberc-core.md` |
| Go 公共库能力 | `backend/go-common.md` | `backend/workspace-overview.md` |
| Vite+ 命令或工作区任务 | `viteplus.md` | 对应前端 skill |
| 新增、归并、重写 skill | `global/skill-authoring-style.md` | `global/work-log.md`、`global/git-record.md` |

## Skill 清单

### 全局与工具链

- `viteplus.md`：Vite+ / `vp` 命令、依赖和任务运行规则。
- `global/skill-authoring-style.md`：skill 写作、拆分和维护规范。
- `global/work-log.md`：`AI_WORK_LOG.md` 里程碑记录规范。
- `global/git-record.md`：`AI_GIT_LOG.md` 提交证据记录规范。

### 前端

- `front/utils-library.md`：`front/packages/utils` 工具库开发。
- `front/admin-react/admin-framework.md`：Admin React 项目结构、路由、菜单、执行路径。
- `front/admin-react/crud-page.md`：CRUD 页面、ProTable、ModalForm、分页、批量操作。
- `front/admin-react/api.md`：Alova 请求层、资源 API、分页类型。
- `front/admin-react/stores.md`：Zustand store 约定。
- `front/admin-react/antd-theme.md`：Ant Design v6、主题 token、主题切换。
- `front/admin-react/playwright-e2e.md`：Playwright E2E 和组件测试。
- `front/admin-react/i18n.md`：i18next 文案资源维护。

### 后端

- `backend/workspace-overview.md`：Go workspace 模块分布和命令入口。
- `backend/admin-service.md`：`backend/admin` 接口开发、业务逻辑、错误处理。
- `backend/router-style.md`：路由注册、handler 包装器、中间件顺序。
- `backend/fiberc-core.md`：FiberC 初始化、全局中间件、错误处理、关闭链路。
- `backend/services-lifecycle.md`：服务装配、启动、健康检查、终止。
- `backend/swagger-style.md`：Swagger 注释和文档生成。
- `backend/orm-models.md`：models、query、gorm/gen、数据库读写。
- `backend/orm-query.md`：分页过滤 `query` 语法。
- `backend/orm-crud.md`：`orm-crud` 基础设施模块。
- `backend/go-common.md`：Go 公共库边界。

## 维护规则

1. 新增、移动或改名 skill 后，同步更新本文件和 `AGENTS.md`。
2. 不把一次性业务过程写入 skill，只沉淀可复用规则和当前项目事实。
3. 大段语法或协议参考单独成文，任务入口文档保持短而可执行。
4. 旧链接仍被历史记录引用时，保留兼容入口并说明新位置。
