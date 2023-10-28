FROM golang:alpine as builder

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the Go module and all Go source code
COPY go.mod go.sum ./
COPY *.go ./

# Download and install any dependencies
RUN go mod download

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -tags urfave_cli_no_docs -o /mastodonctl

FROM scratch

# mastodonctl needs certs to make API calls over TLS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /mastodonctl /mastodonctl

ENTRYPOINT ["/mastodonctl"]
