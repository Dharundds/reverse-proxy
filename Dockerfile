FROM golang:1.24-alpine 

WORKDIR /app

COPY . .

RUN go build -o rp cmd/rp/main.go

EXPOSE 80 5000

CMD ["./rp"]