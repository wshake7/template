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

// @Summary 获取字典类型分页列表
// @Remark 分页查询字典类型信息
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body v1.PagingRequest true "分页参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[models.SysDictType]} "成功"
// @Router /api/dict/type/list [post]
func (*SysDictHandler) TypeList(ctx *handler.Ctx, req *v1.PagingRequest) (*gormc.PagingResult[models.SysDictType], error) {
	pagination, err := query.SysDictType.PageWithPaging(req)
	if err != nil {
		return nil, res.FailDefault
	}
	return pagination, nil
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
		CreatedBy: mixin.CreatedBy{CreatedBy: operationID},
		UpdatedBy: mixin.UpdatedBy{UpdatedBy: operationID},
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
	sysDictType := query.SysDictType

	exprs := []field.AssignExpr{sysDictType.UpdatedBy.Value(operationID)}
	query.ExprAppendSelf(&exprs, req.TypeCode, sysDictType.TypeCode.Value)
	query.ExprAppendSelf(&exprs, req.TypeName, sysDictType.TypeName.Value)
	query.ExprAppendSelf(&exprs, req.IsEnabled, sysDictType.IsEnabled.Value)
	query.ExprAppendSelf(&exprs, req.SortOrder, sysDictType.SortOrder.Value)
	query.ExprAppendSelf(&exprs, req.Remark, sysDictType.Remark.Value)

	_, err := sysDictType.Where(sysDictType.ID.Eq(req.ID)).UpdateSimple(exprs...)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return res.FailMsg("类型编码已存在")
		}
		return res.FailDefault
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
	err := query.Q.Transaction(func(tx *query.Query) error {
		sysDictEntry := tx.SysDictEntry
		_, err := sysDictEntry.
			Where(sysDictEntry.SysDictTypeId.In(req.IDs...)).
			Delete()
		if err != nil {
			return err
		}
		sysDictTypeSub := tx.SysDictType
		_, err = sysDictTypeSub.
			Where(sysDictTypeSub.ID.In(req.IDs...)).
			Delete()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		ctx.L().Error("批量删除字典类型失败", zap.Error(err), zap.Uint64s("ids", req.IDs))
		return res.FailDefault
	}
	return nil
}

// --- 字典数据项 (DictEntry) ---

type ReqDictEntryCreate struct {
	EntryLabel    string `json:"entryLabel" binding:"required,max=255" binding_msg:"required=显示标签不能为空,max=显示标签最多255位"`
	EntryValue    string `json:"entryValue" binding:"required,max=255" binding_msg:"required=数据值不能为空,max=数据值最多255位"`
	NumericValue  int32  `json:"numericValue"`
	LanguageCode  string `json:"languageCode" binding:"max=32" binding_msg:"max=语言代码最多32位"`
	SysDictTypeId uint64 `json:"sysDictTypeId" binding:"required" binding_msg:"required=字典类型ID不能为空"`
	SortOrder     int32  `json:"sortOrder"`
	IsEnabled     bool   `json:"isEnabled"`
	Remark        string `json:"remark" binding:"max=255" binding_msg:"max=备注最多255位"`
}

type ReqDictEntryUpdate struct {
	ID            uint64  `json:"id" binding:"required" binding_msg:"required=请求错误"`
	EntryLabel    *string `json:"entryLabel" binding:"omitempty,max=255" binding_msg:"max=显示标签最多255位"`
	EntryValue    *string `json:"entryValue" binding:"omitempty,max=255" binding_msg:"max=数据值最多255位"`
	NumericValue  *int32  `json:"numericValue"`
	LanguageCode  *string `json:"languageCode" binding:"omitempty,max=32" binding_msg:"max=语言代码最多32位"`
	SysDictTypeId *uint64 `json:"sysDictTypeId"`
	SortOrder     *int32  `json:"sortOrder"`
	IsEnabled     *bool   `json:"isEnabled"`
	Remark        *string `json:"remark" binding:"omitempty,max=255" binding_msg:"max=备注最多255位"`
}

type ReqDictEntryBatchDelete struct {
	IDs []uint64 `json:"ids" binding:"required,min=1" binding_msg:"required=请选择字典项,min=至少选择一项"`
}

type ReqDictEntryBatchCopy struct {
	EntryIds     []uint64 `json:"entryIds" binding:"required,min=1" binding_msg:"required=请选择字典项,min=至少选择一项"`
	TargetTypeId uint64   `json:"targetTypeId" binding:"required" binding_msg:"required=目标字典类型不能为空"`
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
	pagination, err := query.SysDictEntry.PageWithPaging(req)
	if err != nil {
		return nil, res.FailDefault
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
	sysDictType := query.SysDictType
	_, err := sysDictType.
		Select(sysDictType.ID).
		Where(sysDictType.ID.Eq(req.SysDictTypeId)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res.FailMsg("字典类型不存在")
		}
		return res.FailDefault
	}

	err = query.SysDictEntry.Create(&models.SysDictEntry{
		CreatedBy:     mixin.CreatedBy{CreatedBy: operationID},
		UpdatedBy:     mixin.UpdatedBy{UpdatedBy: operationID},
		SortOrder:     mixin.SortOrder{SortOrder: req.SortOrder},
		IsEnabled:     mixin.IsEnabled{IsEnabled: req.IsEnabled},
		Remark:        mixin.Remark{Remark: req.Remark},
		EntryLabel:    req.EntryLabel,
		EntryValue:    req.EntryValue,
		NumericValue:  req.NumericValue,
		LanguageCode:  req.LanguageCode,
		SysDictTypeId: req.SysDictTypeId,
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
	if req.SysDictTypeId != nil {
		sysDictType := query.SysDictType
		_, err := sysDictType.
			Select(sysDictType.ID).
			Where(sysDictType.ID.Eq(*req.SysDictTypeId)).
			First()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return res.FailMsg("字典类型不存在")
			}
			return res.FailDefault
		}
	}

	sysDictEntry := query.SysDictEntry
	exprs := []field.AssignExpr{sysDictEntry.UpdatedBy.Value(operationID)}
	query.ExprAppendSelf(&exprs, req.EntryLabel, sysDictEntry.EntryLabel.Value)
	query.ExprAppendSelf(&exprs, req.EntryValue, sysDictEntry.EntryValue.Value)
	query.ExprAppendSelf(&exprs, req.NumericValue, sysDictEntry.NumericValue.Value)
	query.ExprAppendSelf(&exprs, req.LanguageCode, sysDictEntry.LanguageCode.Value)
	query.ExprAppendSelf(&exprs, req.SysDictTypeId, sysDictEntry.SysDictTypeId.Value)
	query.ExprAppendSelf(&exprs, req.SortOrder, sysDictEntry.SortOrder.Value)
	query.ExprAppendSelf(&exprs, req.IsEnabled, sysDictEntry.IsEnabled.Value)
	query.ExprAppendSelf(&exprs, req.Remark, sysDictEntry.Remark.Value)

	_, err := sysDictEntry.Where(sysDictEntry.ID.Eq(req.ID)).UpdateSimple(exprs...)
	if err != nil {
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
	sysDictEntry := query.SysDictEntry
	_, err := sysDictEntry.Where(sysDictEntry.ID.In(req.IDs...)).Delete()
	if err != nil {
		return res.FailDefault
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
	sysDictType := query.SysDictType
	_, err := sysDictType.
		Select(sysDictType.ID).
		Where(sysDictType.ID.Eq(req.TargetTypeId)).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return res.FailMsg("目标字典类型不存在")
		}
		ctx.L().Error("校验字典类型失败", zap.Error(err), zap.Uint64("targetTypeId", req.TargetTypeId))
		return res.FailDefault
	}

	sourceEntries, err := query.SysDictEntry.
		Where(query.SysDictEntry.ID.In(req.EntryIds...)).
		Find()
	if err != nil {
		ctx.L().Error("查询源字典项失败", zap.Error(err), zap.Uint64s("entryIds", req.EntryIds))
		return res.FailDefault
	}
	if len(sourceEntries) == 0 {
		return res.FailMsg("未找到要复制的字典项")
	}

	var newEntries []*models.SysDictEntry
	for _, entry := range sourceEntries {
		newEntries = append(newEntries, &models.SysDictEntry{
			EntryLabel:    entry.EntryLabel,
			EntryValue:    entry.EntryValue,
			NumericValue:  entry.NumericValue,
			LanguageCode:  entry.LanguageCode,
			SysDictTypeId: req.TargetTypeId,
			SortOrder:     mixin.SortOrder{SortOrder: entry.SortOrder.SortOrder},
			IsEnabled:     mixin.IsEnabled{IsEnabled: entry.IsEnabled.IsEnabled},
			Remark:        mixin.Remark{Remark: entry.Remark.Remark},
		})
	}

	err = query.SysDictEntry.Create(newEntries...)
	if err != nil {
		ctx.L().Error("批量复制字典项失败", zap.Error(err), zap.Uint64s("entryIds", req.EntryIds), zap.Uint64("targetTypeId", req.TargetTypeId))
		return res.FailDefault
	}

	return nil
}
