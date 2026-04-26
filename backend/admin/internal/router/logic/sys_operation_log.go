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

type SysOperationLogHandler struct{}

type ReqLogDetail struct {
	ID uint64 `json:"id" binding:"required" binding_msg:"required=请求错误"`
}

// @Summary 获取操作日志分页列表
// @Description 分页查询操作日志信息
// @Tags Log
// @Accept json
// @Produce json
// @Param req body v1.PagingRequest true "分页参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[models.SysOperationLog]} "成功"
// @Router /api/log/list [post]
func (*SysOperationLogHandler) List(ctx *handler.Ctx, req *v1.PagingRequest) (*gormc.PagingResult[models.SysOperationLog], error) {
	pagination, err := query.SysOperationLog.PageWithPaging(req)
	if err != nil {
		return nil, res.FailDefault
	}
	return pagination, nil
}

// @Summary 获取操作日志详情
// @Description 根据 ID 获取操作日志详情
// @Tags Log
// @Accept json
// @Produce json
// @Param req body ReqLogDetail true "日志ID"
// @Success 200 {object} res.Response{data=models.SysOperationLog} "成功"
// @Router /api/log/detail [post]
func (*SysOperationLogHandler) Detail(ctx *handler.Ctx, req *ReqLogDetail) (*models.SysOperationLog, error) {
	logEntry, err := query.SysOperationLog.
		Where(query.SysOperationLog.ID.Eq(req.ID)).
		First()
	if err != nil {
		ctx.L().Error("查询操作日志失败", zap.Error(err), zap.Uint64("id", req.ID))
		return nil, res.FailDefault
	}
	return logEntry, nil
}
