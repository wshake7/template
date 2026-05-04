bob_dir = ./orm
bobgen_dir = ./etc/bobgen.yaml
bob:
	@mkdir -p ${bob_dir}/queries
	@(cd ${bob_dir} && go run github.com/stephenafamo/bob/gen/bobgen-psql@latest -c ../${bobgen_dir})

.PHONY: script-imports

script-imports:
	go run ./cmd/script/gen_imports

script-orm:
	go run ./cmd/script/orm

run:
	go run ./cmd

swagger:
	swag init -g cmd/main.go --parseDependency
