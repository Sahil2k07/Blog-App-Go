FROM golang

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod download
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

RUN go build -o main .

CMD ["./main"]