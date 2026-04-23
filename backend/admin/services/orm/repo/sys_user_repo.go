package repo

import (
	"admin/services/orm/models"
	"admin/services/orm/query"
	"go-common/mapper"
	"orm-crud/gorm"
)

type sysUserRepo[T, R any] struct {
	*gorm.Repository[T, R]
}

var SysUserRepo *sysUserRepo[models.SysUser, models.SysUser]

func init() {
	repository := gorm.NewRepository(mapper.NewCopierMapper[models.SysUser, models.SysUser]())
	SysUserRepo = &sysUserRepo[models.SysUser, models.SysUser]{
		Repository: repository,
	}
}

func (sysUserRepo[T, R]) ChangePwd(id uint64, encodePwd string) error {
	sysUser := query.SysUser
	_, err := sysUser.Where(sysUser.ID.Eq(id)).Update(sysUser.Password, encodePwd)
	return err
}
