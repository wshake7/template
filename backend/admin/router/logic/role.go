package logic

import (
	"admin/fiberc/handler"
	"admin/fiberc/res"
	"admin/services/orm"
	"admin/services/orm/models"
	"admin/services/orm/repo"
	"orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gorm"
	"orm-crud/gorm/mixin"
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
	Name   string `json:"name" binding:"required,min=1,max=64"`
	Code   string `json:"code" binding:"required,min=1,max=64"`
	Remark string `json:"remark" binding:"max=255"`
}

func (*RoleHandler) Create(ctx *handler.Ctx, req *ReqRoleCreate) error {
	_, err := repo.SysRoleRepo.Create(ctx.Context(), orm.DB(), &models.SysRole{Code: req.Code, Name: req.Name, Remark: mixin.Remark{Remark: req.Remark}})
	if err != nil {
		return res.FailDefault
	}
	return nil
}

type ReqRoleUpdate struct {
	ID     uint64 `json:"id" binding:"required,min=1"`
	Name   string `json:"name" binding:"required,min=1,max=64"`
	Code   string `json:"code" binding:"required,min=1,max=64"`
	Remark string `json:"remark" binding:"max=255"`
}

func (*RoleHandler) Update(ctx *handler.Ctx, req *ReqRoleUpdate) error {
	return nil
}

type ReqRoleSwitchStatus struct {
	ID     uint64 `json:"id" binding:"required,min=1"`
	Status uint8  `json:"status" binding:"required"`
}

func (*RoleHandler) SwitchStatus(ctx *handler.Ctx, req *ReqRoleSwitchStatus) error {
	return nil
}

type ReqRoleDelete struct {
	ID uint64 `json:"id" binding:"required,min=1"`
}

func (*RoleHandler) Delete(ctx *handler.Ctx, req *ReqRoleDelete) error {
	return nil
}
