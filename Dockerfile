FROM golang:1.24.7

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /app/todo-list && chmod +x /app/todo-list

ENTRYPOINT ["/app/todo-list"]