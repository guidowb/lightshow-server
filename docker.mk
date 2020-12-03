.PHONY: docker-image build push deploy

.docker-image-built: ${SRCS}
	docker build -t ${IMAGE} .
	@touch .docker-image-built
	@make next-build

.docker-image-pushed: .docker-image-built VERSION
	$(eval TAG := $(shell cat VERSION 2>/dev/null || echo 0.0.1))
	docker tag ${IMAGE}:latest ${REGISTRY}/${IMAGE}:${TAG}
	docker tag ${IMAGE}:latest ${REGISTRY}/${IMAGE}:latest
	docker push ${REGISTRY}/${IMAGE}:${TAG}
	docker push ${REGISTRY}/${IMAGE}:latest
	@touch .docker-image-pushed

build: .docker-image-built

push: .docker-image-pushed

deploy: .docker-image-pushed
	$(eval TAG := $(shell cat VERSION 2>/dev/null || echo 0.0.1))
	kubectl set image deployment.apps/${IMAGE} ${IMAGE}=${REGISTRY}/${IMAGE}:${TAG}
