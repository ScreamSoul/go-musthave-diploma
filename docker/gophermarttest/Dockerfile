FROM golang:1.22

WORKDIR /app

RUN curl -sSfL -o /gophermarttest https://github.com/Yandex-Practicum/go-autotests/releases/download/v0.10.9/gophermarttest
RUN chmod +x /gophermarttest

COPY docker/gophermarttest/entrypoint.sh /
RUN chmod +x /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]