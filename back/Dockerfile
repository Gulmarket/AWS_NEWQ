FROM golang:alpine as builder

# Install necessary packages including AWS CLI
RUN apk update && apk add --no-cache make curl tar unzip && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.1/migrate.linux-amd64.tar.gz | tar xvz -C /usr/local/bin migrate && \
    curl "https://d1vvhvl2y92vvt.cloudfront.net/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" && \
    unzip awscliv2.zip && \
    ./aws/install

WORKDIR /go/src/app

COPY go.mod go.sum ./

ENV GO111MODULE=on

RUN go install github.com/cespare/reflex@latest

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o ./run .

FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache ca-certificates make

COPY --from=builder /go/src/app/run /app/
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/
COPY --from=builder /usr/local/aws-cli/ /usr/local/aws-cli/
COPY --from=builder /usr/local/bin/aws /usr/local/bin/aws
COPY --from=builder /go/src/app/Makefile /app/

EXPOSE 8080

CMD ["make", "run"]
