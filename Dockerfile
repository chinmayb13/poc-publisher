# syntax=docker/dockerfile:1

FROM golang:1.16 as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server/
RUN echo "hi docker" 
RUN echo $(ls -lrt ./)
FROM alpine
RUN apk add --no-cache ca-certificates
#COPY --from=builder /app/main /main
CMD ["./main"]