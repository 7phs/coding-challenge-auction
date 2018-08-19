FROM golang:1.10-stretch

ADD . /go/src/github.com/7phs/coding-challenge-auction
WORKDIR /go/src/github.com/7phs/coding-challenge-auction

RUN make build

FROM debian:stretch

RUN apt-get update \
    && apt-get clean

EXPOSE 8080
WORKDIR /root/
COPY --from=0 /go/src/github.com/7phs/coding-challenge-auction .

ENV LOG_LEVEL info
ENV STAGE production
ENV ADDR :8080
ENV CORS true

CMD ["./auction", "run"]