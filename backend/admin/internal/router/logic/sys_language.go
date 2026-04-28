package logic

import (
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/res"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"errors"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc"
	"orm-crud/gormc/mixin"

	"gorm.io/gen/field"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SysLanguageHandler struct{}

// --- 语言类型 (LanguageType) ---

type ReqLangTypeCreate struct {
	TypeCode  string `json:"typeCode" binding:"required,max=128" binding_msg:"required=语言编码不能为空,max=语言编码最多128位"`
	TypeName  string `json:"typeName" binding:"required,max=255" binding_msg:"required=语言名称不能为空,max=语言名称最多255位"`
	IsDefault bool   `json:"isDefault"`
	IsEnabled bool   `json:"isEnabled"`
	SortOrder int32  `json:"sortOrder"`
}

type ReqLangTypeUpdate struct {
	ID        uint64  `json:"id" binding:"required" binding_msg:"required=请求错误"`
	TypeCode  *string `json:"typeCode" binding:"omitempty,max=128" binding_msg:"max=语言编码最多128位"`
	TypeName  *string `json:"typeName" binding:"omitempty,max=255" binding_msg:"max=语言名称最多255位"`
	IsDefault *bool   `json:"isDefault"`
	IsEnabled *bool   `json:"isEnabled"`
	SortOrder *int32  `json:"sortOrder"`
}

type ReqLangTypeBatchDelete struct {
	IDs []uint64 `json:"ids" binding:"required,min=1" binding_msg:"required=请选择语言类型,min=至少选择一项"`
}

// @Summary 获取语言类型分页列表
// @Description 分页查询语言类型信息
// @Tags Language
// @Accept json
// @Produce json
// @Param req body v1.PagingRequest true "分页参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[models.SysLanguageType]} "成功"
// @Router /api/sys/language/type/list [post]
func (*SysLanguageHandler) TypeList(ctx *handler.Ctx, req *v1.PagingRequest) (*gormc.PagingResult[models.SysLanguageType], error) {
	pagination, err := query.SysLanguageType.PageWithPaging(req)
	if err != nil {
		return nil, res.FailDefault
	}
	return pagination, nil
}

// @Summary 创建语言类型
// @Description 创建新的语言类型
// @Tags Language
// @Accept json
// @Produce json
// @Param req body ReqLangTypeCreate true "语言类型创建参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/language/type/create [post]
func (*SysLanguageHandler) TypeCreate(ctx *handler.Ctx, req *ReqLangTypeCreate) error {
	sysLanguageType := query.SysLanguageType

	if req.IsDefault {
		_, err := sysLanguageType.Where(sysLanguageType.IsDefault.Is(true)).Update(sysLanguageType.IsDefault, false)
		if err != nil {
			ctx.L().Error("重置默认语言失败", zap.Error(err))
			return res.FailDefault
		}
	}

	operationID := ctx.SessionInfo.Id

	err := sysLanguageType.Create(&models.SysLanguageType{
		OperatorID: mixin.OperatorID{
			CreatedBy: mixin.CreatedBy{CreatedBy: operationID},
			UpdatedBy: mixin.UpdatedBy{UpdatedBy: operationID},
		},
		IsEnabled: mixin.IsEnabled{IsEnabled: req.IsEnabled},
		SortOrder: mixin.SortOrder{SortOrder: req.SortOrder},
		TypeCode:  req.TypeCode,
		TypeName:  req.TypeName,
		IsDefault: req.IsDefault,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return res.FailMsg("语言编码已存在")
		}
		return res.FailDefault
	}
	return nil
}

// @Summary 更新语言类型
// @Description 根据 ID 更新语言类型信息
// @Tags Language
// @Accept json
// @Produce json
// @Param req body ReqLangTypeUpdate true "语言类型更新参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/language/type/update [post]
func (*SysLanguageHandler) TypeUpdate(ctx *handler.Ctx, req *ReqLangTypeUpdate) error {
	sysLanguageType := query.SysLanguageType
	if req.IsEnabled != nil && !*req.IsEnabled {
		target, err := sysLanguageType.
			Select(sysLanguageType.IsDefault).
			Where(sysLanguageType.ID.Eq(req.ID)).
			First()
		if err == nil && target.IsDefault {
			return res.FailMsg("默认语言不能停用")
		}
	}

	if req.IsDefault != nil && *req.IsDefault {
		_, err := sysLanguageType.
			Where(sysLanguageType.ID.Neq(req.ID), sysLanguageType.IsDefault.Is(true)).
			Update(sysLanguageType.IsDefault, false)
		if err != nil {
			ctx.L().Error("重置默认语言失败", zap.Error(err))
			return res.FailDefault
		}
	}

	operationID := ctx.SessionInfo.Id
	exprs := []field.AssignExpr{sysLanguageType.UpdatedBy.Value(operationID)}
	query.ExprAppendSelf(&exprs, req.TypeCode, sysLanguageType.TypeCode.Value)
	query.ExprAppendSelf(&exprs, req.TypeName, sysLanguageType.TypeName.Value)
	query.ExprAppendSelf(&exprs, req.IsDefault, sysLanguageType.IsDefault.Value)
	query.ExprAppendSelf(&exprs, req.IsEnabled, sysLanguageType.IsEnabled.Value)
	query.ExprAppendSelf(&exprs, req.SortOrder, sysLanguageType.SortOrder.Value)

	_, err := sysLanguageType.Where(sysLanguageType.ID.Eq(req.ID)).UpdateSimple(exprs...)

	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return res.FailMsg("语言编码已存在")
		}
		return res.FailDefault
	}
	return nil
}

// @Summary 批量删除语言类型
// @Description 根据 ID 列表批量删除语言类型及其关联的所有语言条目
// @Tags Language
// @Accept json
// @Produce json
// @Param req body ReqLangTypeBatchDelete true "批量删除参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/language/type/del [post]
func (*SysLanguageHandler) TypeDel(ctx *handler.Ctx, req *ReqLangTypeBatchDelete) error {
	sysLanguageType := query.SysLanguageType
	defaultLang, err := sysLanguageType.
		Select(sysLanguageType.ID).
		Where(sysLanguageType.ID.In(req.IDs...), sysLanguageType.IsDefault.Is(true)).
		First()
	if err == nil && defaultLang != nil {
		return res.FailMsg("默认语言不能删除")
	}

	err = query.Q.Transaction(func(tx *query.Query) error {
		sysLanguageEntry := tx.SysLanguageEntry
		_, err = sysLanguageEntry.
			Where(sysLanguageEntry.SysLanguageTypeId.In(req.IDs...)).
			Delete()
		if err != nil {
			return err
		}
		sysLanguageTypeSub := tx.SysLanguageType
		_, err = sysLanguageTypeSub.
			Where(sysLanguageTypeSub.ID.In(req.IDs...)).
			Delete()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		ctx.L().Error("批量删除语言类型失败", zap.Error(err), zap.Uint64s("ids", req.IDs))
		return res.FailDefault
	}
	return nil
}

// --- 语言条目 (LanguageEntry) ---

type ReqLangEntryCreate struct {
	EntryCode         string `json:"entryCode" binding:"required,max=128" binding_msg:"required=条目编码不能为空,max=条目编码最多128位"`
	EntryValue        string `json:"entryValue" binding:"required,max=255" binding_msg:"required=语言值不能为空,max=语言值最多255位"`
	SysLanguageTypeId uint64 `json:"sysLanguageTypeId" binding:"required" binding_msg:"required=语言类型ID不能为空"`
	SortOrder         int32  `json:"sortOrder"`
	IsEnabled         bool   `json:"isEnabled"`
	Remark            string `json:"remark" binding:"max=255" binding_msg:"max=备注最多255位"`
}

type ReqLangEntryUpdate struct {
	ID                *uint64                  `json:"id"`
	EntryCode         *string                  `json:"entryCode" binding:"omitempty,max=128" binding_msg:"max=条目编码最多128位"`
	EntryValue        *string                  `json:"entryValue" binding:"omitempty,max=255" binding_msg:"max=语言值最多255位"`
	SysLanguageTypeId *uint64                  `json:"sysLanguageTypeId"`
	SortOrder         *int32                   `json:"sortOrder"`
	IsEnabled         *bool                    `json:"isEnabled"`
	Remark            *string                  `json:"remark" binding:"omitempty,max=255" binding_msg:"max=备注最多255位"`
	Updates           []ReqLangEntryUpdateItem `json:"updates"`
}

type ReqLangEntryUpdateItem struct {
	ID                uint64  `json:"id" binding:"required" binding_msg:"required=请求错误"`
	EntryCode         *string `json:"entryCode" binding:"omitempty,max=128" binding_msg:"max=条目编码最多128位"`
	EntryValue        *string `json:"entryValue" binding:"omitempty,max=255" binding_msg:"max=语言值最多255位"`
	SysLanguageTypeId *uint64 `json:"sysLanguageTypeId"`
	SortOrder         *int32  `json:"sortOrder"`
	IsEnabled         *bool   `json:"isEnabled"`
	Remark            *string `json:"remark" binding:"omitempty,max=255" binding_msg:"max=备注最多255位"`
}

type ReqLangEntryBatchDelete struct {
	IDs []uint64 `json:"ids" binding:"required,min=1" binding_msg:"required=请选择语言条目,min=至少选择一项"`
}

type ReqLangEntryBatchCreate struct {
	EntryCode string            `json:"entryCode" binding:"required,max=128" binding_msg:"required=条目编码不能为空,max=条目编码最多128位"`
	Values    map[string]string `json:"values" binding:"required,min=1" binding_msg:"required=语言值不能为空"`
	SortOrder int32             `json:"sortOrder"`
	IsEnabled bool              `json:"isEnabled"`
}

// @Summary 获取语言条目分页列表
// @Description 分页查询语言条目信息
// @Tags Language
// @Accept json
// @Produce json
// @Param req body v1.PagingRequest true "分页参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[models.SysLanguageEntry]} "成功"
// @Router /api/sys/language/entry/list [post]
func (*SysLanguageHandler) EntryList(ctx *handler.Ctx, req *v1.PagingRequest) (*gormc.PagingResult[models.SysLanguageEntry], error) {
	pagination, err := query.SysLanguageEntry.PageWithPaging(req)
	if err != nil {
		return nil, res.FailDefault
	}
	return pagination, nil
}

// @Summary 创建语言条目
// @Description 创建新的语言条目
// @Tags Language
// @Accept json
// @Produce json
// @Param req body ReqLangEntryCreate true "语言条目创建参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/language/entry/create [post]
func (*SysLanguageHandler) EntryCreate(ctx *handler.Ctx, req *ReqLangEntryCreate) error {
	operationID := ctx.SessionInfo.Id
	sysLanguageType := query.SysLanguageType
	_, err := sysLanguageType.
		Select(sysLanguageType.ID).
		Where(sysLanguageType.ID.Eq(req.SysLanguageTypeId)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res.FailMsg("语言类型不存在")
		}
		return res.FailDefault
	}
	err = query.SysLanguageEntry.Create(&models.SysLanguageEntry{
		OperatorID: mixin.OperatorID{
			CreatedBy: mixin.CreatedBy{CreatedBy: operationID},
			UpdatedBy: mixin.UpdatedBy{UpdatedBy: operationID},
		},
		SortOrder:         mixin.SortOrder{SortOrder: req.SortOrder},
		IsEnabled:         mixin.IsEnabled{IsEnabled: req.IsEnabled},
		Remark:            mixin.Remark{Remark: req.Remark},
		EntryCode:         req.EntryCode,
		EntryValue:        req.EntryValue,
		SysLanguageTypeId: req.SysLanguageTypeId,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return res.FailMsg("该语言下已存在相同编码的条目")
		}
		return res.FailDefault
	}
	return nil
}

// @Summary 更新语言条目
// @Description 根据 ID 更新语言条目信息
// @Tags Language
// @Accept json
// @Produce json
// @Param req body ReqLangEntryUpdate true "语言条目更新参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/language/entry/update [post]
func (*SysLanguageHandler) EntryUpdate(ctx *handler.Ctx, req *ReqLangEntryUpdate) error {
	operationID := ctx.SessionInfo.Id
	if len(req.Updates) > 0 {
		for _, item := range req.Updates {
			if err := updateLanguageEntry(operationID, &item); err != nil {
				return err
			}
		}
		return nil
	}
	if req.ID == nil {
		return res.FailMsg("请求错误")
	}
	return updateLanguageEntry(operationID, &ReqLangEntryUpdateItem{
		ID:                *req.ID,
		EntryCode:         req.EntryCode,
		EntryValue:        req.EntryValue,
		SysLanguageTypeId: req.SysLanguageTypeId,
		SortOrder:         req.SortOrder,
		IsEnabled:         req.IsEnabled,
		Remark:            req.Remark,
	})
}

func updateLanguageEntry(operationID uint64, req *ReqLangEntryUpdateItem) error {
	if req.SysLanguageTypeId != nil {
		sysLanguageType := query.SysLanguageType
		_, err := sysLanguageType.
			Select(sysLanguageType.ID).
			Where(sysLanguageType.ID.Eq(*req.SysLanguageTypeId)).
			First()

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return res.FailMsg("语言类型不存在")
			}
			return res.FailDefault
		}
	}
	sysLanguageEntry := query.SysLanguageEntry
	exprs := []field.AssignExpr{sysLanguageEntry.UpdatedBy.Value(operationID)}
	query.ExprAppendSelf(&exprs, req.EntryCode, sysLanguageEntry.EntryCode.Value)
	query.ExprAppendSelf(&exprs, req.EntryValue, sysLanguageEntry.EntryValue.Value)
	query.ExprAppendSelf(&exprs, req.SysLanguageTypeId, sysLanguageEntry.SysLanguageTypeId.Value)
	query.ExprAppendSelf(&exprs, req.SortOrder, sysLanguageEntry.SortOrder.Value)
	query.ExprAppendSelf(&exprs, req.IsEnabled, sysLanguageEntry.IsEnabled.Value)
	query.ExprAppendSelf(&exprs, req.Remark, sysLanguageEntry.Remark.Value)
	_, err := sysLanguageEntry.
		Where(sysLanguageEntry.ID.Eq(req.ID)).
		UpdateSimple(exprs...)
	if err != nil {
		return res.FailDefault
	}
	return nil
}

// @Summary 批量删除语言条目
// @Description 根据 ID 列表批量删除语言条目
// @Tags Language
// @Accept json
// @Produce json
// @Param req body ReqLangEntryBatchDelete true "批量删除参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/language/entry/del [post]
func (*SysLanguageHandler) EntryDel(ctx *handler.Ctx, req *ReqLangEntryBatchDelete) error {
	sysLanguageEntry := query.SysLanguageEntry
	_, err := sysLanguageEntry.Where(sysLanguageEntry.ID.In(req.IDs...)).Delete()
	if err != nil {
		return res.FailDefault
	}
	return nil
}

// @Summary 批量创建语言条目
// @Description 根据条目编码和语言值映射批量创建语言条目，已存在的则更新
// @Tags Language
// @Accept json
// @Produce json
// @Param req body ReqLangEntryBatchCreate true "批量创建参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/language/entry/batch/create [post]
func (*SysLanguageHandler) EntryBatchCreate(ctx *handler.Ctx, req *ReqLangEntryBatchCreate) error {
	operationID := ctx.SessionInfo.Id

	var typeCodes []string
	for tc := range req.Values {
		typeCodes = append(typeCodes, tc)
	}

	sysLanguageType := query.SysLanguageType
	typeList, err := sysLanguageType.
		Select(sysLanguageType.ID, sysLanguageType.TypeCode).
		Where(sysLanguageType.TypeCode.In(typeCodes...)).
		Find()
	if err != nil {
		ctx.L().Error("查询语言类型失败", zap.Error(err))
		return res.FailDefault
	}

	typeCodeToID := make(map[string]uint64, len(typeList))
	var typeIDs []uint64
	for _, t := range typeList {
		typeCodeToID[t.TypeCode] = t.ID
		typeIDs = append(typeIDs, t.ID)
	}

	existingEntries, err := query.SysLanguageEntry.
		Where(
			query.SysLanguageEntry.EntryCode.Eq(req.EntryCode),
			query.SysLanguageEntry.SysLanguageTypeId.In(typeIDs...),
		).Find()
	if err != nil {
		ctx.L().Error("查询已有语言条目失败", zap.Error(err))
		return res.FailDefault
	}

	existingMap := make(map[uint64]*models.SysLanguageEntry, len(existingEntries))
	for _, e := range existingEntries {
		existingMap[e.SysLanguageTypeId] = e
	}

	var createEntries []*models.SysLanguageEntry
	var updateEntries []*models.SysLanguageEntry
	for typeCode, entryValue := range req.Values {
		typeID, ok := typeCodeToID[typeCode]
		if !ok {
			continue
		}
		if existing, ok := existingMap[typeID]; ok {
			existing.EntryValue = entryValue
			existing.UpdatedBy = mixin.UpdatedBy{UpdatedBy: operationID}
			updateEntries = append(updateEntries, existing)
		} else {
			createEntries = append(createEntries, &models.SysLanguageEntry{
				OperatorID: mixin.OperatorID{
					CreatedBy: mixin.CreatedBy{CreatedBy: operationID},
					UpdatedBy: mixin.UpdatedBy{UpdatedBy: operationID},
				},
				SortOrder:         mixin.SortOrder{SortOrder: req.SortOrder},
				IsEnabled:         mixin.IsEnabled{IsEnabled: req.IsEnabled},
				EntryCode:         req.EntryCode,
				EntryValue:        entryValue,
				SysLanguageTypeId: typeID,
			})
		}
	}

	if len(createEntries) > 0 {
		err = query.SysLanguageEntry.Create(createEntries...)
		if err != nil {
			ctx.L().Error("批量创建语言条目失败", zap.Error(err))
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return res.FailMsg("语言条目有重复")
			}
			return res.FailDefault
		}
	}

	if len(updateEntries) > 0 {
		err = query.SysLanguageEntry.Save(updateEntries...)
		if err != nil {
			ctx.L().Error("批量更新语言条目失败", zap.Error(err))
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return res.FailMsg("语言条目有重复")
			}
			return res.FailDefault
		}
	}

	return nil
}
