package logic

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/res"
	datapermission "admin/internal/services/orm/data_permission"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"go-common/utils/slices_utils"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc"
	"orm-crud/gormc/mixin"
	paginationFilter "orm-crud/pagination/filter"

	"go.uber.org/zap"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

type SysDictHandler struct{}

// --- 字典类型 (DictType) ---

type ReqDictTypeCreate struct {
	TypeCode  string `json:"typeCode" binding:"required,max=128" binding_msg:"required=字典类型代码不能为空,max=字典类型代码最多128位"`
	TypeName  string `json:"typeName" binding:"required,max=255" binding_msg:"required=字典类型名称不能为空,max=字典类型名称最多255位"`
	IsEnabled bool   `json:"isEnabled"`
	SortOrder int32  `json:"sortOrder"`
	Remark    string `json:"remark" binding:"max=255" binding_msg:"max=描述最多255位"`
}

type ReqDictTypeUpdate struct {
	ID        uint64  `json:"id" binding:"required" binding_msg:"required=请求错误"`
	TypeCode  *string `json:"typeCode" binding:"omitempty,max=128" binding_msg:"max=字典类型代码最多128位"`
	TypeName  *string `json:"typeName" binding:"omitempty,max=255" binding_msg:"max=字典类型名称最多255位"`
	IsEnabled *bool   `json:"isEnabled"`
	SortOrder *int32  `json:"sortOrder"`
	Remark    *string `json:"remark" binding:"omitempty,max=255" binding_msg:"max=描述最多255位"`
}

type ReqDictTypeBatchDelete struct {
	IDs []uint64 `json:"ids" binding:"required,min=1" binding_msg:"required=请选择字典类型,min=至少选择一项"`
}

type RespDictType struct {
	models.SysDictType
	CanWrite  bool `json:"canWrite"`
	CanDelete bool `json:"canDelete"`
}

// @Summary 获取字典类型分页列表
// @Remark 分页查询字典类型信息
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body v1.PagingRequest true "分页参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[RespDictType]} "成功"
// @Router /api/dict/type/list [post]
func (*SysDictHandler) TypeList(ctx *handler.Ctx, req *v1.PagingRequest) (*gormc.PagingResult[RespDictType], error) {
	permissionExprs, err := datapermission.BuildPermissionFilterExprsForCtx(
		ctx,
		models.SysDictType{}.TableName(),
		datapermission.ActionRead,
		datapermission.ActionWrite,
		datapermission.ActionDelete,
	)
	if err != nil {
		ctx.L().Error("build dict type permission expressions failed", zap.Error(err))
		return nil, res.FailDefault
	}
	if err := datapermission.ApplyPagePermissionExpr(req, permissionExprs[datapermission.ActionRead]); err != nil {
		ctx.L().Error("apply dict type read permission failed", zap.Error(err))
		return nil, res.FailDefault
	}

	pagination, err := query.SysDictType.PageWithPaging(req)
	if err != nil {
		return nil, res.FailDefault
	}

	items := make([]*RespDictType, 0, len(pagination.Items))
	ids := make([]uint64, 0, len(pagination.Items))
	for _, item := range pagination.Items {
		ids = append(ids, item.ID)
	}

	writeIDSet, err := queryAllowedDictTypeIDSetByExpr(ids, permissionExprs[datapermission.ActionWrite])
	if err != nil {
		ctx.L().Error("apply dict type write permission failed", zap.Error(err), zap.Uint64s("ids", ids))
		return nil, res.FailDefault
	}
	deleteIDSet := writeIDSet
	if !reflect.DeepEqual(permissionExprs[datapermission.ActionWrite], permissionExprs[datapermission.ActionDelete]) {
		deleteIDSet, err = queryAllowedDictTypeIDSetByExpr(ids, permissionExprs[datapermission.ActionDelete])
		if err != nil {
			ctx.L().Error("apply dict type delete permission failed", zap.Error(err), zap.Uint64s("ids", ids))
			return nil, res.FailDefault
		}
	}

	for _, item := range pagination.Items {
		items = append(items, &RespDictType{
			SysDictType: *item,
			CanWrite:    writeIDSet[item.ID],
			CanDelete:   deleteIDSet[item.ID],
		})
	}
	return &gormc.PagingResult[RespDictType]{
		Items: items,
		Total: pagination.Total,
	}, nil
}

func queryAllowedDictTypeIDSetByExpr(ids []uint64, permissionExpr *v1.FilterExpr) (map[uint64]bool, error) {
	allowedIDSet := make(map[uint64]bool, len(ids))
	if len(ids) == 0 {
		return allowedIDSet, nil
	}

	permissionQuery, err := datapermission.BuildPermissionQueryFromExpr(permissionExpr)
	if err != nil {
		return nil, err
	}

	sysDictType := permissionQuery.SysDictType
	allowedTypes, err := sysDictType.
		Select(sysDictType.ID).
		Where(sysDictType.ID.In(ids...)).
		Find()
	if err != nil {
		return nil, err
	}
	for _, item := range allowedTypes {
		allowedIDSet[item.ID] = true
	}
	return allowedIDSet, nil
}

func queryAllowedDictTypeIDSet(ctx *handler.Ctx, ids []uint64, buildPermissionQuery func(*handler.Ctx, string) (*query.Query, error)) (map[uint64]bool, error) {
	allowedIDSet := make(map[uint64]bool, len(ids))
	if len(ids) == 0 {
		return allowedIDSet, nil
	}

	permissionQuery, err := buildPermissionQuery(ctx, models.SysDictType{}.TableName())
	if err != nil {
		return nil, err
	}

	sysDictType := permissionQuery.SysDictType
	allowedTypes, err := sysDictType.
		Select(sysDictType.ID).
		Where(sysDictType.ID.In(ids...)).
		Find()
	if err != nil {
		return nil, err
	}
	for _, item := range allowedTypes {
		allowedIDSet[item.ID] = true
	}
	return allowedIDSet, nil
}

func queryAllowedDictTypeIDs(ctx *handler.Ctx, buildPermissionQuery func(*handler.Ctx, string) (*query.Query, error)) ([]uint64, error) {
	permissionQuery, err := buildPermissionQuery(ctx, models.SysDictType{}.TableName())
	if err != nil {
		return nil, err
	}

	sysDictType := permissionQuery.SysDictType
	allowedTypes, err := sysDictType.
		Select(sysDictType.ID).
		Find()
	if err != nil {
		return nil, err
	}

	ids := make([]uint64, 0, len(allowedTypes))
	for _, item := range allowedTypes {
		ids = append(ids, item.ID)
	}
	return ids, nil
}

func applyDictEntryTypeIDPageFilter(req *v1.PagingRequest, typeIDs []uint64) error {
	values := make([]string, 0, len(typeIDs))
	for _, id := range typeIDs {
		values = append(values, strconv.FormatUint(id, 10))
	}

	typeFilterExpr := &v1.FilterExpr{
		Type: v1.ExprType_AND,
		Conditions: []*v1.FilterCondition{
			{
				Field:  "sysDictTypeId",
				Op:     v1.Operator_IN,
				Values: values,
			},
		},
	}

	currentExpr, err := paginationFilter.ConvertFilterByPagingRequest(req)
	if err != nil {
		return err
	}
	if currentExpr == nil {
		req.FilteringType = &v1.PagingRequest_FilterExpr{FilterExpr: typeFilterExpr}
		return nil
	}

	req.FilteringType = &v1.PagingRequest_FilterExpr{FilterExpr: &v1.FilterExpr{
		Type:   v1.ExprType_AND,
		Groups: []*v1.FilterExpr{currentExpr, typeFilterExpr},
	}}
	return nil
}

// @Summary 创建字典类型
// @Remark 创建新的字典类型
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body ReqDictTypeCreate true "字典类型创建参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/dict/type/create [post]
func (*SysDictHandler) TypeCreate(ctx *handler.Ctx, req *ReqDictTypeCreate) error {
	operationID := ctx.SessionInfo.Id

	err := query.SysDictType.Create(&models.SysDictType{
		OperatorID: mixin.OperatorID{
			CreatedBy: mixin.CreatedBy{CreatedBy: operationID},
			UpdatedBy: mixin.UpdatedBy{UpdatedBy: operationID},
		},
		IsEnabled: mixin.IsEnabled{IsEnabled: req.IsEnabled},
		SortOrder: mixin.SortOrder{SortOrder: req.SortOrder},
		Remark:    mixin.Remark{Remark: req.Remark},
		TypeCode:  req.TypeCode,
		TypeName:  req.TypeName,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return res.FailMsg("类型编码已存在")
		}
		return res.FailDefault
	}
	return nil
}

// @Summary 更新字典类型
// @Remark 根据 ID 更新字典类型信息
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body ReqDictTypeUpdate true "字典类型更新参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/dict/type/update [post]
func (*SysDictHandler) TypeUpdate(ctx *handler.Ctx, req *ReqDictTypeUpdate) error {
	operationID := ctx.SessionInfo.Id
	permissionQuery, err := datapermission.BuildWritePermissionQuery(ctx, models.SysDictType{}.TableName())
	if err != nil {
		ctx.L().Error("apply dict type write permission failed", zap.Error(err), zap.Uint64("id", req.ID))
		return res.FailDefault
	}

	sysDictType := permissionQuery.SysDictType

	exprs := []field.AssignExpr{sysDictType.UpdatedBy.Value(operationID)}
	query.ExprAppendSelf(&exprs, req.TypeCode, sysDictType.TypeCode.Value)
	query.ExprAppendSelf(&exprs, req.TypeName, sysDictType.TypeName.Value)
	query.ExprAppendSelf(&exprs, req.IsEnabled, sysDictType.IsEnabled.Value)
	query.ExprAppendSelf(&exprs, req.SortOrder, sysDictType.SortOrder.Value)
	query.ExprAppendSelf(&exprs, req.Remark, sysDictType.Remark.Value)

	info, err := sysDictType.Where(sysDictType.ID.Eq(req.ID)).UpdateSimple(exprs...)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return res.FailMsg("类型编码已存在")
		}
		return res.FailDefault
	}
	if info.RowsAffected == 0 {
		return res.FailMsg("无权限或数据不存在")
	}
	return nil
}

// @Summary 批量删除字典类型
// @Remark 根据 ID 列表批量删除字典类型及其关联的所有字典项
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body ReqDictTypeBatchDelete true "批量删除参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/dict/type/del [post]
func (*SysDictHandler) TypeDel(ctx *handler.Ctx, req *ReqDictTypeBatchDelete) error {
	ids := slices_utils.Distinct(req.IDs)
	permissionQuery, err := datapermission.BuildDeletePermissionQuery(ctx, models.SysDictType{}.TableName())
	if err != nil {
		ctx.L().Error("apply dict type delete permission failed", zap.Error(err), zap.Uint64s("ids", ids))
		return res.FailDefault
	}

	sysDictType := permissionQuery.SysDictType
	allowedTypes, err := sysDictType.
		Select(sysDictType.ID).
		Where(sysDictType.ID.In(ids...)).
		Find()
	if err != nil {
		ctx.L().Error("校验字典类型删除权限失败", zap.Error(err), zap.Uint64s("ids", ids))
		return res.FailDefault
	}
	if len(allowedTypes) != len(ids) {
		return res.FailMsg("无权限或数据不存在")
	}

	err = query.Q.Transaction(func(tx *query.Query) error {
		sysDictEntry := tx.SysDictEntry
		_, err = sysDictEntry.
			Where(sysDictEntry.SysDictTypeId.In(ids...)).
			Delete()
		if err != nil {
			return err
		}
		sysDictTypeSub := tx.SysDictType
		_, err = sysDictTypeSub.
			Where(sysDictTypeSub.ID.In(ids...)).
			Delete()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		ctx.L().Error("批量删除字典类型失败", zap.Error(err), zap.Uint64s("ids", ids))
		return res.FailDefault
	}
	return nil
}

// --- 字典数据项 (DictEntry) ---

type ReqDictEntryCreate struct {
	LabelComponent string `json:"labelComponent" binding:"omitempty,max=255" binding_msg:"max=显示标签组件最多255位"`
	EntryLabel     string `json:"entryLabel" binding:"required,max=255" binding_msg:"required=显示标签不能为空,max=显示标签最多255位"`
	EntryValue     string `json:"entryValue" binding:"required,max=255" binding_msg:"required=数据值不能为空,max=数据值最多255位"`
	LanguageCode   string `json:"languageCode" binding:"max=32" binding_msg:"max=语言代码最多32位"`
	SysDictTypeId  uint64 `json:"sysDictTypeId" binding:"required" binding_msg:"required=字典类型ID不能为空"`
	SortOrder      int32  `json:"sortOrder"`
	IsEnabled      bool   `json:"isEnabled"`
	Remark         string `json:"remark" binding:"max=255" binding_msg:"max=备注最多255位"`
}

type ReqDictEntryUpdate struct {
	ID             *uint64                  `json:"id"`
	LabelComponent *string                  `json:"labelComponent" binding:"omitempty,max=255" binding_msg:"max=显示标签组件最多255位"`
	EntryLabel     *string                  `json:"entryLabel" binding:"omitempty,max=255" binding_msg:"max=显示标签最多255位"`
	EntryValue     *string                  `json:"entryValue" binding:"omitempty,max=255" binding_msg:"max=数据值最多255位"`
	LanguageCode   *string                  `json:"languageCode" binding:"omitempty,max=32" binding_msg:"max=语言代码最多32位"`
	SysDictTypeId  *uint64                  `json:"sysDictTypeId"`
	SortOrder      *int32                   `json:"sortOrder"`
	IsEnabled      *bool                    `json:"isEnabled"`
	Remark         *string                  `json:"remark" binding:"omitempty,max=255" binding_msg:"max=备注最多255位"`
	Updates        []ReqDictEntryUpdateItem `json:"updates"`
}

type ReqDictEntryUpdateItem struct {
	ID             uint64  `json:"id" binding:"required" binding_msg:"required=请求错误"`
	LabelComponent *string `json:"labelComponent" binding:"omitempty,max=255" binding_msg:"max=显示标签组件最多255位"`
	EntryLabel     *string `json:"entryLabel" binding:"omitempty,max=255" binding_msg:"max=显示标签最多255位"`
	EntryValue     *string `json:"entryValue" binding:"omitempty,max=255" binding_msg:"max=数据值最多255位"`
	LanguageCode   *string `json:"languageCode" binding:"omitempty,max=32" binding_msg:"max=语言代码最多32位"`
	SysDictTypeId  *uint64 `json:"sysDictTypeId"`
	SortOrder      *int32  `json:"sortOrder"`
	IsEnabled      *bool   `json:"isEnabled"`
	Remark         *string `json:"remark" binding:"omitempty,max=255" binding_msg:"max=备注最多255位"`
}

type ReqDictEntryBatchDelete struct {
	IDs []uint64 `json:"ids" binding:"required,min=1" binding_msg:"required=请选择字典项,min=至少选择一项"`
}

type ReqDictEntryBatchCopy struct {
	EntryIds     []uint64 `json:"entryIds" binding:"required,min=1" binding_msg:"required=请选择字典项,min=至少选择一项"`
	TargetTypeId uint64   `json:"targetTypeId" binding:"required" binding_msg:"required=目标字典类型不能为空"`
}

type ReqDictEntryListByCode struct {
	Code string `json:"code" binding:"required,max=128" binding_msg:"required=字典类型编码不能为空,max=字典类型编码最多128位"`
}

type RespDictEntryByCode struct {
	ID             uint64 `json:"id"`
	LabelComponent string `json:"labelComponent"`
	EntryLabel     string `json:"entryLabel"`
	EntryValue     string `json:"entryValue"`
}

// @Summary 通过字典编码获取启用字典项
// @Remark 根据字典类型编码查询启用字典项；若字典项配置了语言条目编码，则按当前请求语言替换显示标签
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body ReqDictEntryListByCode true "字典类型编码"
// @Success 200 {object} res.Response{data=[]RespDictEntryByCode} "成功"
// @Router /api/sys/dict/entry/match [post]
func (*SysDictHandler) EntryMatch(ctx *handler.Ctx, req *ReqDictEntryListByCode) (*[]*RespDictEntryByCode, error) {
	permissionQuery, err := datapermission.BuildReadPermissionQuery(ctx, models.SysDictType{}.TableName())
	if err != nil {
		ctx.L().Error("apply dict type read permission failed", zap.Error(err), zap.String("code", req.Code))
		return nil, res.FailDefault
	}

	sysDictType := permissionQuery.SysDictType
	dictType, err := sysDictType.
		Select(sysDictType.ID).
		Where(sysDictType.TypeCode.Eq(req.Code), sysDictType.IsEnabled.Is(true)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			items := []*RespDictEntryByCode{}
			return &items, nil
		}
		ctx.L().Error("query dict type by code failed", zap.Error(err), zap.String("code", req.Code))
		return nil, res.FailDefault
	}

	sysDictEntry := query.SysDictEntry
	entries, err := sysDictEntry.
		Where(sysDictEntry.SysDictTypeId.Eq(dictType.ID), sysDictEntry.IsEnabled.Is(true)).
		Order(sysDictEntry.SortOrder.Asc(), sysDictEntry.ID.Asc()).
		Find()
	if err != nil {
		ctx.L().Error("query dict entries by code failed", zap.Error(err), zap.String("code", req.Code), zap.Uint64("typeId", dictType.ID))
		return nil, res.FailDefault
	}
	if len(entries) == 0 {
		items := []*RespDictEntryByCode{}
		return &items, nil
	}

	translationMap, err := queryDictEntryTranslationMap(ctx, entries)
	if err != nil {
		return nil, err
	}

	items := make([]*RespDictEntryByCode, 0, len(entries))
	for _, entry := range entries {
		entryLabel := entry.EntryLabel
		if translation, ok := translationMap[strings.TrimSpace(entry.LanguageCode)]; ok {
			entryLabel = translation
		}
		items = append(items, &RespDictEntryByCode{
			ID:             entry.ID,
			LabelComponent: entry.LabelComponent,
			EntryLabel:     entryLabel,
			EntryValue:     entry.EntryValue,
		})
	}
	return &items, nil
}

func queryDictEntryTranslationMap(ctx *handler.Ctx, entries []*models.SysDictEntry) (map[string]string, error) {
	languageCodes := make([]string, 0, len(entries))
	languageCodeSet := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		languageCode := strings.TrimSpace(entry.LanguageCode)
		if languageCode == "" {
			continue
		}
		if _, exists := languageCodeSet[languageCode]; exists {
			continue
		}
		languageCodeSet[languageCode] = struct{}{}
		languageCodes = append(languageCodes, languageCode)
	}
	language := strings.TrimSpace(ctx.Language)
	if len(languageCodes) == 0 || language == "" {
		return map[string]string{}, nil
	}

	sysLanguageType := query.SysLanguageType
	languageType, err := sysLanguageType.
		Select(sysLanguageType.ID).
		Where(sysLanguageType.TypeCode.Eq(language), sysLanguageType.IsEnabled.Is(true)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return map[string]string{}, nil
		}
		ctx.L().Error("query language type failed", zap.Error(err), zap.String("language", language))
		return nil, res.FailDefault
	}

	sysLanguageEntry := query.SysLanguageEntry
	languageEntries, err := sysLanguageEntry.
		Select(sysLanguageEntry.EntryCode, sysLanguageEntry.EntryValue).
		Where(
			sysLanguageEntry.SysLanguageTypeId.Eq(languageType.ID),
			sysLanguageEntry.EntryCode.In(languageCodes...),
			sysLanguageEntry.IsEnabled.Is(true),
		).
		Find()
	if err != nil {
		ctx.L().Error("query language entries failed", zap.Error(err), zap.String("language", language), zap.Strings("entryCodes", languageCodes))
		return nil, res.FailDefault
	}

	translationMap := make(map[string]string, len(languageEntries))
	for _, entry := range languageEntries {
		translationMap[entry.EntryCode] = entry.EntryValue
	}
	return translationMap, nil
}

// @Summary 获取字典数据项分页列表
// @Remark 分页查询字典数据项信息
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body v1.PagingRequest true "分页参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[models.SysDictEntry]} "成功"
// @Router /api/dict/entry/list [post]
func (*SysDictHandler) EntryList(ctx *handler.Ctx, req *v1.PagingRequest) (*gormc.PagingResult[models.SysDictEntry], error) {
	readableTypeIDs, err := queryAllowedDictTypeIDs(ctx, datapermission.BuildReadPermissionQuery)
	if err != nil {
		ctx.L().Error("apply dict type read permission failed", zap.Error(err))
		return nil, res.FailDefault
	}
	if len(readableTypeIDs) == 0 {
		return &gormc.PagingResult[models.SysDictEntry]{
			Items: []*models.SysDictEntry{},
			Total: 0,
		}, nil
	}
	if err := applyDictEntryTypeIDPageFilter(req, readableTypeIDs); err != nil {
		ctx.L().Error("apply dict entry type read permission failed", zap.Error(err), zap.Uint64s("typeIds", readableTypeIDs))
		return nil, res.FailDefault
	}

	pagination, err := query.SysDictEntry.PageWithPaging(req)
	if err != nil {
		return nil, res.FailDefault
	}
	if len(pagination.Items) == 0 {
		return pagination, nil
	}

	typeIDSet := make(map[uint64]struct{}, len(pagination.Items))
	typeIDs := make([]uint64, 0, len(pagination.Items))
	for _, item := range pagination.Items {
		if _, exists := typeIDSet[item.SysDictTypeId]; exists {
			continue
		}
		typeIDSet[item.SysDictTypeId] = struct{}{}
		typeIDs = append(typeIDs, item.SysDictTypeId)
	}
	if len(typeIDs) == 0 {
		return pagination, nil
	}

	sysDictType := query.SysDictType
	typeList, err := sysDictType.
		Select(sysDictType.ID, sysDictType.TypeCode, sysDictType.TypeName).
		Where(sysDictType.ID.In(typeIDs...)).
		Find()
	if err != nil {
		ctx.L().Error("查询字典类型失败", zap.Error(err), zap.Uint64s("typeIds", typeIDs))
		return nil, res.FailDefault
	}
	typeMap := make(map[uint64]*models.SysDictType, len(typeList))
	for _, item := range typeList {
		typeMap[item.ID] = item
	}
	for _, item := range pagination.Items {
		item.SysDictType = typeMap[item.SysDictTypeId]
	}
	return pagination, nil
}

// @Summary 创建字典数据项
// @Remark 创建新的字典数据项
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body ReqDictEntryCreate true "字典数据项创建参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/dict/entry/create [post]
func (*SysDictHandler) EntryCreate(ctx *handler.Ctx, req *ReqDictEntryCreate) error {
	operationID := ctx.SessionInfo.Id
	permissionQuery, err := datapermission.BuildWritePermissionQuery(ctx, models.SysDictType{}.TableName())
	if err != nil {
		ctx.L().Error("apply dict type write permission failed", zap.Error(err), zap.Uint64("typeId", req.SysDictTypeId))
		return res.FailDefault
	}

	sysDictType := permissionQuery.SysDictType
	_, err = sysDictType.
		Select(sysDictType.ID).
		Where(sysDictType.ID.Eq(req.SysDictTypeId)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res.FailMsg("无权限或字典类型不存在")
		}
		return res.FailDefault
	}

	err = query.SysDictEntry.Create(&models.SysDictEntry{
		OperatorID: mixin.OperatorID{
			CreatedBy: mixin.CreatedBy{CreatedBy: operationID},
			UpdatedBy: mixin.UpdatedBy{UpdatedBy: operationID},
		},
		SortOrder:      mixin.SortOrder{SortOrder: req.SortOrder},
		IsEnabled:      mixin.IsEnabled{IsEnabled: req.IsEnabled},
		Remark:         mixin.Remark{Remark: req.Remark},
		LabelComponent: req.LabelComponent,
		EntryLabel:     req.EntryLabel,
		EntryValue:     req.EntryValue,
		LanguageCode:   req.LanguageCode,
		SysDictTypeId:  req.SysDictTypeId,
	})
	if err != nil {
		return res.FailDefault
	}
	return nil
}

// @Summary 更新字典数据项
// @Remark 根据 ID 更新字典数据项信息
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body ReqDictEntryUpdate true "字典数据项更新参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/dict/entry/update [post]
func (*SysDictHandler) EntryUpdate(ctx *handler.Ctx, req *ReqDictEntryUpdate) error {
	operationID := ctx.SessionInfo.Id
	typeWriteQuery, err := datapermission.BuildWritePermissionQuery(ctx, models.SysDictType{}.TableName())
	if err != nil {
		ctx.L().Error("apply dict type write permission failed", zap.Error(err))
		return res.FailDefault
	}

	if len(req.Updates) > 0 {
		for _, item := range req.Updates {
			if err := updateDictEntry(operationID, typeWriteQuery, &item); err != nil {
				return err
			}
		}
		return nil
	}
	if req.ID == nil {
		return res.FailMsg("请求错误")
	}
	return updateDictEntry(operationID, typeWriteQuery, &ReqDictEntryUpdateItem{
		ID:             *req.ID,
		LabelComponent: req.LabelComponent,
		EntryLabel:     req.EntryLabel,
		EntryValue:     req.EntryValue,
		LanguageCode:   req.LanguageCode,
		SysDictTypeId:  req.SysDictTypeId,
		SortOrder:      req.SortOrder,
		IsEnabled:      req.IsEnabled,
		Remark:         req.Remark,
	})
}

func updateDictEntry(operationID uint64, typeWriteQuery *query.Query, req *ReqDictEntryUpdateItem) error {
	entry, err := query.SysDictEntry.
		Select(query.SysDictEntry.ID, query.SysDictEntry.SysDictTypeId).
		Where(query.SysDictEntry.ID.Eq(req.ID)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res.FailMsg("无权限或数据不存在")
		}
		return res.FailDefault
	}

	if err := ensureDictTypeAllowed(typeWriteQuery, entry.SysDictTypeId); err != nil {
		return err
	}
	if req.SysDictTypeId != nil && *req.SysDictTypeId != entry.SysDictTypeId {
		if err := ensureDictTypeAllowed(typeWriteQuery, *req.SysDictTypeId); err != nil {
			return err
		}
	}

	sysDictEntry := query.SysDictEntry
	exprs := []field.AssignExpr{sysDictEntry.UpdatedBy.Value(operationID)}
	query.ExprAppendSelf(&exprs, req.LabelComponent, sysDictEntry.LabelComponent.Value)
	query.ExprAppendSelf(&exprs, req.EntryLabel, sysDictEntry.EntryLabel.Value)
	query.ExprAppendSelf(&exprs, req.EntryValue, sysDictEntry.EntryValue.Value)
	query.ExprAppendSelf(&exprs, req.LanguageCode, sysDictEntry.LanguageCode.Value)
	query.ExprAppendSelf(&exprs, req.SysDictTypeId, sysDictEntry.SysDictTypeId.Value)
	query.ExprAppendSelf(&exprs, req.SortOrder, sysDictEntry.SortOrder.Value)
	query.ExprAppendSelf(&exprs, req.IsEnabled, sysDictEntry.IsEnabled.Value)
	query.ExprAppendSelf(&exprs, req.Remark, sysDictEntry.Remark.Value)

	info, err := sysDictEntry.Where(sysDictEntry.ID.Eq(req.ID)).UpdateSimple(exprs...)
	if err != nil {
		return res.FailDefault
	}
	if info.RowsAffected == 0 {
		return res.FailMsg("无权限或数据不存在")
	}
	return nil
}

func ensureDictTypeAllowed(permissionQuery *query.Query, typeID uint64) error {
	sysDictType := permissionQuery.SysDictType
	_, err := sysDictType.
		Select(sysDictType.ID).
		Where(sysDictType.ID.Eq(typeID)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res.FailMsg("无权限或字典类型不存在")
		}
		return res.FailDefault
	}
	return nil
}

// @Summary 批量删除字典数据项
// @Remark 根据 ID 列表批量删除字典数据项
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body ReqDictEntryBatchDelete true "批量删除参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/dict/entry/del [post]
func (*SysDictHandler) EntryDel(ctx *handler.Ctx, req *ReqDictEntryBatchDelete) error {
	ids := slices_utils.Distinct(req.IDs)
	sysDictEntry := query.SysDictEntry
	entries, err := sysDictEntry.
		Select(sysDictEntry.ID, sysDictEntry.SysDictTypeId).
		Where(sysDictEntry.ID.In(ids...)).
		Find()
	if err != nil {
		ctx.L().Error("查询字典项删除范围失败", zap.Error(err), zap.Uint64s("ids", ids))
		return res.FailDefault
	}
	if len(entries) != len(ids) {
		return res.FailMsg("无权限或数据不存在")
	}

	typeIDSet := make(map[uint64]struct{}, len(entries))
	typeIDs := make([]uint64, 0, len(entries))
	for _, entry := range entries {
		if _, exists := typeIDSet[entry.SysDictTypeId]; exists {
			continue
		}
		typeIDSet[entry.SysDictTypeId] = struct{}{}
		typeIDs = append(typeIDs, entry.SysDictTypeId)
	}
	deleteIDSet, err := queryAllowedDictTypeIDSet(ctx, typeIDs, datapermission.BuildDeletePermissionQuery)
	if err != nil {
		ctx.L().Error("apply dict type delete permission failed", zap.Error(err), zap.Uint64s("typeIds", typeIDs))
		return res.FailDefault
	}
	if len(deleteIDSet) != len(typeIDs) {
		return res.FailMsg("无权限或数据不存在")
	}

	info, err := sysDictEntry.Where(sysDictEntry.ID.In(ids...)).Delete()
	if err != nil {
		return res.FailDefault
	}
	if info.RowsAffected != int64(len(ids)) {
		return res.FailMsg("无权限或数据不存在")
	}
	return nil
}

// @Summary 批量复制字典数据项
// @Remark 将选中的字典数据项批量复制到指定字典类型下（不支持复制到同一类型）
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body ReqDictEntryBatchCopy true "批量复制参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/dict/entry/batch/copy [post]
func (*SysDictHandler) EntryBatchCopy(ctx *handler.Ctx, req *ReqDictEntryBatchCopy) error {
	writePermissionQuery, err := datapermission.BuildWritePermissionQuery(ctx, models.SysDictType{}.TableName())
	if err != nil {
		ctx.L().Error("apply dict type write permission failed", zap.Error(err), zap.Uint64("targetTypeId", req.TargetTypeId))
		return res.FailDefault
	}

	if err := ensureDictTypeAllowed(writePermissionQuery, req.TargetTypeId); err != nil {
		return err
	}

	entryIDs := slices_utils.Distinct(req.EntryIds)
	sourceEntries, err := query.SysDictEntry.
		Where(query.SysDictEntry.ID.In(entryIDs...)).
		Find()
	if err != nil {
		ctx.L().Error("查询源字典项失败", zap.Error(err), zap.Uint64s("entryIds", entryIDs))
		return res.FailDefault
	}
	if len(sourceEntries) != len(entryIDs) {
		return res.FailMsg("无权限或源字典项不存在")
	}

	sourceTypeIDSet := make(map[uint64]struct{}, len(sourceEntries))
	sourceTypeIDs := make([]uint64, 0, len(sourceEntries))
	for _, entry := range sourceEntries {
		if _, exists := sourceTypeIDSet[entry.SysDictTypeId]; exists {
			continue
		}
		sourceTypeIDSet[entry.SysDictTypeId] = struct{}{}
		sourceTypeIDs = append(sourceTypeIDs, entry.SysDictTypeId)
	}
	readIDSet, err := queryAllowedDictTypeIDSet(ctx, sourceTypeIDs, datapermission.BuildReadPermissionQuery)
	if err != nil {
		ctx.L().Error("apply source dict type read permission failed", zap.Error(err), zap.Uint64s("typeIds", sourceTypeIDs))
		return res.FailDefault
	}
	if len(readIDSet) != len(sourceTypeIDs) {
		return res.FailMsg("无权限或源字典项不存在")
	}

	var newEntries []*models.SysDictEntry
	for _, entry := range sourceEntries {
		newEntries = append(newEntries, &models.SysDictEntry{
			LabelComponent: entry.LabelComponent,
			EntryLabel:     entry.EntryLabel,
			EntryValue:     entry.EntryValue,
			LanguageCode:   entry.LanguageCode,
			SysDictTypeId:  req.TargetTypeId,
			SortOrder:      mixin.SortOrder{SortOrder: entry.SortOrder.SortOrder},
			IsEnabled:      mixin.IsEnabled{IsEnabled: entry.IsEnabled.IsEnabled},
			Remark:         mixin.Remark{Remark: entry.Remark.Remark},
		})
	}

	err = query.SysDictEntry.Create(newEntries...)
	if err != nil {
		ctx.L().Error("批量复制字典项失败", zap.Error(err), zap.Uint64s("entryIds", entryIDs), zap.Uint64("targetTypeId", req.TargetTypeId))
		return res.FailDefault
	}

	return nil
}
