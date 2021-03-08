FROM golang:1.15-alpine AS dev

WORKDIR /app

RUN apk add git

RUN GO111MODULE=on go get github.com/cortesi/modd/cmd/modd

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go install github.com/tnyie/tny-api

CMD ["go", "run", "*.go"]

FROM alpine

WORKDIR /bin

COPY --from=dev /go/bin/tny ./tny

CMD ["sh", "-c", "tny", "-p"]