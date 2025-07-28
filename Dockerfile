FROM oven/bun:latest AS bun-builder

WORKDIR /app

COPY reverse-proxy/frontend/package.json .

RUN bun install

COPY reverse-proxy/frontend/ .

RUN bun vite build


FROM golang:1.24-alpine AS go-builder

WORKDIR /app

COPY reverse-proxy/go.mod .
COPY reverse-proxy/go.sum .

RUN go mod download

COPY reverse-proxy/cmd/ ./cmd/
COPY reverse-proxy/internal/ ./internal/

RUN go build -o rp cmd/rp/main.go

FROM alpine:latest 

WORKDIR /rp

COPY --from=go-builder /app/rp .
COPY --from=bun-builder /app/dist ./dist

EXPOSE 80 5000 3000

CMD ["./rp"]