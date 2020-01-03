.PHONY: lint test

test:
	go test -race -covermode atomic -coverprofile=profile.out -v ./...

lint:
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint:latest golangci-lint run -v
