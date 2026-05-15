package logic

import (
	"fmt"
	"testing"

	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	v1 "orm-crud/api/gen/go/pagination/v1"
	"orm-crud/gormc/mixin"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestSysResourceMenuCreateSavesAPILinks(t *testing.T) {
	setupSysResourceMenuTestDB(t)

	api1 := createTestResourceAPI(t, 1)
	api2 := createTestResourceAPI(t, 2)
	parentID := createTestResourceMenu(t, 1)
	h := SysResourceMenuHandler{}
	if err := h.Create(newTestHandlerCtx(7), &ReqResourceMenuCreate{
		MenuType:  MenuTypeButton,
		ParentID:  &parentID,
		Name:      "button",
		IsEnabled: true,
		ApiIDs:    []uint64{api1, api1, api2},
	}); err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	var links []*models.SysResourceMenuApi
	if err := query.SysResourceMenu.UnderlyingDB().Session(&gorm.Session{NewDB: true}).Model(&models.SysResourceMenuApi{}).Where("deleted_at = 0").Order("api_id asc").Find(&links).Error; err != nil {
		t.Fatalf("query menu api links failed: %v", err)
	}
	if len(links) != 2 || links[0].ApiID != api1 || links[1].ApiID != api2 {
		t.Fatalf("expected deduplicated api links [%d %d], got %#v", api1, api2, links)
	}
}

func TestSysResourceMenuUpdateReplacesAndClearsAPILinks(t *testing.T) {
	setupSysResourceMenuTestDB(t)

	menuID := createTestValidButtonMenu(t, 1)
	api1 := createTestResourceAPI(t, 1)
	api2 := createTestResourceAPI(t, 2)
	h := SysResourceMenuHandler{}
	apis := []uint64{api1, api2}
	if err := h.Update(newTestHandlerCtx(7), &ReqResourceMenuUpdate{
		ID:     menuID,
		ApiIDs: &apis,
	}); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	apis = []uint64{api2}
	if err := h.Update(newTestHandlerCtx(7), &ReqResourceMenuUpdate{
		ID:     menuID,
		ApiIDs: &apis,
	}); err != nil {
		t.Fatalf("replacement Update returned error: %v", err)
	}
	var links []*models.SysResourceMenuApi
	if err := query.SysResourceMenu.UnderlyingDB().Session(&gorm.Session{NewDB: true}).Model(&models.SysResourceMenuApi{}).Where("deleted_at = 0").Find(&links).Error; err != nil {
		t.Fatalf("query menu api links failed: %v", err)
	}
	if len(links) != 1 || links[0].ApiID != api2 {
		t.Fatalf("expected only api2 after replacement, got %#v", links)
	}

	apis = []uint64{}
	if err := h.Update(newTestHandlerCtx(7), &ReqResourceMenuUpdate{
		ID:     menuID,
		ApiIDs: &apis,
	}); err != nil {
		t.Fatalf("clear Update returned error: %v", err)
	}
	links = nil
	if err := query.SysResourceMenu.UnderlyingDB().Session(&gorm.Session{NewDB: true}).Model(&models.SysResourceMenuApi{}).Where("deleted_at = 0").Find(&links).Error; err != nil {
		t.Fatalf("query menu api links after clear failed: %v", err)
	}
	if len(links) != 0 {
		t.Fatalf("expected api links to be cleared, got %#v", links)
	}
}

func TestSysResourceMenuAPILinksOnlyForMenuAndButton(t *testing.T) {
	setupSysResourceMenuTestDB(t)

	apiID := createTestResourceAPI(t, 1)
	h := SysResourceMenuHandler{}
	if err := h.Create(newTestHandlerCtx(7), &ReqResourceMenuCreate{
		MenuType:  MenuTypeCatalog,
		Name:      "catalog",
		IsEnabled: true,
		ApiIDs:    []uint64{apiID},
	}); err != nil {
		t.Fatalf("Create catalog returned error: %v", err)
	}
	var count int64
	if err := query.SysResourceMenu.UnderlyingDB().Session(&gorm.Session{NewDB: true}).Model(&models.SysResourceMenuApi{}).Where("deleted_at = 0").Count(&count).Error; err != nil {
		t.Fatalf("count menu api links failed: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected catalog api links to be ignored, got %d", count)
	}
}

func TestSysResourceMenuListReturnsAPILinks(t *testing.T) {
	setupSysResourceMenuTestDB(t)

	menuID := createTestValidButtonMenu(t, 1)
	apiID := createTestResourceAPI(t, 1)
	apis := []uint64{apiID}
	h := SysResourceMenuHandler{}
	if err := h.Update(newTestHandlerCtx(7), &ReqResourceMenuUpdate{
		ID:     menuID,
		ApiIDs: &apis,
	}); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}

	noPaging := true
	result, err := h.List(newTestHandlerCtx(7), &v1.PagingRequest{NoPaging: &noPaging})
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	node := findRespResourceMenu(result.Items, menuID)
	if node == nil || len(node.ApiIDs) != 1 || node.ApiIDs[0] != apiID {
		t.Fatalf("expected api id in list response, got %#v", result.Items)
	}
}

func findRespResourceMenu(items []*RespSysResourceMenu, id uint64) *RespSysResourceMenu {
	for _, item := range items {
		if item.ID == id {
			return item
		}
		if child := findRespResourceMenu(item.Children, id); child != nil {
			return child
		}
	}
	return nil
}

func createTestValidButtonMenu(t *testing.T, id uint64) uint64 {
	t.Helper()

	parentID := createTestResourceMenu(t, id+1000)
	item := &models.SysResourceMenu{
		AutoIncrementID: mixin.AutoIncrementID{ID: id},
		OperatorID: mixin.OperatorID{
			CreatedBy: mixin.CreatedBy{CreatedBy: 1},
			UpdatedBy: mixin.UpdatedBy{UpdatedBy: 1},
		},
		SortOrder: mixin.SortOrder{SortOrder: int32(id)},
		IsEnabled: mixin.IsEnabled{IsEnabled: true},
		MenuType:  MenuTypeButton,
		Name:      fmt.Sprintf("button-%d", id),
		ParentID:  &parentID,
	}
	if err := query.SysResourceMenu.Create(item); err != nil {
		t.Fatalf("create test button menu failed: %v", err)
	}
	return id
}

func setupSysResourceMenuTestDB(t *testing.T) {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		t.Fatalf("open sqlite db failed: %v", err)
	}
	if err := db.AutoMigrate(
		&models.SysResourceMenu{},
		&models.SysResourceApi{},
		&models.SysResourceMenuApi{},
	); err != nil {
		t.Fatalf("auto migrate failed: %v", err)
	}
	query.SetDefault(db)
}
