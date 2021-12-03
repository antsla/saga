FROM golang:1.17-alpine as builder

ARG SERVICE_NAME

COPY ./${SERVICE_NAME} /app/${SERVICE_NAME}
WORKDIR /app/${SERVICE_NAME}

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.0
RUN apk add --no-cache nano bash postgresql-client shadow build-base

RUN mkdir /.cache
RUN chown nobody:nobody -R /.cache
RUN chown nobody:nobody -R /go/pkg

USER nobody