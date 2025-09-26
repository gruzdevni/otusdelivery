FROM golang:1.24
WORKDIR /app
COPY . .
RUN go mod download
COPY *.go ./

EXPOSE 8080

RUN CGO_ENABLED=0 GOOS=linux go build -o /otusdelivery
CMD [ "/otusdelivery", "rest" ]