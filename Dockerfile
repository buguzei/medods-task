FROM golang:latest

COPY ./ ./

RUN go build -o ./bin/medods-task ./cmd

CMD ["./bin/medods-task"]