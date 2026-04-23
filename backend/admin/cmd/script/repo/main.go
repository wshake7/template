package main

import (
	"admin/config"
	"admin/services/repo/models"
	"admin/services/repo/query"
	"go-common/utils/passwd"
	"go-common/viperc"
	"go.uber.org/zap"
	"gorm.io/gen"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
	gormCrud "orm-crud/gorm"
)

func main() {
	var conf config.Config
	_, err := viperc.ParseFile("./etc/config.yaml", &conf)
	if err != nil {
		panic(err)
	}
	var options []gormCrud.Option
	zapLogger := zap.L()
	logger := zapgorm2.New(zapLogger)
	logger.SetAsDefault()
	options = append(options,
		gormCrud.WithLogger(zapLogger.Sugar()),
		gormCrud.WithGormConfig(&gorm.Config{Logger: logger}),
		gormCrud.WithDriverName(conf.Repo.DriverName),
		gormCrud.WithDSN(conf.Repo.DataSource),
		gormCrud.WithEnableTrace(true),
		gormCrud.WithEnableMetrics(true),
	)

	if conf.Repo.IsAutoMigrate {
		options = append(options, gormCrud.WithAutoMigrate(models.Models...))
	}

	client, err := gormCrud.NewClient(options...)
	if err != nil {
		panic(err)
	}

	query.SetDefault(client.DB)
	codeGenCode(client.DB, models.Models)
	genUserAdd()
}

// 从数据库生成代码
func dbGenCode(db *gorm.DB) {
	cfg := gen.Config{
		OutPath:           "./services/repo/query",
		OutFile:           "",
		ModelPkgPath:      "",
		WithUnitTest:      false,
		FieldNullable:     true,
		FieldCoverable:    true,
		FieldSignable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
		Mode:              gen.WithDefaultQuery | gen.WithoutContext | gen.WithQueryInterface,
	}
	// 驼峰
	//cfg.WithJSONTagNameStrategy(func(columnName string) (tagContent string) {
	//	return strcase.LowerCamelCase(columnName)
	//})
	g := gen.NewGenerator(cfg)
	g.UseDB(db)
	m := g.GenerateAllTable()
	g.ApplyBasic(m...)
	g.Execute()
}

func codeGenCode(db *gorm.DB, models []any) {
	cfg := gen.Config{
		OutPath:           "./services/repo/query",
		OutFile:           "",
		ModelPkgPath:      "",
		WithUnitTest:      false,
		FieldNullable:     true,
		FieldCoverable:    true,
		FieldSignable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
		Mode:              gen.WithDefaultQuery | gen.WithoutContext | gen.WithQueryInterface,
	}
	// 驼峰
	//cfg.WithJSONTagNameStrategy(func(columnName string) (tagContent string) {
	//	return strcase.LowerCamelCase(columnName)
	//})
	g := gen.NewGenerator(cfg)
	g.UseDB(db)
	g.ApplyBasic(models...)
	g.Execute()
}

func genUserAdd() {
	sysUser := query.SysUser
	pwd, _ := passwd.Encode("123456")
	_ = sysUser.Create(&models.SysUser{Username: "admin", Password: pwd})
}
