.DEFAULT_GOAL := build

.PHONY:fmt vet build
fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: notes calendar

notes: vet
	go build -o cmd/notes/notes ./cmd/notes

calendar: vet
	go build -o cmd/calendar/calendar ./cmd/calendar

run: build
	DD_SERVICE=notes DD_ENV=dev DD_VERSION=0.0.1 ./cmd/notes/notes &
	DD_SERVICE=calendar DD_ENV=dev DD_VERSION=0.0.1 ./cmd/calendar/calendar &

exit:
	curl -X POST localhost:8080/notes/quit
	curl -X POST localhost:9090/calendar/quit

clean:
	go clean