GOBIN=$(GOPATH)/bin
APP_DIR_LIST=$(shell go list ./... | grep -v /vendor/)

docker_build: build_anywhere
	docker build -t tapng-gateway .

push_docker: docker_build
	docker tag tapng-gateway $(REPOSITORY_URL)/tapng-gateway:latest
	docker push $(REPOSITORY_URL)/tapng-gateway:latest

kubernetes_deploy: docker_build
	kubectl create -f configmap.yaml
	kubectl create -f service.yaml
	kubectl create -f deployment.yaml

kubernetes_update: docker_build
	kubectl delete -f deployment.yaml
	kubectl create -f deployment.yaml

prepare_dirs:
	mkdir -p ./temp/src/github.com/trustedanalytics/gateway
	$(eval REPOFILES=$(shell pwd)/*)
	ln -sf $(REPOFILES) temp/src/github.com/trustedanalytics/gateway

build_anywhere: prepare_dirs
	$(eval GOPATH=$(shell cd ./temp; pwd))
	$(eval APP_DIR_LIST=$(shell GOPATH=$(GOPATH) go list ./temp/src/github.com/trustedanalytics/gateway/... | grep -v /vendor/))
	GOPATH=$(GOPATH) CGO_ENABLED=0 go build $(APP_DIR_LIST)
	rm -Rf application && mkdir application
	cp ./gateway ./application/gateway
	rm -Rf ./temp
