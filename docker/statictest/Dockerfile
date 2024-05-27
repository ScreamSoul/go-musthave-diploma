FROM golang:1.22

WORKDIR /app

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.59.0

COPY docker/statictest/entrypoint.sh /
RUN chmod +x /entrypoint.sh

COPY docker/statictest/statictest /
RUN chmod +x /statictest

ENTRYPOINT [ "/entrypoint.sh" ]