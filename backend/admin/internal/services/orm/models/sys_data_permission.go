package models

import (
	"orm-crud/gormc/mixin"

	"gorm.io/datatypes"
	"gorm.io/plugin/soft_delete"
)

func init() {
	Models = append(Models, &SysDataPermission{})
}

type SysDataPermission struct {
	mixin.AutoIncrementID
	mixin.CreatedAt
	mixin.UpdatedAt
	mixin.Remark
	mixin.OperatorID
	mixin.IsEnabled
	DeletedAt     soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;not null;default:0;uniqueIndex:idx_sys_data_permission_subject_resource_action_deleted_at,priority:4" json:"deletedAt"`
	SubjectType   string                `gorm:"column:subject_type;type:varchar(16);not null;uniqueIndex:idx_sys_data_permission_subject_resource_action_deleted_at,priority:1;comment:主体类型(USER/ROLE/ANY_USER/ANY_ROLE)" json:"subjectType"`
	SubjectID     uint64                `gorm:"column:subject_id;type:bigint;not null;index:idx_sys_data_permission_subject;uniqueIndex:idx_sys_data_permission_subject_resource_action_deleted_at,priority:2;comment:主体ID，ANY_*时为0" json:"subjectID"`
	ResourceTable string                `gorm:"column:resource_table;type:varchar(32);not null;uniqueIndex:idx_sys_data_permission_subject_resource_action_deleted_at,priority:3;comment:资源表名" json:"resourceTable"`
	Action        datatypes.JSON        `gorm:"column:action;not null;default:'[\"read\"]';uniqueIndex:idx_sys_data_permission_subject_resource_action_deleted_at,priority:5;comment:操作列表(all/read/write/delete)" json:"action"`
	ScopeType     string                `gorm:"column:scope_type;type:varchar(32);not null;default:none;comment:作用域类型(all/none/include/exclude/owner/custom)" json:"scopeType"`
	ScopeField    string                `gorm:"column:scope_field;type:varchar(64);not null;default:id;comment:用于匹配scope_values的字段" json:"scopeField"`
	ScopeValues   datatypes.JSON        `gorm:"column:scope_values;not null;default:'[]';comment:作用域值列表" json:"scopeValues"`
	Conditions    datatypes.JSONMap     `gorm:"column:conditions;not null;default:'{}';comment:行过滤条件" json:"conditions"`
	Priority      int                   `gorm:"column:priority;type:integer;not null;default:0;comment:多角色冲突时的优先级" json:"priority"`
}

func (*SysDataPermission) TableName() string {
	return "sys_data_permission"
}
