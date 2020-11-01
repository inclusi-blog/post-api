WORK_DIR = $(shell pwd)

PROJECT := post-api
REVISION := latest
RELEASE_SCRIPTS_VERSION := latest

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

create_user:
	docker exec -it gola-db /bin/sh /sql/create_user.sh

create_user_dev:
	docker exec -it gola-db /bin/sh /sql/create_user_dev.sh

run_migration:
	docker-compose -f docker-compose.db.yml up -d post-migration

run_test_migration:
	docker-compose -f docker-compose.test.yml up -d post-test-migration

safesql: install_deps
	docker-compose -f infrastructure/build.yml --project-name $(PROJECT) \
	run --rm build-env /bin/sh -c "go get github.com/stripe/safesql && safesql main.go"

vet: install_deps
	docker-compose -f infrastructure/build.yml --project-name $(PROJECT) \
	run --rm build-env /bin/sh -c "go vet -mod=vendor ./..."

clean:
	chmod -R +w ./.gopath vendor || true

create-db:
	docker network prune -f && docker volume prune -f && \
	docker-compose -f docker-compose.db.yml --project-name $(PROJECT) up -d gola-db && \
	sleep 100

start-db: create-db create_user run_migration run_test_migration

stop-db:
	docker-compose -f docker-compose.db.yml -f docker-compose.test.yml down -v

publish: docker_login
	docker tag post-api $(ARTIFACTORY_USER)/post-api:$(REVISION); \
	docker push $(ARTIFACTORY_USER)/post-api:$(REVISION)

start: create-db create_user run_migration run_test_migration build
	docker-compose -f docker-compose.db.yml -f docker-compose.test.yml -f docker-compose.local-app.yml up -d

stop:
	docker-compose -f docker-compose.db.yml -f docker-compose.test.yml -f docker-compose.local-app.yml down -v

test: install_deps
	docker-compose -f docker-compose.test.yml \
	--project-name $(PROJECT) \
	run --rm post-test

dockerize: docker_login
	docker-compose -f docker-compose.db.yml -f docker-compose.test.yml -f docker-compose.local-app.yml build --no-cache

hadolint: docker_login
	docker run --rm -i hadolint/hadolint:latest hadolint --ignore DL3007 --ignore DL3008 --ignore SC2016 - < infrastructure/Dockerfile
	docker run --rm -i hadolint/hadolint:latest hadolint --ignore DL3007 --ignore DL3008 --ignore SC2016 - < infrastructure/Migrate.Dockerfile

golangci-lint: install_deps
	docker run --rm \
	    -e GOLANGCI_LINT_CACHE=/tmp/.cache \
        -v $(WORK_DIR):/post-api \
        golangci/golangci-lint:v1.21 /bin/sh -c "cd /post-api && mkdir -p /tmp/.cache && golangci-lint run -v ./... "

dev_migration:
	docker-compose -f docker-compose-db.dev.migration.yml up -d

healthcheck: start
	EXIT_CODE=$(shell ./docker-compose-scripts/test-scripts/verify_healthcheck.sh http://localhost:30003/api/post/healthz > /dev/null 2>&1; echo $$?); \
	docker-compose -f docker-compose.db.yml -f docker-compose.test.yml -f docker-compose.local-app.yml --project-name $(PROJECT) down -v; \
	exit $$EXIT_CODE

##Generates metadata for k8s
generate_metadata:
	echo "$$METADATA" > metadata

##generates pipeline template for release
generate_template: docker_login
	docker run --rm \
	-v $(WORD_DIR):/home/gola \
	gola_release_scripts:$(RELEASE_SCRIPTS_VERSION) \
	/bin/bash -c "/scripts/export_template.sh $(GO_PIPELINE_NAME)"

generate_release_artifacts: generate_template
	tar zcf release_artifacts.tar.gz post-api.gocd.yaml template-post-api.gocd.json

format:
	go fmt ./...

pre_commit:
	go mod tidy
	go vet ./...
	go fmt ./...

pre_push:
	true

install_hooks: ## Dev: Install pre-commit and pre-push hooks
	if [ -f ${WORK_DIR}/.git/hooks/pre-commit ]; then mv ${WORK_DIR}/.git/hooks/pre-commit ${WORK_DIR}/.git/hooks/old-pre-commit; fi
	if [ -f ${WORK_DIR}/.git/hooks/pre-push ]; then mv ${WORK_DIR}/.git/hooks/pre-push ${WORK_DIR}/.git/hooks/old-pre-push; fi
	ln -s ${WORK_DIR}/infrastructure/hooks/pre-push ${WORK_DIR}/.git/hooks/pre-push
	ln -s ${WORK_DIR}/infrastructure/hooks/pre-commit ${WORK_DIR}/.git/hooks/pre-commit
	chmod +x ${WORK_DIR}/.git/hooks/pre-push ${WORK_DIR}/.git/hooks/pre-commit

uninstall_hooks: ## Dev: Uninstall pre-commit and pre-push hooks
	rm ${WORK_DIR}/.git/hooks/pre-commit
	rm ${WORK_DIR}/.git/hooks/pre-push;
	if [ -f ${WORK_DIR}/.git/hooks/old-pre-commit ]; then mv ${WORK_DIR}/.git/hooks/old-pre-commit ${WORK_DIR}/.git/hooks/pre-commit; fi
	if [ -f ${WORK_DIR}/.git/hooks/old-pre-push ]; then mv ${WORK_DIR}/.git/hooks/old-pre-push ${WORK_DIR}/.git/hooks/pre-push; fi
