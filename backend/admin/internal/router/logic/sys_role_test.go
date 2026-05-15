package logic

import (
	"errors"
	"fmt"
	"testing"

	"admin/internal/fiberc/res"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc/mixin"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSysRoleCreateSuccess(t *testing.T) {
	setupSysRoleTestDB(t)

	h := SysRoleHandler{}
	if err := h.Create(newTestHandlerCtx(100), &ReqSysRoleCreate{
		Name:      "Admin",
		Code:      "admin",
		IsEnabled: true,
		Remark:    "admin role",
	}); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	sysRole := query.SysRole
	role, err := sysRole.
		Select(sysRole.ID, sysRole.Name, sysRole.Code, sysRole.IsEnabled, sysRole.Remark, sysRole.CreatedBy).
		Where(sysRole.Code.Eq("admin")).
		First()
	if err != nil {
		t.Fatalf("query created role failed: %v", err)
	}
	if role.Name != "Admin" || role.Code != "admin" {
		t.Fatalf("unexpected role values: %#v", role)
	}
	if !role.IsEnabled.IsEnabled {
		t.Fatal("expected role to be enabled")
	}
	if role.CreatedBy.CreatedBy != 100 {
		t.Fatalf("expected created_by to be 100, got %d", role.CreatedBy.CreatedBy)
	}
}

func TestSysRoleCreateDuplicateCode(t *testing.T) {
	setupSysRoleTestDB(t)

	h := SysRoleHandler{}
	req := &ReqSysRoleCreate{Name: "Admin", Code: "admin", IsEnabled: true}
	if err := h.Create(newTestHandlerCtx(1), req); err != nil {
		t.Fatalf("initial Create returned error: %v", err)
	}
	err := h.Create(newTestHandlerCtx(2), &ReqSysRoleCreate{Name: "Other", Code: "admin", IsEnabled: true})
	if err == nil {
		t.Fatal("expected duplicate Create to fail")
	}
	var resp res.Response
	if !errors.As(err, &resp) {
		t.Fatalf("expected res.Response error, got %T", err)
	}
	if resp.Msg != "角色标识已存在" {
		t.Fatalf("expected duplicate code message, got %q", resp.Msg)
	}
}

func TestSysRoleUpdatePreservesUnsetFields(t *testing.T) {
	setupSysRoleTestDB(t)

	roleID := createTestSysRole(t, 1, "before", "before", nil)
	h := SysRoleHandler{}
	name := "After"
	remark := "updated"
	if err := h.Update(newTestHandlerCtx(10), &ReqSysRoleUpdate{
		ID:     roleID,
		Name:   &name,
		Remark: &remark,
	}); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	sysRole := query.SysRole
	role, err := sysRole.
		Select(sysRole.ID, sysRole.Name, sysRole.Code, sysRole.Remark, sysRole.UpdatedBy).
		Where(sysRole.ID.Eq(roleID)).
		First()
	if err != nil {
		t.Fatalf("query updated role failed: %v", err)
	}
	if role.Name != "After" {
		t.Fatalf("expected name to be updated, got %q", role.Name)
	}
	if role.Code != "before" {
		t.Fatalf("expected code to remain before, got %q", role.Code)
	}
	if role.Remark.Remark != "updated" {
		t.Fatalf("expected remark to be updated, got %q", role.Remark.Remark)
	}
	if role.UpdatedBy.UpdatedBy != 10 {
		t.Fatalf("expected updated_by to be 10, got %d", role.UpdatedBy.UpdatedBy)
	}
}

func TestSysRoleUpdateRejectsSelfOrDescendantParent(t *testing.T) {
	setupSysRoleTestDB(t)

	rootID := createTestSysRole(t, 1, "root", "root", nil)
	childID := createTestSysRole(t, 2, "child", "child", &rootID)
	h := SysRoleHandler{}

	selfParent := rootID
	if err := h.Update(newTestHandlerCtx(10), &ReqSysRoleUpdate{ID: rootID, ParentID: &selfParent}); err == nil {
		t.Fatal("expected self parent update to fail")
	}

	descendantParent := childID
	if err := h.Update(newTestHandlerCtx(10), &ReqSysRoleUpdate{ID: rootID, ParentID: &descendantParent}); err == nil {
		t.Fatal("expected descendant parent update to fail")
	}
}

func TestSysRoleDeleteRejectsBindings(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T, roleID uint64)
		want  string
	}{
		{
			name: "child role",
			setup: func(t *testing.T, roleID uint64) {
				createTestSysRole(t, 2, "child", "child", &roleID)
			},
			want: "存在子角色，不能删除",
		},
		{
			name: "user role",
			setup: func(t *testing.T, roleID uint64) {
				createTestSysUser(t, 1, "u1", "U1", "secret123")
				if err := query.SysUserRole.Create(&models.SysUserRole{UserID: 1, RoleID: roleID}); err != nil {
					t.Fatalf("create user role failed: %v", err)
				}
			},
			want: "角色已绑定用户，不能删除",
		},
		{
			name: "menu permission",
			setup: func(t *testing.T, roleID uint64) {
				menuID := createTestResourceMenu(t, 1)
				if err := query.SysRoleMenu.Create(&models.SysRoleMenu{RoleID: roleID, MenuID: menuID}); err != nil {
					t.Fatalf("create role menu failed: %v", err)
				}
			},
			want: "角色已绑定菜单权限，不能删除",
		},
		{
			name: "api permission",
			setup: func(t *testing.T, roleID uint64) {
				apiID := createTestResourceAPI(t, 1)
				if err := query.SysRoleApi.Create(&models.SysRoleApi{RoleID: roleID, ApiID: apiID}); err != nil {
					t.Fatalf("create role api failed: %v", err)
				}
			},
			want: "角色已绑定API权限，不能删除",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupSysRoleTestDB(t)
			roleID := createTestSysRole(t, 1, "role", "role", nil)
			tt.setup(t, roleID)

			h := SysRoleHandler{}
			err := h.Del(newTestHandlerCtx(20), &ReqSysRoleBatchDelete{IDs: []uint64{roleID}})
			if err == nil {
				t.Fatal("expected Del to fail")
			}
			var resp res.Response
			if !errors.As(err, &resp) {
				t.Fatalf("expected res.Response error, got %T", err)
			}
			if resp.Msg != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, resp.Msg)
			}
		})
	}
}

func TestSysRoleSavePermissionsReplacesAndDeduplicates(t *testing.T) {
	setupSysRoleTestDB(t)

	roleID := createTestSysRole(t, 1, "role", "role", nil)
	menu1 := createTestResourceMenu(t, 1)
	menu2 := createTestResourceMenu(t, 2)
	api1 := createTestResourceAPI(t, 1)
	api2 := createTestResourceAPI(t, 2)

	h := SysRoleHandler{}
	if err := h.SavePermissions(newTestHandlerCtx(9), &ReqSysRolePermissionSave{
		ID:      roleID,
		MenuIDs: []uint64{menu1, menu1, menu2},
		ApiIDs:  []uint64{api1, api1, api2},
	}); err != nil {
		t.Fatalf("SavePermissions returned error: %v", err)
	}

	result, err := h.Permissions(newTestHandlerCtx(9), &ReqSysRolePermissionQuery{ID: roleID})
	if err != nil {
		t.Fatalf("Permissions returned error: %v", err)
	}
	if len(result.MenuIDs) != 2 || len(result.ApiIDs) != 2 {
		t.Fatalf("expected deduplicated permissions, got menus=%v apis=%v", result.MenuIDs, result.ApiIDs)
	}

	if err := h.SavePermissions(newTestHandlerCtx(9), &ReqSysRolePermissionSave{
		ID:      roleID,
		MenuIDs: []uint64{menu2},
		ApiIDs:  []uint64{api2},
	}); err != nil {
		t.Fatalf("replacement SavePermissions returned error: %v", err)
	}
	result, err = h.Permissions(newTestHandlerCtx(9), &ReqSysRolePermissionQuery{ID: roleID})
	if err != nil {
		t.Fatalf("Permissions after replacement returned error: %v", err)
	}
	if len(result.MenuIDs) != 1 || result.MenuIDs[0] != menu2 {
		t.Fatalf("expected only menu2 after replacement, got %v", result.MenuIDs)
	}
	if len(result.ApiIDs) != 1 || result.ApiIDs[0] != api2 {
		t.Fatalf("expected only api2 after replacement, got %v", result.ApiIDs)
	}
}

func TestSysRoleSavePermissionsRejectsMissingResource(t *testing.T) {
	setupSysRoleTestDB(t)

	roleID := createTestSysRole(t, 1, "role", "role", nil)
	h := SysRoleHandler{}
	err := h.SavePermissions(newTestHandlerCtx(9), &ReqSysRolePermissionSave{
		ID:      roleID,
		MenuIDs: []uint64{999},
	})
	if err == nil {
		t.Fatal("expected SavePermissions to reject missing menu")
	}
}

func TestSysRoleListDefaultOrderByIDDESC(t *testing.T) {
	setupSysRoleTestDB(t)

	firstID := createTestSysRole(t, 1, "first", "first", nil)
	secondID := createTestSysRole(t, 2, "second", "second", nil)
	page := uint32(1)
	pageSize := uint32(10)
	h := SysRoleHandler{}
	result, err := h.List(newTestHandlerCtx(0), &v1.PagingRequest{
		Page:     &page,
		PageSize: &pageSize,
	})
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(result.Items) != 2 {
		t.Fatalf("expected 2 roles, got %d", len(result.Items))
	}
	if result.Items[0].ID != secondID || result.Items[1].ID != firstID {
		t.Fatalf("expected roles ordered by id desc, got [%d %d]", result.Items[0].ID, result.Items[1].ID)
	}
}

func setupSysRoleTestDB(t *testing.T) {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		t.Fatalf("open sqlite db failed: %v", err)
	}
	if err := db.AutoMigrate(
		&models.SysRole{},
		&models.SysUser{},
		&models.SysUserRole{},
		&models.SysResourceMenu{},
		&models.SysResourceApi{},
		&models.SysResourceMenuApi{},
		&models.SysRoleMenu{},
		&models.SysRoleApi{},
	); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}
	query.SetDefault(db)
}

func createTestSysRole(t *testing.T, id uint64, name, code string, parentID *uint64) uint64 {
	t.Helper()

	role := &models.SysRole{
		AutoIncrementID: mixin.AutoIncrementID{ID: id},
		OperatorID: mixin.OperatorID{
			CreatedBy: mixin.CreatedBy{CreatedBy: 1},
			UpdatedBy: mixin.UpdatedBy{UpdatedBy: 1},
		},
		IsEnabled: mixin.IsEnabled{IsEnabled: true},
		Name:      name,
		Code:      code,
		ParentID:  parentID,
	}
	if err := query.SysRole.Create(role); err != nil {
		t.Fatalf("create test role failed: %v", err)
	}
	return id
}

func createTestResourceMenu(t *testing.T, id uint64) uint64 {
	t.Helper()

	item := &models.SysResourceMenu{
		AutoIncrementID: mixin.AutoIncrementID{ID: id},
		OperatorID: mixin.OperatorID{
			CreatedBy: mixin.CreatedBy{CreatedBy: 1},
			UpdatedBy: mixin.UpdatedBy{UpdatedBy: 1},
		},
		SortOrder: mixin.SortOrder{SortOrder: int32(id)},
		IsEnabled: mixin.IsEnabled{IsEnabled: true},
		MenuType:  MenuTypeMenu,
		Path:      fmt.Sprintf("/menu/%d", id),
		Name:      fmt.Sprintf("menu-%d", id),
		Component: fmt.Sprintf("/menu/%d.tsx", id),
	}
	if err := query.SysResourceMenu.Create(item); err != nil {
		t.Fatalf("create test resource menu failed: %v", err)
	}
	return id
}

func createTestResourceAPI(t *testing.T, id uint64) uint64 {
	t.Helper()

	item := &models.SysResourceApi{
		AutoIncrementID: mixin.AutoIncrementID{ID: id},
		OperatorID: mixin.OperatorID{
			CreatedBy: mixin.CreatedBy{CreatedBy: 1},
			UpdatedBy: mixin.UpdatedBy{UpdatedBy: 1},
		},
		SortOrder: mixin.SortOrder{SortOrder: int32(id)},
		IsEnabled: mixin.IsEnabled{IsEnabled: true},
		Method:    HttpMethodGet,
		Path:      fmt.Sprintf("/api/test/%d", id),
	}
	if err := query.SysResourceApi.Create(item); err != nil {
		t.Fatalf("create test resource api failed: %v", err)
	}
	return id
}
