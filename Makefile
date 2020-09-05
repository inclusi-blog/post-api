WORK_DIR = $(shell pwd)

PROJECT := post-api

BUILD_VENDOR := git config --global url."https://gola-glitch:2f139c1997392434c4acfd282d8d91d70325ac8f@github.com".insteadOf "https://github.com" && \
                go env -w GOPRIVATE=github.com/gola-glitch && go mod vendor && chmod -R +w vendor

docker_login:
	@docker login -u $(ARTIFACTORY_USER) -p $(ARTIFACTORY_PASSWORD)

install_deps: docker_login
	docker-compose -f infrastructure/build.yml --project-name $(PROJECT) \
	run --rm build-env /bin/sh -c "$(BUILD_VENDOR)"

safesql: install_deps
	docker-compose -f infrastructure/build.yml --project-name $(PROJECT) \
	run --rm build-env /bin/sh -c "go get github.com/stripe/safesql && safesql main.go"
