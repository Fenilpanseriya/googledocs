FROM golang:1.24.3-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o main ./main.go
RUN chmod +x main
EXPOSE 5000
CMD ["./main"]

