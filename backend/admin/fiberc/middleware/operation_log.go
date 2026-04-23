package middleware

import (
	"admin/fiberc/handler"
	"admin/services/orm"
	"admin/services/orm/models"
	"admin/services/orm/repo"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v3"
	"go-common/utils/coroutine"
	"go-common/utils/function"
	"go-common/utils/ip_util"
	"go.uber.org/zap"
	"reflect"
	"strings"
	"time"
)

func DiffChange(before, after any) string {
	if after == nil {
		return ""
	}

	// 反射查找 ChangeString 方法
	afterVal := reflect.ValueOf(after)
	method := afterVal.MethodByName("ChangeString")
	if method.IsValid() {
		var beforeVal reflect.Value
		if before == nil {
			// before 为 nil，传零值
			beforeVal = reflect.Zero(afterVal.Type())
		} else {
			beforeVal = reflect.ValueOf(before)
		}
		results := method.Call([]reflect.Value{beforeVal})
		return results[0].String()
	}

	if before == nil {
		return diffNewByTag(after)
	}
	return diffByTag(before, after)
}

func diffByTag(before, after any) string {
	bVal := reflect.ValueOf(before)
	aVal := reflect.ValueOf(after)
	// 解指针
	for bVal.Kind() == reflect.Pointer {
		bVal = bVal.Elem()
	}
	for aVal.Kind() == reflect.Pointer {
		aVal = aVal.Elem()
	}
	if bVal.Kind() != reflect.Struct || aVal.Kind() != reflect.Struct {
		return ""
	}
	t := bVal.Type()
	var parts []string
	for i := range t.NumField() {
		field := t.Field(i)
		tag := field.Tag.Get("change")
		if tag == "" {
			continue
		}
		bField := bVal.Field(i)
		aField := aVal.Field(i)
		if !reflect.DeepEqual(bField.Interface(), aField.Interface()) {
			parts = append(parts, fmt.Sprintf("%s：%v->%v", tag, bField.Interface(), aField.Interface()))
		}
	}
	return strings.Join(parts, ", ")
}

func diffNewByTag(after any) string {
	aVal := reflect.ValueOf(after)
	for aVal.Kind() == reflect.Pointer {
		aVal = aVal.Elem()
	}
	if aVal.Kind() != reflect.Struct {
		return ""
	}
	t := aVal.Type()
	var parts []string
	for i := range t.NumField() {
		field := t.Field(i)
		tag := field.Tag.Get("change")
		if tag == "" {
			continue
		}
		parts = append(parts, fmt.Sprintf("%s：%v", tag, aVal.Field(i).Interface()))
	}
	return strings.Join(parts, ", ")
}

type Changeable[T any] interface {
	ChangeString(before T) string
}

type ChangeQueryHandler func(ctx *handler.Ctx) (any, error)

func ChangeQueryNilFn[R any](fn func(ctx *handler.Ctx) (*R, error)) ChangeQueryHandler {
	return func(ctx *handler.Ctx) (any, error) {
		r, err := fn(ctx)
		if err != nil {
			return nil, err
		}
		return r, nil
	}

}

func ChangeQueryParamsFn[T any, R any](fn func(ctx *handler.Ctx, t *T) (*R, error)) ChangeQueryHandler {
	return func(ctx *handler.Ctx) (any, error) {
		t := new(T)
		err := ctx.DefaultCtx.Bind().All(t)
		if err != nil {
			return nil, err
		}
		r, err := fn(ctx, t)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
}

type OperationConfig struct {
	BeforeChangeQuery ChangeQueryHandler
	AfterChangeQuery  ChangeQueryHandler
	Module            string
}

type Option func(*OperationConfig)

func WithChangeQuery(fn ChangeQueryHandler) Option {
	return func(config *OperationConfig) {
		config.BeforeChangeQuery = fn
		config.AfterChangeQuery = fn
	}
}

func WithModule(module string) Option {
	return func(config *OperationConfig) {
		config.Module = module
	}
}

func OperationLogMiddleware(options ...Option) fiber.Handler {
	conf := &OperationConfig{}
	for _, option := range options {
		option(conf)
	}
	return handler.CtxNilMiddlewareFunc(func(ctx *handler.Ctx) error {
		now := time.Now()
		var beforeChange string
		var beforeChangeData any
		if conf.BeforeChangeQuery != nil {
			function.RecFn(func() {
				var err error
				beforeChangeData, err = conf.BeforeChangeQuery(ctx)
				if err != nil {
					ctx.L().Error("BeforeChangeQuery fail", zap.Error(err))
					return
				}
				beforeChange, err = sonic.MarshalString(beforeChangeData)
				if err != nil {
					ctx.L().Error("BeforeChangeQuery marshal fail", zap.Error(err))
				}
			})
		}
		defer func() {
			nextRecoverErr := recover()
			logger := ctx.L()
			var afterChange string
			var afterChangeData any
			if conf.AfterChangeQuery != nil {
				function.RecFn(func() {
					var err error
					afterChangeData, err = conf.AfterChangeQuery(ctx)
					if err != nil {
						logger.Error("AfterChangeQuery fail", zap.Error(err))
						return
					}
					afterChange, err = sonic.MarshalString(afterChangeData)
					if err != nil {
						logger.Error("AfterChangeQuery marshal fail", zap.Error(err))
					}
				})
			}

			ip := ctx.IP()
			method := ctx.Method()
			path := ctx.Path()
			requestID := ctx.RequestID()
			requestBody := string(ctx.Request().Body())
			responseBody := string(ctx.Response().Body())
			statusCode := ctx.Response().StatusCode()
			headers, _ := sonic.MarshalString(ctx.GetReqHeaders())
			referer := ctx.Referer()
			requestURI := ctx.Request().URI().String()
			userAgent := ctx.UserAgent()
			userId := ctx.SessionInfo.Id
			username := ctx.SessionInfo.Username
			errMsg := ctx.ErrMsg
			errCode := ctx.ErrCode
			costTime := time.Since(now).Milliseconds()
			var formatChange string
			function.RecFn(func() {
				formatChange = DiffChange(beforeChangeData, afterChangeData)
			})
			coroutine.Launch(func() {
				result, err := ip_util.Client.Query(ip)
				if err != nil {
					logger.Error("Query Ip fail", zap.Error(err))
				}
				m := &models.SysOperationLog{
					RequestID:      requestID,
					Method:         method,
					Module:         conf.Module,
					Path:           path,
					Referer:        referer,
					BeforeChange:   beforeChange,
					AfterChange:    afterChange,
					FormatChange:   formatChange,
					RequestURI:     requestURI,
					RequestBody:    requestBody,
					RequestHeader:  headers,
					Response:       responseBody,
					CostTime:       costTime,
					UserID:         userId,
					Username:       username,
					ClientIP:       ip,
					StatusCode:     statusCode,
					Reason:         errMsg,
					Success:        errCode == nil,
					Location:       result.String(),
					UserAgent:      userAgent,
					BrowserName:    "",
					BrowserVersion: "",
					ClientID:       "",
					ClientName:     "",
					OSName:         "",
					OSVersion:      "",
				}
				_, err = repo.SysOperationLogRepo.Create(ctx.Context(), orm.DB(), m)
				if err != nil {
					logger.Error("SysOperationLog.Create fail", zap.Error(err))
				}
			})

			if nextRecoverErr != nil {
				panic(nextRecoverErr)
			}
		}()
		return ctx.Next()
	})
}
