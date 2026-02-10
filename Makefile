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
	@go test -v -short -count=1 -race ./... -coverprofile=coverage.out -covermode=atomic

test.integration :: docker-up
	@echo "Waiting for PostgreSQL to be ready..."
	@until docker compose exec -T postgres pg_isready -U postgres > /dev/null 2>&1; do sleep 1; done
	@go run cmd/seed/main.go
	@go test -v -count=1 ./...

docker-up ::
	docker compose up -d

docker-down ::
	docker compose down
