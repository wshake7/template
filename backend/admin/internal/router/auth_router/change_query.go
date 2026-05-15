package auth_router

import (
	"errors"
	"strings"
	"time"

	"admin/internal/fiberc/handler"
	"admin/internal/router/logic"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"go-common/utils/pt"
	"gorm.io/gorm"
)

type jobExecutionChange struct {
	ID           uint64     `json:"id"`
	Status       string     `json:"status" change:"状态"`
	EndTime      *time.Time `json:"endTime" change:"结束时间"`
	ErrorMessage string     `json:"errorMessage" change:"错误信息"`
}

func bindChangeReq[T any](ctx *handler.Ctx) (*T, error) {
	req := new(T)
	if err := ctx.DefaultCtx.Bind().All(req); err != nil {
		return nil, err
	}
	return req, nil
}

func nilOnRecordNotFound[T any](item *T, err error) (*T, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return item, err
}

func sysUserCreateChangeQuery(ctx *handler.Ctx) (any, error) {
	req, err := bindChangeReq[logic.ReqSysUserCreate](ctx)
	if err != nil {
		return nil, err
	}
	m, err := nilOnRecordNotFound(query.SysUser.Where(query.SysUser.Username.Eq(strings.TrimSpace(req.Username))).First())
	if m == nil || err != nil {
		return m, err
	}
	return &logic.ReqSysUserCreate{
		Username:     m.Username,
		Nickname:     m.Nickname,
		LanguageCode: m.LanguageCode,
		IsEnabled:    m.IsEnabled.IsEnabled,
		Remark:       m.Remark.Remark,
	}, nil
}

func sysUserUpdateChangeQuery(ctx *handler.Ctx) (any, error) {
	req, err := bindChangeReq[logic.ReqSysUserUpdate](ctx)
	if err != nil {
		return nil, err
	}
	m, err := nilOnRecordNotFound(query.SysUser.Where(query.SysUser.ID.Eq(req.ID)).First())
	if m == nil || err != nil {
		return m, err
	}
	return &logic.ReqSysUserUpdate{
		ID:           m.ID,
		Username:     pt.Ptr(m.Username),
		Nickname:     pt.Ptr(m.Nickname),
		LanguageCode: pt.Ptr(m.LanguageCode),
		IsEnabled:    pt.Ptr(m.IsEnabled.IsEnabled),
		Remark:       pt.Ptr(m.Remark.Remark),
	}, nil
}

func sysRoleCreateChangeQuery(ctx *handler.Ctx) (any, error) {
	req, err := bindChangeReq[logic.ReqSysRoleCreate](ctx)
	if err != nil {
		return nil, err
	}
	m, err := nilOnRecordNotFound(query.SysRole.Where(query.SysRole.Code.Eq(strings.TrimSpace(req.Code))).First())
	if m == nil || err != nil {
		return m, err
	}
	return &logic.ReqSysRoleCreate{
		Name:      m.Name,
		Code:      m.Code,
		ParentID:  m.ParentID,
		IsEnabled: m.IsEnabled.IsEnabled,
		Remark:    m.Remark.Remark,
	}, nil
}

func sysRoleUpdateChangeQuery(ctx *handler.Ctx) (any, error) {
	req, err := bindChangeReq[logic.ReqSysRoleUpdate](ctx)
	if err != nil {
		return nil, err
	}
	m, err := nilOnRecordNotFound(query.SysRole.Where(query.SysRole.ID.Eq(req.ID)).First())
	if m == nil || err != nil {
		return m, err
	}
	return &logic.ReqSysRoleUpdate{
		ID:        m.ID,
		Name:      pt.Ptr(m.Name),
		Code:      pt.Ptr(m.Code),
		ParentID:  m.ParentID,
		IsEnabled: pt.Ptr(m.IsEnabled.IsEnabled),
		Remark:    pt.Ptr(m.Remark.Remark),
	}, nil
}

func sysRolePermissionChangeQuery(ctx *handler.Ctx) (any, error) {
	req, err := bindChangeReq[logic.ReqSysRolePermissionSave](ctx)
	if err != nil {
		return nil, err
	}
	menus, err := query.SysRoleMenu.Select(query.SysRoleMenu.MenuID).Where(query.SysRoleMenu.RoleID.Eq(req.ID)).Find()
	if err != nil {
		return nil, err
	}
	apis, err := query.SysRoleApi.Select(query.SysRoleApi.ApiID).Where(query.SysRoleApi.RoleID.Eq(req.ID)).Find()
	if err != nil {
		return nil, err
	}
	result := &logic.ReqSysRolePermissionSave{ID: req.ID}
	for _, item := range menus {
		result.MenuIDs = append(result.MenuIDs, item.MenuID)
	}
	for _, item := range apis {
		result.ApiIDs = append(result.ApiIDs, item.ApiID)
	}
	return result, nil
}

func resourceApiCreateChangeQuery(ctx *handler.Ctx) (any, error) {
	req, err := bindChangeReq[logic.ReqResourceApiCreate](ctx)
	if err != nil {
		return nil, err
	}
	path := normalizeResourceApiPathChange(strings.TrimSpace(req.Path))
	method := strings.ToUpper(strings.TrimSpace(req.Method))
	m, err := nilOnRecordNotFound(query.SysResourceApi.Where(query.SysResourceApi.Path.Eq(path), query.SysResourceApi.Method.Eq(method)).First())
	if m == nil || err != nil {
		return m, err
	}
	return &logic.ReqResourceApiCreate{
		Module:    m.Module,
		Path:      m.Path,
		Method:    m.Method,
		SortOrder: m.SortOrder.SortOrder,
		IsEnabled: m.IsEnabled.IsEnabled,
		Remark:    m.Remark.Remark,
	}, nil
}

func resourceApiUpdateChangeQuery(ctx *handler.Ctx) (any, error) {
	req, err := bindChangeReq[logic.ReqResourceApiUpdate](ctx)
	if err != nil {
		return nil, err
	}
	m, err := nilOnRecordNotFound(query.SysResourceApi.Where(query.SysResourceApi.ID.Eq(req.ID)).First())
	if m == nil || err != nil {
		return m, err
	}
	return &logic.ReqResourceApiUpdate{
		ID:        m.ID,
		Module:    pt.Ptr(m.Module),
		Path:      pt.Ptr(m.Path),
		Method:    pt.Ptr(m.Method),
		SortOrder: pt.Ptr(m.SortOrder.SortOrder),
		IsEnabled: pt.Ptr(m.IsEnabled.IsEnabled),
		Remark:    pt.Ptr(m.Remark.Remark),
	}, nil
}

func resourceMenuCreateChangeQuery(ctx *handler.Ctx) (any, error) {
	req, err := bindChangeReq[logic.ReqResourceMenuCreate](ctx)
	if err != nil {
		return nil, err
	}
	m, err := nilOnRecordNotFound(query.SysResourceMenu.
		Where(
			query.SysResourceMenu.MenuType.Eq(strings.ToUpper(strings.TrimSpace(req.MenuType))),
			query.SysResourceMenu.Name.Eq(strings.TrimSpace(req.Name)),
			query.SysResourceMenu.Path.Eq(strings.TrimSpace(req.Path)),
		).
		Order(query.SysResourceMenu.ID.Desc()).
		First())
	if m == nil || err != nil {
		return m, err
	}
	return resourceMenuCreateFromModel(m), nil
}

func resourceMenuUpdateChangeQuery(ctx *handler.Ctx) (any, error) {
	req, err := bindChangeReq[logic.ReqResourceMenuUpdate](ctx)
	if err != nil {
		return nil, err
	}
	m, err := nilOnRecordNotFound(query.SysResourceMenu.Where(query.SysResourceMenu.ID.Eq(req.ID)).First())
	if m == nil || err != nil {
		return m, err
	}
	return &logic.ReqResourceMenuUpdate{
		ID:        m.ID,
		ParentID:  m.ParentID,
		MenuType:  pt.Ptr(m.MenuType),
		Path:      pt.Ptr(m.Path),
		Redirect:  pt.Ptr(m.Redirect),
		Alias:     pt.Ptr(m.Alias),
		Name:      pt.Ptr(m.Name),
		Component: pt.Ptr(m.Component),
		Metadata:  pt.Ptr(m.Metadata.Metadata),
		SortOrder: pt.Ptr(m.SortOrder.SortOrder),
		IsEnabled: pt.Ptr(m.IsEnabled.IsEnabled),
		Remark:    pt.Ptr(m.Remark.Remark),
	}, nil
}

func resourceMenuCreateFromModel(m *models.SysResourceMenu) *logic.ReqResourceMenuCreate {
	return &logic.ReqResourceMenuCreate{
		ParentID:  m.ParentID,
		MenuType:  m.MenuType,
		Path:      m.Path,
		Redirect:  m.Redirect,
		Alias:     m.Alias,
		Name:      m.Name,
		Component: m.Component,
		Metadata:  m.Metadata.Metadata,
		SortOrder: m.SortOrder.SortOrder,
		IsEnabled: m.IsEnabled.IsEnabled,
		Remark:    m.Remark.Remark,
	}
}

func jobScheduleCreateChangeQuery(ctx *handler.Ctx) (any, error) {
	req, err := bindChangeReq[logic.ReqJobScheduleCreate](ctx)
	if err != nil {
		return nil, err
	}
	m, err := nilOnRecordNotFound(query.JobSchedule.Where(query.JobSchedule.JobCode.Eq(strings.TrimSpace(req.JobCode))).First())
	if m == nil || err != nil {
		return m, err
	}
	return jobScheduleCreateFromModel(m), nil
}

func jobScheduleUpdateChangeQuery(ctx *handler.Ctx) (any, error) {
	req, err := bindChangeReq[logic.ReqJobScheduleUpdate](ctx)
	if err != nil {
		return nil, err
	}
	m, err := nilOnRecordNotFound(query.JobSchedule.Where(query.JobSchedule.ID.Eq(req.ID)).First())
	if m == nil || err != nil {
		return m, err
	}
	return &logic.ReqJobScheduleUpdate{
		ID:                       m.ID,
		JobName:                  pt.Ptr(m.JobName),
		WorkflowType:             pt.Ptr(m.WorkflowType),
		TaskQueue:                pt.Ptr(m.TaskQueue),
		ScheduleType:             pt.Ptr(m.ScheduleType),
		CronExpr:                 pt.Ptr(m.CronExpr),
		IntervalSeconds:          m.IntervalSeconds,
		StartTime:                m.StartTime,
		EndTime:                  m.EndTime,
		InputJSON:                pt.Ptr(string(m.InputJSON)),
		Status:                   pt.Ptr(m.Status),
		TemporalWorkflowIDPrefix: pt.Ptr(m.TemporalWorkflowIDPrefix),
		Description:              pt.Ptr(m.Description),
	}, nil
}

func jobScheduleSwitchChangeQuery(ctx *handler.Ctx) (any, error) {
	req, err := bindChangeReq[logic.ReqJobScheduleSwitch](ctx)
	if err != nil {
		return nil, err
	}
	m, err := nilOnRecordNotFound(query.JobSchedule.Where(query.JobSchedule.ID.Eq(req.ID)).First())
	if m == nil || err != nil {
		return m, err
	}
	return &logic.ReqJobScheduleSwitch{ID: m.ID, Enabled: m.Status == models.JobScheduleStatusEnabled}, nil
}

func jobScheduleCreateFromModel(m *models.JobSchedule) *logic.ReqJobScheduleCreate {
	return &logic.ReqJobScheduleCreate{
		JobCode:                  m.JobCode,
		JobName:                  m.JobName,
		WorkflowType:             m.WorkflowType,
		TaskQueue:                m.TaskQueue,
		ScheduleType:             m.ScheduleType,
		CronExpr:                 m.CronExpr,
		IntervalSeconds:          m.IntervalSeconds,
		StartTime:                m.StartTime,
		EndTime:                  m.EndTime,
		InputJSON:                string(m.InputJSON),
		Status:                   m.Status,
		TemporalScheduleID:       m.TemporalScheduleID,
		TemporalWorkflowIDPrefix: m.TemporalWorkflowIDPrefix,
		Description:              m.Description,
	}
}

func normalizeResourceApiPathChange(path string) string {
	if path == "" {
		return ""
	}
	segments := strings.Split(path, "/")
	for i, segment := range segments {
		if len(segment) >= 3 && strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}") {
			name := strings.TrimSpace(segment[1 : len(segment)-1])
			if name != "" {
				segments[i] = ":" + name
			}
		}
	}
	return strings.Join(segments, "/")
}

func jobExecutionChangeQuery(ctx *handler.Ctx) (any, error) {
	req, err := bindChangeReq[logic.ReqJobExecutionID](ctx)
	if err != nil {
		return nil, err
	}
	m, err := nilOnRecordNotFound(query.JobExecution.Where(query.JobExecution.ID.Eq(req.ID)).First())
	if m == nil || err != nil {
		return m, err
	}
	return &jobExecutionChange{
		ID:           m.ID,
		Status:       m.Status,
		EndTime:      m.EndTime,
		ErrorMessage: m.ErrorMessage,
	}, nil
}

func dictTypeCreateChangeQuery(ctx *handler.Ctx) (any, error) {
	req, err := bindChangeReq[logic.ReqDictTypeCreate](ctx)
	if err != nil {
		return nil, err
	}
	m, err := nilOnRecordNotFound(query.SysDictType.Where(query.SysDictType.TypeCode.Eq(req.TypeCode)).First())
	if m == nil || err != nil {
		return m, err
	}
	return &logic.ReqDictTypeCreate{
		TypeCode:  m.TypeCode,
		TypeName:  m.TypeName,
		IsEnabled: m.IsEnabled.IsEnabled,
		SortOrder: m.SortOrder.SortOrder,
		Remark:    m.Remark.Remark,
	}, nil
}

func dictTypeUpdateChangeQuery(ctx *handler.Ctx) (any, error) {
	req, err := bindChangeReq[logic.ReqDictTypeUpdate](ctx)
	if err != nil {
		return nil, err
	}
	m, err := nilOnRecordNotFound(query.SysDictType.Where(query.SysDictType.ID.Eq(req.ID)).First())
	if m == nil || err != nil {
		return m, err
	}
	return &logic.ReqDictTypeUpdate{
		ID:        m.ID,
		TypeCode:  pt.Ptr(m.TypeCode),
		TypeName:  pt.Ptr(m.TypeName),
		IsEnabled: pt.Ptr(m.IsEnabled.IsEnabled),
		SortOrder: pt.Ptr(m.SortOrder.SortOrder),
		Remark:    pt.Ptr(m.Remark.Remark),
	}, nil
}

func dictEntryCreateChangeQuery(ctx *handler.Ctx) (any, error) {
	req, err := bindChangeReq[logic.ReqDictEntryCreate](ctx)
	if err != nil {
		return nil, err
	}
	m, err := nilOnRecordNotFound(query.SysDictEntry.
		Where(
			query.SysDictEntry.SysDictTypeId.Eq(req.SysDictTypeId),
			query.SysDictEntry.EntryValue.Eq(req.EntryValue),
			query.SysDictEntry.EntryLabel.Eq(req.EntryLabel),
		).
		Order(query.SysDictEntry.ID.Desc()).
		First())
	if m == nil || err != nil {
		return m, err
	}
	return &logic.ReqDictEntryCreate{
		LabelComponent: m.LabelComponent,
		EntryLabel:     m.EntryLabel,
		EntryValue:     m.EntryValue,
		LanguageCode:   m.LanguageCode,
		SysDictTypeId:  m.SysDictTypeId,
		SortOrder:      m.SortOrder.SortOrder,
		IsEnabled:      m.IsEnabled.IsEnabled,
		Remark:         m.Remark.Remark,
	}, nil
}

func dictEntryUpdateChangeQuery(ctx *handler.Ctx) (any, error) {
	req, err := bindChangeReq[logic.ReqDictEntryUpdate](ctx)
	if err != nil {
		return nil, err
	}
	ids := make([]uint64, 0, len(req.Updates)+1)
	if req.ID != nil {
		ids = append(ids, *req.ID)
	}
	for _, item := range req.Updates {
		ids = append(ids, item.ID)
	}
	if len(ids) == 0 {
		return &logic.ReqDictEntryUpdate{}, nil
	}
	items, err := query.SysDictEntry.Where(query.SysDictEntry.ID.In(ids...)).Order(query.SysDictEntry.ID.Asc()).Find()
	if err != nil {
		return nil, err
	}
	result := &logic.ReqDictEntryUpdate{}
	for _, m := range items {
		result.Updates = append(result.Updates, logic.ReqDictEntryUpdateItem{
			ID:             m.ID,
			LabelComponent: pt.Ptr(m.LabelComponent),
			EntryLabel:     pt.Ptr(m.EntryLabel),
			EntryValue:     pt.Ptr(m.EntryValue),
			LanguageCode:   pt.Ptr(m.LanguageCode),
			SysDictTypeId:  &m.SysDictTypeId,
			SortOrder:      pt.Ptr(m.SortOrder.SortOrder),
			IsEnabled:      pt.Ptr(m.IsEnabled.IsEnabled),
			Remark:         pt.Ptr(m.Remark.Remark),
		})
	}
	return result, nil
}

func dictEntryBatchCopyChangeQuery(ctx *handler.Ctx) (any, error) {
	req, err := bindChangeReq[logic.ReqDictEntryBatchCopy](ctx)
	if err != nil {
		return nil, err
	}
	sourceEntries, err := query.SysDictEntry.Where(query.SysDictEntry.ID.In(req.EntryIds...)).Find()
	if err != nil {
		return nil, err
	}
	values := make([]string, 0, len(sourceEntries))
	for _, item := range sourceEntries {
		values = append(values, item.EntryValue)
	}
	result := &logic.ReqDictEntryUpdate{}
	if len(values) == 0 {
		return result, nil
	}
	items, err := query.SysDictEntry.
		Where(query.SysDictEntry.SysDictTypeId.Eq(req.TargetTypeId), query.SysDictEntry.EntryValue.In(values...)).
		Order(query.SysDictEntry.ID.Asc()).
		Find()
	if err != nil {
		return nil, err
	}
	for _, m := range items {
		result.Updates = append(result.Updates, logic.ReqDictEntryUpdateItem{
			ID:             m.ID,
			LabelComponent: pt.Ptr(m.LabelComponent),
			EntryLabel:     pt.Ptr(m.EntryLabel),
			EntryValue:     pt.Ptr(m.EntryValue),
			LanguageCode:   pt.Ptr(m.LanguageCode),
			SysDictTypeId:  &m.SysDictTypeId,
			SortOrder:      pt.Ptr(m.SortOrder.SortOrder),
			IsEnabled:      pt.Ptr(m.IsEnabled.IsEnabled),
			Remark:         pt.Ptr(m.Remark.Remark),
		})
	}
	return result, nil
}

func langTypeCreateChangeQuery(ctx *handler.Ctx) (any, error) {
	req, err := bindChangeReq[logic.ReqLangTypeCreate](ctx)
	if err != nil {
		return nil, err
	}
	m, err := nilOnRecordNotFound(query.SysLanguageType.Where(query.SysLanguageType.TypeCode.Eq(req.TypeCode)).First())
	if m == nil || err != nil {
		return m, err
	}
	return &logic.ReqLangTypeCreate{
		TypeCode:  m.TypeCode,
		TypeName:  m.TypeName,
		IsDefault: m.IsDefault,
		IsEnabled: m.IsEnabled.IsEnabled,
		SortOrder: m.SortOrder.SortOrder,
	}, nil
}

func langTypeUpdateChangeQuery(ctx *handler.Ctx) (any, error) {
	req, err := bindChangeReq[logic.ReqLangTypeUpdate](ctx)
	if err != nil {
		return nil, err
	}
	m, err := nilOnRecordNotFound(query.SysLanguageType.Where(query.SysLanguageType.ID.Eq(req.ID)).First())
	if m == nil || err != nil {
		return m, err
	}
	return &logic.ReqLangTypeUpdate{
		ID:        m.ID,
		TypeCode:  pt.Ptr(m.TypeCode),
		TypeName:  pt.Ptr(m.TypeName),
		IsDefault: pt.Ptr(m.IsDefault),
		IsEnabled: pt.Ptr(m.IsEnabled.IsEnabled),
		SortOrder: pt.Ptr(m.SortOrder.SortOrder),
	}, nil
}

func langEntryCreateChangeQuery(ctx *handler.Ctx) (any, error) {
	req, err := bindChangeReq[logic.ReqLangEntryCreate](ctx)
	if err != nil {
		return nil, err
	}
	m, err := nilOnRecordNotFound(query.SysLanguageEntry.
		Where(query.SysLanguageEntry.SysLanguageTypeId.Eq(req.SysLanguageTypeId), query.SysLanguageEntry.EntryCode.Eq(req.EntryCode)).
		First())
	if m == nil || err != nil {
		return m, err
	}
	return &logic.ReqLangEntryCreate{
		EntryCode:         m.EntryCode,
		EntryValue:        m.EntryValue,
		SysLanguageTypeId: m.SysLanguageTypeId,
		SortOrder:         m.SortOrder.SortOrder,
		IsEnabled:         m.IsEnabled.IsEnabled,
		Remark:            m.Remark.Remark,
	}, nil
}

func langEntryUpdateChangeQuery(ctx *handler.Ctx) (any, error) {
	req, err := bindChangeReq[logic.ReqLangEntryUpdate](ctx)
	if err != nil {
		return nil, err
	}
	ids := make([]uint64, 0, len(req.Updates)+1)
	if req.ID != nil {
		ids = append(ids, *req.ID)
	}
	for _, item := range req.Updates {
		ids = append(ids, item.ID)
	}
	if len(ids) == 0 {
		return &logic.ReqLangEntryUpdate{}, nil
	}
	items, err := query.SysLanguageEntry.Where(query.SysLanguageEntry.ID.In(ids...)).Order(query.SysLanguageEntry.ID.Asc()).Find()
	if err != nil {
		return nil, err
	}
	result := &logic.ReqLangEntryUpdate{}
	for _, m := range items {
		result.Updates = append(result.Updates, logic.ReqLangEntryUpdateItem{
			ID:                m.ID,
			EntryCode:         pt.Ptr(m.EntryCode),
			EntryValue:        pt.Ptr(m.EntryValue),
			SysLanguageTypeId: &m.SysLanguageTypeId,
			SortOrder:         pt.Ptr(m.SortOrder.SortOrder),
			IsEnabled:         pt.Ptr(m.IsEnabled.IsEnabled),
			Remark:            pt.Ptr(m.Remark.Remark),
		})
	}
	return result, nil
}

func langEntryBatchCreateChangeQuery(ctx *handler.Ctx) (any, error) {
	req, err := bindChangeReq[logic.ReqLangEntryBatchCreate](ctx)
	if err != nil {
		return nil, err
	}
	typeCodes := make([]string, 0, len(req.Values))
	for code := range req.Values {
		typeCodes = append(typeCodes, code)
	}
	result := &logic.ReqLangEntryBatchCreate{
		EntryCode: req.EntryCode,
		Values:    map[string]string{},
		SortOrder: req.SortOrder,
		IsEnabled: req.IsEnabled,
	}
	if len(typeCodes) == 0 {
		return result, nil
	}
	types, err := query.SysLanguageType.Where(query.SysLanguageType.TypeCode.In(typeCodes...)).Find()
	if err != nil {
		return nil, err
	}
	typeIDToCode := make(map[uint64]string, len(types))
	typeIDs := make([]uint64, 0, len(types))
	for _, item := range types {
		typeIDToCode[item.ID] = item.TypeCode
		typeIDs = append(typeIDs, item.ID)
	}
	if len(typeIDs) == 0 {
		return result, nil
	}
	entries, err := query.SysLanguageEntry.
		Where(query.SysLanguageEntry.EntryCode.Eq(req.EntryCode), query.SysLanguageEntry.SysLanguageTypeId.In(typeIDs...)).
		Find()
	if err != nil {
		return nil, err
	}
	for _, item := range entries {
		if code, ok := typeIDToCode[item.SysLanguageTypeId]; ok {
			result.Values[code] = item.EntryValue
		}
	}
	return result, nil
}
