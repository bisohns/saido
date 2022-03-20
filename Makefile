# Example:
#   make
#   make prep-ci
.PHONY: prep-ci
# Creates the ssh keys and docker container for running test
prep-ci:
	./scripts/make-ci-test.sh
	cat config-ci.yaml
