.PHONY: proto clean deps

# Generate protobuf files
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/user/v1/user_service.proto

# Install dependencies
deps:
	go mod tidy
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Clean generated files
clean:
	find . -name "*.pb.go" -delete

# Build and test
build:
	go build ./...

test:
	go test -v ./...

# Tag a new version
tag:
	@echo "Current tags:"
	@git tag -l | sort -V | tail -5
	@echo "Create new tag with: git tag v1.0.0 && git push origin v1.0.0"
