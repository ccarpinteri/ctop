FROM golang:1.23-alpine AS builder

ARG VERSION=dev
ARG BUILD=unknown

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -tags release -ldflags="-w -s -X main.version=${VERSION} -X main.build=${BUILD}" -o ctop .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
ENV TERM=linux
COPY --from=builder /app/ctop /ctop
ENTRYPOINT ["/ctop"]
