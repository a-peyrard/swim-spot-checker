set shell := ["bash", "-uc"]
set positional-arguments := true

project := 'swim-spot-checker'
docker-account := 'augustinpeyrard'

# show this help
help:
    @just --list

# configure the dev environment
build:
    @go build -o swim-spot-checker main.go

# run unit tests
test *ARGS:
    @go test -- -v -race "$@" ./pkg/...

# publish new version
publish:
    @docker buildx build --platform linux/amd64 -t {{docker-account}}/{{project}}:latest --push .

