FROM golang:1.26-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /ground .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates git
COPY --from=builder /ground /usr/local/bin/ground
ENTRYPOINT ["ground"]
