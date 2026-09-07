.PHONY: api-doc lint

SWAG := github.com/swaggo/swag/v2/cmd/swag@v2.0.0-rc5
ASYNGO := github.com/polanski13/asyngo/cmd/asyngo@v1.0.0

ifeq (lint,$(firstword $(MAKECMDGOALS)))
  LINT_TARGET := $(word 2,$(MAKECMDGOALS))
  $(eval $(LINT_TARGET):;@true)
endif

ifeq (api-doc,$(firstword $(MAKECMDGOALS)))
  API_DOC_TARGET := $(word 2,$(MAKECMDGOALS))
  $(eval $(API_DOC_TARGET):;@true)
endif

# Линтер в одноразовом контейнере (--rm), код монтируется (--volume).
# Не docker compose и не volume БД.
LINTER = docker run --rm \
	--volume $(CURDIR)/apps/$(LINT_TARGET):/app \
	--workdir /app \
	golangci/golangci-lint:v2.6.2 \
	golangci-lint run

api-doc:
ifeq ($(API_DOC_TARGET),)
	rm -rf apps/api-docs/docs
	cd apps/api-docs && go run $(SWAG) init -g main.go -o ../../schemas --outputTypes yaml --v3.1
	mv schemas/swagger.yaml schemas/openapi.yaml
	cd apps/api-docs && go run ./tools/fixspec ../../schemas/openapi.yaml
	cd apps/api-docs && go run $(ASYNGO) init --dir . --main events.go --output ../../schemas --outputTypes yaml
else
	@test -d "apps/$(API_DOC_TARGET)" || (echo "no such service: apps/$(API_DOC_TARGET)"; exit 1)
	cd apps/$(API_DOC_TARGET) && go run $(SWAG) init -g cmd/api/main.go -o docs --parseInternal --outputTypes yaml,json --v3.1
endif

lint:
	@test -n "$(LINT_TARGET)" || (echo "usage: make lint <service>"; exit 1)
	@test -d "apps/$(LINT_TARGET)" || (echo "no such service: apps/$(LINT_TARGET)"; exit 1)
	golangci-lint cache clean >/dev/null 2>&1 || true
	$(LINTER)
