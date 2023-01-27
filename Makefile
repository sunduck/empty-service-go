SHELL:=/bin/bash
IMAGE_TAG:=empty-service-go

all: update

build:
	docker build -t $(IMAGE_TAG) .

push:
	docker push $(IMAGE_TAG)

local-recreate: build
	docker rm --force empty-service-go
	docker run --network paxful --name empty-service-go --env-file .env.local -p 8890:9001 $(IMAGE_TAG)

update: build push

local-restart:
	docker rm --force empty-service-go
	docker run --network paxful --name empty-service-go --env-file .env.local -p 8890:9001 $(IMAGE_TAG)

local-recreate: build local-restart

k8s-recreate: update
	kubectl delete deployment empty-service-go
	kubectl apply -f deploy/k8s/deploy.yaml

swagger:
	swag init --md .
	cp api.md README.md