.PHONY: clean cluster compose context coverage install lint test update vet web
.SILENT:

kind.conf: context
	kubectl config view --raw | sed -E 's/127.0.0.1|localhost/host.docker.internal/' > kind.conf

clean:
	kind delete cluster
	rm -f kind.conf &>/dev/null

cluster: clean
	kind create cluster --config=kind.yaml

compose: kind.conf
	docker-compose up --build

context:
	kubectl config use-context kind-kind

coverage: lint
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	rm -f coverage.out &>/dev/null

install: context
	kubectl create namespace argo
	kubectl apply -n argo -f https://github.com/argoproj/argo-workflows/releases/download/v3.3.9/install.yaml
	kubectl create rolebinding default-admin --clusterrole=admin --serviceaccount=argo:default --namespace=argo

lint: vet
	golangci-lint run ./...

test: lint
	go test ./...

update:
	go get -u -t -d -v ./...
	go mod tidy

vet:
	go vet ./...

web:
	argo server --auth-mode=server