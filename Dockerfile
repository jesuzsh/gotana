FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY cmd ./cmd/
COPY gotana-client ./gotana-client/

RUN go build ./cmd/gotana/main.go

CMD [ "./main" ]
