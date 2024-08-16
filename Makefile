IMAGE_NAME=choisunguk/admission_controller_demo
IMAGE_TAG=v1

up:
	kind create cluster --config kind-cluster/config.yaml

down:
	kind delete cluster --name admission-controller

create-builder:
	docker buildx create --name mybuilder --use

build-push:
	docker buildx build --platform linux/amd64,linux/arm64 -t $(IMAGE_NAME):${IMAGE_TAG} --push .
