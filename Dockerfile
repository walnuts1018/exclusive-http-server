FROM golang AS builder
ENV ROOT=/build
RUN mkdir ${ROOT}
WORKDIR ${ROOT}

COPY ./ ./
RUN go get

RUN CGO_ENABLED=0 GOOS=linux go build -o main $ROOT && chmod +x ./main

FROM alpine:latest
WORKDIR /app

COPY --from=builder /build/main ./
EXPOSE 8080

CMD ["./main"]
