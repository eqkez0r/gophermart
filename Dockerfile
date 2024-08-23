FROM golang:1.22.1-alpine AS builder
LABEL authors="Kozlyk-VA"

WORKDIR /usr/local/src

RUN apk --no-cache add bash git go-task gcc gettext musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go build -o ./build ./cmd/gophermart

FROM alpine AS gophermart_runner

COPY --from=builder /usr/local/src/build/gophermart /

CMD ["./gophermart"]

FROM alpine AS accrual_runner

COPY ./cmd/accrual/accrual_linux_amd64 /

CMD ["./accrual_linux_amd64"]
