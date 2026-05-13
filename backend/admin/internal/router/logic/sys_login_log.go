package logic

import (
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/res"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc"

	"go.uber.org/zap"
)

type SysLoginLogHandler struct{}

type ReqLoginLogDetail struct {
	ID uint64 `json:"id" binding:"required" binding_msg:"required=请求错误"`
}

// @Summary 获取登录日志分页列表
// @Description 分页查询登录日志信息
// @Tags Log
// @Accept json
// @Produce json
// @Param req body v1.PagingRequest true "分页参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[models.SysLoginLog]} "成功"
// @Router /api/sys/login/log/list [post]
func (*SysLoginLogHandler) List(ctx *handler.Ctx, req *v1.PagingRequest) (*gormc.PagingResult[models.SysLoginLog], error) {
	pagination, err := query.SysLoginLog.PageWithPaging(req)
	if err != nil {
		ctx.L().Error("query login log list fail", zap.Error(err))
		return nil, res.FailDefault
	}
	if err = reloadSysLoginLogUsers(pagination.Items); err != nil {
		ctx.L().Error("query login log users fail", zap.Error(err))
		return nil, res.FailDefault
	}
	return pagination, nil
}

func reloadSysLoginLogUsers(items []*models.SysLoginLog) error {
	if len(items) == 0 {
		return nil
	}

	userIDSet := make(map[uint64]struct{})
	for _, item := range items {
		if item == nil || item.SysUserID == nil {
			continue
		}
		userIDSet[*item.SysUserID] = struct{}{}
	}
	if len(userIDSet) == 0 {
		return nil
	}

	userIDs := make([]uint64, 0, len(userIDSet))
	for userID := range userIDSet {
		userIDs = append(userIDs, userID)
	}

	sysUser := query.SysUser
	users, err := sysUser.Select(sysUser.ID, sysUser.Username, sysUser.Nickname).Where(sysUser.ID.In(userIDs...)).Find()
	if err != nil {
		return err
	}

	userMap := make(map[uint64]*models.SysUser, len(users))
	for _, user := range users {
		if user != nil {
			userMap[user.ID] = user
		}
	}
	for _, item := range items {
		if item == nil || item.SysUserID == nil {
			continue
		}
		item.SysUser = userMap[*item.SysUserID]
	}

	return nil
}

// @Summary 获取登录日志详情
// @Description 根据 ID 获取登录日志详情
// @Tags Log
// @Accept json
// @Produce json
// @Param req body ReqLoginLogDetail true "日志ID"
// @Success 200 {object} res.Response{data=models.SysLoginLog} "成功"
// @Router /api/sys/login/log/detail [post]
func (*SysLoginLogHandler) Detail(ctx *handler.Ctx, req *ReqLoginLogDetail) (*models.SysLoginLog, error) {
	sysLoginLog := query.SysLoginLog
	sysUser := query.SysUser
	logEntry, err := sysLoginLog.
		Preload(sysLoginLog.SysUser.Select(sysUser.ID, sysUser.Username, sysUser.Nickname)).
		Where(sysLoginLog.ID.Eq(req.ID)).
		First()
	if err != nil {
		ctx.L().Error("query login log detail fail", zap.Error(err), zap.Uint64("id", req.ID))
		return nil, res.FailDefault
	}
	return logEntry, nil
}
