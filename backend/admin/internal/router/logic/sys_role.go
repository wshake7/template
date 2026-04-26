package logic

import (
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/res"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc"
	"orm-crud/gormc/mixin"
)

type SysRoleHandler struct{}

// @Summary 获取角色分页列表
// @Description 分页查询角色信息
// @Tags Role
// @Accept json
// @Produce json
// @Param req body v1.PagingRequest true "分页参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[models.SysRole]} "成功"
// @Router /api/role/list [get]
func (*SysRoleHandler) List(ctx *handler.Ctx, req *v1.PagingRequest) (*gormc.PagingResult[models.SysRole], error) {
	pagination, err := query.SysRole.PageWithPaging(req)
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

// @Summary 创建角色
// @Description 创建新的角色
// @Tags Role
// @Accept json
// @Produce json
// @Param req body ReqRoleCreate true "角色创建参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/role/create [post]
func (*SysRoleHandler) Create(ctx *handler.Ctx, req *ReqRoleCreate) error {
	err := query.SysRole.Create(&models.SysRole{Code: req.Code, Name: req.Name, Remark: mixin.Remark{Remark: req.Remark}})
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

// @Summary 更新角色
// @Description 根据角色 ID 更新角色信息
// @Tags Role
// @Accept json
// @Produce json
// @Param req body ReqRoleUpdate true "角色更新参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/role/update [post]
func (*SysRoleHandler) Update(ctx *handler.Ctx, req *ReqRoleUpdate) error {
	sysRole := query.SysRole
	_, err := sysRole.Where(sysRole.ID.Eq(req.ID)).UpdateSimple(
		sysRole.Name.Value(req.Name),
		sysRole.Code.Value(req.Code),
		sysRole.Remark.Value(req.Remark),
	)
	if err != nil {
		return res.FailDefault
	}
	return nil
}

type ReqRoleSwitchStatus struct {
	ID        uint64 `json:"id" binding:"required" binding_msg:"required=请求错误"`
	IsEnabled uint8  `json:"isEnabled"`
}

// @Summary 切换角色状态
// @Description 根据角色 ID 修改启用状态
// @Tags Role
// @Accept json
// @Produce json
// @Param req body ReqRoleSwitchStatus true "角色状态参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/role/switch [post]
func (*SysRoleHandler) Switch(ctx *handler.Ctx, req *ReqRoleSwitchStatus) error {
	sysRole := query.SysRole
	_, err := sysRole.Where(sysRole.ID.Eq(req.ID)).UpdateSimple(
		sysRole.IsEnabled.Value(req.IsEnabled != 0),
	)
	if err != nil {
		return res.FailDefault
	}
	return nil
}

type ReqRoleDelete struct {
	ID uint64 `json:"id" binding:"required" binding_msg:"required=请求错误"`
}

// @Summary 删除角色
// @Description 根据角色 ID 删除角色
// @Tags Role
// @Accept json
// @Produce json
// @Param req body ReqRoleDelete true "角色删除参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/role/del [post]
func (*SysRoleHandler) Del(ctx *handler.Ctx, req *ReqRoleDelete) error {
	_, err := query.SysRole.Where(query.SysRole.ID.Eq(req.ID)).Delete()
	if err != nil {
		return res.FailDefault
	}
	return nil
}
