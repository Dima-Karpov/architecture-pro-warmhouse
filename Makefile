.PHONY: api-doc

SWAG := github.com/swaggo/swag/v2/cmd/swag@v2.0.0-rc5
ASYNGO := github.com/polanski13/asyngo/cmd/asyngo@v1.0.0

api-doc:
	rm -rf apps/api-docs/docs
	cd apps/api-docs && go run $(SWAG) init -g main.go -o ../../schemas --outputTypes yaml --v3.1
	mv schemas/swagger.yaml schemas/openapi.yaml
	cd apps/api-docs && go run ./tools/fixspec ../../schemas/openapi.yaml
	cd apps/api-docs && go run $(ASYNGO) init --dir . --main events.go --output ../../schemas --outputTypes yaml
