FROM golang:1.19-alpine3.16 as builder

ENV APP_NAME=companysvc
WORKDIR /${APP_NAME}

COPY ./api ./api
COPY ./config ./config
COPY ./db ./db
COPY ./domain ./domain
COPY ./service ./service
COPY ./main.go .
COPY ./go.mod .
COPY ./go.sum .

RUN go build -v -o ${APP_NAME} .

FROM alpine:3.16
COPY --from=builder /companysvc/companysvc /bin/
COPY ./config.yaml /etc/companysvc/config.yaml
COPY ./migrations/ /etc/migrations/

ENV CONFIG="/etc/companysvc/config"

CMD "companysvc"
