FROM golang:1.21 as service-post-notification

WORKDIR /project

COPY go.mod .
RUN go mod download

COPY . /project
RUN go build -o /bin/service-post-notification -v ./cmd/service

RUN rm -rf /project

CMD ["/bin/service-post-notification"]