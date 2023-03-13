FROM golang:1.20-alpine

WORKDIR /app

COPY ./ ./

RUN apk add --update --no-cache git make
RUN make build
RUN go build

EXPOSE 8000
ENTRYPOINT ["/app/simple-memorizer"]
