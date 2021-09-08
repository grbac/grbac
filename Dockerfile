FROM golang:1.16-alpine AS builder

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o grbac ./cmd

FROM alpine

WORKDIR /usr/local/grbac

COPY --from=builder /build/grbac bin/grbac
COPY scripts/docker-compose.sh docker-compose.sh

ENV PATH=/usr/local/grbac/bin:$PATH

ENTRYPOINT [ "grbac" ]
CMD [ "version" ]
