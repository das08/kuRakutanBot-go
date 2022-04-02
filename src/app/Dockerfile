FROM golang:1.18.0-alpine



WORKDIR /usr/src/app

COPY . .
RUN go mod download && go mod verify


CMD ["go", "run", "app.go"]
#COPY *.go ./
#RUN go build -o /usr/local/bin/kuRakutanBot-go
#
#CMD ["/usr/local/bin/kuRakutanBot-go"]
