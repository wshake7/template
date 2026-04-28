package datapermission

import (
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"errors"
	"maps"
	v1 "orm-crud/api/gen/go/pagination/v1"
	paginationFilter "orm-crud/pagination/filter"
	"slices"

	"github.com/bytedance/sonic"
	"gorm.io/gen/field"
)

type FilterAction string

const (
	actionRead   = FilterAction("read")
	actionWrite  = FilterAction("write")
	actionDelete = FilterAction("delete")
)

type Subject struct {
	Type string
	ID   uint64
}

type permissionFilter struct {
	query      map[string]any
	fullAccess bool
}

func ApplyPageFilter(req *v1.PagingRequest, resourceTable string, subjects []Subject) error {
	if req == nil {
		return errors.New("paging request is nil")
	}

	permissionExpr, err := BuildFilterExpr(resourceTable, actionRead, subjects)
	if err != nil {
		return err
	}
	if permissionExpr == nil {
		return nil
	}

	currentExpr, err := paginationFilter.ConvertFilterByPagingRequest(req)
	if err != nil {
		return err
	}
	if currentExpr == nil {
		req.FilteringType = &v1.PagingRequest_FilterExpr{FilterExpr: permissionExpr}
		return nil
	}

	req.FilteringType = &v1.PagingRequest_FilterExpr{FilterExpr: &v1.FilterExpr{
		Type:   v1.ExprType_AND,
		Groups: []*v1.FilterExpr{currentExpr, permissionExpr},
	}}
	return nil
}

func BuildFilterExpr(resourceTable string, action FilterAction, subjects []Subject) (*v1.FilterExpr, error) {
	queryMap, err := buildActionFilterMap(resourceTable, action, subjects)
	if err != nil {
		return nil, err
	}
	if queryMap == nil {
		return nil, nil
	}
	return filterMapToExpr(queryMap)
}

func BuildSubjects(userID uint64, roleIDs []uint64) []Subject {
	subjects := []Subject{
		{Type: "USER", ID: userID},
		{Type: "ANY_USER", ID: 0},
		{Type: "ANY_ROLE", ID: 0},
	}
	for _, roleID := range roleIDs {
		subjects = append(subjects, Subject{Type: "ROLE", ID: roleID})
	}
	return subjects
}

func buildActionFilterMap(resourceTable string, action FilterAction, subjects []Subject) (map[string]any, error) {
	permissions, err := findCandidatePermissions(resourceTable, subjects)
	if err != nil {
		return nil, err
	}

	allowBranches := make([]map[string]any, 0, len(permissions))
	for _, permission := range permissions {
		if !permissionAllowsAction(permission, action) {
			continue
		}

		filter, ok := buildPermissionFilter(permission)
		if !ok {
			continue
		}
		if filter.fullAccess {
			return nil, nil
		}
		allowBranches = append(allowBranches, filter.query)
	}

	return joinAllowBranches(allowBranches), nil
}

func findCandidatePermissions(resourceTable string, subjects []Subject) ([]*models.SysDataPermission, error) {
	subjectPredicates := buildSubjectPredicates(subjects)
	if len(subjectPredicates) == 0 {
		return nil, nil
	}

	sysDataPermission := query.SysDataPermission
	return sysDataPermission.
		Where(
			sysDataPermission.IsEnabled.Is(true),
			sysDataPermission.ResourceTable.Eq(resourceTable),
			field.Or(subjectPredicates...),
		).
		Find()
}

func buildSubjectPredicates(subjects []Subject) []field.Expr {
	predicates := make([]field.Expr, 0, len(subjects))
	for _, subject := range subjects {
		predicates = append(predicates, field.And(
			query.SysDataPermission.SubjectType.Eq(subject.Type),
			query.SysDataPermission.SubjectID.Eq(subject.ID),
		))
	}
	return predicates
}

func permissionAllowsAction(permission *models.SysDataPermission, action FilterAction) bool {
	if len(permission.Action) == 0 {
		return false
	}

	var actions []FilterAction
	if err := sonic.Unmarshal(permission.Action, &actions); err != nil {
		return false
	}
	return slices.Contains(actions, action)
}

func buildPermissionFilter(permission *models.SysDataPermission) (permissionFilter, bool) {
	condition := map[string]any{}
	maps.Copy(condition, permission.Conditions)

	switch permission.ScopeType {
	case "all":
		return permissionFilter{query: condition, fullAccess: len(condition) == 0}, true
	case "custom":
		return permissionFilter{query: condition, fullAccess: len(condition) == 0}, true
	case "include":
		values, ok := parseScopeValues(permission)
		if !ok {
			return permissionFilter{}, false
		}
		condition[permission.ScopeField+"__in"] = values
		return permissionFilter{query: condition}, true
	case "exclude":
		values, ok := parseScopeValues(permission)
		if !ok {
			return permissionFilter{}, false
		}
		condition[permission.ScopeField+"__not_in"] = values
		return permissionFilter{query: condition}, true
	default:
		return permissionFilter{}, false
	}
}

func parseScopeValues(permission *models.SysDataPermission) ([]any, bool) {
	if len(permission.ScopeValues) == 0 {
		return nil, false
	}

	var values []any
	if err := sonic.Unmarshal(permission.ScopeValues, &values); err != nil {
		return nil, false
	}
	if len(values) == 0 {
		return nil, false
	}
	return values, true
}

func joinAllowBranches(allowBranches []map[string]any) map[string]any {
	if len(allowBranches) == 0 {
		return map[string]any{"id": 0}
	}
	if len(allowBranches) == 1 {
		return allowBranches[0]
	}
	return map[string]any{"$or": allowBranches}
}

func filterMapToExpr(queryMap map[string]any) (*v1.FilterExpr, error) {
	queryBytes, err := sonic.Marshal(queryMap)
	if err != nil {
		return nil, err
	}
	return paginationFilter.ConvertFilterByPagingRequest(&v1.PagingRequest{
		FilteringType: &v1.PagingRequest_Query{Query: string(queryBytes)},
	})
}
