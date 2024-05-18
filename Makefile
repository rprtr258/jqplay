GOBIN ?= $(CURDIR)/build
.PHONY: build
build:
	docker run -ti -v .:/data node sh -c 'cd data && yarn && rm -rf node_modules'
	go build .

.PHONY: test
test:
	docker \
		buildx \
		build \
		--rm \
		--build-arg TIMESTAMP=$$(date +%s) \
		--target gotest \
		.

.PHONY: vet
vet:
	docker \
		run \
		--rm \
		-v $(CURDIR):/app \
		-w /app \
		golangci/golangci-lint:latest \
		golangci-lint run --timeout 5m -v

TAG ?= latest
REPO ?= ghcr.io/owenthereal/jqplay
.PHONY: docker_build
docker_build:
	docker buildx build --rm -t $(REPO):$(TAG) --load .

.PHONY: docker_push
docker_push: docker_build
	docker buildx build --rm -t $(REPO):$(TAG) --push .

.PHONY: start
start:
	docker compose up --build --force-recreate

.PHONY: watch
watch:
	docker compose watch
