package datapermission

import (
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"encoding/json"
	"errors"
	"maps"
	v1 "orm-crud/api/gen/go/pagination/v1"
	paginationFilter "orm-crud/pagination/filter"
	"slices"

	"gorm.io/gen/field"
)

const (
	effectAllow = "allow"
	actionRead  = "read"
)

type Subject struct {
	Type string
	ID   uint64
}

func ApplyPageFilter(req *v1.PagingRequest, resourceTable string, subjects []Subject) error {
	if req == nil {
		return errors.New("paging request is nil")
	}

	permissionExpr, err := CommonFilterExpr(resourceTable, subjects)
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

func CommonFilterExpr(resourceTable string, subjects []Subject) (*v1.FilterExpr, error) {
	filterJSON, err := buildReadFilterJSON(resourceTable, subjects)
	if err != nil {
		return nil, err
	}
	if filterJSON == "" {
		return nil, nil
	}
	return filterJSONToExpr(filterJSON)
}

func BuildSessionSubjects(userID uint64, roleIDs []uint64) []Subject {
	subjects := []Subject{
		{Type: "USER", ID: userID},
		{Type: "ANY_USER", ID: 0},
	}
	if len(roleIDs) > 0 {
		subjects = append(subjects, Subject{Type: "ANY_ROLE", ID: 0})
	}
	for _, roleID := range roleIDs {
		subjects = append(subjects, Subject{Type: "ROLE", ID: roleID})
	}
	return subjects
}

func buildReadFilterJSON(resourceTable string, subjects []Subject) (string, error) {
	permissions, err := findReadCandidatePermissions(resourceTable, subjects)
	if err != nil {
		return "", err
	}

	subjectSet := buildSubjectSet(subjects)
	allowBranches := make([]any, 0)
	for _, permission := range permissions {
		if !permissionAllowsSubject(permission, subjectSet) {
			continue
		}
		if !permissionAllowsAction(permission, actionRead) {
			continue
		}

		condition, allScope, ok := buildPermissionFilterMap(permission)
		if !ok {
			continue
		}
		if allScope {
			return "", nil
		}
		allowBranches = append(allowBranches, map[string]any{
			"$and": []any{condition},
		})
	}

	return joinAllowBranches(allowBranches)
}

func findReadCandidatePermissions(resourceTable string, subjects []Subject) ([]*models.SysDataPermission, error) {
	subjectPredicates := buildSubjectPredicates(subjects)
	if len(subjectPredicates) == 0 {
		return nil, nil
	}

	sysDataPermission := query.SysDataPermission
	return sysDataPermission.
		Where(
			sysDataPermission.IsEnabled.Is(true),
			sysDataPermission.ResourceTable.Eq(resourceTable),
			field.Or(sysDataPermission.SubjectEffect.Eq(""), sysDataPermission.SubjectEffect.Eq(effectAllow)),
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

func buildSubjectSet(subjects []Subject) map[Subject]struct{} {
	subjectSet := make(map[Subject]struct{}, len(subjects))
	for _, subject := range subjects {
		subjectSet[subject] = struct{}{}
	}
	return subjectSet
}

func permissionAllowsSubject(permission *models.SysDataPermission, subjectSet map[Subject]struct{}) bool {
	if permission.SubjectEffect != "" && permission.SubjectEffect != effectAllow {
		return false
	}
	_, ok := subjectSet[Subject{
		Type: permission.SubjectType,
		ID:   permission.SubjectID,
	}]
	return ok
}

func permissionAllowsAction(permission *models.SysDataPermission, action string) bool {
	if len(permission.Action) == 0 {
		return false
	}

	var actions []string
	if err := json.Unmarshal(permission.Action, &actions); err != nil {
		return false
	}
	return slices.Contains(actions, action)
}

func buildPermissionFilterMap(permission *models.SysDataPermission) (map[string]any, bool, bool) {
	condition := map[string]any{}
	maps.Copy(condition, permission.Conditions)

	switch permission.ScopeType {
	case "all":
		return condition, len(condition) == 0, true
	case "custom":
		return condition, len(condition) == 0, true
	case "include":
		values, ok := parseScopeValues(permission)
		if !ok {
			return nil, false, false
		}
		condition[permission.ScopeField+"__in"] = values
		return condition, false, true
	case "exclude":
		values, ok := parseScopeValues(permission)
		if !ok {
			return nil, false, false
		}
		condition[permission.ScopeField+"__not_in"] = values
		return condition, false, true
	default:
		return nil, false, false
	}
}

func parseScopeValues(permission *models.SysDataPermission) ([]any, bool) {
	if len(permission.ScopeValues) == 0 {
		return nil, false
	}

	var values []any
	if err := json.Unmarshal(permission.ScopeValues, &values); err != nil {
		return nil, false
	}
	if len(values) == 0 {
		return nil, false
	}
	return values, true
}

func joinAllowBranches(allowBranches []any) (string, error) {
	if len(allowBranches) == 0 {
		return mapToFilterJSON(map[string]any{"id": 0})
	}
	if len(allowBranches) == 1 {
		if branch, ok := allowBranches[0].(map[string]any); ok {
			return mapToFilterJSON(branch)
		}
	}
	return mapToFilterJSON(map[string]any{"$or": allowBranches})
}

func mapToFilterJSON(queryMap map[string]any) (string, error) {
	queryBytes, err := json.Marshal(queryMap)
	if err != nil {
		return "", err
	}
	return string(queryBytes), nil
}

func filterJSONToExpr(queryString string) (*v1.FilterExpr, error) {
	return paginationFilter.ConvertFilterByPagingRequest(&v1.PagingRequest{
		FilteringType: &v1.PagingRequest_Query{Query: queryString},
	})
}
