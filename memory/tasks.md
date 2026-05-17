# Tasks

- [x] 整理 template 仓库项目级 skills，并放入 `.agents/skills/` 分层目录。
- [x] 添加 backend/front 变更驱动的 AI 技能自动优化 workflow。
- [x] 2026-05-08 移除 SysResource 后端接口、ORM 模型/生成代码、Swagger 文档和前端资源管理页 API。
- [x] 2026-05-08 实现 SysResourceMenu 动态菜单接口、前端动态侧栏和菜单资源树表管理页。
- [x] 2026-05-09 优化菜单资源页为右侧抽屉表单，并按菜单类型维护 `metadata` 字段。
- [x] 2026-05-10 实现 SysResourceApi 前后端 CRUD 管理页，并支持 `{id}` 到 `:id` 的路径参数模板归一化。
- [x] 2026-05-10 实现 `SysUser` 后端 CRUD：新增 `/api/sys/user/*` 列表/创建/更新/删除接口、Swagger 注释与 `logic` 层真实 ORM 测试，覆盖密码加密、重名校验、默认排序和软删除。
- [x] 2026-05-10 接入 `SysUser` 前端管理页：新增 `/system/user` 路由、`sysUser` API 封装与抽屉表单 CRUD 页面，并完成前端构建验证。
- [x] 2026-05-11 移除 SysUser.SysRoles 关联字段，登录角色逻辑改为直接查询 SysUserRole。
- [x] 2026-05-11 移除 SysRole.ChildIDs 字段，并同步初始化 SQL、迁移脚本与 ORM query。
- [x] 2026-05-17 将 `admin-react`、`app-react`、`app-react-ssr` 的可复用前端代码抽取到 `front/packages`，优先保留 app wrapper 降低迁移风险。
