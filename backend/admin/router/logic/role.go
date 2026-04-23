package logic

import (
	"admin/fiberc/handler"
	"admin/services/orm/models"
	"orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gorm"
)

type RoleHandler struct{}

func (h *RoleHandler) List(ctx *handler.Ctx, req *v1.PaginationRequest) (*gorm.PagingResult[models.SysRole], error) {
	return nil, nil
}

type ReqRoleCreate struct {
	Name   string `json:"name" binding:"required,min=1,max=64"`
	Code   string `json:"code" binding:"required,min=1,max=64"`
	Remark string `json:"remark" binding:"max=255"`
}

func (*RoleHandler) Create(ctx *handler.Ctx, req *ReqRoleCreate) error {
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
