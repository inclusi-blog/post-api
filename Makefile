WORK_DIR = $(shell pwd)

PROJECT := post-api
REVISION := latest

BUILD_VENDOR := git config --global url."https://inclusi-blog:$(GITHUB_ACCESS_TOKEN)@github.com".insteadOf "https://github.com" && \
                go env -w GOPRIVATE=github.com/inclusi-blog && go mod vendor && chmod -R +w vendor

install_deps:
	docker compose -f infrastructure/build.yml --project-name $(PROJECT) \
	run --rm build-env /bin/sh -c "apk update && apk add git && $(BUILD_VENDOR)"

build: install_deps
	docker compose -f infrastructure/build.yml --project-name $(PROJECT) \
	run --rm build-env /bin/sh -c "go build -mod=vendor -o ./bin/$(PROJECT)"

vet: install_deps
	docker-compose -f infrastructure/build.yml --project-name $(PROJECT) \
	run --rm build-env /bin/sh -c "go vet -mod=vendor ./..."

start: build
	docker-compose -f docker-compose.local-app.yml up -d

stop:
	docker-compose -f docker-compose.local-app.yml down -v

clean:
	chmod -R +w ./.gopath vendor || true
