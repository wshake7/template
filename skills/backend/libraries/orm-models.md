# Backend ORM Models 编写与维护

## 场景与目标
- 适用场景：修改 `backend/admin/services/orm/models` 下实体模型、字段、关联关系、表结构时
- 目标：保持模型定义、自动迁移、`gorm/gen` 查询代码三者一致，降低回归风险

## 目录/文件位置
- 模型目录：`backend/admin/services/orm/models/`
- 模型注册入口：`backend/admin/services/orm/models/models.go`
- ORM 初始化：`backend/admin/services/orm/orm.go`
- 查询代码生成目录：`backend/admin/services/orm/query/`
- 生成脚本：`backend/admin/cmd/scripts/orm/main.go`

## 当前模型组织方式（按现状）
1. 每个模型文件通过 `init()` 执行 `Models = append(Models, &Xxx{})`
2. `orm.New()` 中通过 `gormCrud.WithAutoMigrate(models.Models...)` 执行自动迁移
3. `orm.New()` / `cmd/scripts/orm/main.go` 可触发 `gen` 生成 `query/*.gen.go`
4. 业务层统一通过 `query.SysUser`、`query.SysRole` 等访问

## 字段风格与约定
- 通用字段优先复用 `orm-crud/gorm/mixin`（如 `CreatedAt`、`Status`、`OperatorID`）
- 软删除字段统一使用 `gorm.DeletedAt` 或 `soft_delete.DeletedAt`，并配合唯一索引
- 表名必须实现 `TableName()`，保持显式可控
- 关联字段显式声明 `foreignKey` / `references` / `constraint`

## 常见模型关系（当前）
- 用户与角色：`SysUser` <-> `SysRole`（`sys_user_role` 多对多中间表）
- 字典类型与字典项：`SysDictType` -> `[]SysDictEntry`（一对多）
- 角色树：`SysRole` 的 `ParentSysRole` / `Children`（自关联）

## 修改步骤（推荐）
1. 改模型结构：编辑 `models/*.go`
2. 确认模型已注册：检查对应 `init()` 是否 append 到 `Models`
3. 若涉及迁移：确认 `config.Orm.IsAutoMigrate` 或手工迁移策略
4. 重新生成查询代码（必要时）：运行脚本更新 `query/*.gen.go`
5. 回归调用点：检查 `router/logic` 和 `services` 中的 query 使用

## 常用命令
```bash
cd backend/admin

# 启动服务（带自动迁移配置时可触发表结构更新）
go run ./cmd/main.go -f ./etc/config.yaml

# 生成 orm query 代码
go run ./cmd/scripts/orm

# 全量测试
go test ./...
```

## 注意事项
1. 新增模型但未注册到 `Models`，会导致迁移与代码生成遗漏
2. 改字段名/类型后，必须同步关注生成的 `query` 与业务使用处
3. 索引与唯一约束变更要先评估线上数据兼容性
4. `sys_login_log.go` 当前为占位文件，若计划启用请先补完整模型并注册
