FROM golang:1.18.0-alpine

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -o /usr/local/bin/kuRakutanBot-go
CMD ["/usr/local/bin/kuRakutanBot-go"]
#CMD ["go", "run", "app.go"]
