IMAGE = github.com/7phs/coding-challenge-auction
VERSION = latest

build: export GO111MODULE=on
build:
	go get
	go build -o auction -ldflags "-X github.com/7phs/coding-challenge-auction/cmd.BuildTime=`date +%Y-%m-%d:%H:%M:%S` -X github.com/7phs/coding-challenge-auction/cmd.GitHash=`git rev-parse --short HEAD`"

testing: export GO111MODULE=on
testing:
	go get
	LOG_LEVEL=error ADDR=:8080 STAGE=testing go test ./...

bench: export GO111MODULE=on
bench:
	go get
	LOG_LEVEL=error ADDR=:8080 STAGE=testing go test ./... -bench . -run ^$$

run: export GO111MODULE=on
run:
	go get
	LOG_LEVEL=info ADDR=:8080 STAGE=production go run main.go run

image:
	docker build -t $(IMAGE):$(VERSION) .

image-run:
	docker run --rm -it -p 8080:8080 $(IMAGE):$(VERSION)

all: build
