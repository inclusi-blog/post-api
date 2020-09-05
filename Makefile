WORK_DIR = $(shell pwd)

PROJECT := post-api

BUILD_VENDOR := git config --global url."https://gola-glitch:2f139c1997392434c4acfd282d8d91d70325ac8f@github.com".insteadOf "https://github.com" && \
                go env -w GOPRIVATE=github.com/gola-glitch && go mod vendor && chmod -R +w vendor

docker_login:
	@docker login -u $(ARTIFACTORY_USER) -p $(ARTIFACTORY_PASSWORD)

install_deps: docker_login
	docker-compose -f infrastructure/build.yml --project-name $(PROJECT) \
	run --rm build-env /bin/sh -c "$(BUILD_VENDOR)"

build: install_deps
	docker-compose -f infrastructure/build.yml --project-name $(PROJECT) \
	run --rm build-env /bin/sh -c "go build -mod=vendor -o ./bin/post-api"

safesql: install_deps
	docker-compose -f infrastructure/build.yml --project-name $(PROJECT) \
	run --rm build-env /bin/sh -c "go get github.com/stripe/safesql && safesql main.go"

vet: install_deps
	docker-compose -f infrastructure/build.yml --project-name $(PROJECT) \
	run --rm build-env /bin/sh -c "go vet -mod=vendor ./..."

clean:
	chmod -R +w ./.gopath vendor || true

start-db:
	docker stop gola-db-test && docker rm gola-db-test && docker network prune -f && docker volume prune -f
	docker-compose -f docker-compose.db.yml --project-name $(PROJECT) up -d

start-test-db:
	docker stop gola-db && docker rm gola-db && docker network prune -f && docker volume prune -f
	docker-compose -f docker-compose.test.yml --project-name $(PROJECT) up -d
