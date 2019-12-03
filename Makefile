# Highlight
HL = @printf "\033[36m>> $1\033[0m\n"

default: help

help:
	@echo "Usage: make <TARGET>\n\nTargets:"
	@grep -E "^[\. a-zA-Z_-]+:.*?## .*$$" $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' |sort

test: ## Run unit tests with ARGS=<go_test_args>
	$(call HL,test)
	@DB_URL=db-url go test -count=1 $(ARGS) ./...
