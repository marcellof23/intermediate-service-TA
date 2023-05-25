FROM golang:1.18-alpine AS BuildStage

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN  go build -o ./tmp/main ./cmd/intermediate

EXPOSE 8080

CMD [ "./tmp/main", "api" ]

FROM alpine:latest
WORKDIR /
COPY --from=BuildStage /app/tmp/main /main
COPY --from=BuildStage /app/config.yaml /config.yaml

EXPOSE 8080

CMD ["/main", "api"]