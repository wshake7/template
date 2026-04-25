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
	gormCrud "orm-crud/gormc"
	"orm-crud/gormc/mixin"
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
		gormCrud.WithDriverName(conf.Orm.DriverName),
		gormCrud.WithDSN(conf.Orm.DataSource),
		gormCrud.WithEnableTrace(true),
		gormCrud.WithEnableMetrics(true),
	)

	if conf.Orm.IsAutoMigrate {
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

	sysCasbinModel := query.SysCasbinModel
	_ = sysCasbinModel.Create(&models.SysCasbinModel{
		IsEnabled: mixin.IsEnabled{IsEnabled: true},
		Name:      "pbac",
		Content:   "[request_definition]\nr = sub, obj, act\n\n[policy_definition]\np = sub_rule, obj_rule, act\n\n[policy_effect]\ne = some(where (p.eft == allow))\n\n[matchers]\nm = eval(p.sub_rule) && eval(p.obj_rule) && r.act == p.act",
	})
}

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

	templateDir := filepath.Join(wd, "cmd/scripts/orm/templates")
	outDir := filepath.Join(wd, "services/orm/repo")

	entries, err := os.ReadDir(templateDir)
	if err != nil {
		panic(fmt.Sprintf("failed to read template dir: %v", err))
	}

	if err = os.MkdirAll(outDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create output dir: %v", err))
	}

	moduleName, err := getModuleName()
	if err != nil {
		panic(fmt.Sprintf("failed to get module name: %v", err))
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

		tmpl, err := template.New(entry.Name()).Funcs(template.FuncMap{
			"toLower": func(s string) string {
				if len(s) == 0 {
					return s
				}
				return strings.ToLower(s[:1]) + s[1:]
			},
		}).Parse(string(tplContent))
		if err != nil {
			panic(fmt.Sprintf("failed to parse template %s: %v", tplPath, err))
		}

		baseName := strings.TrimSuffix(entry.Name(), ".tpl")
		if !strings.HasSuffix(baseName, ".go") {
			baseName += ".go"
		}

		for _, model := range models {
			t := reflect.TypeOf(model)
			for t.Kind() == reflect.Pointer {
				t = t.Elem()
			}
			modelName := t.Name()

			data := struct {
				ModelName  string
				ModuleName string
			}{
				ModelName:  modelName,
				ModuleName: moduleName,
			}

			fileName := toSnakeCase(modelName) + "_" + baseName
			filePath := filepath.Join(outDir, fileName)

			// 文件已存在则跳过
			if _, err = os.Stat(filePath); err == nil {
				fmt.Printf("skipping existing file: %s\n", filePath)
				continue
			}

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

func getModuleName() (string, error) {
	wd, _ := os.Getwd()
	goModPath := filepath.Join(wd, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return "", fmt.Errorf("failed to read go.mod: %v", err)
	}
	for line := range strings.SplitSeq(string(content), "\n") {
		line = strings.TrimSpace(line)
		if after, ok := strings.CutPrefix(line, "module "); ok {
			return after, nil
		}
	}
	return "", fmt.Errorf("module name not found in go.mod")
}
