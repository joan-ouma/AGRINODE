FROM golang:latest
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o command-svc ./cmd/command-svc
CMD ["./command-svc"]
