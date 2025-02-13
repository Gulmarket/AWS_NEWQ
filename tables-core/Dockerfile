FROM nexus.sberbank.kz:5001/golang:latest AS builder
ARG CERT_LOCATION=/usr/local/share/ca-certificates
RUN openssl s_client -showcerts -connect nexus.sberbank.kz:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${CERT_LOCATION}/nexus.sberbank.crt
RUN openssl s_client -showcerts -connect bitbucket.sberbank.kz:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${CERT_LOCATION}/bitbucket.sberbank.crt
RUN update-ca-certificates
RUN go env -w GOPROXY=https://nexus.sberbank.kz/repository/golang
RUN go env -w 'GOSUMDB=sum.golang.org https://nexus.sberbank.kz/repository/gosum'
RUN go env -w GONOSUMDB=bitbucket.sberbank.kz/*
RUN go env -w GONOPROXY=bitbucket.sberbank.kz/*
WORKDIR /app
COPY . .
RUN GO111MODULE="on" CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd

FROM nexus.sberbank.kz:5001/alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
COPY --from=builder /app/internal/config /app/internal/config
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
COPY ./Almaty .
ENV TZ=Asia/Almaty
ENV ZONEINFO=/zoneinfo.zip
CMD ["./app"]
