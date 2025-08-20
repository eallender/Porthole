VERSION 0.8

all:
    BUILD +test
    BUILD +lint
    BUILD +build-all

golang:
    FROM golang:1.24.5-alpine

deps:
    FROM +golang
    COPY go.mod ./
    RUN go mod download
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum

build:
    FROM +deps
    COPY . .
    ARG GOOS=linux
    ARG GOARCH=amd64
    ARG VERSION=dev
    RUN CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build \
        -ldflags="-w -s -X main.version=$VERSION" \
        -o porthole-$GOOS-$GOARCH \
        ./cmd/porthole
    SAVE ARTIFACT porthole-$GOOS-$GOARCH AS LOCAL bin/porthole-$GOOS-$GOARCH

build-all:
    BUILD +build --GOOS=linux --GOARCH=amd64
    BUILD +build --GOOS=darwin --GOARCH=amd64
    BUILD +build --GOOS=darwin --GOARCH=arm64
    BUILD +build --GOOS=windows --GOARCH=amd64

test:
    FROM +deps
    COPY . .
    RUN go test -v ./...

lint:
    FROM +deps
    RUN go install honnef.co/go/tools/cmd/staticcheck@latest
    COPY . .
    RUN go vet ./...
    RUN staticcheck ./...

