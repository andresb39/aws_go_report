FROM golang:1.16-alpine

LABEL Version="0.1" \
        Autor="J. Andres Bergano" \
        Date="Oct 2022"

WORKDIR /app

COPY . ./

RUN go mod download


RUN go build -o /aws_cost_export

CMD [ "/aws_cost_export" ]