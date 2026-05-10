package logic

import (
	"errors"
	"fmt"
	"testing"

	"admin/internal/auth"
	"admin/internal/fiberc/handler"
	"admin/internal/fiberc/res"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"go-common/utils/passwd"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc/mixin"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSysUserCreateEncodesPassword(t *testing.T) {
	setupSysUserTestDB(t)

	h := SysUserHandler{}
	if err := h.Create(newTestHandlerCtx(100), &ReqSysUserCreate{
		Username:     "alice",
		Nickname:     "Alice",
		Password:     "secret123",
		LanguageCode: "zh-CN",
		IsEnabled:    true,
		Remark:       "first user",
	}); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	sysUser := query.SysUser
	user, err := sysUser.
		Select(sysUser.ID, sysUser.Username, sysUser.Password, sysUser.LanguageCode, sysUser.Remark, sysUser.IsEnabled).
		Where(sysUser.Username.Eq("alice")).
		First()
	if err != nil {
		t.Fatalf("query created user failed: %v", err)
	}
	if user.Password == "secret123" {
		t.Fatal("expected stored password to be encoded")
	}
	if !passwd.Match("secret123", user.Password) {
		t.Fatal("expected stored password to match original password")
	}
	if user.LanguageCode != "zh-CN" {
		t.Fatalf("expected language code zh-CN, got %q", user.LanguageCode)
	}
	if !user.IsEnabled.IsEnabled {
		t.Fatal("expected user to be enabled")
	}
}

func TestSysUserCreateDuplicateUsername(t *testing.T) {
	setupSysUserTestDB(t)

	h := SysUserHandler{}
	createReq := &ReqSysUserCreate{
		Username:  "duplicate",
		Nickname:  "first",
		Password:  "secret123",
		IsEnabled: true,
	}
	if err := h.Create(newTestHandlerCtx(1), createReq); err != nil {
		t.Fatalf("initial Create returned error: %v", err)
	}

	err := h.Create(newTestHandlerCtx(2), &ReqSysUserCreate{
		Username:  "duplicate",
		Nickname:  "second",
		Password:  "secret456",
		IsEnabled: true,
	})
	if err == nil {
		t.Fatal("expected duplicate Create to fail")
	}

	var resp res.Response
	if !errors.As(err, &resp) {
		t.Fatalf("expected res.Response error, got %T", err)
	}
	if resp.Msg != "用户名已存在" {
		t.Fatalf("expected duplicate username message, got %q", resp.Msg)
	}
}

func TestSysUserUpdatePreservesUnsetFields(t *testing.T) {
	setupSysUserTestDB(t)

	userID := createTestSysUser(t, 1, "before", "Before", "secret123")
	h := SysUserHandler{}
	nickname := "After"
	remark := "updated"
	if err := h.Update(newTestHandlerCtx(10), &ReqSysUserUpdate{
		ID:       userID,
		Nickname: &nickname,
		Remark:   &remark,
	}); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	sysUser := query.SysUser
	user, err := sysUser.
		Select(sysUser.ID, sysUser.Username, sysUser.Nickname, sysUser.LanguageCode, sysUser.Remark, sysUser.UpdatedBy).
		Where(sysUser.ID.Eq(userID)).
		First()
	if err != nil {
		t.Fatalf("query updated user failed: %v", err)
	}
	if user.Username != "before" {
		t.Fatalf("expected username to remain %q, got %q", "before", user.Username)
	}
	if user.Nickname != "After" {
		t.Fatalf("expected nickname to be updated, got %q", user.Nickname)
	}
	if user.LanguageCode != "" {
		t.Fatalf("expected language code to remain empty, got %q", user.LanguageCode)
	}
	if user.Remark.Remark != "updated" {
		t.Fatalf("expected remark to be updated, got %q", user.Remark.Remark)
	}
	if user.UpdatedBy.UpdatedBy != 10 {
		t.Fatalf("expected updated_by to be 10, got %d", user.UpdatedBy.UpdatedBy)
	}
}

func TestSysUserUpdateMissingUser(t *testing.T) {
	setupSysUserTestDB(t)

	h := SysUserHandler{}
	nickname := "After"
	err := h.Update(newTestHandlerCtx(10), &ReqSysUserUpdate{
		ID:       999,
		Nickname: &nickname,
	})
	if err == nil {
		t.Fatal("expected Update to fail for missing user")
	}

	var resp res.Response
	if !errors.As(err, &resp) {
		t.Fatalf("expected res.Response error, got %T", err)
	}
	if resp.Msg != "用户不存在" {
		t.Fatalf("expected missing user message, got %q", resp.Msg)
	}
}

func TestSysUserDeleteSuccess(t *testing.T) {
	setupSysUserTestDB(t)

	id1 := createTestSysUser(t, 1, "u1", "U1", "secret123")
	id2 := createTestSysUser(t, 2, "u2", "U2", "secret123")
	h := SysUserHandler{}
	if err := h.Del(newTestHandlerCtx(20), &ReqSysUserBatchDelete{
		IDs: []uint64{id1, id2, id1},
	}); err != nil {
		t.Fatalf("Del returned error: %v", err)
	}

	sysUser := query.SysUser
	if _, err := sysUser.Where(sysUser.ID.Eq(id1)).First(); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected first deleted user to be hidden by soft delete, got %v", err)
	}
	if _, err := sysUser.Where(sysUser.ID.Eq(id2)).First(); !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected second deleted user to be hidden by soft delete, got %v", err)
	}
}

func TestSysUserDeleteMissingUser(t *testing.T) {
	setupSysUserTestDB(t)

	id := createTestSysUser(t, 1, "u1", "U1", "secret123")
	h := SysUserHandler{}
	err := h.Del(newTestHandlerCtx(20), &ReqSysUserBatchDelete{
		IDs: []uint64{id, 999},
	})
	if err == nil {
		t.Fatal("expected Del to fail when any user is missing")
	}

	var resp res.Response
	if !errors.As(err, &resp) {
		t.Fatalf("expected res.Response error, got %T", err)
	}
	if resp.Msg != "用户不存在" {
		t.Fatalf("expected missing user message, got %q", resp.Msg)
	}
}

func TestSysUserListDefaultOrderByIDDESC(t *testing.T) {
	setupSysUserTestDB(t)

	firstID := createTestSysUser(t, 1, "first", "First", "secret123")
	secondID := createTestSysUser(t, 2, "second", "Second", "secret123")
	h := SysUserHandler{}
	page := uint32(1)
	pageSize := uint32(10)
	result, err := h.List(newTestHandlerCtx(0), &v1.PagingRequest{
		Page:     &page,
		PageSize: &pageSize,
	})
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(result.Items) != 2 {
		t.Fatalf("expected 2 users, got %d", len(result.Items))
	}
	if result.Items[0].ID != secondID || result.Items[1].ID != firstID {
		t.Fatalf("expected users ordered by id desc, got [%d %d]", result.Items[0].ID, result.Items[1].ID)
	}
}

func setupSysUserTestDB(t *testing.T) {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		t.Fatalf("open sqlite db failed: %v", err)
	}
	if err := db.AutoMigrate(&models.SysRole{}, &models.SysUser{}, &models.SysUserRole{}); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}
	query.SetDefault(db)
}

func createTestSysUser(t *testing.T, id uint64, username, nickname, rawPassword string) uint64 {
	t.Helper()

	encodedPwd, err := passwd.Encode(rawPassword)
	if err != nil {
		t.Fatalf("encode password failed: %v", err)
	}

	user := &models.SysUser{
		OperatorID: mixin.OperatorID{
			CreatedBy: mixin.CreatedBy{CreatedBy: 1},
			UpdatedBy: mixin.UpdatedBy{UpdatedBy: 1},
		},
		IsEnabled: mixin.IsEnabled{IsEnabled: true},
		AutoIncrementID: mixin.AutoIncrementID{
			ID: id,
		},
		Username: username,
		Nickname: nickname,
		Password: encodedPwd,
	}
	if err := query.SysUser.Create(user); err != nil {
		t.Fatalf("create test user failed: %v", err)
	}
	return id
}

func newTestHandlerCtx(userID uint64) *handler.Ctx {
	return &handler.Ctx{
		SessionInfo: &auth.SessionInfo{Id: userID},
	}
}
