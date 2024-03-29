FROM golang:1.20-alpine

WORKDIR /app

COPY ./ ./

RUN apk add --update --no-cache git make
RUN make build
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2

EXPOSE 8000
