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
	mixin.CreatedBy
	mixin.UpdatedBy
	mixin.IsEnabled
	DeletedAt     soft_delete.DeletedAt `gorm:"column:deleted_at;softDelete:milli;default:0;uniqueIndex:idx_sys_data_permission_subject_resource_action_deleted_at,priority:4" json:"deletedAt"`
	SubjectType   string                `gorm:"column:subject_type;type:varchar(16);not null;uniqueIndex:idx_sys_data_permission_subject_resource_action_deleted_at,priority:1;comment:主体类型(USER/ROLE)" json:"subjectType"`
	SubjectID     uint64                `gorm:"column:subject_id;type:bigint;not null;index:idx_sys_data_permission_subject;uniqueIndex:idx_sys_data_permission_subject_resource_action_deleted_at,priority:2;comment:主体ID" json:"subjectID"`
	ResourceTable string                `gorm:"column:resource_table;type:varchar(32);not null;uniqueIndex:idx_sys_data_permission_subject_resource_action_deleted_at,priority:3;comment:资源表名" json:"resourceTable"`
	Action        datatypes.JSON        `gorm:"column:action;not null;default:'[\"read\"]';uniqueIndex:idx_sys_data_permission_subject_resource_action_deleted_at,priority:5;comment:操作(read/write/delete)" json:"action"`
	ScopeType     string                `gorm:"column:scope_type;type:varchar(32);not null;default:none;comment:scope类型(all/none/include/exclude/owner/custom)" json:"scopeType"`
	ScopeValues   datatypes.JSON        `gorm:"column:scope_values;not null;default:'[]';comment:scope值(ID集合)" json:"scopeValues"`
	Conditions    datatypes.JSONMap     `gorm:"column:conditions;not null;default:'{}';comment:行级条件" json:"conditions"`
	Effect        string                `gorm:"column:effect;type:varchar(16);not null;default:allow;comment:生效类型(allow/deny)" json:"effect"`
	Priority      int                   `gorm:"column:priority;type:integer;not null;default:0;comment:优先级，多角色冲突时使用" json:"priority"`
}

func (*SysDataPermission) TableName() string {
	return "sys_data_permission"
}
