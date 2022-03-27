build: go-mod-tidy build-gateway build-op-processor build-wlt-balancer copy-migrations copy-settings

go-mod-tidy:
	go mod tidy

build-gateway:
	go build -ldflags "-s -w" -o bin/gateway gateway/main.go

build-op-processor:
	go build -ldflags "-s -w" -o bin/op_processor op_processor/main.go

build-wlt-balancer:
	go build -ldflags "-s -w" -o bin/wlt_balancer wlt_balancer/main.go

copy-migrations:
	rsync -rup migrations bin/

copy-settings:
	rsync -up settings.yml bin/
