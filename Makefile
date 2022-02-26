CONFIG_PATH=${HOME}/.proglog/

.PHONY: init
init:
	mkdir -p ${CONFIG_PATH}

.PHONY: gencert
gencert:
	cfssl gencert \
		-initca test/ca-csr.json | cfssljson -bare ca

	cfssl gencert \
		-ca=ca.pem \
		-ca-key=ca-key.pem \
		-config=test/ca-config.json \
		-profile=server \
		test/server-csr.json | cfssljson -bare server

		cfssl gencert \
			-ca=ca.pem \
			-ca-key=ca-key.pem \
			-config=test/ca-config.json \
			-profile=client \
			-cn="root" \
			test/client-csr.json | cfssljson -bare root-client

		cfssl gencert \
			-ca=ca.pem \
			-ca-key=ca-key.pem \
			-config=test/ca-config.json \
			-profile=client \
			-cn="nobody" \
			test/client-csr.json | cfssljson -bare nobody-client

	mv *.pem *.csr ${CONFIG_PATH}

.PHONY: proto
proto:
	protoc api/v1/*.proto \
	--go_out=. \
	--go_opt=paths=source_relative \
	--proto_path=.

.PHONY: compile
compile:
	protoc api/v1/*.proto \
	--go_out=. \
	--go-grpc_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_opt=paths=source_relative \
	--proto_path=.

$(CONFIG_PATH)/model.conf:
	cp test/model.conf $(CONFIG_PATH)/model.conf

$(CONFIG_PATH)/policy.csv:
	cp test/policy.csv $(CONFIG_PATH)/policy.csv

.PHONY: test
test: $(CONFIG_PATH)/policy.csv $(CONFIG_PATH)/model.conf
	go test -race ./...

.PHONY: docker-build
docker-build:
	docker buildx build . --platform linux/amd64 --tag caladreas/proglog:latest --load

.PHONY: dbuild
docker-run:
	docker run --name proglog caladreas/proglog:latest

.PHONY: dbuild
dbuild:
	docker buildx build . --platform linux/arm64,linux/amd64 --tag caladreas/proglog:latest --push

.PHONY: multiarch
multiarch:
	CGO_ENABLE=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o bin/$(TARGETARCH)/proglog ./cmd/proglog

.PHONY: linux
linux:
	CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -o bin/proglog ./cmd/proglog

.PHONY: clean
clean:
	rm -rf /tmp/proglog