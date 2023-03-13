FROM golang:1.18.0-buster AS build_base

WORKDIR /store

COPY go.mod .
COPY go.sum .

# Download all depeendencies
RUN go mod download

FROM build_base AS builder

WORKDIR /store

# Copy all the relevant files
COPY main.go main.go
COPY pkg pkg

# Create a static build of the package
RUN CGO_ENABLED=0 go install -ldflags '-extldflags "-static"' .

# Used for the runtime. Nice and light.
FROM alpine:3.11

# Install certificates that we'll add to the production env
RUN apk --no-cache add ca-certificates tzdata

# Add the binary
COPY --from=builder /go/bin/store /store

ENTRYPOINT ["/store"]


