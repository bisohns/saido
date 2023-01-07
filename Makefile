version=fake
ifeq ($(OS),Windows_NT)
bin=main.exe
build_bin=tmp\main.exe
else
bin=main
build_bin=tmp/main
endif
# Example:
#   make
#   make prep-ci-ssh
.PHONY: prep-ci-ssh
# Creates the ssh keys and docker container for running test
prep-ci-ssh:
	./scripts/prep-test-ssh.sh
	cat config-test.yaml

.PHONY: prep-ci-locals
prep-ci-local:
	cp ./scripts/config.local.yaml config-test.yaml
	cat config-test.yaml

.PHONY: prep-ci-local-windows
# Because we are using make, we can use linux cp and cat commands
prep-ci-local-windows:
	cp ".\scripts\config.local.yaml" ".\config-test.yaml"
	cat config-test.yaml

ifneq ($(findstring fake, $(version)), fake)
upgrade:
	@echo ">>> Recreating version_num.go"
	@echo 'package cmd\n\nconst Version = "$(version)"' > cmd/version_num.go
	@go run main.go version
	@git tag v$(version)
	@echo ">>> Pushing Tag to Remote..."
	@git push origin v$(version)
else
upgrade:
	@echo "Version not set - use syntax \`make upgrade version=0.x.x\`"
endif

dependencies:
ifeq ($(bin),main.exe)
	@make prep-ci-local-windows
	cd web && yarn add react-scripts@latest
else
	@make prep-ci-local
endif
	go get .
	cd web && yarn install && cd ..

.PHONY: build-frontend
build-frontend:
	cd web && export BUILD_PATH=../cmd/build && CI=false yarn build && cd ..

.PHONY: dev-backend
dev-backend:
	air --build.cmd "go build -o ./tmp/$(bin) ." --build.bin "$(build_bin)" --build.exclude_dir "assets,docs,tmp,web,scripts,ssh-key,.github,.git" --build.include_ext "go,yaml,html,js" --build.exclude_file "config.example.yaml" --build.args_bin "--config,config-test.yaml,--verbose" --build.stop_on_error true --misc.clean_on_exit true --log.time true

prod-monolith:
	go build -tags prod -o ./tmp/$(bin) . && $(build_bin) --config config-test.yaml -b --verbose

app: build-frontend prod-monolith
