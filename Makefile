#!/bin/bash

update_dependency_injection:
	@if !command -v wire >/dev/null 2>&1 ; then \
		echo "Go Wire is not installed. Installing..."; \
		go install github.com/google/wire/cmd/wire@latest; \
	fi

	@echo "Updating dependency injection";
	@wire di/wire.go;

