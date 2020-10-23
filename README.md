#post-api (story-api): Repository for gola microservice

## Download source
    git clone git@github.com:gola-glitch/post-api.git
    
## Scripts
### Prerequisites
* You need docker and docker-compose installed.

### Format source code
    make format

### Install dependencies
    make install_deps
    
### Build go binaries
    make build
    
### Go vet linter
    make vet
    
### Build docker image
    make dockerize
    
### Getting started as a service
###Steps to run post-api in local environment using docker compose
    make start

### Stop the app and db
    make stop
    
### Lint new docker files
    make hadolint

### healthcheckz
    make healthcheck
    
### DEV guidelines
### Install pre commit and pre push hooks for code quality
    make install_hooks
    
### To uninstall hooks
    make uninstall_hooks
    
### Update swagger documentation
####Setup swaggo
    go get -u github.com/swaggo/swag/cmd/swag

####Update
    swag init
