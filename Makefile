fmt:
	go fmt ./...

vet:
	go vet ./...

test:
	go test -v -race -cover ./...

lint:
	golangci-lint run \
		-v \
		--no-config \
		-E goconst \
		-E goimports \
		-E gocritic \
		-E golint \
		-E interfacer \
		-E maligned \
		-E misspell \
		-E stylecheck \
		-E unconvert \
		-E unparam \
		-D errcheck \
		--skip-dirs vendor ./...
	golint -set_exit_status ./...
