# container for building the app
FROM golang:1.22-bullseye AS builder
WORKDIR /app

# copy files
COPY go.mod ./
COPY go.sum ./
COPY main.go ./
COPY vendor/ ./vendor/

# build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s" -o /kdc

# create an app user
RUN echo "appuser:x:10001:10001:App User:/:/sbin/nologin" > /etc/minimal-passwd

# runtime container
FROM scratch

# copy output from builder
COPY --from=builder /kdc /kdc
COPY --from=builder /etc/minimal-passwd /etc/passwd

USER appuser
ENTRYPOINT ["/kdc"]
