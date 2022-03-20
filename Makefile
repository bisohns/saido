# Example:
#   make
#   make prep-ci-ssh
.PHONY: prep-ci-ssh
# Creates the ssh keys and docker container for running test
prep-ci-ssh:
	./scripts/prep-test-ssh.sh
	cat config-ci.yaml

.PHONY: prep-ci-locals
prep-ci-local:
	./scripts/prep-test-local.sh
	cat config-ci.yaml

.PHONY: prep-ci-local-windows
prep-ci-local-windows:
	# DO SOMETHING
