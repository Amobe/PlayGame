FROM --platform=$BUILDPLATFORM golang:1.20-alpine3.17 AS builder
ENV GO111MODULE=on CGO_ENABLED=0
RUN go install honnef.co/go/tools/cmd/staticcheck@v0.4.3
WORKDIR /workspace/PlayGame/server
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
ARG TARGETOS
ARG TARGETARCH
ENV GOOS $TARGETOS
ENV GOARCH $TARGETARCH
RUN go vet ./... && staticcheck ./... && go test ./...
RUN go build -o /bin/server ./cmd/server

FROM alpine:3.17 AS release
COPY --from=builder /bin/server /bin/server
ENTRYPOINT ["./bin/server"]
