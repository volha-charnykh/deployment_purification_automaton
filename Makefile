

GIT_COMMIT     := $(shell git describe --dirty=-unsupported --always --tags || echo pre-commit)
IMAGE_VERSION  ?= $(GIT_COMMIT)

.PHONY: db-up
db-up:
	docker run --name pur-db -e POSTGRES_PASSWORD=infoblox123 -e POSTGRES_USER=postgres -e POSTGRES_DB=postgress -p 5432:5432 -d postgres:10.13

.PHONY: db-down
db-down:
	if docker inspect pur-db &>/dev/null; then \
		docker rm pur-db -fv; \
	fi

.PHONY: docker
docker:
	docker build -t automaton-$(IMAGE_VERSION) -f Dockerfile .
