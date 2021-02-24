.PHONY: migration.create migration.up migration.down

all: build

build: bin/migrate

bin/migrate: $(./.../*.go)
	go build -mod vendor -o bin/migrate ./cmd/migrate

migration.create:
	go run ./cmd/migrate create "$(name)"
migration.up:
	bin/migrate up $(count)
migration.down:
	bin/migrate down $(count)
