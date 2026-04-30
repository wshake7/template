package datapermission

import (
	"testing"

	"admin/internal/services/orm/models"
	"gorm.io/datatypes"
)

func TestPermissionIncludesActionAll(t *testing.T) {
	permission := &models.SysDataPermission{
		Action: datatypes.JSON(`["all"]`),
	}

	for _, action := range []permissionAction{actionRead, actionWrite, actionDelete} {
		if !permissionIncludesAction(permission, action) {
			t.Fatalf("expected all action to include %s", action)
		}
	}
}

func TestPermissionIncludesActionSpecific(t *testing.T) {
	permission := &models.SysDataPermission{
		Action: datatypes.JSON(`["read"]`),
	}

	if !permissionIncludesAction(permission, actionRead) {
		t.Fatal("expected read action to include read")
	}
	if permissionIncludesAction(permission, actionWrite) {
		t.Fatal("did not expect read action to include write")
	}
}
