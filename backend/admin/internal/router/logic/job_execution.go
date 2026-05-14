package logic

import (
	"fmt"
	"strings"
	"time"

	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/res"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"admin/internal/services/temporaljob"

	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc"

	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type JobExecutionHandler struct{}

type RespJobExecution struct {
	models.JobExecution
	CanWrite  bool `json:"canWrite"`
	CanDelete bool `json:"canDelete"`
}

type ReqJobExecutionID struct {
	ID uint64 `json:"id" binding:"required" binding_msg:"required=请求错误"`
}

// @Summary 获取任务执行记录分页列表
// @Tags JobExecution
// @Accept json
// @Produce json
// @Param req body v1.PagingRequest true "分页参数"
// @Success 200 {object} res.Response{data=gormc.PagingResult[RespJobExecution]} "成功"
// @Router /api/sys/job/execution/list [post]
func (*JobExecutionHandler) List(ctx *handler.Ctx, req *v1.PagingRequest) (*gormc.PagingResult[RespJobExecution], error) {
	if req.GetOrderBy() == "" {
		orderBy := "id desc"
		req.OrderBy = &orderBy
	}

	pagination, err := query.JobExecution.PageWithPaging(req)
	if err != nil {
		ctx.L().Error("query job execution list fail", zap.Error(err))
		return nil, res.FailDefault
	}

	items := make([]*RespJobExecution, 0, len(pagination.Items))
	for _, item := range pagination.Items {
		items = append(items, &RespJobExecution{
			JobExecution: *item,
			CanWrite:     true,
			CanDelete:    true,
		})
	}
	return &gormc.PagingResult[RespJobExecution]{
		Items: items,
		Total: pagination.Total,
	}, nil
}

// @Summary 获取任务执行记录详情
// @Tags JobExecution
// @Accept json
// @Produce json
// @Param req body ReqJobExecutionID true "详情参数"
// @Success 200 {object} res.Response{data=models.JobExecution} "成功"
// @Router /api/sys/job/execution/detail [post]
func (*JobExecutionHandler) Detail(ctx *handler.Ctx, req *ReqJobExecutionID) (*models.JobExecution, error) {
	current, err := findJobExecution(req.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, res.FailMsg("执行记录不存在")
		}
		ctx.L().Error("query job execution detail fail", zap.Error(err), zap.Uint64("id", req.ID))
		return nil, res.FailDefault
	}
	return current, nil
}

// @Summary 取消运行中的任务执行
// @Tags JobExecution
// @Accept json
// @Produce json
// @Param req body ReqJobExecutionID true "取消参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/job/execution/cancel [post]
func (*JobExecutionHandler) Cancel(ctx *handler.Ctx, req *ReqJobExecutionID) error {
	current, err := findJobExecution(req.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return res.FailMsg("执行记录不存在")
		}
		return res.FailDefault
	}
	if current.Status != models.JobExecutionStatusRunning {
		return res.FailMsg("仅运行中的执行记录可取消")
	}
	service, err := requireTemporal()
	if err != nil {
		return err
	}
	if err = service.Client.CancelWorkflow(ctx.Context(), current.TemporalWorkflowID, current.TemporalRunID); err != nil {
		ctx.L().Error("cancel workflow fail", zap.Error(err), zap.Uint64("id", req.ID), zap.String("workflowID", current.TemporalWorkflowID), zap.String("runID", current.TemporalRunID))
		return res.FailMsg("取消 Workflow 失败")
	}

	now := time.Now()
	jobExecution := query.JobExecution
	_, err = jobExecution.Where(jobExecution.ID.Eq(req.ID), jobExecution.Status.Eq(models.JobExecutionStatusRunning)).Updates(map[string]any{
		"status":        models.JobExecutionStatusCanceled,
		"end_time":      &now,
		"error_message": "canceled by admin",
	})
	if err != nil {
		return res.FailDefault
	}
	return nil
}

// @Summary 重试失败态任务执行
// @Tags JobExecution
// @Accept json
// @Produce json
// @Param req body ReqJobExecutionID true "重试参数"
// @Success 200 {object} res.Response "成功"
// @Router /api/sys/job/execution/retry [post]
func (*JobExecutionHandler) Retry(ctx *handler.Ctx, req *ReqJobExecutionID) error {
	current, err := findJobExecution(req.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return res.FailMsg("执行记录不存在")
		}
		return res.FailDefault
	}
	if !isRetryableJobExecution(current.Status) {
		return res.FailMsg("仅失败、取消、超时的执行记录可重试")
	}
	schedule, err := findJobScheduleByCode(current.JobCode)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return res.FailMsg("任务配置不存在")
		}
		return res.FailDefault
	}
	inputValue, err := parseJSONValue(current.InputJSON)
	if err != nil {
		return res.FailMsg("执行记录输入参数不是合法 JSON")
	}
	service, err := requireTemporal()
	if err != nil {
		return err
	}

	workflowIDPrefix := schedule.TemporalWorkflowIDPrefix
	if workflowIDPrefix == "" {
		workflowIDPrefix = schedule.JobCode
	}
	workflowID := fmt.Sprintf("%s-retry-dispatch-%d", workflowIDPrefix, time.Now().UnixNano())
	_, err = service.Client.ExecuteWorkflow(ctx.Context(), client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: schedule.TaskQueue,
	}, temporaljob.DispatchWorkflowName, temporaljob.DispatchInput{
		JobCode:          schedule.JobCode,
		WorkflowType:     schedule.WorkflowType,
		TaskQueue:        schedule.TaskQueue,
		WorkflowIDPrefix: workflowIDPrefix,
		Input:            inputValue,
		RetryCount:       current.RetryCount + 1,
	})
	if err != nil {
		ctx.L().Error("retry workflow fail", zap.Error(err), zap.Uint64("id", req.ID), zap.String("workflowID", workflowID))
		return res.FailMsg("重试 Workflow 失败")
	}
	return nil
}

func findJobExecution(id uint64) (*models.JobExecution, error) {
	jobExecution := query.JobExecution
	return jobExecution.Where(jobExecution.ID.Eq(id)).First()
}

func findJobScheduleByCode(jobCode string) (*models.JobSchedule, error) {
	jobSchedule := query.JobSchedule
	return jobSchedule.Where(
		jobSchedule.JobCode.Eq(strings.TrimSpace(jobCode)),
		jobSchedule.Status.Neq(models.JobScheduleStatusDeleted),
	).First()
}

func isRetryableJobExecution(status string) bool {
	switch status {
	case models.JobExecutionStatusFailed, models.JobExecutionStatusCanceled, models.JobExecutionStatusTimeout:
		return true
	default:
		return false
	}
}
