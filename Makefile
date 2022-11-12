version=fake
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
