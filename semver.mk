SHA := $(shell git rev-parse --short HEAD)
TAG := $(shell cat VERSION 2>/dev/null || echo 0.0.1)
VERSION := $(shell echo ${TAG} | cut -d'-' -f1)
MAJOR := $(shell echo ${VERSION} | cut -d'.' -f1)
MINOR := $(shell echo ${VERSION} | cut -d'.' -f2)
PATCH := $(shell echo ${VERSION} | cut -d'.' -f3)
BUILD := $(shell echo ${TAG}-build.0 | cut -d'-' -f2 | sed 's/^build.//' )

.PHONY: version next-build patch minor major

version:
	$(eval TAG := $(shell cat VERSION 2>/dev/null || echo 0.0.1))
	@echo ${TAG}

next-build: ${SRCS}
	$(eval BUILD := $(shell expr ${BUILD} + 1))
	@echo ${VERSION}-build.${BUILD} >VERSION

patch:
	$(eval PATCH := $(shell expr ${PATCH} + 1))
	@echo ${MAJOR}.${MINOR}.${PATCH} | tee VERSION

minor:
	$(eval MINOR := $(shell expr ${MINOR} + 1))
	@echo ${MAJOR}.${MINOR}.0 | tee VERSION

major:
	$(eval MAJOR := $(shell expr ${MAJOR} + 1))
	@echo ${MAJOR}.0.0 | tee VERSION
