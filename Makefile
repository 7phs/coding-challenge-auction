IMAGE = github.com/7phs/coding-challenge-auction
VERSION = latest

build:
	GO111MODULE=on go build -o auction -ldflags "-X github.com/7phs/coding-challenge-auction/cmd.BuildTime=`date +%Y-%m-%d:%H:%M:%S` -X github.com/7phs/coding-challenge-auction/cmd.GitHash=`git rev-parse --short HEAD`"

testing:
	GO111MODULE=on LOG_LEVEL=error ADDR=:8080 STAGE=testing go test ./...

run:
    GO111MODULE=on LOG_LEVEL=info ADDR=:8080 STAGE=production go run main.go

image:
	docker build -t $(IMAGE):$(VERSION)

image-run:
	docker run --rm -it -p 8080:8080 $(IMAGE):$(VERSION)

all: dep_update build
