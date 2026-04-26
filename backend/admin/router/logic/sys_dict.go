package logic

import (
	"admin/fiberc/handler"
	"admin/fiberc/res"
	"admin/services/orm"
	"admin/services/orm/models"
	"admin/services/orm/query"
	"admin/services/orm/repo"
	"errors"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc"
	"orm-crud/gormc/mixin"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SysDictHandler struct{}

// --- 字典类型 (DictType) ---

type ReqDictTypeCreate struct {
	TypeCode    string `json:"typeCode" binding:"required,max=128" binding_msg:"required=字典类型代码不能为空,max=字典类型代码最多128位"`
	TypeName    string `json:"typeName" binding:"required,max=255" binding_msg:"required=字典类型名称不能为空,max=字典类型名称最多255位"`
	IsEnabled   bool   `json:"isEnabled"`
	SortOrder   int32  `json:"sortOrder"`
	Description string `json:"description" binding:"max=255" binding_msg:"max=描述最多255位"`
}

type ReqDictTypeUpdate struct {
	ID          uint64  `json:"id" binding:"required" binding_msg:"required=请求错误"`
	TypeCode    *string `json:"typeCode" binding:"omitempty,max=128" binding_msg:"max=字典类型代码最多128位"`
	TypeName    *string `json:"typeName" binding:"omitempty,max=255" binding_msg:"max=字典类型名称最多255位"`
	IsEnabled   *bool   `json:"isEnabled"`
	SortOrder   *int32  `json:"sortOrder"`
	Description *string `json:"description" binding:"omitempty,max=255" binding_msg:"max=描述最多255位"`
}

type ReqDictTypeSwitchStatus struct {
	ID        uint64 `json:"id" binding:"required" binding_msg:"required=请求错误"`
	IsEnabled bool   `json:"isEnabled"`
}

type ReqDictTypeBatchDelete struct {
	IDs []uint64 `json:"ids" binding:"required,min=1" binding_msg:"required=请选择字典类型,min=至少选择一项"`
}

// @Summary 获取字典类型分页列表
// @Description 分页查询字典类型信息
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body v1.PagingRequest true "分页参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[models.SysDictType]} "成功"
// @Router /api/dict/type/list [post]
func (*SysDictHandler) TypeList(ctx *handler.Ctx, req *v1.PagingRequest) (*gormc.PagingResult[models.SysDictType], error) {
	pagination, err := repo.SysDictTypeRepo.ListWithPaging(ctx.Context(), orm.DB(), req)
	if err != nil {
		return nil, res.FailDefault
	}
	return pagination, nil
}

// @Summary 创建字典类型
// @Description 创建新的字典类型
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body ReqDictTypeCreate true "字典类型创建参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/dict/type/create [post]
func (*SysDictHandler) TypeCreate(ctx *handler.Ctx, req *ReqDictTypeCreate) error {
	operationID := ctx.SessionInfo.Id
	_, err := repo.SysDictTypeRepo.Create(ctx.Context(), orm.DB(), &models.SysDictType{
		CreatedBy:   mixin.CreatedBy{CreatedBy: operationID},
		UpdatedBy:   mixin.UpdatedBy{UpdatedBy: operationID},
		IsEnabled:   mixin.IsEnabled{IsEnabled: req.IsEnabled},
		SortOrder:   mixin.SortOrder{SortOrder: req.SortOrder},
		Description: mixin.Description{Description: req.Description},
		TypeCode:    req.TypeCode,
		TypeName:    req.TypeName,
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
// @Description 根据 ID 更新字典类型信息
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body ReqDictTypeUpdate true "字典类型更新参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/dict/type/update [post]
func (*SysDictHandler) TypeUpdate(ctx *handler.Ctx, req *ReqDictTypeUpdate) error {
	operationID := ctx.SessionInfo.Id
	sysDictType := query.SysDictType
	m := map[string]any{
		sysDictType.UpdatedBy.ColumnName().String():   operationID,
		sysDictType.TypeCode.ColumnName().String():    req.TypeCode,
		sysDictType.TypeName.ColumnName().String():    req.TypeName,
		sysDictType.IsEnabled.ColumnName().String():   req.IsEnabled,
		sysDictType.SortOrder.ColumnName().String():   req.SortOrder,
		sysDictType.Description.ColumnName().String(): req.Description,
	}
	_, err := repo.SysDictTypeRepo.UpdateNoNilMap(m, sysDictType.ID.Eq(req.ID))
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return res.FailMsg("类型编码已存在")
		}
		return res.FailDefault
	}
	return nil
}

// @Summary 切换字典类型状态
// @Description 根据 ID 修改启用状态
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body ReqDictTypeSwitchStatus true "状态参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/dict/type/switch [post]
func (*SysDictHandler) TypeSwitch(ctx *handler.Ctx, req *ReqDictTypeSwitchStatus) error {
	operationID := ctx.SessionInfo.Id
	sysDictType := query.SysDictType
	_, err := repo.SysDictTypeRepo.UpdateMap(map[string]any{
		sysDictType.UpdatedBy.ColumnName().String(): operationID,
		sysDictType.IsEnabled.ColumnName().String(): req.IsEnabled,
	}, sysDictType.ID.Eq(req.ID))
	if err != nil {
		return res.FailDefault
	}
	return nil
}

// @Summary 批量删除字典类型
// @Description 根据 ID 列表批量删除字典类型及其关联的所有字典项
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body ReqDictTypeBatchDelete true "批量删除参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/dict/type/del [post]
func (*SysDictHandler) TypeDel(ctx *handler.Ctx, req *ReqDictTypeBatchDelete) error {
	err := orm.DB().Transaction(func(tx *gorm.DB) error {
		// 1. 删除关联的字典项
		_, err := repo.SysDictEntryRepo.SoftDelete(ctx.Context(), tx.Where(query.SysDictEntry.SysDictTypeId.In(req.IDs...)))
		if err != nil {
			return err
		}

		// 2. 删除字典类型
		_, err = repo.SysDictTypeRepo.SoftDelete(ctx.Context(), tx.Where(query.SysDictType.ID.In(req.IDs...)))
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

type ReqDictEntrySwitchStatus struct {
	ID        uint64 `json:"id" binding:"required" binding_msg:"required=请求错误"`
	IsEnabled bool   `json:"isEnabled"`
}

type ReqDictEntryBatchDelete struct {
	IDs []uint64 `json:"ids" binding:"required,min=1" binding_msg:"required=请选择字典项,min=至少选择一项"`
}

type ReqDictEntryBatchCopy struct {
	EntryIds     []uint64 `json:"entryIds" binding:"required,min=1" binding_msg:"required=请选择字典项,min=至少选择一项"`
	TargetTypeId uint64   `json:"targetTypeId" binding:"required" binding_msg:"required=目标字典类型不能为空"`
}

// @Summary 获取字典数据项分页列表
// @Description 分页查询字典数据项信息
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body v1.PagingRequest true "分页参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[models.SysDictEntry]} "成功"
// @Router /api/dict/entry/list [post]
func (*SysDictHandler) EntryList(ctx *handler.Ctx, req *v1.PagingRequest) (*gormc.PagingResult[models.SysDictEntry], error) {
	pagination, err := repo.SysDictEntryRepo.ListWithPaging(ctx.Context(), orm.DB(), req)
	if err != nil {
		return nil, res.FailDefault
	}
	return pagination, nil
}

// @Summary 创建字典数据项
// @Description 创建新的字典数据项
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body ReqDictEntryCreate true "字典数据项创建参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/dict/entry/create [post]
func (*SysDictHandler) EntryCreate(ctx *handler.Ctx, req *ReqDictEntryCreate) error {
	operationID := ctx.SessionInfo.Id
	// 校验字典类型是否存在
	exists, err := repo.SysDictTypeRepo.Exists(ctx.Context(), orm.DB().Where(query.SysDictType.ID.Eq(req.SysDictTypeId), query.SysDictType.IsEnabled.Is(req.IsEnabled)))
	if err != nil {
		return res.FailDefault
	}
	if !exists {
		return res.FailMsg("字典类型不存在")
	}

	_, err = repo.SysDictEntryRepo.Create(ctx.Context(), orm.DB(), &models.SysDictEntry{
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
		SysDictType:   nil,
	})
	if err != nil {
		return res.FailDefault
	}
	return nil
}

// @Summary 更新字典数据项
// @Description 根据 ID 更新字典数据项信息
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body ReqDictEntryUpdate true "字典数据项更新参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/dict/entry/update [post]
func (*SysDictHandler) EntryUpdate(ctx *handler.Ctx, req *ReqDictEntryUpdate) error {
	operationID := ctx.SessionInfo.Id
	// 如果更新了 SysDictTypeId，校验其是否存在
	if req.SysDictTypeId != nil {
		exists, err := repo.SysDictTypeRepo.Exists(ctx.Context(), orm.DB().Where(query.SysDictType.ID.Eq(*req.SysDictTypeId)))
		if err != nil {
			return res.FailDefault
		}
		if !exists {
			return res.FailMsg("字典类型不存在")
		}
	}

	sysDictEntry := query.SysDictEntry
	_, err := repo.SysDictEntryRepo.UpdateNoNilMap(map[string]any{
		sysDictEntry.UpdatedBy.ColumnName().String():     operationID,
		sysDictEntry.EntryLabel.ColumnName().String():    req.EntryLabel,
		sysDictEntry.EntryValue.ColumnName().String():    req.EntryValue,
		sysDictEntry.NumericValue.ColumnName().String():  req.NumericValue,
		sysDictEntry.LanguageCode.ColumnName().String():  req.LanguageCode,
		sysDictEntry.SysDictTypeId.ColumnName().String(): req.SysDictTypeId,
		sysDictEntry.SortOrder.ColumnName().String():     req.SortOrder,
		sysDictEntry.IsEnabled.ColumnName().String():     req.IsEnabled,
		sysDictEntry.Remark.ColumnName().String():        req.Remark,
	}, sysDictEntry.ID.Eq(req.ID))
	if err != nil {
		return res.FailDefault
	}
	return nil
}

// @Summary 切换字典数据项状态
// @Description 根据 ID 修改启用状态
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body ReqDictEntrySwitchStatus true "状态参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/dict/entry/switch [post]
func (*SysDictHandler) EntrySwitch(ctx *handler.Ctx, req *ReqDictEntrySwitchStatus) error {
	operationID := ctx.SessionInfo.Id
	sysDictEntry := query.SysDictEntry
	_, err := repo.SysDictEntryRepo.UpdateMap(map[string]any{
		sysDictEntry.UpdatedBy.ColumnName().String(): operationID,
		sysDictEntry.IsEnabled.ColumnName().String(): req.IsEnabled,
	}, sysDictEntry.ID.Eq(req.ID))
	if err != nil {
		return res.FailDefault
	}
	return nil
}

// @Summary 批量删除字典数据项
// @Description 根据 ID 列表批量删除字典数据项
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body ReqDictEntryBatchDelete true "批量删除参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/dict/entry/del [post]
func (*SysDictHandler) EntryDel(ctx *handler.Ctx, req *ReqDictEntryBatchDelete) error {
	_, err := repo.SysDictEntryRepo.SoftDelete(ctx.Context(), orm.DB().Where(query.SysDictEntry.ID.In(req.IDs...)))
	if err != nil {
		return res.FailDefault
	}
	return nil
}

// @Summary 批量复制字典数据项
// @Description 将选中的字典数据项批量复制到指定字典类型下（不支持复制到同一类型）
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body ReqDictEntryBatchCopy true "批量复制参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/dict/entry/batch/copy [post]
func (*SysDictHandler) EntryBatchCopy(ctx *handler.Ctx, req *ReqDictEntryBatchCopy) error {
	// 校验目标字典类型是否存在
	exists, err := repo.SysDictTypeRepo.Exists(ctx.Context(), orm.DB().Where(query.SysDictType.ID.Eq(req.TargetTypeId)))
	if err != nil {
		ctx.L().Error("校验字典类型失败", zap.Error(err), zap.Uint64("targetTypeId", req.TargetTypeId))
		return res.FailDefault
	}
	if !exists {
		return res.FailMsg("目标字典类型不存在")
	}

	// 查询源字典项
	var sourceEntries []models.SysDictEntry
	err = orm.DB().WithContext(ctx.Context()).Where("id IN ?", req.EntryIds).Find(&sourceEntries).Error
	if err != nil {
		ctx.L().Error("查询源字典项失败", zap.Error(err), zap.Uint64s("entryIds", req.EntryIds))
		return res.FailDefault
	}
	if len(sourceEntries) == 0 {
		return res.FailMsg("未找到要复制的字典项")
	}

	// 创建新字典项（复制到目标类型）
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

	_, err = repo.SysDictEntryRepo.BatchCreate(ctx.Context(), orm.DB(), newEntries)
	if err != nil {
		ctx.L().Error("批量复制字典项失败", zap.Error(err), zap.Uint64s("entryIds", req.EntryIds), zap.Uint64("targetTypeId", req.TargetTypeId))
		return res.FailDefault
	}

	return nil
}
