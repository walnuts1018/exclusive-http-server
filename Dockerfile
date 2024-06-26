FROM golang AS builder
ENV ROOT=/build
RUN mkdir ${ROOT}
WORKDIR ${ROOT}

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main $ROOT && chmod +x ./main

FROM alpine:latest
WORKDIR /app

COPY --from=builder /build/main ./
EXPOSE 8080 8081

CMD ["./main"]
