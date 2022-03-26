build: build-gateway build-op-processor copy-migrations copy-settings

build-gateway:
	go build -ldflags "-s -w" -o bin/gateway gateway/main.go

build-op-processor:
	go build -ldflags "-s -w" -o bin/op_processor op_processor/main.go

copy-migrations:
	rsync -rup migrations bin/

copy-settings:
	rsync -up settings.yml bin/
