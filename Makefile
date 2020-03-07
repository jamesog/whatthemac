IMAGE := mac.jog.li
VERSION := 0.1.0

GCR_REGION  := eu
GCR_HOST    := $(GCR_REGION).gcr.io
GCP_PROJECT := whatthemac

all: build

build:
	docker build -t $(IMAGE):$(VERSION) .

run:
	docker run -e PORT=8000 -p 8000:8000 $(IMAGE):$(VERSION)

push:
	docker tag $(IMAGE):$(VERSION) $(GCR_HOST)/$(GCP_PROJECT)/$(IMAGE):$(VERSION)
	docker push $(GCR_HOST)/$(GCP_PROJECT)/$(IMAGE):$(VERSION)
	gsutil acl ch -r -u AllUsers:READ gs://$(GCR_REGION).artifacts.$(GCP_PROJECT).appspot.com/

deploy:
	kubectl apply -f deployment.yaml
