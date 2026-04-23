package logic

import (
	"admin/fiberc/handler"
	"admin/fiberc/res"
	"admin/services/orm"
	"admin/services/orm/models"
	"admin/services/orm/query"
	"admin/services/orm/repo"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gorm"
	"orm-crud/gorm/mixin"

	"gorm.io/gen/field"
)

type RoleHandler struct{}

func (h *RoleHandler) List(ctx *handler.Ctx, req *v1.PaginationRequest) (*gorm.PagingResult[models.SysRole], error) {
	pagination, err := repo.SysRoleRepo.ListWithPagination(ctx.Context(), orm.DB(), req)
	if err != nil {
		return nil, res.FailDefault
	}
	return pagination, nil
}

type ReqRoleCreate struct {
	Name   string `json:"name" binding:"required,min=1,max=64" binding_msg:"required=角色名称不能为空,min=角色名称最少1位,max=角色名称最多64位"`
	Code   string `json:"code" binding:"required,min=1,max=64" binding_msg:"required=角色编码不能为空,min=角色编码最少1位,max=角色编码最多64位"`
	Remark string `json:"remark" binding:"max=255" binding_msg:"max=备注最多255位"`
}

func (*RoleHandler) Create(ctx *handler.Ctx, req *ReqRoleCreate) error {
	_, err := repo.SysRoleRepo.Create(ctx.Context(), orm.DB(), &models.SysRole{Code: req.Code, Name: req.Name, Remark: mixin.Remark{Remark: req.Remark}})
	if err != nil {
		return res.FailDefault
	}
	return nil
}

type ReqRoleUpdate struct {
	ID     uint64 `json:"id" binding:"required" binding_msg:"required=请求错误"`
	Name   string `json:"name" binding:"required,min=1,max=64" binding_msg:"required=角色名称不能为空,min=角色名称最少1位,max=角色名称最多64位"`
	Code   string `json:"code" binding:"required,min=1,max=64" binding_msg:"required=角色编码不能为空,min=角色编码最少1位,max=角色编码最多64位"`
	Remark string `json:"remark" binding:"max=255" binding_msg:"max=备注最多255位"`
}

func (*RoleHandler) Update(ctx *handler.Ctx, req *ReqRoleUpdate) error {
	sysRole := query.SysRole
	_, err := repo.SysRoleRepo.UpdateMap(map[field.Expr]any{
		sysRole.Name:   req.Name,
		sysRole.Code:   req.Code,
		sysRole.Remark: mixin.Remark{Remark: req.Remark},
	}, sysRole.ID.Eq(req.ID))
	if err != nil {
		return res.FailDefault
	}
	return nil
}

type ReqRoleSwitchStatus struct {
	ID     uint64 `json:"id" binding:"required" binding_msg:"required=请求错误"`
	Status uint8  `json:"status" binding:"required" binding_msg:"required=角色状态不能为空"`
}

func (*RoleHandler) SwitchStatus(ctx *handler.Ctx, req *ReqRoleSwitchStatus) error {
	sysRole := query.SysRole
	_, err := repo.SysRoleRepo.UpdateMap(map[field.Expr]any{
		sysRole.Status: req.Status,
	}, sysRole.ID.Eq(req.ID))
	if err != nil {
		return res.FailDefault
	}
	return nil
}

type ReqRoleDelete struct {
	ID uint64 `json:"id" binding:"required" binding_msg:"required=请求错误"`
}

func (*RoleHandler) Delete(ctx *handler.Ctx, req *ReqRoleDelete) error {
	_, err := repo.SysRoleRepo.SoftDelete(ctx.Context(), orm.DB().Where(query.SysUser.ID.Eq(req.ID)))
	if err != nil {
		return res.FailDefault
	}
	return nil
}
