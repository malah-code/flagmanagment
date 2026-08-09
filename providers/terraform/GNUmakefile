TEST?=$$(go list ./... | grep -v /vendor/)
HOSTNAME=malah-code
NAMESPACE=flagmanagment
NAME=flagmanagment
BINARY=terraform-provider-${NAME}
VERSION=1.0.0
OS_ARCH=linux_amd64

default: build

build:
	go build -o bin/${BINARY} main.go

test:
	go test -i $(TEST) || true
	go test $(TEST) -v $(TESTARGS)

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv bin/${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
