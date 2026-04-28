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
	"gorm.io/gorm"
)

type SysResourceHandler struct{}

type ReqResourceCreate struct {
	Type      string `json:"type" binding:"required,min=1,max=32" binding_msg:"required=资源类型不能为空,min=资源类型不能为空,max=资源类型最多32位"`
	Code      string `json:"code" binding:"required,min=1,max=255" binding_msg:"required=资源编码不能为空,min=资源编码不能为空,max=资源编码最多255位"`
	Name      string `json:"name" binding:"required,min=1,max=255" binding_msg:"required=资源名称不能为空,min=资源名称不能为空,max=资源名称最多255位"`
	IsEnabled bool   `json:"isEnabled"`
	Remark    string `json:"remark" binding:"max=255" binding_msg:"max=备注最多255位"`
}

type ReqResourceUpdate struct {
	ID        uint64  `json:"id" binding:"required" binding_msg:"required=请求错误"`
	Type      *string `json:"type" binding:"omitempty,max=32" binding_msg:"max=资源类型最多32位"`
	Code      *string `json:"code" binding:"omitempty,max=255" binding_msg:"max=资源编码最多255位"`
	Name      *string `json:"name" binding:"omitempty,max=255" binding_msg:"max=资源名称最多255位"`
	IsEnabled *bool   `json:"isEnabled"`
	Remark    *string `json:"remark" binding:"omitempty,max=255" binding_msg:"max=备注最多255位"`
}

type ReqResourceBatchDelete struct {
	IDs []uint64 `json:"ids" binding:"required,min=1" binding_msg:"required=请选择资源,min=至少选择一项"`
}

// @Summary 获取资源分页列表
// @Remark 分页查询资源信息
// @Tags Resource
// @Accept json
// @Produce json
// @Param req body v1.PagingRequest true "分页参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[models.SysResource]} "成功"
// @Router /api/sys/resource/list [post]
func (*SysResourceHandler) List(ctx *handler.Ctx, req *v1.PagingRequest) (*gormc.PagingResult[models.SysResource], error) {
	pagination, err := query.SysResource.PageWithPaging(req)
	if err != nil {
		return nil, res.FailDefault
	}
	return pagination, nil
}

// @Summary 创建资源
// @Remark 创建新的资源
// @Tags Resource
// @Accept json
// @Produce json
// @Param req body ReqResourceCreate true "资源创建参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/resource/create [post]
func (*SysResourceHandler) Create(ctx *handler.Ctx, req *ReqResourceCreate) error {
	operationID := ctx.SessionInfo.Id
	err := query.SysResource.Create(&models.SysResource{
		OperatorID: mixin.OperatorID{
			CreatedBy: mixin.CreatedBy{CreatedBy: operationID},
			UpdatedBy: mixin.UpdatedBy{UpdatedBy: operationID},
		},
		IsEnabled: mixin.IsEnabled{IsEnabled: req.IsEnabled},
		Remark:    mixin.Remark{Remark: req.Remark},
		Type:      req.Type,
		Code:      req.Code,
		Name:      req.Name,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return res.FailMsg("资源编码已存在")
		}
		return res.FailDefault
	}
	return nil
}

// @Summary 更新资源
// @Remark 根据 ID 更新资源信息
// @Tags Resource
// @Accept json
// @Produce json
// @Param req body ReqResourceUpdate true "资源更新参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/resource/update [post]
func (*SysResourceHandler) Update(ctx *handler.Ctx, req *ReqResourceUpdate) error {
	sysResource := query.SysResource
	operationID := ctx.SessionInfo.Id

	exprs := []field.AssignExpr{sysResource.UpdatedBy.Value(operationID)}
	query.ExprAppendSelf(&exprs, req.Type, sysResource.Type.Value)
	query.ExprAppendSelf(&exprs, req.Code, sysResource.Code.Value)
	query.ExprAppendSelf(&exprs, req.Name, sysResource.Name.Value)
	query.ExprAppendSelf(&exprs, req.IsEnabled, sysResource.IsEnabled.Value)
	query.ExprAppendSelf(&exprs, req.Remark, sysResource.Remark.Value)

	_, err := sysResource.Where(sysResource.ID.Eq(req.ID)).UpdateSimple(exprs...)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return res.FailMsg("资源编码已存在")
		}
		return res.FailDefault
	}
	return nil
}

// @Summary 批量删除资源
// @Remark 根据 ID 列表批量删除资源
// @Tags Resource
// @Accept json
// @Produce json
// @Param req body ReqResourceBatchDelete true "批量删除参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/resource/del [post]
func (*SysResourceHandler) Del(ctx *handler.Ctx, req *ReqResourceBatchDelete) error {
	_, err := query.SysResource.Where(query.SysResource.ID.In(req.IDs...)).Delete()
	if err != nil {
		return res.FailDefault
	}
	return nil
}
