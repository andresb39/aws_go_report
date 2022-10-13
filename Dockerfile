FROM golang:1.16-alpine

LABEL Version="0.1" \
        Autor="J. Andres Bergano" \
        Date="Oct 2022"

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /aws_cost_export

CMD [ "/aws_cost_export" ]
