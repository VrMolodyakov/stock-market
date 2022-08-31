FROM golang:alpine

WORKDIR /app/jwt-auth

COPY go.mod .
COPY go.sum .
ENV GOPATH=/
RUN go mod download

#build appliction
COPY . .
RUN go build -o jwt-auth ./cmd/main/app.go

CMD ["./jwt-auth"]