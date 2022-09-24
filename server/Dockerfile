FROM golang:alpine

WORKDIR /app/stock-server

COPY go.mod .
COPY go.sum .
ENV GOPATH=/
RUN go mod download

#build appliction
COPY . .
RUN go build -o stock-server ./cmd/main/app.go

CMD ["./stock-server"]