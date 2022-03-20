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
prep-ci-local-windows:
	cp ".\scripts\config.local.yaml" ".\config-test.yaml"
	type config-test.yaml

