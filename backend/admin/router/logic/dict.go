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

type DictHandler struct{}

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

type ReqDictTypeDelete struct {
	ID uint64 `json:"id" binding:"required" binding_msg:"required=请求错误"`
}

// @Summary 获取字典类型分页列表
// @Description 分页查询字典类型信息
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body v1.PagingRequest true "分页参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[models.SysDictType]} "成功"
// @Router /api/dict/type/list [post]
func (*DictHandler) TypeList(ctx *handler.Ctx, req *v1.PagingRequest) (*gormc.PagingResult[models.SysDictType], error) {
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
func (*DictHandler) TypeCreate(ctx *handler.Ctx, req *ReqDictTypeCreate) error {
	_, err := repo.SysDictTypeRepo.Create(ctx.Context(), orm.DB(), &models.SysDictType{
		TypeCode:    req.TypeCode,
		TypeName:    req.TypeName,
		IsEnabled:   mixin.IsEnabled{IsEnabled: req.IsEnabled},
		SortOrder:   mixin.SortOrder{SortOrder: req.SortOrder},
		Description: mixin.Description{Description: req.Description},
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
func (*DictHandler) TypeUpdate(ctx *handler.Ctx, req *ReqDictTypeUpdate) error {
	sysDictType := query.SysDictType
	m := map[string]any{
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
func (*DictHandler) TypeSwitch(ctx *handler.Ctx, req *ReqDictTypeSwitchStatus) error {
	sysDictType := query.SysDictType
	_, err := repo.SysDictTypeRepo.UpdateMap(map[string]any{
		sysDictType.IsEnabled.ColumnName().String(): req.IsEnabled,
	}, sysDictType.ID.Eq(req.ID))
	if err != nil {
		return res.FailDefault
	}
	return nil
}

// @Summary 删除字典类型
// @Description 根据 ID 删除字典类型及其关联的所有字典项
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body ReqDictTypeDelete true "删除参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/dict/type/del [post]
func (*DictHandler) TypeDel(ctx *handler.Ctx, req *ReqDictTypeDelete) error {
	err := orm.DB().Transaction(func(tx *gorm.DB) error {
		// 1. 删除关联的字典项
		_, err := repo.SysDictEntryRepo.SoftDelete(ctx.Context(), tx.Where(query.SysDictEntry.SysDictTypeId.Eq(req.ID)))
		if err != nil {
			return err
		}

		// 2. 删除字典类型
		_, err = repo.SysDictTypeRepo.SoftDelete(ctx.Context(), tx.Where(query.SysDictType.ID.Eq(req.ID)))
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		ctx.L().Error("删除字典类型失败", zap.Error(err), zap.Uint64("id", req.ID))
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

type ReqDictEntryDelete struct {
	ID uint64 `json:"id" binding:"required" binding_msg:"required=请求错误"`
}

// @Summary 获取字典数据项分页列表
// @Description 分页查询字典数据项信息
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body v1.PagingRequest true "分页参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[models.SysDictEntry]} "成功"
// @Router /api/dict/entry/list [post]
func (*DictHandler) EntryList(ctx *handler.Ctx, req *v1.PagingRequest) (*gormc.PagingResult[models.SysDictEntry], error) {
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
func (*DictHandler) EntryCreate(ctx *handler.Ctx, req *ReqDictEntryCreate) error {
	// 校验字典类型是否存在
	exists, err := repo.SysDictTypeRepo.Exists(ctx.Context(), orm.DB().Where(query.SysDictType.ID.Eq(req.SysDictTypeId), query.SysDictType.IsEnabled.Is(req.IsEnabled)))
	if err != nil {
		return res.FailDefault
	}
	if !exists {
		return res.FailMsg("字典类型不存在")
	}

	_, err = repo.SysDictEntryRepo.Create(ctx.Context(), orm.DB(), &models.SysDictEntry{
		EntryLabel:    req.EntryLabel,
		EntryValue:    req.EntryValue,
		NumericValue:  req.NumericValue,
		LanguageCode:  req.LanguageCode,
		SysDictTypeId: req.SysDictTypeId,
		SortOrder:     mixin.SortOrder{SortOrder: req.SortOrder},
		IsEnabled:     mixin.IsEnabled{IsEnabled: req.IsEnabled},
		Remark:        mixin.Remark{Remark: req.Remark},
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
func (*DictHandler) EntryUpdate(ctx *handler.Ctx, req *ReqDictEntryUpdate) error {
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
func (*DictHandler) EntrySwitch(ctx *handler.Ctx, req *ReqDictEntrySwitchStatus) error {
	sysDictEntry := query.SysDictEntry
	_, err := repo.SysDictEntryRepo.UpdateMap(map[string]any{
		sysDictEntry.IsEnabled.ColumnName().String(): req.IsEnabled,
	}, sysDictEntry.ID.Eq(req.ID))
	if err != nil {
		return res.FailDefault
	}
	return nil
}

// @Summary 删除字典数据项
// @Description 根据 ID 删除字典数据项
// @Tags Dict
// @Accept json
// @Produce json
// @Param req body ReqDictEntryDelete true "删除参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/dict/entry/del [post]
func (*DictHandler) EntryDel(ctx *handler.Ctx, req *ReqDictEntryDelete) error {
	_, err := repo.SysDictEntryRepo.SoftDelete(ctx.Context(), orm.DB().Where(query.SysDictEntry.ID.Eq(req.ID)))
	if err != nil {
		return res.FailDefault
	}
	return nil
}
