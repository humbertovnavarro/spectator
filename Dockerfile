FROM golang:1.21-alpine AS build

RUN apk add --no-cache \
    git \
    gcc \
    musl-dev

WORKDIR /app
COPY . .
ARG CGO_ENABLED=1
RUN CGOENABLED=1 go build -o main .

FROM alpine:latest AS runtime
RUN apk add --no-cache \
    ca-certificates
WORKDIR /app
COPY --from=build /app/main .
CMD ["/app/main"]
