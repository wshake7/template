package main

import (
	"admin/config"
	"admin/services/orm/models"
	"admin/services/orm/query"
	"fmt"
	"go-common/utils/passwd"
	"go-common/viperc"
	"go.uber.org/zap"
	"gorm.io/gen"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
	gormCrud "orm-crud/gorm"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
	"unicode"
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

func genUserAdd() {
	sysUser := query.SysUser
	pwd, _ := passwd.Encode("123456")
	_ = sysUser.Create(&models.SysUser{Username: "admin", Password: pwd})
}

// 从数据库生成代码
func dbGenCode(db *gorm.DB, models []any) {
	cfg := gen.Config{
		OutPath:           "./services/orm/query",
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

	generateExtraFiles(models)
}

func codeGenCode(db *gorm.DB, models []any) {
	cfg := gen.Config{
		OutPath:           "./services/orm/query",
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
	g.WithFileNameStrategy(func(tableName string) string {
		return tableName + "_repo"
	})
	g.Execute()

	generateExtraFiles(models)
}

func generateExtraFiles(models []any) {
	wd, _ := os.Getwd()
	fmt.Println("current working dir:", wd)

	// 确认后再硬改路径，或者：
	templateDir := filepath.Join(wd, "cmd/orm/templates")
	outDir := filepath.Join(wd, "services/orm")

	var err error
	// 读取 templates 目录下所有模板文件
	entries, err := os.ReadDir(templateDir)
	if err != nil {
		panic(fmt.Sprintf("failed to read template dir: %v", err))
	}

	if err = os.MkdirAll(outDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create output dir: %v", err))
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".tpl") {
			continue
		}

		tplPath := filepath.Join(templateDir, entry.Name())
		tplContent, err := os.ReadFile(tplPath)
		if err != nil {
			panic(fmt.Sprintf("failed to read template %s: %v", tplPath, err))
		}

		tmpl, err := template.New(entry.Name()).Parse(string(tplContent))
		if err != nil {
			panic(fmt.Sprintf("failed to parse template %s: %v", tplPath, err))
		}

		// 模板文件名去掉 .tpl 作为生成文件的后缀规则，例如 repo.go.tpl -> {model}_repo.go
		baseName := strings.TrimSuffix(entry.Name(), ".tpl") // e.g. "repo.go"

		for _, model := range models {
			t := reflect.TypeOf(model)
			for t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			modelName := t.Name()

			data := struct {
				ModelName string
			}{
				ModelName: modelName,
			}

			// 生成文件名：sys_user_repo.go
			fileName := toSnakeCase(modelName) + "_" + baseName
			filePath := filepath.Join(outDir, fileName)

			f, err := os.Create(filePath)
			if err != nil {
				panic(fmt.Sprintf("failed to create file %s: %v", filePath, err))
			}

			if err = tmpl.Execute(f, data); err != nil {
				_ = f.Close()
				panic(fmt.Sprintf("failed to execute template for %s: %v", modelName, err))
			}
			_ = f.Close()
			fmt.Printf("generated: %s\n", filePath)
		}
	}
}

func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}
