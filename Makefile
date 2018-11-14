
.PHONY: ci_deps vendor

default: ci_deps vendor

vendor:
	vndr
ci_deps:
	@echo "fetching vndr"
	@go get -u github.com/LK4D4/vndr