tidy ::
	@go mod tidy && go mod vendor

seed ::
	@go run cmd/seed/main.go

run ::
	@go run cmd/server/main.go

mocks ::
	@for f in app/*/contracts.go; do \
		name=$$(basename $$(dirname $$f)); \
		PATH="$$PATH:$$(go env GOPATH)/bin" mockgen -source=$$f -destination=app/mocks/$${name}_mock.go -package=mocks; \
	done

test ::
	@go test -v -count=1 -race ./... -coverprofile=coverage.out -covermode=atomic

docker-up ::
	docker compose up -d

docker-down ::
	docker compose down
