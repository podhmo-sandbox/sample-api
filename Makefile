# ----------------------------------------
# management
# ----------------------------------------
test:
	go test ./... # -coverprofile=./test_results/cover.out
# 	go tool cover -html=./test_results/cover.out -o ./test_results/cover.html
.PHONY: test
format:
	go list -f '{{.Dir}}' ./... | xargs goimports -w
.PHONY: format
lint:
	staticcheck ./...
	go vet ./...
.PHONY: lint

# ----------------------------------------
# local development
# ----------------------------------------
serve:
	docker-compose -f build/docker-compose.yml up
.PHONY: serve
