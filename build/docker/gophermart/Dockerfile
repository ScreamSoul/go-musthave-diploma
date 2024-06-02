FROM golang:1.22 AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0

RUN apt-get update && apt-get install tzdata

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN cd cmd/gophermart && go build -buildvcs=false  -o /app/main

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Russia/Moskow

WORKDIR /app
COPY --from=builder /app/main /app/main

CMD ["./main"]
