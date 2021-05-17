FROM golang:1.15.2-alpine as builder

WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# Cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN GOPROXY=https://goproxy.io,direct  go mod download

# Copy the go source code
COPY main.go main.go
#COPY pkg/ pkg/

# Build
RUN GOPROXY=https://goproxy.io,direct CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o http-server main.go


FROM alpine
COPY --from=builder /workspace/http-server /app/
COPY config.yaml /app/
COPY static /app/static
COPY dist  /app/dist

WORKDIR /app
ENTRYPOINT ["/app/http-server"]