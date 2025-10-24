IMAGE ?= shimit/faulti_pod
TAG ?= latest
RELEASE ?= faulty
NAMESPACE ?= default
CHART_DIR := helm

.PHONY: build push deploy-cluster deploy-e2e

build:
	docker build -t $(IMAGE):$(TAG) .

push:
	docker push $(IMAGE):$(TAG)

deploy-cluster:
	helm upgrade --install $(RELEASE) $(CHART_DIR) \
	  --namespace $(NAMESPACE) --create-namespace \
	  --set image.repository=$(IMAGE) --set image.tag=$(TAG)

deploy-e2e: build push deploy-cluster


