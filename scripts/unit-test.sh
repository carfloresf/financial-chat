#!/usr/bin/env bash

# tools needed
tools=(
	gotest.tools/gotestsum@latest
)

# Install missing tools
for tool in ${tools[@]}; do
	which $(basename ${tool}) > /dev/null || go install ${tool}
done

echo "Running unit tests."

# Generate tests report
go clean -testcache
gotestsum --format pkgname -- -cover ./...; test ${PIPESTATUS[0]} -eq 0 || status=${PIPESTATUS[0]}

exit ${status:-0}
