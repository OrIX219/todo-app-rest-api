FROM golang

ENV GOPATH=/

COPY . .

RUN go mod download
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN go build -o main ./cmd/main.go

CMD ["make", "prod"]